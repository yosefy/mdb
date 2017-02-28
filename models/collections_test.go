package models

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/vattle/sqlboiler/boil"
	"github.com/vattle/sqlboiler/randomize"
	"github.com/vattle/sqlboiler/strmangle"
)

func testCollections(t *testing.T) {
	t.Parallel()

	query := Collections(nil)

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}
func testCollectionsDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	collection := &Collection{}
	if err = randomize.Struct(seed, collection, collectionDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Collection struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = collection.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = collection.Delete(tx); err != nil {
		t.Error(err)
	}

	count, err := Collections(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testCollectionsQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	collection := &Collection{}
	if err = randomize.Struct(seed, collection, collectionDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Collection struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = collection.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = Collections(tx).DeleteAll(); err != nil {
		t.Error(err)
	}

	count, err := Collections(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testCollectionsSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	collection := &Collection{}
	if err = randomize.Struct(seed, collection, collectionDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Collection struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = collection.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := CollectionSlice{collection}

	if err = slice.DeleteAll(tx); err != nil {
		t.Error(err)
	}

	count, err := Collections(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}
func testCollectionsExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	collection := &Collection{}
	if err = randomize.Struct(seed, collection, collectionDBTypes, true, collectionColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Collection struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = collection.Insert(tx); err != nil {
		t.Error(err)
	}

	e, err := CollectionExists(tx, collection.ID)
	if err != nil {
		t.Errorf("Unable to check if Collection exists: %s", err)
	}
	if !e {
		t.Errorf("Expected CollectionExistsG to return true, but got false.")
	}
}
func testCollectionsFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	collection := &Collection{}
	if err = randomize.Struct(seed, collection, collectionDBTypes, true, collectionColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Collection struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = collection.Insert(tx); err != nil {
		t.Error(err)
	}

	collectionFound, err := FindCollection(tx, collection.ID)
	if err != nil {
		t.Error(err)
	}

	if collectionFound == nil {
		t.Error("want a record, got nil")
	}
}
func testCollectionsBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	collection := &Collection{}
	if err = randomize.Struct(seed, collection, collectionDBTypes, true, collectionColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Collection struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = collection.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = Collections(tx).Bind(collection); err != nil {
		t.Error(err)
	}
}

func testCollectionsOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	collection := &Collection{}
	if err = randomize.Struct(seed, collection, collectionDBTypes, true, collectionColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Collection struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = collection.Insert(tx); err != nil {
		t.Error(err)
	}

	if x, err := Collections(tx).One(); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testCollectionsAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	collectionOne := &Collection{}
	collectionTwo := &Collection{}
	if err = randomize.Struct(seed, collectionOne, collectionDBTypes, false, collectionColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Collection struct: %s", err)
	}
	if err = randomize.Struct(seed, collectionTwo, collectionDBTypes, false, collectionColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Collection struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = collectionOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = collectionTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := Collections(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testCollectionsCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	collectionOne := &Collection{}
	collectionTwo := &Collection{}
	if err = randomize.Struct(seed, collectionOne, collectionDBTypes, false, collectionColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Collection struct: %s", err)
	}
	if err = randomize.Struct(seed, collectionTwo, collectionDBTypes, false, collectionColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Collection struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = collectionOne.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = collectionTwo.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Collections(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func testCollectionsInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	collection := &Collection{}
	if err = randomize.Struct(seed, collection, collectionDBTypes, true, collectionColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Collection struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = collection.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Collections(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testCollectionsInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	collection := &Collection{}
	if err = randomize.Struct(seed, collection, collectionDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Collection struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = collection.Insert(tx, collectionColumns...); err != nil {
		t.Error(err)
	}

	count, err := Collections(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testCollectionToManyCollectionsContentUnits(t *testing.T) {
	var err error
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Collection
	var b, c CollectionsContentUnit

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, collectionDBTypes, true, collectionColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Collection struct: %s", err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	randomize.Struct(seed, &b, collectionsContentUnitDBTypes, false, collectionsContentUnitColumnsWithDefault...)
	randomize.Struct(seed, &c, collectionsContentUnitDBTypes, false, collectionsContentUnitColumnsWithDefault...)

	b.CollectionID = a.ID
	c.CollectionID = a.ID
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx); err != nil {
		t.Fatal(err)
	}

	collectionsContentUnit, err := a.CollectionsContentUnits(tx).All()
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range collectionsContentUnit {
		if v.CollectionID == b.CollectionID {
			bFound = true
		}
		if v.CollectionID == c.CollectionID {
			cFound = true
		}
	}

	if !bFound {
		t.Error("expected to find b")
	}
	if !cFound {
		t.Error("expected to find c")
	}

	slice := CollectionSlice{&a}
	if err = a.L.LoadCollectionsContentUnits(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.CollectionsContentUnits); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.CollectionsContentUnits = nil
	if err = a.L.LoadCollectionsContentUnits(tx, true, &a); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.CollectionsContentUnits); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", collectionsContentUnit)
	}
}

func testCollectionToManyAddOpCollectionsContentUnits(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Collection
	var b, c, d, e CollectionsContentUnit

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, collectionDBTypes, false, strmangle.SetComplement(collectionPrimaryKeyColumns, collectionColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*CollectionsContentUnit{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, collectionsContentUnitDBTypes, false, strmangle.SetComplement(collectionsContentUnitPrimaryKeyColumns, collectionsContentUnitColumnsWithoutDefault)...); err != nil {
			t.Fatal(err)
		}
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx); err != nil {
		t.Fatal(err)
	}

	foreignersSplitByInsertion := [][]*CollectionsContentUnit{
		{&b, &c},
		{&d, &e},
	}

	for i, x := range foreignersSplitByInsertion {
		err = a.AddCollectionsContentUnits(tx, i != 0, x...)
		if err != nil {
			t.Fatal(err)
		}

		first := x[0]
		second := x[1]

		if a.ID != first.CollectionID {
			t.Error("foreign key was wrong value", a.ID, first.CollectionID)
		}
		if a.ID != second.CollectionID {
			t.Error("foreign key was wrong value", a.ID, second.CollectionID)
		}

		if first.R.Collection != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}
		if second.R.Collection != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}

		if a.R.CollectionsContentUnits[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.CollectionsContentUnits[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.CollectionsContentUnits(tx).Count()
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}
func testCollectionToOneStringTranslationUsingDescription(t *testing.T) {
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var local Collection
	var foreign StringTranslation

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, collectionDBTypes, true, collectionColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Collection struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, stringTranslationDBTypes, true, stringTranslationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize StringTranslation struct: %s", err)
	}

	local.DescriptionID.Valid = true

	if err := foreign.Insert(tx); err != nil {
		t.Fatal(err)
	}

	local.DescriptionID.Int64 = foreign.ID
	if err := local.Insert(tx); err != nil {
		t.Fatal(err)
	}

	check, err := local.Description(tx).One()
	if err != nil {
		t.Fatal(err)
	}

	if check.ID != foreign.ID {
		t.Errorf("want: %v, got %v", foreign.ID, check.ID)
	}

	slice := CollectionSlice{&local}
	if err = local.L.LoadDescription(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if local.R.Description == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.Description = nil
	if err = local.L.LoadDescription(tx, true, &local); err != nil {
		t.Fatal(err)
	}
	if local.R.Description == nil {
		t.Error("struct should have been eager loaded")
	}
}

func testCollectionToOneStringTranslationUsingName(t *testing.T) {
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var local Collection
	var foreign StringTranslation

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, collectionDBTypes, true, collectionColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Collection struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, stringTranslationDBTypes, true, stringTranslationColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize StringTranslation struct: %s", err)
	}

	if err := foreign.Insert(tx); err != nil {
		t.Fatal(err)
	}

	local.NameID = foreign.ID
	if err := local.Insert(tx); err != nil {
		t.Fatal(err)
	}

	check, err := local.Name(tx).One()
	if err != nil {
		t.Fatal(err)
	}

	if check.ID != foreign.ID {
		t.Errorf("want: %v, got %v", foreign.ID, check.ID)
	}

	slice := CollectionSlice{&local}
	if err = local.L.LoadName(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if local.R.Name == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.Name = nil
	if err = local.L.LoadName(tx, true, &local); err != nil {
		t.Fatal(err)
	}
	if local.R.Name == nil {
		t.Error("struct should have been eager loaded")
	}
}

func testCollectionToOneContentTypeUsingType(t *testing.T) {
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var local Collection
	var foreign ContentType

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, collectionDBTypes, true, collectionColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Collection struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, contentTypeDBTypes, true, contentTypeColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ContentType struct: %s", err)
	}

	if err := foreign.Insert(tx); err != nil {
		t.Fatal(err)
	}

	local.TypeID = foreign.ID
	if err := local.Insert(tx); err != nil {
		t.Fatal(err)
	}

	check, err := local.Type(tx).One()
	if err != nil {
		t.Fatal(err)
	}

	if check.ID != foreign.ID {
		t.Errorf("want: %v, got %v", foreign.ID, check.ID)
	}

	slice := CollectionSlice{&local}
	if err = local.L.LoadType(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if local.R.Type == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.Type = nil
	if err = local.L.LoadType(tx, true, &local); err != nil {
		t.Fatal(err)
	}
	if local.R.Type == nil {
		t.Error("struct should have been eager loaded")
	}
}

func testCollectionToOneSetOpStringTranslationUsingDescription(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Collection
	var b, c StringTranslation

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, collectionDBTypes, false, strmangle.SetComplement(collectionPrimaryKeyColumns, collectionColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, stringTranslationDBTypes, false, strmangle.SetComplement(stringTranslationPrimaryKeyColumns, stringTranslationColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, stringTranslationDBTypes, false, strmangle.SetComplement(stringTranslationPrimaryKeyColumns, stringTranslationColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}

	for i, x := range []*StringTranslation{&b, &c} {
		err = a.SetDescription(tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.Description != x {
			t.Error("relationship struct not set to correct value")
		}

		if x.R.DescriptionCollections[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if a.DescriptionID.Int64 != x.ID {
			t.Error("foreign key was wrong value", a.DescriptionID.Int64)
		}

		zero := reflect.Zero(reflect.TypeOf(a.DescriptionID.Int64))
		reflect.Indirect(reflect.ValueOf(&a.DescriptionID.Int64)).Set(zero)

		if err = a.Reload(tx); err != nil {
			t.Fatal("failed to reload", err)
		}

		if a.DescriptionID.Int64 != x.ID {
			t.Error("foreign key was wrong value", a.DescriptionID.Int64, x.ID)
		}
	}
}

func testCollectionToOneRemoveOpStringTranslationUsingDescription(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Collection
	var b StringTranslation

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, collectionDBTypes, false, strmangle.SetComplement(collectionPrimaryKeyColumns, collectionColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, stringTranslationDBTypes, false, strmangle.SetComplement(stringTranslationPrimaryKeyColumns, stringTranslationColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err = a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	if err = a.SetDescription(tx, true, &b); err != nil {
		t.Fatal(err)
	}

	if err = a.RemoveDescription(tx, &b); err != nil {
		t.Error("failed to remove relationship")
	}

	count, err := a.Description(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 0 {
		t.Error("want no relationships remaining")
	}

	if a.R.Description != nil {
		t.Error("R struct entry should be nil")
	}

	if a.DescriptionID.Valid {
		t.Error("foreign key value should be nil")
	}

	if len(b.R.DescriptionCollections) != 0 {
		t.Error("failed to remove a from b's relationships")
	}
}

func testCollectionToOneSetOpStringTranslationUsingName(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Collection
	var b, c StringTranslation

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, collectionDBTypes, false, strmangle.SetComplement(collectionPrimaryKeyColumns, collectionColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, stringTranslationDBTypes, false, strmangle.SetComplement(stringTranslationPrimaryKeyColumns, stringTranslationColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, stringTranslationDBTypes, false, strmangle.SetComplement(stringTranslationPrimaryKeyColumns, stringTranslationColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}

	for i, x := range []*StringTranslation{&b, &c} {
		err = a.SetName(tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.Name != x {
			t.Error("relationship struct not set to correct value")
		}

		if x.R.NameCollections[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if a.NameID != x.ID {
			t.Error("foreign key was wrong value", a.NameID)
		}

		zero := reflect.Zero(reflect.TypeOf(a.NameID))
		reflect.Indirect(reflect.ValueOf(&a.NameID)).Set(zero)

		if err = a.Reload(tx); err != nil {
			t.Fatal("failed to reload", err)
		}

		if a.NameID != x.ID {
			t.Error("foreign key was wrong value", a.NameID, x.ID)
		}
	}
}
func testCollectionToOneSetOpContentTypeUsingType(t *testing.T) {
	var err error

	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a Collection
	var b, c ContentType

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, collectionDBTypes, false, strmangle.SetComplement(collectionPrimaryKeyColumns, collectionColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, contentTypeDBTypes, false, strmangle.SetComplement(contentTypePrimaryKeyColumns, contentTypeColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, contentTypeDBTypes, false, strmangle.SetComplement(contentTypePrimaryKeyColumns, contentTypeColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}

	for i, x := range []*ContentType{&b, &c} {
		err = a.SetType(tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.Type != x {
			t.Error("relationship struct not set to correct value")
		}

		if x.R.TypeCollections[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if a.TypeID != x.ID {
			t.Error("foreign key was wrong value", a.TypeID)
		}

		zero := reflect.Zero(reflect.TypeOf(a.TypeID))
		reflect.Indirect(reflect.ValueOf(&a.TypeID)).Set(zero)

		if err = a.Reload(tx); err != nil {
			t.Fatal("failed to reload", err)
		}

		if a.TypeID != x.ID {
			t.Error("foreign key was wrong value", a.TypeID, x.ID)
		}
	}
}
func testCollectionsReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	collection := &Collection{}
	if err = randomize.Struct(seed, collection, collectionDBTypes, true, collectionColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Collection struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = collection.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = collection.Reload(tx); err != nil {
		t.Error(err)
	}
}

func testCollectionsReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	collection := &Collection{}
	if err = randomize.Struct(seed, collection, collectionDBTypes, true, collectionColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Collection struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = collection.Insert(tx); err != nil {
		t.Error(err)
	}

	slice := CollectionSlice{collection}

	if err = slice.ReloadAll(tx); err != nil {
		t.Error(err)
	}
}
func testCollectionsSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	collection := &Collection{}
	if err = randomize.Struct(seed, collection, collectionDBTypes, true, collectionColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Collection struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = collection.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := Collections(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	collectionDBTypes = map[string]string{`CreatedAt`: `timestamp with time zone`, `DescriptionID`: `bigint`, `ExternalID`: `character varying`, `ID`: `bigint`, `NameID`: `bigint`, `Properties`: `jsonb`, `TypeID`: `bigint`, `UID`: `character`}
	_                 = bytes.MinRead
)

func testCollectionsUpdate(t *testing.T) {
	t.Parallel()

	if len(collectionColumns) == len(collectionPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	collection := &Collection{}
	if err = randomize.Struct(seed, collection, collectionDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Collection struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = collection.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Collections(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, collection, collectionDBTypes, true, collectionColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Collection struct: %s", err)
	}

	if err = collection.Update(tx); err != nil {
		t.Error(err)
	}
}

func testCollectionsSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(collectionColumns) == len(collectionPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	collection := &Collection{}
	if err = randomize.Struct(seed, collection, collectionDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Collection struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = collection.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := Collections(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, collection, collectionDBTypes, true, collectionPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Collection struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(collectionColumns, collectionPrimaryKeyColumns) {
		fields = collectionColumns
	} else {
		fields = strmangle.SetComplement(
			collectionColumns,
			collectionPrimaryKeyColumns,
		)
	}

	value := reflect.Indirect(reflect.ValueOf(collection))
	updateMap := M{}
	for _, col := range fields {
		updateMap[col] = value.FieldByName(strmangle.TitleCase(col)).Interface()
	}

	slice := CollectionSlice{collection}
	if err = slice.UpdateAll(tx, updateMap); err != nil {
		t.Error(err)
	}
}
func testCollectionsUpsert(t *testing.T) {
	t.Parallel()

	if len(collectionColumns) == len(collectionPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	collection := Collection{}
	if err = randomize.Struct(seed, &collection, collectionDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Collection struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = collection.Upsert(tx, false, nil, nil); err != nil {
		t.Errorf("Unable to upsert Collection: %s", err)
	}

	count, err := Collections(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &collection, collectionDBTypes, false, collectionPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Collection struct: %s", err)
	}

	if err = collection.Upsert(tx, true, nil, nil); err != nil {
		t.Errorf("Unable to upsert Collection: %s", err)
	}

	count, err = Collections(tx).Count()
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
