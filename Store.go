package settingstore

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"log/slog"
	"strconv"
	"strings"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"     // importing mysql dialect
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"  // importing postgres dialect
	_ "github.com/doug-martin/goqu/v9/dialect/sqlite3"   // importing sqlite3 dialect
	_ "github.com/doug-martin/goqu/v9/dialect/sqlserver" // importing sqlserver dialect
	"github.com/dromara/carbon/v2"
	"github.com/gouniverse/base/database"
	"github.com/gouniverse/sb"
	"github.com/samber/lo"
)

// == INTERFACE ===============================================================

var _ StoreInterface = (*store)(nil) // verify it extends the store interface

// == TYPE ====================================================================

// Store defines a setting store
type store struct {
	settingTableName   string
	db                 *sql.DB
	dbDriverName       string
	timeoutSeconds     int64
	automigrateEnabled bool
	debugEnabled       bool
	sqlLogger          *slog.Logger
}

// PUBLIC METHODS ============================================================

// AutoMigrate creates the settings table if it does not exist
//
// Parameters:
// - ctx: the context
//
// Returns:
// - error - nil if no error, error otherwise
func (store *store) AutoMigrate(ctx context.Context) error {
	sqlStr := store.SQLCreateTable()

	if sqlStr == "" {
		return errors.New("setting store: table create sql is empty")
	}

	if store.db == nil {
		return errors.New("setting store: database is nil")
	}

	_, err := database.Execute(database.Context(ctx, store.db), sqlStr)

	if err != nil {
		return err
	}

	return nil
}

// EnableDebug - enables the debug option
//
// # If enabled will log the SQL statements to the provided logger
//
// Parameters:
// - debug: true to enable, false otherwise
//
// Returns:
// - void
func (st *store) EnableDebug(debug bool) {
	st.debugEnabled = debug
}

// Delete is a shortcut method to delete a value by key
//
// # It is a convenience method which wraps SettingDeleteByKey
//
// Parameters:
// - ctx: the context
// - settingKey: the key of the setting to delete
//
// Returns:
// - error - nil if no error, error otherwise
func (st *store) Delete(ctx context.Context, settingKey string) error {
	return st.SettingDeleteByKey(ctx, settingKey)
}

// Get is a shortcut method to get a value by key, or a default, if not found
//
// It is a convenience method which wraps SettingFindByKey and returns
// the value directly
//
// Parameters:
// - ctx: the context
// - settingKey: the key of the setting to get
// - valueDefault: the default value to return if the setting is not found
//
// Returns:
// - string - the value of the setting, or the default value if not found
// - error - nil if no error, error otherwise
func (st *store) Get(ctx context.Context, settingKey string, valueDefault string) (string, error) {
	setting, errFindByKey := st.SettingFindByKey(ctx, settingKey)

	if errFindByKey != nil {
		return "", errFindByKey
	}

	if setting != nil {
		return setting.GetValue(), nil
	}

	return valueDefault, nil
}

// GetAny is a shortcut method to get a value by key as an interface, or a default if not found
//
// It is a convenience method which wraps SettingFindByKey, gets the value
// directly and attempts to parse the value as interface
//
// Parameters:
// - ctx: the context
// - settingKey: the key of the setting to get
// - valueDefault: the default value to return if the setting is not found
//
// Returns:
// - interface{}, error
func (st *store) GetAny(ctx context.Context, key string, valueDefault any) (any, error) {
	setting, errFindByKey := st.SettingFindByKey(ctx, key)

	if errFindByKey != nil {
		return valueDefault, errFindByKey
	}

	if setting != nil {
		jsonValue := setting.GetValue()
		var val interface{}
		jsonError := json.Unmarshal([]byte(jsonValue), &val)
		if jsonError != nil {
			return valueDefault, jsonError
		}

		return val, nil
	}

	return valueDefault, nil
}

// GetMap is a shortcut method to get a value by key as a map, or a default if not found
//
// It is a convenience method which wraps SettingFindByKey, and attempts
// to parse the value as map[string]any
//
// Parameters:
// - ctx: the context
// - settingKey: the key of the setting to get
// - valueDefault: the default value to return if the setting is not found
//
// Returns:
// - map[string]any - the value of the setting, or the default value if not found
// - error - nil if no error, error otherwise
func (st *store) GetMap(ctx context.Context, key string, valueDefault map[string]any) (map[string]any, error) {
	setting, errFindByKey := st.SettingFindByKey(ctx, key)

	if errFindByKey != nil {
		return valueDefault, errFindByKey
	}

	if setting != nil {
		jsonValue := setting.GetValue()
		var val map[string]any
		jsonError := json.Unmarshal([]byte(jsonValue), &val)
		if jsonError != nil {
			return valueDefault, jsonError
		}

		return val, nil
	}

	return valueDefault, nil
}

