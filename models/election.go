// Code generated by SQLBoiler 4.2.0 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/jackc/pgtype"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
	"github.com/volatiletech/sqlboiler/v4/types"
	"github.com/volatiletech/strmangle"
)

// Election is an object representing the database table.
type Election struct {
	ID        int64              `boil:"id" json:"id" toml:"id" yaml:"id"`
	Key       string             `boil:"key" json:"key" toml:"key" yaml:"key"`
	Question  string             `boil:"question" json:"question" toml:"question" yaml:"question"`
	Choices   types.StringArray  `boil:"choices" json:"choices" toml:"choices" yaml:"choices"`
	CreatedAt pgtype.Timestamptz `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	BallotKey string             `boil:"ballot_key" json:"ballot_key" toml:"ballot_key" yaml:"ballot_key"`
	Close     pgtype.Timestamptz `boil:"close" json:"close,omitempty" toml:"close" yaml:"close,omitempty"`
	CloseTZ   null.String        `boil:"close_tz" json:"close_tz,omitempty" toml:"close_tz" yaml:"close_tz,omitempty"`
	Flags     types.StringArray  `boil:"flags" json:"flags" toml:"flags" yaml:"flags"`

	R *electionR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L electionL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var ElectionColumns = struct {
	ID        string
	Key       string
	Question  string
	Choices   string
	CreatedAt string
	BallotKey string
	Close     string
	CloseTZ   string
	Flags     string
}{
	ID:        "id",
	Key:       "key",
	Question:  "question",
	Choices:   "choices",
	CreatedAt: "created_at",
	BallotKey: "ballot_key",
	Close:     "close",
	CloseTZ:   "close_tz",
	Flags:     "flags",
}

// Generated where

type whereHelperint64 struct{ field string }

func (w whereHelperint64) EQ(x int64) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.EQ, x) }
func (w whereHelperint64) NEQ(x int64) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.NEQ, x) }
func (w whereHelperint64) LT(x int64) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.LT, x) }
func (w whereHelperint64) LTE(x int64) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.LTE, x) }
func (w whereHelperint64) GT(x int64) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.GT, x) }
func (w whereHelperint64) GTE(x int64) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.GTE, x) }
func (w whereHelperint64) IN(slice []int64) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereIn(fmt.Sprintf("%s IN ?", w.field), values...)
}
func (w whereHelperint64) NIN(slice []int64) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereNotIn(fmt.Sprintf("%s NOT IN ?", w.field), values...)
}

type whereHelperstring struct{ field string }

func (w whereHelperstring) EQ(x string) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.EQ, x) }
func (w whereHelperstring) NEQ(x string) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.NEQ, x) }
func (w whereHelperstring) LT(x string) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.LT, x) }
func (w whereHelperstring) LTE(x string) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.LTE, x) }
func (w whereHelperstring) GT(x string) qm.QueryMod  { return qmhelper.Where(w.field, qmhelper.GT, x) }
func (w whereHelperstring) GTE(x string) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.GTE, x) }
func (w whereHelperstring) IN(slice []string) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereIn(fmt.Sprintf("%s IN ?", w.field), values...)
}
func (w whereHelperstring) NIN(slice []string) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereNotIn(fmt.Sprintf("%s NOT IN ?", w.field), values...)
}

type whereHelpertypes_StringArray struct{ field string }

func (w whereHelpertypes_StringArray) EQ(x types.StringArray) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.EQ, x)
}
func (w whereHelpertypes_StringArray) NEQ(x types.StringArray) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.NEQ, x)
}
func (w whereHelpertypes_StringArray) LT(x types.StringArray) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LT, x)
}
func (w whereHelpertypes_StringArray) LTE(x types.StringArray) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LTE, x)
}
func (w whereHelpertypes_StringArray) GT(x types.StringArray) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GT, x)
}
func (w whereHelpertypes_StringArray) GTE(x types.StringArray) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GTE, x)
}

type whereHelperpgtype_Timestamptz struct{ field string }

func (w whereHelperpgtype_Timestamptz) EQ(x pgtype.Timestamptz) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.EQ, x)
}
func (w whereHelperpgtype_Timestamptz) NEQ(x pgtype.Timestamptz) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.NEQ, x)
}
func (w whereHelperpgtype_Timestamptz) LT(x pgtype.Timestamptz) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LT, x)
}
func (w whereHelperpgtype_Timestamptz) LTE(x pgtype.Timestamptz) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LTE, x)
}
func (w whereHelperpgtype_Timestamptz) GT(x pgtype.Timestamptz) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GT, x)
}
func (w whereHelperpgtype_Timestamptz) GTE(x pgtype.Timestamptz) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GTE, x)
}

type whereHelpernull_String struct{ field string }

func (w whereHelpernull_String) EQ(x null.String) qm.QueryMod {
	return qmhelper.WhereNullEQ(w.field, false, x)
}
func (w whereHelpernull_String) NEQ(x null.String) qm.QueryMod {
	return qmhelper.WhereNullEQ(w.field, true, x)
}
func (w whereHelpernull_String) IsNull() qm.QueryMod    { return qmhelper.WhereIsNull(w.field) }
func (w whereHelpernull_String) IsNotNull() qm.QueryMod { return qmhelper.WhereIsNotNull(w.field) }
func (w whereHelpernull_String) LT(x null.String) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LT, x)
}
func (w whereHelpernull_String) LTE(x null.String) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.LTE, x)
}
func (w whereHelpernull_String) GT(x null.String) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GT, x)
}
func (w whereHelpernull_String) GTE(x null.String) qm.QueryMod {
	return qmhelper.Where(w.field, qmhelper.GTE, x)
}

var ElectionWhere = struct {
	ID        whereHelperint64
	Key       whereHelperstring
	Question  whereHelperstring
	Choices   whereHelpertypes_StringArray
	CreatedAt whereHelperpgtype_Timestamptz
	BallotKey whereHelperstring
	Close     whereHelperpgtype_Timestamptz
	CloseTZ   whereHelpernull_String
	Flags     whereHelpertypes_StringArray
}{
	ID:        whereHelperint64{field: "\"rv\".\"election\".\"id\""},
	Key:       whereHelperstring{field: "\"rv\".\"election\".\"key\""},
	Question:  whereHelperstring{field: "\"rv\".\"election\".\"question\""},
	Choices:   whereHelpertypes_StringArray{field: "\"rv\".\"election\".\"choices\""},
	CreatedAt: whereHelperpgtype_Timestamptz{field: "\"rv\".\"election\".\"created_at\""},
	BallotKey: whereHelperstring{field: "\"rv\".\"election\".\"ballot_key\""},
	Close:     whereHelperpgtype_Timestamptz{field: "\"rv\".\"election\".\"close\""},
	CloseTZ:   whereHelpernull_String{field: "\"rv\".\"election\".\"close_tz\""},
	Flags:     whereHelpertypes_StringArray{field: "\"rv\".\"election\".\"flags\""},
}

// ElectionRels is where relationship names are stored.
var ElectionRels = struct {
	Votes string
}{
	Votes: "Votes",
}

// electionR is where relationships are stored.
type electionR struct {
	Votes VoteSlice `boil:"Votes" json:"Votes" toml:"Votes" yaml:"Votes"`
}

// NewStruct creates a new relationship struct
func (*electionR) NewStruct() *electionR {
	return &electionR{}
}

// electionL is where Load methods for each relationship are stored.
type electionL struct{}

var (
	electionAllColumns            = []string{"id", "key", "question", "choices", "created_at", "ballot_key", "close", "close_tz", "flags"}
	electionColumnsWithoutDefault = []string{"key", "question", "choices", "created_at", "ballot_key", "close", "close_tz"}
	electionColumnsWithDefault    = []string{"id", "flags"}
	electionPrimaryKeyColumns     = []string{"id"}
)

type (
	// ElectionSlice is an alias for a slice of pointers to Election.
	// This should generally be used opposed to []Election.
	ElectionSlice []*Election

	electionQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	electionType                 = reflect.TypeOf(&Election{})
	electionMapping              = queries.MakeStructMapping(electionType)
	electionPrimaryKeyMapping, _ = queries.BindMapping(electionType, electionMapping, electionPrimaryKeyColumns)
	electionInsertCacheMut       sync.RWMutex
	electionInsertCache          = make(map[string]insertCache)
	electionUpdateCacheMut       sync.RWMutex
	electionUpdateCache          = make(map[string]updateCache)
	electionUpsertCacheMut       sync.RWMutex
	electionUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

// One returns a single election record from the query.
func (q electionQuery) One(ctx context.Context, exec boil.ContextExecutor) (*Election, error) {
	o := &Election{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for election")
	}

	return o, nil
}

// All returns all Election records from the query.
func (q electionQuery) All(ctx context.Context, exec boil.ContextExecutor) (ElectionSlice, error) {
	var o []*Election

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to Election slice")
	}

	return o, nil
}

// Count returns the count of all Election records in the query.
func (q electionQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count election rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q electionQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if election exists")
	}

	return count > 0, nil
}

// Votes retrieves all the vote's Votes with an executor.
func (o *Election) Votes(mods ...qm.QueryMod) voteQuery {
	var queryMods []qm.QueryMod
	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"rv\".\"vote\".\"election_id\"=?", o.ID),
	)

	query := Votes(queryMods...)
	queries.SetFrom(query.Query, "\"rv\".\"vote\"")

	if len(queries.GetSelect(query.Query)) == 0 {
		queries.SetSelect(query.Query, []string{"\"rv\".\"vote\".*"})
	}

	return query
}

// LoadVotes allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for a 1-M or N-M relationship.
func (electionL) LoadVotes(ctx context.Context, e boil.ContextExecutor, singular bool, maybeElection interface{}, mods queries.Applicator) error {
	var slice []*Election
	var object *Election

	if singular {
		object = maybeElection.(*Election)
	} else {
		slice = *maybeElection.(*[]*Election)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &electionR{}
		}
		args = append(args, object.ID)
	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &electionR{}
			}

			for _, a := range args {
				if a == obj.ID {
					continue Outer
				}
			}

			args = append(args, obj.ID)
		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(
		qm.From(`rv.vote`),
		qm.WhereIn(`rv.vote.election_id in ?`, args...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load vote")
	}

	var resultSlice []*Vote
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice vote")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results in eager load on vote")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for vote")
	}

	if singular {
		object.R.Votes = resultSlice
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.ID == foreign.ElectionID {
				local.R.Votes = append(local.R.Votes, foreign)
				break
			}
		}
	}

	return nil
}

// AddVotes adds the given related objects to the existing relationships
// of the election, optionally inserting them as new records.
// Appends related to o.R.Votes.
// Sets related.R.Election appropriately.
func (o *Election) AddVotes(ctx context.Context, exec boil.ContextExecutor, insert bool, related ...*Vote) error {
	var err error
	for _, rel := range related {
		if insert {
			rel.ElectionID = o.ID
			if err = rel.Insert(ctx, exec, boil.Infer()); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"rv\".\"vote\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"election_id"}),
				strmangle.WhereClause("\"", "\"", 2, votePrimaryKeyColumns),
			)
			values := []interface{}{o.ID, rel.ID}

			if boil.IsDebug(ctx) {
				writer := boil.DebugWriterFrom(ctx)
				fmt.Fprintln(writer, updateQuery)
				fmt.Fprintln(writer, values)
			}
			if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			rel.ElectionID = o.ID
		}
	}

	if o.R == nil {
		o.R = &electionR{
			Votes: related,
		}
	} else {
		o.R.Votes = append(o.R.Votes, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &voteR{
				Election: o,
			}
		} else {
			rel.R.Election = o
		}
	}
	return nil
}

// Elections retrieves all the records using an executor.
func Elections(mods ...qm.QueryMod) electionQuery {
	mods = append(mods, qm.From("\"rv\".\"election\""))
	return electionQuery{NewQuery(mods...)}
}

// FindElection retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindElection(ctx context.Context, exec boil.ContextExecutor, iD int64, selectCols ...string) (*Election, error) {
	electionObj := &Election{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"rv\".\"election\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, electionObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from election")
	}

	return electionObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *Election) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no election provided for insertion")
	}

	var err error

	nzDefaults := queries.NonZeroDefaultSet(electionColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	electionInsertCacheMut.RLock()
	cache, cached := electionInsertCache[key]
	electionInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			electionAllColumns,
			electionColumnsWithDefault,
			electionColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(electionType, electionMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(electionType, electionMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"rv\".\"election\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"rv\".\"election\" %sDEFAULT VALUES%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			queryReturning = fmt.Sprintf(" RETURNING \"%s\"", strings.Join(returnColumns, "\",\""))
		}

		cache.query = fmt.Sprintf(cache.query, queryOutput, queryReturning)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}

	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}

	if err != nil {
		return errors.Wrap(err, "models: unable to insert into election")
	}

	if !cached {
		electionInsertCacheMut.Lock()
		electionInsertCache[key] = cache
		electionInsertCacheMut.Unlock()
	}

	return nil
}

// Update uses an executor to update the Election.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *Election) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	key := makeCacheKey(columns, nil)
	electionUpdateCacheMut.RLock()
	cache, cached := electionUpdateCache[key]
	electionUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			electionAllColumns,
			electionPrimaryKeyColumns,
		)

		if len(wl) == 0 {
			return 0, errors.New("models: unable to update election, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"rv\".\"election\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, electionPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(electionType, electionMapping, append(wl, electionPrimaryKeyColumns...))
		if err != nil {
			return 0, err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, values)
	}
	var result sql.Result
	result, err = exec.ExecContext(ctx, cache.query, values...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update election row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for election")
	}

	if !cached {
		electionUpdateCacheMut.Lock()
		electionUpdateCache[key] = cache
		electionUpdateCacheMut.Unlock()
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values.
func (q electionQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for election")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for election")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o ElectionSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	ln := int64(len(o))
	if ln == 0 {
		return 0, nil
	}

	if len(cols) == 0 {
		return 0, errors.New("models: update all requires at least one column argument")
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), electionPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"rv\".\"election\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, electionPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in election slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all election")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *Election) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no election provided for upsert")
	}

	nzDefaults := queries.NonZeroDefaultSet(electionColumnsWithDefault, o)

	// Build cache key in-line uglily - mysql vs psql problems
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
	buf.WriteString(strconv.Itoa(updateColumns.Kind))
	for _, c := range updateColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(insertColumns.Kind))
	for _, c := range insertColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzDefaults {
		buf.WriteString(c)
	}
	key := buf.String()
	strmangle.PutBuffer(buf)

	electionUpsertCacheMut.RLock()
	cache, cached := electionUpsertCache[key]
	electionUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			electionAllColumns,
			electionColumnsWithDefault,
			electionColumnsWithoutDefault,
			nzDefaults,
		)
		update := updateColumns.UpdateColumnSet(
			electionAllColumns,
			electionPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert election, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(electionPrimaryKeyColumns))
			copy(conflict, electionPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"rv\".\"election\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(electionType, electionMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(electionType, electionMapping, ret)
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

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}
	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(returns...)
		if err == sql.ErrNoRows {
			err = nil // Postgres doesn't return anything when there's no update
		}
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}
	if err != nil {
		return errors.Wrap(err, "models: unable to upsert election")
	}

	if !cached {
		electionUpsertCacheMut.Lock()
		electionUpsertCache[key] = cache
		electionUpsertCacheMut.Unlock()
	}

	return nil
}

// Delete deletes a single Election record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Election) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no Election provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), electionPrimaryKeyMapping)
	sql := "DELETE FROM \"rv\".\"election\" WHERE \"id\"=$1"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from election")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for election")
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q electionQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no electionQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from election")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for election")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o ElectionSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), electionPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"rv\".\"election\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, electionPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from election slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for election")
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *Election) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindElection(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *ElectionSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := ElectionSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), electionPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"rv\".\"election\".* FROM \"rv\".\"election\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, electionPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in ElectionSlice")
	}

	*o = slice

	return nil
}

// ElectionExists checks if the Election row exists.
func ElectionExists(ctx context.Context, exec boil.ContextExecutor, iD int64) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"rv\".\"election\" where \"id\"=$1 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, iD)
	}
	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if election exists")
	}

	return exists, nil
}
