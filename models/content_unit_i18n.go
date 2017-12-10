// Code generated by SQLBoiler (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import (
	"bytes"
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"github.com/volatiletech/sqlboiler/strmangle"
	"gopkg.in/volatiletech/null.v6"
)

// ContentUnitI18n is an object representing the database table.
type ContentUnitI18n struct {
	ContentUnitID    int64       `boil:"content_unit_id" json:"content_unit_id" toml:"content_unit_id" yaml:"content_unit_id"`
	Language         string      `boil:"language" json:"language" toml:"language" yaml:"language"`
	OriginalLanguage null.String `boil:"original_language" json:"original_language,omitempty" toml:"original_language" yaml:"original_language,omitempty"`
	Name             null.String `boil:"name" json:"name,omitempty" toml:"name" yaml:"name,omitempty"`
	Description      null.String `boil:"description" json:"description,omitempty" toml:"description" yaml:"description,omitempty"`
	UserID           null.Int64  `boil:"user_id" json:"user_id,omitempty" toml:"user_id" yaml:"user_id,omitempty"`
	CreatedAt        time.Time   `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`

	R *contentUnitI18nR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L contentUnitI18nL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var ContentUnitI18nColumns = struct {
	ContentUnitID    string
	Language         string
	OriginalLanguage string
	Name             string
	Description      string
	UserID           string
	CreatedAt        string
}{
	ContentUnitID:    "content_unit_id",
	Language:         "language",
	OriginalLanguage: "original_language",
	Name:             "name",
	Description:      "description",
	UserID:           "user_id",
	CreatedAt:        "created_at",
}

// contentUnitI18nR is where relationships are stored.
type contentUnitI18nR struct {
	ContentUnit *ContentUnit
	User        *User
}

// contentUnitI18nL is where Load methods for each relationship are stored.
type contentUnitI18nL struct{}

var (
	contentUnitI18nColumns               = []string{"content_unit_id", "language", "original_language", "name", "description", "user_id", "created_at"}
	contentUnitI18nColumnsWithoutDefault = []string{"content_unit_id", "language", "original_language", "name", "description", "user_id"}
	contentUnitI18nColumnsWithDefault    = []string{"created_at"}
	contentUnitI18nPrimaryKeyColumns     = []string{"content_unit_id", "language"}
)

type (
	// ContentUnitI18nSlice is an alias for a slice of pointers to ContentUnitI18n.
	// This should generally be used opposed to []ContentUnitI18n.
	ContentUnitI18nSlice []*ContentUnitI18n

	contentUnitI18nQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	contentUnitI18nType                 = reflect.TypeOf(&ContentUnitI18n{})
	contentUnitI18nMapping              = queries.MakeStructMapping(contentUnitI18nType)
	contentUnitI18nPrimaryKeyMapping, _ = queries.BindMapping(contentUnitI18nType, contentUnitI18nMapping, contentUnitI18nPrimaryKeyColumns)
	contentUnitI18nInsertCacheMut       sync.RWMutex
	contentUnitI18nInsertCache          = make(map[string]insertCache)
	contentUnitI18nUpdateCacheMut       sync.RWMutex
	contentUnitI18nUpdateCache          = make(map[string]updateCache)
	contentUnitI18nUpsertCacheMut       sync.RWMutex
	contentUnitI18nUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force bytes in case of primary key column that uses []byte (for relationship compares)
	_ = bytes.MinRead
)

// OneP returns a single contentUnitI18n record from the query, and panics on error.
func (q contentUnitI18nQuery) OneP() *ContentUnitI18n {
	o, err := q.One()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// One returns a single contentUnitI18n record from the query.
func (q contentUnitI18nQuery) One() (*ContentUnitI18n, error) {
	o := &ContentUnitI18n{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for content_unit_i18n")
	}

	return o, nil
}

// AllP returns all ContentUnitI18n records from the query, and panics on error.
func (q contentUnitI18nQuery) AllP() ContentUnitI18nSlice {
	o, err := q.All()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// All returns all ContentUnitI18n records from the query.
func (q contentUnitI18nQuery) All() (ContentUnitI18nSlice, error) {
	var o []*ContentUnitI18n

	err := q.Bind(&o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to ContentUnitI18n slice")
	}

	return o, nil
}

// CountP returns the count of all ContentUnitI18n records in the query, and panics on error.
func (q contentUnitI18nQuery) CountP() int64 {
	c, err := q.Count()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return c
}

// Count returns the count of all ContentUnitI18n records in the query.
func (q contentUnitI18nQuery) Count() (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count content_unit_i18n rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table, and panics on error.
func (q contentUnitI18nQuery) ExistsP() bool {
	e, err := q.Exists()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// Exists checks if the row exists in the table.
func (q contentUnitI18nQuery) Exists() (bool, error) {
	var count int64

	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if content_unit_i18n exists")
	}

	return count > 0, nil
}

// ContentUnitG pointed to by the foreign key.
func (o *ContentUnitI18n) ContentUnitG(mods ...qm.QueryMod) contentUnitQuery {
	return o.ContentUnit(boil.GetDB(), mods...)
}

// ContentUnit pointed to by the foreign key.
func (o *ContentUnitI18n) ContentUnit(exec boil.Executor, mods ...qm.QueryMod) contentUnitQuery {
	queryMods := []qm.QueryMod{
		qm.Where("id=?", o.ContentUnitID),
	}

	queryMods = append(queryMods, mods...)

	query := ContentUnits(exec, queryMods...)
	queries.SetFrom(query.Query, "\"content_units\"")

	return query
}

// UserG pointed to by the foreign key.
func (o *ContentUnitI18n) UserG(mods ...qm.QueryMod) userQuery {
	return o.User(boil.GetDB(), mods...)
}

// User pointed to by the foreign key.
func (o *ContentUnitI18n) User(exec boil.Executor, mods ...qm.QueryMod) userQuery {
	queryMods := []qm.QueryMod{
		qm.Where("id=?", o.UserID),
	}

	queryMods = append(queryMods, mods...)

	query := Users(exec, queryMods...)
	queries.SetFrom(query.Query, "\"users\"")

	return query
} // LoadContentUnit allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (contentUnitI18nL) LoadContentUnit(e boil.Executor, singular bool, maybeContentUnitI18n interface{}) error {
	var slice []*ContentUnitI18n
	var object *ContentUnitI18n

	count := 1
	if singular {
		object = maybeContentUnitI18n.(*ContentUnitI18n)
	} else {
		slice = *maybeContentUnitI18n.(*[]*ContentUnitI18n)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &contentUnitI18nR{}
		}
		args[0] = object.ContentUnitID
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &contentUnitI18nR{}
			}
			args[i] = obj.ContentUnitID
		}
	}

	query := fmt.Sprintf(
		"select * from \"content_units\" where \"id\" in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, count, 1, 1),
	)

	if boil.DebugMode {
		fmt.Fprintf(boil.DebugWriter, "%s\n%v\n", query, args)
	}

	results, err := e.Query(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to eager load ContentUnit")
	}
	defer results.Close()

	var resultSlice []*ContentUnit
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice ContentUnit")
	}

	if len(resultSlice) == 0 {
		return nil
	}

	if singular {
		object.R.ContentUnit = resultSlice[0]
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.ContentUnitID == foreign.ID {
				local.R.ContentUnit = foreign
				break
			}
		}
	}

	return nil
}

