package settingstore

import (
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
	"github.com/georgysavva/scany/sqlscan"
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

// AutoMigrate auto migrate
func (store *store) AutoMigrate() error {
	sqlStr := store.SQLCreateTable()

	if sqlStr == "" {
		return errors.New("setting store: table create sql is empty")
	}

	if store.db == nil {
		return errors.New("setting store: database is nil")
	}

	_, err := store.db.Exec(sqlStr)

	if err != nil {
		return err
	}

	return nil
}

// EnableDebug - enables the debug option
func (st *store) EnableDebug(debug bool) {
	st.debugEnabled = debug
}

// Delete deletes a setting
func (st *store) Delete(settingKey string) error {
	wheres := []goqu.Expression{
		goqu.C(COLUMN_SETTING_KEY).Eq(settingKey),
	}

	sqlStr, sqlParams, err := goqu.Dialect(st.dbDriverName).
		From(st.settingTableName).
		Where(wheres...).
		Delete().
		Prepared(true).
		ToSQL()

	if err != nil {
		return err
	}

	st.logSql("delete", sqlStr, sqlParams)

	_, err = st.db.Exec(sqlStr, sqlParams...)

	if err != nil {
		if err == sql.ErrNoRows {
			// Looks like this is now outdated for sqlscan
			return nil
		}

		if sqlscan.NotFound(err) {
			return nil
		}

		return err
	}

	return nil
}

// FindByKey finds a setting by key
func (store *store) FindByKey(settingKey string) (SettingInterface, error) {
	if settingKey == "" {
		return nil, errors.New("setting store > find by key: setting key is required")
	}

	query := SettingQuery().
		SetKey(settingKey).
		SetLimit(1)

	list, err := store.SettingList(query)

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return list[0], nil
	}

	return nil, nil
}

// Gets the setting value as a string
func (st *store) Get(settingKey string, valueDefault string) (string, error) {
	setting, errFindByKey := st.FindByKey(settingKey)

	if errFindByKey != nil {
		return "", errFindByKey
	}

	if setting != nil {
		return setting.GetValue(), nil
	}

	return valueDefault, nil
}

// GetAny attempts to parse the value as interface, use with SetAny
func (st *store) GetAny(key string, valueDefault interface{}) (interface{}, error) {
	setting, errFindByKey := st.FindByKey(key)

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

// GetMap attempts to parse the value as map[string]any, use with SetMap
func (st *store) GetMap(key string, valueDefault map[string]any) (map[string]any, error) {
	setting, errFindByKey := st.FindByKey(key)

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

// Has finds if a setting by key exists
func (store *store) Has(settingKey string) (bool, error) {
	if settingKey == "" {
		return false, errors.New("setting store > find by key: setting key is required")
	}

	query := SettingQuery().
		SetKey(settingKey).
		SetLimit(1)

	count, err := store.SettingCount(query)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (st *store) MergeMap(key string, mergeMap map[string]any, seconds int64) error {
	currentMap, err := st.GetMap(key, nil)

	if err != nil {
		return err
	}

	if currentMap == nil {
		return errors.New("settingstore. nil found")
	}

	for mapKey, mapValue := range mergeMap {
		currentMap[mapKey] = mapValue
	}

	return st.SetMap(key, currentMap, seconds)
}

func (store *store) SettingCount(options SettingQueryInterface) (int64, error) {
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

func (st *store) SettingCreate(setting SettingInterface) error {
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
func (store *store) SettingDelete(setting SettingInterface) error {
	if setting == nil {
		return errors.New("setting is nil")
	}

	return store.SettingDeleteByID(setting.GetID())
}

// SettingDeleteByID deletes a setting by id
func (store *store) SettingDeleteByID(id string) error {
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
func (store *store) SettingDeleteByKey(settingKey string) error {
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
func (store *store) SettingFindByID(settingID string) (SettingInterface, error) {
	if settingID == "" {
		return nil, errors.New("setting store > find by id: setting id is required")
	}

	query := SettingQuery().
		SetID(settingID).
		SetLimit(1)

	list, err := store.SettingList(query)

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return list[0], nil
	}

	return nil, nil
}

// SettingFindByKey finds a setting by key
func (store *store) SettingFindByKey(settingKey string) (SettingInterface, error) {
	if settingKey == "" {
		return nil, errors.New("setting store > find by key: setting key is required")
	}

	query := SettingQuery().
		SetKey(settingKey).
		SetLimit(1)

	list, err := store.SettingList(query)

	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return list[0], nil
	}

	return nil, nil
}

func (store *store) SettingList(query SettingQueryInterface) ([]SettingInterface, error) {
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

func (store *store) SettingSoftDelete(setting SettingInterface) error {
	if setting == nil {
		return errors.New("setting is nil")
	}

	setting.SetSoftDeletedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

	return store.SettingUpdate(setting)
}

func (store *store) SettingSoftDeleteByID(id string) error {
	setting, err := store.SettingFindByID(id)

	if err != nil {
		return err
	}

	return store.SettingSoftDelete(setting)
}

func (store *store) SettingUpdate(setting SettingInterface) error {
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

// Set sets a key in store
func (st *store) Set(settingKey string, value string, seconds int64) error {
	setting, errFindByKey := st.FindByKey(settingKey)

	if errFindByKey != nil {
		return errFindByKey
	}

	if setting == nil {
		newSetting := NewSetting().
			SetKey(settingKey).
			SetValue(value)

		return st.SettingCreate(newSetting)
	} else {
		setting.SetValue(value)
		setting.SetUpdatedAt(carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC))

		return st.SettingUpdate(setting)
	}
}

// SetAny convenience method which saves the supplied interface value, use GetAny to extract
// Internally it serializes the data to JSON
func (st *store) SetAny(key string, value interface{}, seconds int64) error {
	jsonValue, jsonError := json.Marshal(value)
	if jsonError != nil {
		return jsonError
	}

	return st.Set(key, string(jsonValue), seconds)
}

// SetMap convenience method which saves the supplied map, use GetMap to extract
func (st *store) SetMap(key string, value map[string]any, seconds int64) error {
	jsonValue, jsonError := json.Marshal(value)
	if jsonError != nil {
		return jsonError
	}

	return st.Set(key, string(jsonValue), seconds)
}

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

func (store *store) logSql(sqlOperationType string, sql string, params ...interface{}) {
	if !store.debugEnabled {
		return
	}

	if store.sqlLogger != nil {
		store.sqlLogger.Debug("sql: "+sqlOperationType, slog.String("sql", sql), slog.Any("params", params))
	}
}
