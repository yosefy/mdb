package kmedia

import (
	"database/sql"
	"runtime/debug"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/queries/qm"

	"github.com/Bnei-Baruch/mdb/api"
	"github.com/Bnei-Baruch/mdb/importer/kmedia/kmodels"
	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
)

func ImportimportFlatCatalogs() {
	clock := Init()

	stats = NewImportStatistics()
	utils.Must(importFlatCatalog(120, api.CT_FRIENDS_GATHERING))
	utils.Must(importFlatCatalog(4791, api.CT_MEAL))
	stats.dump()

	Shutdown()
	log.Info("Success")
	log.Infof("Total run time: %s", time.Now().Sub(clock).String())
}

func importFlatCatalog(catalogID int, cuType string) error {
	catalog, err := kmodels.Catalogs(kmdb,
		qm.Where("id=?", catalogID),
		qm.Load("Containers")).One()
	if err != nil {
		return errors.Wrapf(err, "Load catalog %d", catalogID)
	}

	stats.CatalogsProcessed.Inc(1)

	for i := range catalog.R.Containers {
		tx, err := mdb.Begin()
		utils.Must(err)

		if err = importContainerWOCollection(tx, catalog.R.Containers[i], cuType); err != nil {
			utils.Must(tx.Rollback())
			stats.TxRolledBack.Inc(1)
			log.Error(err)
			debug.PrintStack()
			continue
		} else {
			utils.Must(tx.Commit())
			stats.TxCommitted.Inc(1)
		}
	}

	return nil
}

func importContainerWOCollection(exec boil.Executor, container *kmodels.Container, cuType string) error {
	stats.ContainersProcessed.Inc(1)

	unit, err := models.ContentUnits(mdb, qm.Where("(properties->>'kmedia_id')::int = ?", container.ID)).One()
	if err != nil {
		if err == sql.ErrNoRows {
			log.Infof("New CU %d %s", container.ID, container.Name.String)
			return importContainerWOCollectionNewCU(exec, container, cuType)
		}
		return errors.Wrapf(err, "Lookup content unit kmid %s", container.ID)
	}

	log.Infof("CU exists [%d] container: %s %d", unit.ID, container.Name.String, container.ID)
	_, err = importContainer(exec, container, nil, cuType, "", 0)
	if err != nil {
		return errors.Wrapf(err, "Import container %d", container.ID)
	}

	if cuType != api.CONTENT_TYPE_REGISTRY.ByID[unit.TypeID].Name {
		log.Infof("Overriding CU Type to %s", cuType)
		unit.TypeID = api.CONTENT_TYPE_REGISTRY.ByName[cuType].ID
		err = unit.Update(exec, "type_id")
		if err != nil {
			return errors.Wrapf(err, "Update CU type %d", unit.ID)
		}
	}

	return nil
}

func importContainerWOCollectionNewCU(exec boil.Executor, container *kmodels.Container, cuType string) error {
	err := container.L.LoadFileAssets(kmdb, true, container)
	if err != nil {
		return errors.Wrapf(err, "Load kmedia file assets %d", container.ID)
	}

	// Create import operation
	operation, err := api.CreateOperation(exec, api.OP_IMPORT_KMEDIA,
		api.Operation{WorkflowID: strconv.Itoa(container.ID)}, nil)
	if err != nil {
		return errors.Wrapf(err, "Create operation %d", container.ID)
	}
	stats.OperationsCreated.Inc(1)

	// import container
	unit, err := importContainer(exec, container, nil, cuType, "", 0)
	if err != nil {
		return errors.Wrapf(err, "Import container %d", container.ID)
	}

	// import container files
	var file *models.File
	for _, fileAsset := range container.R.FileAssets {
		log.Infof("Processing file_asset %d", fileAsset.ID)
		stats.FileAssetsProcessed.Inc(1)

		// Create or update MDB file
		file, err = importFileAsset(exec, fileAsset, unit, operation)
		if err != nil {
			log.Error(err)
			debug.PrintStack()
			break
		}
		if file != nil && file.Published {
			unit.Published = true
		}
	}
	if err != nil {
		return errors.Wrapf(err, "Import container files %d", container.ID)
	}

	if unit.Published {
		err = unit.Update(exec, "published")
		if err != nil {
			return errors.Wrapf(err, "Update unit published column %d", container.ID)
		}
	}

	return nil
}
