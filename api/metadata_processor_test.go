package api

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/queries/qm"
	"gopkg.in/nullbio/null.v6"

	"github.com/Bnei-Baruch/mdb/models"
	"github.com/Bnei-Baruch/mdb/utils"
)

type MetadataProcessorSuite struct {
	suite.Suite
	utils.TestDBManager
	tx boil.Transactor
}

func (suite *MetadataProcessorSuite) SetupSuite() {
	suite.Require().Nil(suite.InitTestDB())
	suite.Require().Nil(OPERATION_TYPE_REGISTRY.Init())
	suite.Require().Nil(CONTENT_TYPE_REGISTRY.Init())
}

func (suite *MetadataProcessorSuite) TearDownSuite() {
	suite.Require().Nil(suite.DestroyTestDB())
}

func (suite *MetadataProcessorSuite) SetupTest() {
	var err error
	suite.tx, err = boil.Begin()
	suite.Require().Nil(err)
}

func (suite *MetadataProcessorSuite) TearDownTest() {
	err := suite.tx.Rollback()
	suite.Require().Nil(err)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestMetadataProcessor(t *testing.T) {
	suite.Run(t, new(MetadataProcessorSuite))
}

func (suite *MetadataProcessorSuite) TestDailyLesson() {
	chain := suite.simulateLessonChain()

	// send parts
	// send full
	// send kitei makor of part 1

	metadata := CITMetadata{
		ContentType:    CT_LESSON_PART,
		AutoName:       "auto_name",
		FinalName:      "final_name",
		CaptureDate:    Date{time.Now()},
		Language:       LANG_HEBREW,
		HasTranslation: true,
		Lecturer:       "rav",
		Number:         null.IntFrom(1),
		Part:           null.IntFrom(0),
		Sources:        suite.someSources(),
		Tags:           suite.someTags(),
		RequireTest:    false,
	}
	original, proxy := chain["part0"].Original, chain["part0"].Proxy

	err := ProcessCITMetadata(suite.tx, metadata, original, proxy)
	suite.Require().Nil(err)

	err = original.Reload(suite.tx)
	suite.Require().Nil(err)
	err = proxy.Reload(suite.tx)
	suite.Require().Nil(err)

	suite.assertFiles(metadata, original, proxy)
	suite.assertContentUnit(metadata, original, proxy)

	// collection association
	err = original.L.LoadContentUnit(suite.tx, true, original)
	suite.Require().Nil(err)
	cu := original.R.ContentUnit
	err = cu.L.LoadCollectionsContentUnits(suite.tx, true, cu)
	suite.Require().Nil(err)
	suite.Equal(1, len(cu.R.CollectionsContentUnits), "len(cu.R.CollectionsContentUnits)")
	ccu := cu.R.CollectionsContentUnits[0]
	suite.Equal("0", ccu.Name, "ccu.Name")

	// collection
	err = ccu.L.LoadCollection(suite.tx, true, ccu)
	suite.Require().Nil(err)
	c := ccu.R.Collection
	suite.Equal(CONTENT_TYPE_REGISTRY.ByName[CT_DAILY_LESSON].ID, c.TypeID, "c.TypeID")
	suite.True(c.Properties.Valid, "c.Properties.Valid")
	var props map[string]interface{}
	err = json.Unmarshal(c.Properties.JSON, &props)
	suite.Require().Nil(err)
	suite.Equal(metadata.CaptureDate.Format(time.RFC3339Nano), props["capture_date"], "c.Properties[\"capture_date\"]")
	suite.Equal(metadata.CaptureDate.Format(time.RFC3339Nano), props["film_date"], "c.Properties[\"film_date\"]")
	suite.Equal("c12356789", props["capture_id"], "c.Properties[\"capture_id\"]")
	suite.EqualValues(metadata.Number.Int, props["number"], "c.Properties[\"number\"]")

	// process other parts
	for i := 1; i < 4; i++ {
		metadata.Part = null.IntFrom(i)
		metadata.Sources = suite.someSources()
		metadata.Tags = suite.someTags()
		tf := chain[fmt.Sprintf("part%d", i)]
		original, proxy := tf.Original, tf.Proxy

		err := ProcessCITMetadata(suite.tx, metadata, original, proxy)
		suite.Require().Nil(err)

		err = original.Reload(suite.tx)
		suite.Require().Nil(err)
		err = proxy.Reload(suite.tx)
		suite.Require().Nil(err)

		suite.assertFiles(metadata, original, proxy)
		suite.assertContentUnit(metadata, original, proxy)

		// collection association
		err = original.L.LoadContentUnit(suite.tx, true, original)
		suite.Require().Nil(err)
		cu := original.R.ContentUnit
		err = cu.L.LoadCollectionsContentUnits(suite.tx, true, cu)
		suite.Require().Nil(err)
		suite.Equal(1, len(cu.R.CollectionsContentUnits), "len(cu.R.CollectionsContentUnits)")
		ccu := cu.R.CollectionsContentUnits[0]
		suite.Equal(strconv.Itoa(i), ccu.Name, "ccu.Name")
		suite.Equal(c.ID, ccu.CollectionID, "ccu.CollectionID")
	}

	// process full
	metadata.ContentType = CT_FULL_LESSON
	metadata.Part = null.NewInt(0, false)
	metadata.Sources = nil
	metadata.Tags = nil
	tf := chain["full"]
	original, proxy = tf.Original, tf.Proxy

	err = ProcessCITMetadata(suite.tx, metadata, original, proxy)
	suite.Require().Nil(err)

	err = original.Reload(suite.tx)
	suite.Require().Nil(err)
	err = proxy.Reload(suite.tx)
	suite.Require().Nil(err)

	suite.assertFiles(metadata, original, proxy)
	suite.assertContentUnit(metadata, original, proxy)

	// collection association
	err = original.L.LoadContentUnit(suite.tx, true, original)
	suite.Require().Nil(err)
	cu = original.R.ContentUnit
	err = cu.L.LoadCollectionsContentUnits(suite.tx, true, cu)
	suite.Require().Nil(err)
	suite.Equal(1, len(cu.R.CollectionsContentUnits), "len(cu.R.CollectionsContentUnits)")
	ccu = cu.R.CollectionsContentUnits[0]
	suite.Equal("full", ccu.Name, "ccu.Name")
	suite.Equal(c.ID, ccu.CollectionID, "ccu.CollectionID")

	// full with week_date different from capture_date
	metadata.WeekDate = &Date{Time: time.Now().AddDate(1, 0, 0)}
	err = ProcessCITMetadata(suite.tx, metadata, original, proxy)
	suite.Require().Nil(err)
	err = c.Reload(suite.tx)
	suite.Require().Nil(err)
	suite.Equal(CONTENT_TYPE_REGISTRY.ByName[CT_SATURDAY_LESSON].ID, c.TypeID, "c.TypeID")
	suite.True(c.Properties.Valid, "c.Properties.Valid")
	err = json.Unmarshal(c.Properties.JSON, &props)
	suite.Require().Nil(err)
	suite.Equal(metadata.CaptureDate.Format(time.RFC3339Nano), props["capture_date"], "c.Properties[\"capture_date\"]")
	suite.Equal(metadata.WeekDate.Format(time.RFC3339Nano), props["film_date"], "c.Properties[\"film_date\"]")
	suite.Equal("c12356789", props["capture_id"], "c.Properties[\"capture_id\"]")
	suite.EqualValues(metadata.Number.Int, props["number"], "c.Properties[\"number\"]")

	// process kitei makor for part 1
	metadata.ContentType = CT_LESSON_PART
	metadata.Part = null.IntFrom(1)
	metadata.ArtifactType = null.StringFrom("kitei_makor")
	metadata.WeekDate = nil
	tf = chain["part1_kitei-makor"]
	original, proxy = tf.Original, tf.Proxy
	err = ProcessCITMetadata(suite.tx, metadata, original, proxy)
	suite.Require().Nil(err)

	err = original.Reload(suite.tx)
	suite.Require().Nil(err)
	err = proxy.Reload(suite.tx)
	suite.Require().Nil(err)

	suite.assertFiles(metadata, original, proxy)
	suite.assertContentUnit(metadata, original, proxy)

	// not associated with collection
	err = original.L.LoadContentUnit(suite.tx, true, original)
	suite.Require().Nil(err)
	cu = original.R.ContentUnit
	err = cu.L.LoadCollectionsContentUnits(suite.tx, true, cu)
	suite.Require().Nil(err)
	suite.Empty(cu.R.CollectionsContentUnits, "cu.R.CollectionsContentUnits empty")

	// no changes to collection
	err = c.Reload(suite.tx)
	suite.Require().Nil(err)
	suite.Equal(CONTENT_TYPE_REGISTRY.ByName[CT_SATURDAY_LESSON].ID, c.TypeID, "c.TypeID")
	suite.True(c.Properties.Valid, "c.Properties.Valid")
	err = json.Unmarshal(c.Properties.JSON, &props)
	suite.Require().Nil(err)
	suite.Equal(metadata.CaptureDate.Format(time.RFC3339Nano), props["capture_date"], "c.Properties[\"capture_date\"]")
	suite.NotEqual(metadata.CaptureDate.Format(time.RFC3339Nano), props["film_date"], "c.Properties[\"film_date\"]")
	suite.Equal("c12356789", props["capture_id"], "c.Properties[\"capture_id\"]")
	suite.EqualValues(metadata.Number.Int, props["number"], "c.Properties[\"number\"]")
}

func (suite *MetadataProcessorSuite) TestVideoProgram() {
	tf := suite.simulateSimpleChain()
	original, proxy := tf.Original, tf.Proxy

	c, err := CreateCollection(suite.tx, CT_VIDEO_PROGRAM, nil)
	suite.Require().Nil(err)

	metadata := CITMetadata{
		ContentType:    CT_VIDEO_PROGRAM_CHAPTER,
		AutoName:       "auto_name",
		FinalName:      "final_name",
		CaptureDate:    Date{time.Now()},
		Language:       LANG_HEBREW,
		HasTranslation: false,
		Lecturer:       "norav",
		CollectionUID:  null.StringFrom(c.UID),
		Episode:        null.StringFrom("827"),
		RequireTest:    true,
	}

	err = ProcessCITMetadata(suite.tx, metadata, original, proxy)
	suite.Require().Nil(err)

	err = original.Reload(suite.tx)
	suite.Require().Nil(err)
	err = proxy.Reload(suite.tx)
	suite.Require().Nil(err)

	suite.assertFiles(metadata, original, proxy)
	suite.assertContentUnit(metadata, original, proxy)

	// collection association
	err = original.L.LoadContentUnit(suite.tx, true, original)
	suite.Require().Nil(err)
	cu := original.R.ContentUnit
	err = cu.L.LoadCollectionsContentUnits(suite.tx, true, cu)
	suite.Require().Nil(err)
	suite.Equal(1, len(cu.R.CollectionsContentUnits), "len(cu.R.CollectionsContentUnits)")
	ccu := cu.R.CollectionsContentUnits[0]
	suite.Equal(metadata.Episode.String, ccu.Name, "ccu.Name")

}

// Helpers

type TrimFiles struct {
	Original *models.File
	Proxy    *models.File
}

func (suite *MetadataProcessorSuite) simulateSimpleChain() TrimFiles {
	CS_SHA1 := utils.RandomSHA1()
	DMX_O_SHA1 := utils.RandomSHA1()
	DMX_P_SHA1 := utils.RandomSHA1()
	TRM_O_SHA1 := utils.RandomSHA1()
	TRM_P_SHA1 := utils.RandomSHA1()

	// capture_start
	_, err := handleCaptureStart(suite.tx, CaptureStartRequest{
		Operation: Operation{
			Station:    "Capture station",
			User:       "operator@dev.com",
			WorkflowID: "c12356788",
		},
		FileName:      "capture_start_simple",
		CaptureSource: "mltcap",
		CollectionUID: "c12356788",
	})
	suite.Require().Nil(err)

	// capture_stop
	_, err = handleCaptureStop(suite.tx, CaptureStopRequest{
		Operation: Operation{
			Station:    "Capture station",
			User:       "operator@dev.com",
			WorkflowID: "c12356789",
		},
		File: File{
			FileName:  "capture_stop_simple.mp4",
			Sha1:      CS_SHA1,
			Size:      98737,
			CreatedAt: &Timestamp{Time: time.Now()},
			Language:  LANG_HEBREW,
		},
		CaptureSource: "mltcap",
		CollectionUID: "c12356789",
		Part:          "false",
	})
	suite.Require().Nil(err)

	// demux
	_, err = handleDemux(suite.tx, DemuxRequest{
		Operation: Operation{
			Station: "Trimmer station",
			User:    "operator@dev.com",
		},
		Sha1: CS_SHA1,
		Original: AVFile{
			File: File{
				FileName:  "demux_simple_original.mp4",
				Sha1:      DMX_O_SHA1,
				Size:      98737,
				CreatedAt: &Timestamp{Time: time.Now()},
			},
			Duration: 892.1900,
		},
		Proxy: AVFile{
			File: File{
				FileName:  "demux_simple_proxy.mp4",
				Sha1:      DMX_P_SHA1,
				Size:      9878,
				CreatedAt: &Timestamp{Time: time.Now()},
			},
			Duration: 891.8800,
		},
		CaptureSource: "mltcap",
	})
	suite.Require().Nil(err)

	// trim
	op, err := handleTrim(suite.tx, TrimRequest{
		Operation: Operation{
			Station: "Trimmer station",
			User:    "operator@dev.com",
		},
		OriginalSha1: DMX_O_SHA1,
		ProxySha1:    DMX_P_SHA1,
		Original: AVFile{
			File: File{
				FileName:  "trim_simple_original.mp4",
				Sha1:      TRM_O_SHA1,
				Size:      98000,
				CreatedAt: &Timestamp{Time: time.Now()},
			},
			Duration: 892.1900,
		},
		Proxy: AVFile{
			File: File{
				FileName:  "trim_simple_proxy.mp4",
				Sha1:      TRM_P_SHA1,
				Size:      9800,
				CreatedAt: &Timestamp{Time: time.Now()},
			},
			Duration: 891.8800,
		},
		CaptureSource: "mltcap",
		In:            []float64{10.05, 249.43},
		Out:           []float64{240.51, 899.27},
	})
	suite.Require().Nil(err)

	files := suite.opFilesBySHA1(op)
	return TrimFiles{
		Original: files[TRM_O_SHA1],
		Proxy:    files[TRM_P_SHA1],
	}
}

func (suite *MetadataProcessorSuite) simulateLessonChain() map[string]TrimFiles {
	CS_SHA1 := [5]string{}
	DMX_O_SHA1 := [5]string{}
	DMX_P_SHA1 := [5]string{}
	TRM_O_SHA1 := [6]string{}
	TRM_P_SHA1 := [6]string{}
	for i := range CS_SHA1 {
		CS_SHA1[i] = utils.RandomSHA1()
		DMX_O_SHA1[i] = utils.RandomSHA1()
		DMX_P_SHA1[i] = utils.RandomSHA1()
		TRM_O_SHA1[i] = utils.RandomSHA1()
		TRM_P_SHA1[i] = utils.RandomSHA1()
	}
	TRM_O_SHA1[5] = utils.RandomSHA1()
	TRM_P_SHA1[5] = utils.RandomSHA1()

	// capture_start
	_, err := handleCaptureStart(suite.tx, CaptureStartRequest{
		Operation: Operation{
			Station:    "Capture station",
			User:       "operator@dev.com",
			WorkflowID: "c12356789",
		},
		FileName:      "capture_start_full",
		CaptureSource: "mltbackup",
		CollectionUID: "c12356789",
	})
	suite.Require().Nil(err)

	for i := 0; i < 4; i++ {
		part := strconv.Itoa(i)
		_, err := handleCaptureStart(suite.tx, CaptureStartRequest{
			Operation: Operation{
				Station:    "Capture station",
				User:       "operator@dev.com",
				WorkflowID: "c" + strings.Repeat(part, 8),
			},
			FileName:      "capture_start_part" + part,
			CaptureSource: "mltcap",
			CollectionUID: "c12356789",
		})
		suite.Require().Nil(err)
	}

	// capture_stop
	_, err = handleCaptureStop(suite.tx, CaptureStopRequest{
		Operation: Operation{
			Station:    "Capture station",
			User:       "operator@dev.com",
			WorkflowID: "c12356789",
		},
		File: File{
			FileName:  "capture_stop_full.mp4",
			Sha1:      CS_SHA1[4],
			Size:      98737,
			CreatedAt: &Timestamp{Time: time.Now()},
			Language:  LANG_MULTI,
		},
		CaptureSource: "mltbackup",
		CollectionUID: "c12356789",
		Part:          "full",
	})
	suite.Require().Nil(err)

	for i := 0; i < 4; i++ {
		part := strconv.Itoa(i)
		_, err := handleCaptureStop(suite.tx, CaptureStopRequest{
			Operation: Operation{
				Station:    "Capture station",
				User:       "operator@dev.com",
				WorkflowID: "c" + strings.Repeat(part, 8),
			},
			File: File{
				FileName:  "capture_stop_part" + part + ".mp4",
				Sha1:      CS_SHA1[i],
				Size:      int64(i),
				CreatedAt: &Timestamp{Time: time.Now()},
				Language:  LANG_MULTI,
			},
			CaptureSource: "mltcap",
			CollectionUID: "c12356789",
			Part:          part,
		})
		suite.Require().Nil(err)
	}

	// demux
	_, err = handleDemux(suite.tx, DemuxRequest{
		Operation: Operation{
			Station: "Trimmer station",
			User:    "operator@dev.com",
		},
		Sha1: CS_SHA1[4],
		Original: AVFile{
			File: File{
				FileName:  "demux_full_original.mp4",
				Sha1:      DMX_O_SHA1[4],
				Size:      98737,
				CreatedAt: &Timestamp{Time: time.Now()},
			},
			Duration: 892.1900,
		},
		Proxy: AVFile{
			File: File{
				FileName:  "demux_full_proxy.mp4",
				Sha1:      DMX_P_SHA1[4],
				Size:      9878,
				CreatedAt: &Timestamp{Time: time.Now()},
			},
			Duration: 891.8800,
		},
		CaptureSource: "mltbackup",
	})
	suite.Require().Nil(err)

	for i := 0; i < 4; i++ {
		part := strconv.Itoa(i)
		_, err := handleDemux(suite.tx, DemuxRequest{
			Operation: Operation{
				Station: "Trimmer station",
				User:    "operator@dev.com",
			},
			Sha1: CS_SHA1[4],
			Original: AVFile{
				File: File{
					FileName:  "demux_part" + part + "_original.mp4",
					Sha1:      DMX_O_SHA1[i],
					Size:      98737,
					CreatedAt: &Timestamp{Time: time.Now()},
				},
				Duration: 892.1900,
			},
			Proxy: AVFile{
				File: File{
					FileName:  "demux_part" + part + "_proxy.mp4",
					Sha1:      DMX_P_SHA1[i],
					Size:      9878,
					CreatedAt: &Timestamp{Time: time.Now()},
				},
				Duration: 891.8800,
			},
			CaptureSource: "mltcap",
		})
		suite.Require().Nil(err)
	}

	trimFiles := make(map[string]TrimFiles)

	// trim
	op, err := handleTrim(suite.tx, TrimRequest{
		Operation: Operation{
			Station: "Trimmer station",
			User:    "operator@dev.com",
		},
		OriginalSha1: DMX_O_SHA1[4],
		ProxySha1:    DMX_P_SHA1[4],
		Original: AVFile{
			File: File{
				FileName:  "trim_full_original.mp4",
				Sha1:      TRM_O_SHA1[4],
				Size:      98000,
				CreatedAt: &Timestamp{Time: time.Now()},
			},
			Duration: 892.1900,
		},
		Proxy: AVFile{
			File: File{
				FileName:  "trim_full_proxy.mp4",
				Sha1:      TRM_P_SHA1[4],
				Size:      9800,
				CreatedAt: &Timestamp{Time: time.Now()},
			},
			Duration: 891.8800,
		},
		CaptureSource: "mltbackup",
		In:            []float64{10.05, 249.43},
		Out:           []float64{240.51, 899.27},
	})
	suite.Require().Nil(err)
	files := suite.opFilesBySHA1(op)
	trimFiles["full"] = TrimFiles{
		Original: files[TRM_O_SHA1[4]],
		Proxy:    files[TRM_P_SHA1[4]],
	}

	for i := 0; i < 4; i++ {
		part := strconv.Itoa(i)
		op, err := handleTrim(suite.tx, TrimRequest{
			Operation: Operation{
				Station: "Trimmer station",
				User:    "operator@dev.com",
			},
			OriginalSha1: DMX_O_SHA1[4],
			ProxySha1:    DMX_P_SHA1[4],
			Original: AVFile{
				File: File{
					FileName:  "trim_part" + part + "_original.mp4",
					Sha1:      TRM_O_SHA1[i],
					Size:      98000,
					CreatedAt: &Timestamp{Time: time.Now()},
				},
				Duration: 892.1900,
			},
			Proxy: AVFile{
				File: File{
					FileName:  "trim_part" + part + "_proxy.mp4",
					Sha1:      TRM_P_SHA1[i],
					Size:      9800,
					CreatedAt: &Timestamp{Time: time.Now()},
				},
				Duration: 891.8800,
			},
			CaptureSource: "mltbackup",
			In:            []float64{10.05, 249.43},
			Out:           []float64{240.51, 899.27},
		})
		suite.Require().Nil(err)
		files := suite.opFilesBySHA1(op)
		trimFiles["part"+part] = TrimFiles{
			Original: files[TRM_O_SHA1[i]],
			Proxy:    files[TRM_P_SHA1[i]],
		}
	}

	// trim kitei makor from part1
	op, err = handleTrim(suite.tx, TrimRequest{
		Operation: Operation{
			Station: "Trimmer station",
			User:    "operator@dev.com",
		},
		OriginalSha1: DMX_O_SHA1[1],
		ProxySha1:    DMX_P_SHA1[1],
		Original: AVFile{
			File: File{
				FileName:  "trim_part1_kitei_makor_original.mp4",
				Sha1:      TRM_O_SHA1[5],
				Size:      6700,
				CreatedAt: &Timestamp{Time: time.Now()},
			},
			Duration: 92.1900,
		},
		Proxy: AVFile{
			File: File{
				FileName:  "trim_part1_kitei_makor_proxy.mp4",
				Sha1:      TRM_P_SHA1[5],
				Size:      6700,
				CreatedAt: &Timestamp{Time: time.Now()},
			},
			Duration: 91.8800,
		},
		CaptureSource: "mltcap",
		In:            []float64{10.05, 249.43, 253.83, 312.23, 463.3, 512.3},
		Out:           []float64{240.51, 250.31, 282.13, 441.03, 483.39, 899.81},
	})
	suite.Require().Nil(err)
	files = suite.opFilesBySHA1(op)
	trimFiles["part1_kitei-makor"] = TrimFiles{
		Original: files[TRM_O_SHA1[5]],
		Proxy:    files[TRM_P_SHA1[5]],
	}

	return trimFiles
}

func (suite *MetadataProcessorSuite) opFilesBySHA1(o *models.Operation) map[string]*models.File {
	files := make(map[string]*models.File)
	for _, f := range o.R.Files {
		files[hex.EncodeToString(f.Sha1.Bytes)] = f
	}
	return files
}

func (suite *MetadataProcessorSuite) someSources() []string {
	items, err := models.Sources(suite.tx, qm.Limit(1+rand.Intn(10))).All()
	suite.Require().Nil(err)
	uids := make([]string, len(items))
	for i, x := range items {
		uids[i] = x.UID
	}
	return uids
}

func (suite *MetadataProcessorSuite) someTags() []string {
	items, err := models.Tags(suite.tx, qm.Limit(1+rand.Intn(10))).All()
	suite.Require().Nil(err)
	uids := make([]string, len(items))
	for i, x := range items {
		uids[i] = x.UID
	}
	return uids
}

func (suite *MetadataProcessorSuite) assertFiles(metadata CITMetadata, original, proxy *models.File) {
	capDate := metadata.CaptureDate
	filmDate := metadata.CaptureDate
	if metadata.WeekDate != nil {
		filmDate = *metadata.WeekDate
	}
	var lang string
	if metadata.HasTranslation {
		lang = LANG_MULTI
	} else {
		lang = StdLang(metadata.Language)
	}

	var props map[string]interface{}

	// original properties
	suite.Require().True(original.Properties.Valid)
	err := original.Properties.Unmarshal(&props)
	suite.Require().Nil(err)
	suite.Equal(capDate.Format(time.RFC3339Nano), props["capture_date"], "original.Properties[\"capture_date\"]")
	suite.Equal(filmDate.Format(time.RFC3339Nano), props["film_date"], "original.Properties[\"film_date\"]")

	// proxy properties
	suite.Require().True(proxy.Properties.Valid)
	err = proxy.Properties.Unmarshal(&props)
	suite.Require().Nil(err)
	suite.Equal(capDate.Format(time.RFC3339Nano), props["capture_date"], "proxy.Properties[\"capture_date\"]")
	suite.Equal(filmDate.Format(time.RFC3339Nano), props["film_date"], "proxy.Properties[\"film_date\"]")

	// original language
	suite.True(original.Language.Valid, "original.Language.Valid")
	suite.Equal(lang, original.Language.String, "original.Language")
}

func (suite *MetadataProcessorSuite) assertContentUnit(metadata CITMetadata, original, proxy *models.File) {
	// reload unit
	err := original.L.LoadContentUnit(suite.tx, true, original)
	suite.Require().Nil(err)
	cu := original.R.ContentUnit

	// properties
	suite.Equal(cu.TypeID, CONTENT_TYPE_REGISTRY.ByName[metadata.ContentType].ID, "cu.type_id")

	capDate := metadata.CaptureDate
	filmDate := metadata.CaptureDate
	if metadata.WeekDate != nil {
		filmDate = *metadata.WeekDate
	}
	suite.Require().True(cu.Properties.Valid)
	var props map[string]interface{}
	err = cu.Properties.Unmarshal(&props)
	suite.Require().Nil(err)
	suite.Equal(capDate.Format(time.RFC3339Nano), props["capture_date"], "cu.Properties[\"capture_date\"]")
	suite.Equal(filmDate.Format(time.RFC3339Nano), props["film_date"], "cu.Properties[\"film_date\"]")

	// files in unit
	err = cu.L.LoadFiles(suite.tx, true, cu)
	suite.Require().Nil(err)
	suite.True(original.ContentUnitID.Valid, "original.ContentUnitID.Valid")
	suite.True(proxy.ContentUnitID.Valid, "proxy.ContentUnitID.Valid")
	suite.Equal(original.ContentUnitID.Int64, proxy.ContentUnitID.Int64, "original.cuid = proxy.cuid")

	// ancestors ?
	if metadata.ArtifactType.Valid && metadata.ArtifactType.String != "main" {
		suite.Equal(2, len(cu.R.Files), "len(cu.R.Files)")
	} else {
		ancestors, err := FindFileAncestors(suite.tx, original.ID)
		suite.Require().Nil(err)
		suite.Equal(2+len(ancestors), len(cu.R.Files), "len(cu.R.Files)")
		for i, f := range ancestors {
			suite.True(f.ContentUnitID.Valid, "Ancestor[%d].ContentUnitID.Valid", i)
			suite.Equal(original.ContentUnitID.Int64, f.ContentUnitID.Int64, "Ancestor[%d]ContentUnitID.Int64", i)
		}
	}

	// sources
	err = cu.L.LoadSources(suite.tx, true, cu)
	suite.Require().Nil(err)
	suite.Equal(len(metadata.Sources), len(cu.R.Sources), "len(cu.R.Sources)")
	for _, x := range metadata.Sources {
		missing := true
		for _, y := range cu.R.Sources {
			if x == y.UID {
				missing = false
				break
			}
		}
		suite.False(missing, "Missing source %s", x)
	}

	// tags
	err = cu.L.LoadTags(suite.tx, true, cu)
	suite.Require().Nil(err)
	suite.Equal(len(metadata.Tags), len(cu.R.Tags), "len(cu.R.Tags)")
	for _, x := range metadata.Tags {
		missing := true
		for _, y := range cu.R.Tags {
			if x == y.UID {
				missing = false
				break
			}
		}
		suite.False(missing, "Missing tag %s", x)
	}
}