// Has is a shortcut method to check if a setting exists by key
//
// It is a convenience method which wraps SettingCount and returns
// true if the count is greater than 0
//
// Parameters:
// - ctx: the context
// - settingKey: the key of the setting to check
//
// Returns:
// - bool - true if the setting exists, false otherwise
// - error - nil if no error, error otherwise
func (store *store) Has(ctx context.Context, settingKey string) (bool, error) {
	if settingKey == "" {
		return false, errors.New("setting store > find by key: setting key is required")
	}

	query := SettingQuery().
		SetKey(settingKey).
		SetLimit(1)

	count, err := store.SettingCount(ctx, query)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// MergeMap is a shortcut method to merge a map with an existing map
//
// It is a convenience method which wraps GetMap and SetMap to merge
// a map with an existing map
//
// Parameters:
// - ctx: the context
// - key: the key of the setting to merge
// - mergeMap: the map to merge with the existing map
//
// Returns:
// - error - nil if no error, error otherwise
func (st *store) MergeMap(ctx context.Context, key string, mergeMap map[string]any) error {
	currentMap, err := st.GetMap(ctx, key, nil)

	if err != nil {
		return err
	}

	if currentMap == nil {
		return errors.New("settingstore. nil found")
	}

	for mapKey, mapValue := range mergeMap {
		currentMap[mapKey] = mapValue
	}

	return st.SetMap(ctx, key, currentMap)
}

func (store *store) SettingCount(ctx context.Context, options SettingQueryInterface) (int64, error) {
	options.SetCountOnly(true)

	q, _, err := store.settingSelectQuery(options)

	if err != nil {
		return -1, err
	}

	sqlStr, params, errSql := q.Prepared(true).
		Limit(1).
		Select(goqu.COUNT(goqu.Star()).As("count")).
		ToSQL()

	if errSql != nil {
		return -1, errSql
	}

	if store.debugEnabled {
		log.Println(sqlStr)
	}

	db := sb.NewDatabase(store.db, store.dbDriverName)
	mapped, err := db.SelectToMapString(sqlStr, params...)
	if err != nil {
		return -1, err
	}

	if len(mapped) < 1 {
		return -1, nil
	}

	countStr := mapped[0]["count"]

	i, err := strconv.ParseInt(countStr, 10, 64)

	if err != nil {
		return -1, err

	}

	return i, nil
}

func (st *store) SettingCreate(ctx context.Context, setting SettingInterface) error {
	if setting == nil {
		return errors.New("settingstore > setting create. setting cannot be nil")
	}

	if setting.GetKey() == "" {
		return errors.New("settingstore > setting create. key cannot be empty")
	}

	if setting.GetCreatedAt() == "" {
		setting.SetCreatedAt(carbon.Now(carbon.UTC).ToDateTimeString())
	}

	if setting.GetUpdatedAt() == "" {
		setting.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString())
	}

	if setting.GetSoftDeletedAt() == "" {
		setting.SetSoftDeletedAt(sb.MAX_DATETIME)
	}

	data := setting.Data()

	sqlStr, sqlParams, sqlErr := goqu.Dialect(st.dbDriverName).
		Insert(st.settingTableName).
		Prepared(true).
		Rows(data).
		ToSQL()

	if sqlErr != nil {
		return sqlErr
	}

	st.logSql("create", sqlStr, sqlParams)

	_, err := st.db.Exec(sqlStr, sqlParams...)

	if err != nil {
		return err
	}

	setting.MarkAsNotDirty()

	return nil
}

// SettingDelete deletes a setting
func (store *store) SettingDelete(ctx context.Context, setting SettingInterface) error {
	if setting == nil {
		return errors.New("setting is nil")
	}

	return store.SettingDeleteByID(ctx, setting.GetID())
}

// SettingDeleteByID deletes a setting by id
func (store *store) SettingDeleteByID(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("setting id is empty")
	}

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Delete(store.settingTableName).
		Prepared(true).
		Where(goqu.C(COLUMN_ID).Eq(id)).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	store.logSql("delete", sqlStr, params...)

	_, err := store.db.Exec(sqlStr, params...)

	return err
}

// SettingDeleteByID deletes a setting by id
func (store *store) SettingDeleteByKey(ctx context.Context, settingKey string) error {
	if settingKey == "" {
		return errors.New("setting id is empty")
	}

	sqlStr, params, errSql := goqu.Dialect(store.dbDriverName).
		Delete(store.settingTableName).
		Prepared(true).
		Where(goqu.C(COLUMN_SETTING_KEY).Eq(settingKey)).
		ToSQL()

	if errSql != nil {
		return errSql
	}

	store.logSql("delete", sqlStr, params...)

	_, err := store.db.Exec(sqlStr, params...)

	return err
}