// LoadUser allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (contentUnitI18nL) LoadUser(e boil.Executor, singular bool, maybeContentUnitI18n interface{}) error {
	var slice []*ContentUnitI18n
	var object *ContentUnitI18n

	count := 1
	if singular {
		object = maybeContentUnitI18n.(*ContentUnitI18n)
	} else {
		slice = *maybeContentUnitI18n.(*[]*ContentUnitI18n)
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &contentUnitI18nR{}
		}
		args[0] = object.UserID
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &contentUnitI18nR{}
			}
			args[i] = obj.UserID
		}
	}

	query := fmt.Sprintf(
		"select * from \"users\" where \"id\" in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, count, 1, 1),
	)

	if boil.DebugMode {
		fmt.Fprintf(boil.DebugWriter, "%s\n%v\n", query, args)
	}

	results, err := e.Query(query, args...)
	if err != nil {
		return errors.Wrap(err, "failed to eager load User")
	}
	defer results.Close()

	var resultSlice []*User
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice User")
	}

	if len(resultSlice) == 0 {
		return nil
	}

	if singular {
		object.R.User = resultSlice[0]
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.UserID.Int64 == foreign.ID {
				local.R.User = foreign
				break
			}
		}
	}

	return nil
}

// SetContentUnitG of the content_unit_i18n to the related item.
// Sets o.R.ContentUnit to related.
// Adds o to related.R.ContentUnitI18ns.
// Uses the global database handle.
func (o *ContentUnitI18n) SetContentUnitG(insert bool, related *ContentUnit) error {
	return o.SetContentUnit(boil.GetDB(), insert, related)
}