// SettingFindByID finds a setting by id
func (store *store) SettingFindByID(ctx context.Context, settingID string) (SettingInterface, error) {
	if settingID == "" {
		return nil, errors.New("setting store > find by id: setting id is required")
	}

	query := SettingQuery().
		SetID(settingID).
		SetLimit(1)

	list, err := store.SettingList(ctx, query)

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return list[0], nil
	}

	return nil, nil
}

// SettingFindByKey finds a setting by key
func (store *store) SettingFindByKey(ctx context.Context, settingKey string) (SettingInterface, error) {
	if settingKey == "" {
		return nil, errors.New("setting store > find by key: setting key is required")
	}

	query := SettingQuery().
		SetKey(settingKey).
		SetLimit(1)

	list, err := store.SettingList(ctx, query)

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return list[0], nil
	}

	return nil, nil
}

func (store *store) SettingList(ctx context.Context, query SettingQueryInterface) ([]SettingInterface, error) {
	if query == nil {
		return []SettingInterface{}, errors.New("at setting list > setting query is nil")
	}

	if err := query.Validate(); err != nil {
		return []SettingInterface{}, err
	}

	q, columns, err := store.settingSelectQuery(query)

	if err != nil {
		return []SettingInterface{}, err
	}

	sqlStr, sqlParams, errSql := q.Prepared(true).Select(columns...).ToSQL()

	if errSql != nil {
		return []SettingInterface{}, errSql
	}

	store.logSql("list", sqlStr, sqlParams...)

	if store.db == nil {
		return []SettingInterface{}, errors.New("settingstore: database is nil")
	}

	db := sb.NewDatabase(store.db, store.dbDriverName)

	if db == nil {
		return []SettingInterface{}, errors.New("settingstore: database is nil")
	}

	modelMaps, err := db.SelectToMapString(sqlStr, sqlParams...)

	if err != nil {
		return []SettingInterface{}, err
	}

	list := []SettingInterface{}

	lo.ForEach(modelMaps, func(modelMap map[string]string, index int) {
		model := NewSettingFromExistingData(modelMap)
		list = append(list, model)
	})

	return list, nil
}

func (store *store) SettingSoftDelete(ctx context.Context, setting SettingInterface) error {
	if setting == nil {
		return errors.New("setting is nil")
	}

	setting.SetSoftDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	return store.SettingUpdate(ctx, setting)
}

func (store *store) SettingSoftDeleteByID(ctx context.Context, id string) error {
	setting, err := store.SettingFindByID(ctx, id)

	if err != nil {
		return err
	}

	return store.SettingSoftDelete(ctx, setting)
}

func (store *store) SettingUpdate(ctx context.Context, setting SettingInterface) error {
	if setting == nil {
		return errors.New("settingstore > setting update. setting cannot be nil")
	}

	if store.db == nil {
		return errors.New("settingstore > setting update. db cannot be nil")
	}

	setting.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	dataChanged := setting.DataChanged()

	if len(dataChanged) == 0 {
		return nil
	}

	delete(dataChanged, COLUMN_ID) // ID cannot be updated

	// fields := map[string]interface{}{}
	// fields[COLUMN_SETTING_VALUE] = setting.GetValue()
	// fields[COLUMN_EXPIRES_AT] = setting.GetExpiresAt()
	// fields[COLUMN_UPDATED_AT] = carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)

	// wheres := []goqu.Expression{
	// 	goqu.C(COLUMN_SETTING_KEY).Eq(setting.GetKey()),
	// 	goqu.C(COLUMN_EXPIRES_AT).Gte(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)),
	// 	goqu.C(COLUMN_SOFT_DELETED_AT).Gte(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)),
	// 	goqu.C(COLUMN_USER_AGENT).Eq(options.UserAgent),
	// 	goqu.C(COLUMN_IP_ADDRESS).Eq(options.IPAddress),
	// }

	// // Only add the condition, if specifically requested
	// if len(options.UserID) > 0 {
	// 	wheres = append(wheres, goqu.C(COLUMN_USER_ID).Eq(options.UserID))
	// }

	sqlStr, sqlParams, sqlErr := goqu.Dialect(store.dbDriverName).
		Update(store.settingTableName).
		Prepared(true).
		Where(goqu.C(COLUMN_SETTING_KEY).Eq(setting.GetKey())).
		Where(goqu.C(COLUMN_ID).Eq(setting.GetID())).
		Set(dataChanged).
		ToSQL()

	if sqlErr != nil {
		return sqlErr
	}

	store.logSql("update", sqlStr, sqlParams...)

	_, err := store.db.Exec(sqlStr, sqlParams...)

	if err != nil {
		return err
	}

	return nil
}