// SetContentUnitP of the content_unit_i18n to the related item.
// Sets o.R.ContentUnit to related.
// Adds o to related.R.ContentUnitI18ns.
// Panics on error.
func (o *ContentUnitI18n) SetContentUnitP(exec boil.Executor, insert bool, related *ContentUnit) {
	if err := o.SetContentUnit(exec, insert, related); err != nil {
		panic(boil.WrapErr(err))
	}
}

// SetContentUnitGP of the content_unit_i18n to the related item.
// Sets o.R.ContentUnit to related.
// Adds o to related.R.ContentUnitI18ns.
// Uses the global database handle and panics on error.
func (o *ContentUnitI18n) SetContentUnitGP(insert bool, related *ContentUnit) {
	if err := o.SetContentUnit(boil.GetDB(), insert, related); err != nil {
		panic(boil.WrapErr(err))
	}
}

// SetContentUnit of the content_unit_i18n to the related item.
// Sets o.R.ContentUnit to related.
// Adds o to related.R.ContentUnitI18ns.
func (o *ContentUnitI18n) SetContentUnit(exec boil.Executor, insert bool, related *ContentUnit) error {
	var err error
	if insert {
		if err = related.Insert(exec); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"content_unit_i18n\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"content_unit_id"}),
		strmangle.WhereClause("\"", "\"", 2, contentUnitI18nPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.ContentUnitID, o.Language}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, updateQuery)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	if _, err = exec.Exec(updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.ContentUnitID = related.ID

	if o.R == nil {
		o.R = &contentUnitI18nR{
			ContentUnit: related,
		}
	} else {
		o.R.ContentUnit = related
	}

	if related.R == nil {
		related.R = &contentUnitR{
			ContentUnitI18ns: ContentUnitI18nSlice{o},
		}
	} else {
		related.R.ContentUnitI18ns = append(related.R.ContentUnitI18ns, o)
	}

	return nil
}

// SetUserG of the content_unit_i18n to the related item.
// Sets o.R.User to related.
// Adds o to related.R.ContentUnitI18ns.
// Uses the global database handle.
func (o *ContentUnitI18n) SetUserG(insert bool, related *User) error {
	return o.SetUser(boil.GetDB(), insert, related)
}

// SetUserP of the content_unit_i18n to the related item.
// Sets o.R.User to related.
// Adds o to related.R.ContentUnitI18ns.
// Panics on error.
func (o *ContentUnitI18n) SetUserP(exec boil.Executor, insert bool, related *User) {
	if err := o.SetUser(exec, insert, related); err != nil {
		panic(boil.WrapErr(err))
	}
}

// SetUserGP of the content_unit_i18n to the related item.
// Sets o.R.User to related.
// Adds o to related.R.ContentUnitI18ns.
// Uses the global database handle and panics on error.
func (o *ContentUnitI18n) SetUserGP(insert bool, related *User) {
	if err := o.SetUser(boil.GetDB(), insert, related); err != nil {
		panic(boil.WrapErr(err))
	}
}

// SetUser of the content_unit_i18n to the related item.
// Sets o.R.User to related.
// Adds o to related.R.ContentUnitI18ns.
func (o *ContentUnitI18n) SetUser(exec boil.Executor, insert bool, related *User) error {
	var err error
	if insert {
		if err = related.Insert(exec); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"content_unit_i18n\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"user_id"}),
		strmangle.WhereClause("\"", "\"", 2, contentUnitI18nPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.ContentUnitID, o.Language}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, updateQuery)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	if _, err = exec.Exec(updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.UserID.Int64 = related.ID
	o.UserID.Valid = true

	if o.R == nil {
		o.R = &contentUnitI18nR{
			User: related,
		}
	} else {
		o.R.User = related
	}

	if related.R == nil {
		related.R = &userR{
			ContentUnitI18ns: ContentUnitI18nSlice{o},
		}
	} else {
		related.R.ContentUnitI18ns = append(related.R.ContentUnitI18ns, o)
	}

	return nil
}

// RemoveUserG relationship.
// Sets o.R.User to nil.
// Removes o from all passed in related items' relationships struct (Optional).
// Uses the global database handle.
func (o *ContentUnitI18n) RemoveUserG(related *User) error {
	return o.RemoveUser(boil.GetDB(), related)
}

// RemoveUserP relationship.
// Sets o.R.User to nil.
// Removes o from all passed in related items' relationships struct (Optional).
// Panics on error.
func (o *ContentUnitI18n) RemoveUserP(exec boil.Executor, related *User) {
	if err := o.RemoveUser(exec, related); err != nil {
		panic(boil.WrapErr(err))
	}
}

// RemoveUserGP relationship.
// Sets o.R.User to nil.
// Removes o from all passed in related items' relationships struct (Optional).
// Uses the global database handle and panics on error.
func (o *ContentUnitI18n) RemoveUserGP(related *User) {
	if err := o.RemoveUser(boil.GetDB(), related); err != nil {
		panic(boil.WrapErr(err))
	}
}

// RemoveUser relationship.
// Sets o.R.User to nil.
// Removes o from all passed in related items' relationships struct (Optional).
func (o *ContentUnitI18n) RemoveUser(exec boil.Executor, related *User) error {
	var err error

	o.UserID.Valid = false
	if err = o.Update(exec, "user_id"); err != nil {
		o.UserID.Valid = true
		return errors.Wrap(err, "failed to update local table")
	}

	o.R.User = nil
	if related == nil || related.R == nil {
		return nil
	}

	for i, ri := range related.R.ContentUnitI18ns {
		if o.UserID.Int64 != ri.UserID.Int64 {
			continue
		}

		ln := len(related.R.ContentUnitI18ns)
		if ln > 1 && i < ln-1 {
			related.R.ContentUnitI18ns[i] = related.R.ContentUnitI18ns[ln-1]
		}
		related.R.ContentUnitI18ns = related.R.ContentUnitI18ns[:ln-1]
		break
	}
	return nil
}

// ContentUnitI18nsG retrieves all records.
func ContentUnitI18nsG(mods ...qm.QueryMod) contentUnitI18nQuery {
	return ContentUnitI18ns(boil.GetDB(), mods...)
}

// ContentUnitI18ns retrieves all the records using an executor.
func ContentUnitI18ns(exec boil.Executor, mods ...qm.QueryMod) contentUnitI18nQuery {
	mods = append(mods, qm.From("\"content_unit_i18n\""))
	return contentUnitI18nQuery{NewQuery(exec, mods...)}
}

// FindContentUnitI18nG retrieves a single record by ID.
func FindContentUnitI18nG(contentUnitID int64, language string, selectCols ...string) (*ContentUnitI18n, error) {
	return FindContentUnitI18n(boil.GetDB(), contentUnitID, language, selectCols...)
}

// FindContentUnitI18nGP retrieves a single record by ID, and panics on error.
func FindContentUnitI18nGP(contentUnitID int64, language string, selectCols ...string) *ContentUnitI18n {
	retobj, err := FindContentUnitI18n(boil.GetDB(), contentUnitID, language, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// FindContentUnitI18n retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindContentUnitI18n(exec boil.Executor, contentUnitID int64, language string, selectCols ...string) (*ContentUnitI18n, error) {
	contentUnitI18nObj := &ContentUnitI18n{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"content_unit_i18n\" where \"content_unit_id\"=$1 AND \"language\"=$2", sel,
	)

	q := queries.Raw(exec, query, contentUnitID, language)

	err := q.Bind(contentUnitI18nObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from content_unit_i18n")
	}

	return contentUnitI18nObj, nil
}

// FindContentUnitI18nP retrieves a single record by ID with an executor, and panics on error.
func FindContentUnitI18nP(exec boil.Executor, contentUnitID int64, language string, selectCols ...string) *ContentUnitI18n {
	retobj, err := FindContentUnitI18n(exec, contentUnitID, language, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// InsertG a single record. See Insert for whitelist behavior description.
func (o *ContentUnitI18n) InsertG(whitelist ...string) error {
	return o.Insert(boil.GetDB(), whitelist...)
}

// InsertGP a single record, and panics on error. See Insert for whitelist
// behavior description.
func (o *ContentUnitI18n) InsertGP(whitelist ...string) {
	if err := o.Insert(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// InsertP a single record using an executor, and panics on error. See Insert
// for whitelist behavior description.
func (o *ContentUnitI18n) InsertP(exec boil.Executor, whitelist ...string) {
	if err := o.Insert(exec, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Insert a single record using an executor.
// Whitelist behavior: If a whitelist is provided, only those columns supplied are inserted
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns without a default value are included (i.e. name, age)
// - All columns with a default, but non-zero are included (i.e. health = 75)
func (o *ContentUnitI18n) Insert(exec boil.Executor, whitelist ...string) error {
	if o == nil {
		return errors.New("models: no content_unit_i18n provided for insertion")
	}

	var err error

	nzDefaults := queries.NonZeroDefaultSet(contentUnitI18nColumnsWithDefault, o)

	key := makeCacheKey(whitelist, nzDefaults)
	contentUnitI18nInsertCacheMut.RLock()
	cache, cached := contentUnitI18nInsertCache[key]
	contentUnitI18nInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := strmangle.InsertColumnSet(
			contentUnitI18nColumns,
			contentUnitI18nColumnsWithDefault,
			contentUnitI18nColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)

		cache.valueMapping, err = queries.BindMapping(contentUnitI18nType, contentUnitI18nMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(contentUnitI18nType, contentUnitI18nMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"content_unit_i18n\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.IndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"content_unit_i18n\" DEFAULT VALUES"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			queryReturning = fmt.Sprintf(" RETURNING \"%s\"", strings.Join(returnColumns, "\",\""))
		}

		if len(wl) != 0 {
			cache.query = fmt.Sprintf(cache.query, queryOutput, queryReturning)
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, vals)
	}

	if len(cache.retMapping) != 0 {
		err = exec.QueryRow(cache.query, vals...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	} else {
		_, err = exec.Exec(cache.query, vals...)
	}

	if err != nil {
		return errors.Wrap(err, "models: unable to insert into content_unit_i18n")
	}

	if !cached {
		contentUnitI18nInsertCacheMut.Lock()
		contentUnitI18nInsertCache[key] = cache
		contentUnitI18nInsertCacheMut.Unlock()
	}

	return nil
}

// UpdateG a single ContentUnitI18n record. See Update for
// whitelist behavior description.
func (o *ContentUnitI18n) UpdateG(whitelist ...string) error {
	return o.Update(boil.GetDB(), whitelist...)
}

// UpdateGP a single ContentUnitI18n record.
// UpdateGP takes a whitelist of column names that should be updated.
// Panics on error. See Update for whitelist behavior description.
func (o *ContentUnitI18n) UpdateGP(whitelist ...string) {
	if err := o.Update(boil.GetDB(), whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateP uses an executor to update the ContentUnitI18n, and panics on error.
// See Update for whitelist behavior description.
func (o *ContentUnitI18n) UpdateP(exec boil.Executor, whitelist ...string) {
	err := o.Update(exec, whitelist...)
	if err != nil {
		panic(boil.WrapErr(err))
	}
}

// Update uses an executor to update the ContentUnitI18n.
// Whitelist behavior: If a whitelist is provided, only the columns given are updated.
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns are inferred to start with
// - All primary keys are subtracted from this set
// Update does not automatically update the record in case of default values. Use .Reload()
// to refresh the records.
func (o *ContentUnitI18n) Update(exec boil.Executor, whitelist ...string) error {
	var err error
	key := makeCacheKey(whitelist, nil)
	contentUnitI18nUpdateCacheMut.RLock()
	cache, cached := contentUnitI18nUpdateCache[key]
	contentUnitI18nUpdateCacheMut.RUnlock()

	if !cached {
		wl := strmangle.UpdateColumnSet(
			contentUnitI18nColumns,
			contentUnitI18nPrimaryKeyColumns,
			whitelist,
		)

		if len(wl) == 0 {
			return errors.New("models: unable to update content_unit_i18n, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"content_unit_i18n\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, contentUnitI18nPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(contentUnitI18nType, contentUnitI18nMapping, append(wl, contentUnitI18nPrimaryKeyColumns...))
		if err != nil {
			return err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	_, err = exec.Exec(cache.query, values...)
	if err != nil {
		return errors.Wrap(err, "models: unable to update content_unit_i18n row")
	}

	if !cached {
		contentUnitI18nUpdateCacheMut.Lock()
		contentUnitI18nUpdateCache[key] = cache
		contentUnitI18nUpdateCacheMut.Unlock()
	}

	return nil
}

// UpdateAllP updates all rows with matching column names, and panics on error.
func (q contentUnitI18nQuery) UpdateAllP(cols M) {
	if err := q.UpdateAll(cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values.
func (q contentUnitI18nQuery) UpdateAll(cols M) error {
	queries.SetUpdate(q.Query, cols)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "models: unable to update all for content_unit_i18n")
	}

	return nil
}

// UpdateAllG updates all rows with the specified column values.
func (o ContentUnitI18nSlice) UpdateAllG(cols M) error {
	return o.UpdateAll(boil.GetDB(), cols)
}

// UpdateAllGP updates all rows with the specified column values, and panics on error.
func (o ContentUnitI18nSlice) UpdateAllGP(cols M) {
	if err := o.UpdateAll(boil.GetDB(), cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAllP updates all rows with the specified column values, and panics on error.
func (o ContentUnitI18nSlice) UpdateAllP(exec boil.Executor, cols M) {
	if err := o.UpdateAll(exec, cols); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o ContentUnitI18nSlice) UpdateAll(exec boil.Executor, cols M) error {
	ln := int64(len(o))
	if ln == 0 {
		return nil
	}

	if len(cols) == 0 {
		return errors.New("models: update all requires at least one column argument")
	}

	colNames := make([]string, len(cols))
	args := make([]interface{}, len(cols))

	i := 0
	for name, value := range cols {
		colNames[i] = name
		args[i] = value
		i++
	}

	// Append all of the primary key values for each column
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), contentUnitI18nPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"content_unit_i18n\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, contentUnitI18nPrimaryKeyColumns, len(o)))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to update all in contentUnitI18n slice")
	}

	return nil
}

// UpsertG attempts an insert, and does an update or ignore on conflict.
func (o *ContentUnitI18n) UpsertG(updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	return o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, whitelist...)
}

// UpsertGP attempts an insert, and does an update or ignore on conflict. Panics on error.
func (o *ContentUnitI18n) UpsertGP(updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(boil.GetDB(), updateOnConflict, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpsertP attempts an insert using an executor, and does an update or ignore on conflict.
// UpsertP panics on error.
func (o *ContentUnitI18n) UpsertP(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) {
	if err := o.Upsert(exec, updateOnConflict, conflictColumns, updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
func (o *ContentUnitI18n) Upsert(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
	if o == nil {
		return errors.New("models: no content_unit_i18n provided for upsert")
	}

	nzDefaults := queries.NonZeroDefaultSet(contentUnitI18nColumnsWithDefault, o)

	// Build cache key in-line uglily - mysql vs postgres problems
	buf := strmangle.GetBuffer()

	if updateOnConflict {
		buf.WriteByte('t')
	} else {
		buf.WriteByte('f')
	}
	buf.WriteByte('.')
	for _, c := range conflictColumns {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range updateColumns {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range whitelist {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzDefaults {
		buf.WriteString(c)
	}
	key := buf.String()
	strmangle.PutBuffer(buf)

	contentUnitI18nUpsertCacheMut.RLock()
	cache, cached := contentUnitI18nUpsertCache[key]
	contentUnitI18nUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := strmangle.InsertColumnSet(
			contentUnitI18nColumns,
			contentUnitI18nColumnsWithDefault,
			contentUnitI18nColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)

		update := strmangle.UpdateColumnSet(
			contentUnitI18nColumns,
			contentUnitI18nPrimaryKeyColumns,
			updateColumns,
		)
		if len(update) == 0 {
			return errors.New("models: unable to upsert content_unit_i18n, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(contentUnitI18nPrimaryKeyColumns))
			copy(conflict, contentUnitI18nPrimaryKeyColumns)
		}
		cache.query = queries.BuildUpsertQueryPostgres(dialect, "\"content_unit_i18n\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(contentUnitI18nType, contentUnitI18nMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(contentUnitI18nType, contentUnitI18nMapping, ret)
			if err != nil {
				return err
			}
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)
	var returns []interface{}
	if len(cache.retMapping) != 0 {
		returns = queries.PtrsFromMapping(value, cache.retMapping)
	}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, vals)
	}

	if len(cache.retMapping) != 0 {
		err = exec.QueryRow(cache.query, vals...).Scan(returns...)
		if err == sql.ErrNoRows {
			err = nil // Postgres doesn't return anything when there's no update
		}
	} else {
		_, err = exec.Exec(cache.query, vals...)
	}
	if err != nil {
		return errors.Wrap(err, "models: unable to upsert content_unit_i18n")
	}

	if !cached {
		contentUnitI18nUpsertCacheMut.Lock()
		contentUnitI18nUpsertCache[key] = cache
		contentUnitI18nUpsertCacheMut.Unlock()
	}

	return nil
}

// DeleteP deletes a single ContentUnitI18n record with an executor.
// DeleteP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *ContentUnitI18n) DeleteP(exec boil.Executor) {
	if err := o.Delete(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteG deletes a single ContentUnitI18n record.
// DeleteG will match against the primary key column to find the record to delete.
func (o *ContentUnitI18n) DeleteG() error {
	if o == nil {
		return errors.New("models: no ContentUnitI18n provided for deletion")
	}

	return o.Delete(boil.GetDB())
}

// DeleteGP deletes a single ContentUnitI18n record.
// DeleteGP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *ContentUnitI18n) DeleteGP() {
	if err := o.DeleteG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// Delete deletes a single ContentUnitI18n record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *ContentUnitI18n) Delete(exec boil.Executor) error {
	if o == nil {
		return errors.New("models: no ContentUnitI18n provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), contentUnitI18nPrimaryKeyMapping)
	sql := "DELETE FROM \"content_unit_i18n\" WHERE \"content_unit_id\"=$1 AND \"language\"=$2"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete from content_unit_i18n")
	}

	return nil
}

// DeleteAllP deletes all rows, and panics on error.
func (q contentUnitI18nQuery) DeleteAllP() {
	if err := q.DeleteAll(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all matching rows.
func (q contentUnitI18nQuery) DeleteAll() error {
	if q.Query == nil {
		return errors.New("models: no contentUnitI18nQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	_, err := q.Query.Exec()
	if err != nil {
		return errors.Wrap(err, "models: unable to delete all from content_unit_i18n")
	}

	return nil
}

// DeleteAllGP deletes all rows in the slice, and panics on error.
func (o ContentUnitI18nSlice) DeleteAllGP() {
	if err := o.DeleteAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAllG deletes all rows in the slice.
func (o ContentUnitI18nSlice) DeleteAllG() error {
	if o == nil {
		return errors.New("models: no ContentUnitI18n slice provided for delete all")
	}
	return o.DeleteAll(boil.GetDB())
}

// DeleteAllP deletes all rows in the slice, using an executor, and panics on error.
func (o ContentUnitI18nSlice) DeleteAllP(exec boil.Executor) {
	if err := o.DeleteAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o ContentUnitI18nSlice) DeleteAll(exec boil.Executor) error {
	if o == nil {
		return errors.New("models: no ContentUnitI18n slice provided for delete all")
	}

	if len(o) == 0 {
		return nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), contentUnitI18nPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"content_unit_i18n\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, contentUnitI18nPrimaryKeyColumns, len(o))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "models: unable to delete all from contentUnitI18n slice")
	}

	return nil
}

// ReloadGP refetches the object from the database and panics on error.
func (o *ContentUnitI18n) ReloadGP() {
	if err := o.ReloadG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadP refetches the object from the database with an executor. Panics on error.
func (o *ContentUnitI18n) ReloadP(exec boil.Executor) {
	if err := o.Reload(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadG refetches the object from the database using the primary keys.
func (o *ContentUnitI18n) ReloadG() error {
	if o == nil {
		return errors.New("models: no ContentUnitI18n provided for reload")
	}

	return o.Reload(boil.GetDB())
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *ContentUnitI18n) Reload(exec boil.Executor) error {
	ret, err := FindContentUnitI18n(exec, o.ContentUnitID, o.Language)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAllGP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *ContentUnitI18nSlice) ReloadAllGP() {
	if err := o.ReloadAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *ContentUnitI18nSlice) ReloadAllP(exec boil.Executor) {
	if err := o.ReloadAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllG refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *ContentUnitI18nSlice) ReloadAllG() error {
	if o == nil {
		return errors.New("models: empty ContentUnitI18nSlice provided for reload all")
	}

	return o.ReloadAll(boil.GetDB())
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *ContentUnitI18nSlice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	contentUnitI18ns := ContentUnitI18nSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), contentUnitI18nPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"content_unit_i18n\".* FROM \"content_unit_i18n\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, contentUnitI18nPrimaryKeyColumns, len(*o))

	q := queries.Raw(exec, sql, args...)

	err := q.Bind(&contentUnitI18ns)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in ContentUnitI18nSlice")
	}

	*o = contentUnitI18ns

	return nil
}

// ContentUnitI18nExists checks if the ContentUnitI18n row exists.
func ContentUnitI18nExists(exec boil.Executor, contentUnitID int64, language string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"content_unit_i18n\" where \"content_unit_id\"=$1 AND \"language\"=$2 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, contentUnitID, language)
	}

	row := exec.QueryRow(sql, contentUnitID, language)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if content_unit_i18n exists")
	}

	return exists, nil
}

// ContentUnitI18nExistsG checks if the ContentUnitI18n row exists.
func ContentUnitI18nExistsG(contentUnitID int64, language string) (bool, error) {
	return ContentUnitI18nExists(boil.GetDB(), contentUnitID, language)
}

// ContentUnitI18nExistsGP checks if the ContentUnitI18n row exists. Panics on error.
func ContentUnitI18nExistsGP(contentUnitID int64, language string) bool {
	e, err := ContentUnitI18nExists(boil.GetDB(), contentUnitID, language)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// ContentUnitI18nExistsP checks if the ContentUnitI18n row exists. Panics on error.
func ContentUnitI18nExistsP(exec boil.Executor, contentUnitID int64, language string) bool {
	e, err := ContentUnitI18nExists(exec, contentUnitID, language)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}