// Set is a shortcut method to save a value by key, use Get to extract
//
// It is a convenience method which wraps SettingFindByKey,
// and then SettingCreate or SettingUpdate
//
// Parameters:
// - ctx: the context
// - settingKey: the key of the setting to save
// - value: the value to save
//
// Returns:
// - error - nil if no error, error otherwise
func (st *store) Set(ctx context.Context, settingKey string, value string) error {
	setting, errFindByKey := st.SettingFindByKey(ctx, settingKey)

	if errFindByKey != nil {
		return errFindByKey
	}

	if setting == nil {
		newSetting := NewSetting().
			SetKey(settingKey).
			SetValue(value)

		return st.SettingCreate(ctx, newSetting)
	} else {
		setting.SetValue(value)

		return st.SettingUpdate(ctx, setting)
	}
}

// SetAny is a shortcut method to save any value by key, use GetAny to extract
//
// It is a convenience method which wraps SettingCreate or SettingUpdate
// and uses JSON to serialize the data
//
// Parameters:
// - ctx: the context
// - settingKey: the key of the setting to save
// - value: the value to save
//
// Returns:
// - error - nil if no error, error otherwise
func (st *store) SetAny(ctx context.Context, key string, value interface{}, seconds int64) error {
	jsonValue, jsonError := json.Marshal(value)
	if jsonError != nil {
		return jsonError
	}

	return st.Set(ctx, key, string(jsonValue))
}

// SetMap is a shortcut method to save a map by key, use GetMap to extract
//
// It is a convenience method which wraps SettingCreate or SettingUpdate
// to save a map by key
//
// Parameters:
// - ctx: the context
// - settingKey: the key of the setting to save
// - value: the value to save
//
// Returns:
// - error - nil if no error, error otherwise
func (st *store) SetMap(ctx context.Context, key string, value map[string]any) error {
	jsonValue, jsonError := json.Marshal(value)

	if jsonError != nil {
		return jsonError
	}

	return st.Set(ctx, key, string(jsonValue))
}

// settingSelectQuery builds the select query
//
// Parameters:
// - options: the options for the query
//
// Returns:
// - selectDataset: the select dataset
// - columns: the columns to select
// - err: the error
func (store *store) settingSelectQuery(options SettingQueryInterface) (selectDataset *goqu.SelectDataset, columns []any, err error) {
	if options == nil {
		return nil, []any{}, errors.New("setting query: cannot be nil")
	}

	if err := options.Validate(); err != nil {
		return nil, []any{}, err
	}

	q := goqu.Dialect(store.dbDriverName).From(store.settingTableName)

	if options.HasID() {
		q = q.Where(goqu.C(COLUMN_ID).Eq(options.ID()))
	}

	if options.HasIDIn() {
		q = q.Where(goqu.C(COLUMN_ID).In(options.IDIn()))
	}

	if options.HasKey() {
		q = q.Where(goqu.C(COLUMN_SETTING_KEY).Eq(options.Key()))
	}

	if !options.IsCountOnly() {
		if options.HasLimit() {
			q = q.Limit(uint(options.Limit()))
		}

		if options.HasOffset() {
			q = q.Offset(uint(options.Offset()))
		}
	}

	sortOrder := sb.DESC
	if options.HasSortOrder() && options.SortOrder() != "" {
		sortOrder = options.SortOrder()
	}

	if options.HasOrderBy() && options.OrderBy() != "" {
		if strings.EqualFold(sortOrder, sb.ASC) {
			q = q.Order(goqu.I(options.OrderBy()).Asc())
		} else {
			q = q.Order(goqu.I(options.OrderBy()).Desc())
		}
	}

	columns = []any{}

	for _, column := range options.Columns() {
		columns = append(columns, column)
	}

	if options.SoftDeletedIncluded() {
		return q, columns, nil // soft deleted settings requested specifically
	}

	softDeleted := goqu.C(COLUMN_SOFT_DELETED_AT).
		Gt(carbon.Now(carbon.UTC).ToDateTimeString())

	return q.Where(softDeleted), columns, nil
}

// logSql logs the sql statement to the provided logger
//
// # It only logs the sql if debug is enabled, otherwise it does nothing
//
// Parameters:
// - sqlOperationType: the type of sql operation
// - sqlString: the sql statement
// - sqlParams: the sql parameters
func (store *store) logSql(sqlOperationType string, sqlString string, sqlParams ...interface{}) {
	if !store.debugEnabled {
		return
	}

	if store.sqlLogger != nil {
		store.sqlLogger.Debug("sql: "+sqlOperationType, slog.String("sql", sqlString), slog.Any("params", sqlParams))
	}
}
