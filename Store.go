package settingstore

import (
	"database/sql"
	"errors"
	"log"
	"strconv"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"     // importing mysql dialect
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"  // importing postgres dialect
	_ "github.com/doug-martin/goqu/v9/dialect/sqlite3"   // importing sqlite3 dialect
	_ "github.com/doug-martin/goqu/v9/dialect/sqlserver" // importing sqlserver dialect
	"github.com/dromara/carbon/v2"
	"github.com/gouniverse/sb"
	"github.com/samber/lo"
)

// == INTERFACE ===============================================================

var _ StoreInterface = (*store)(nil) // verify it extends the store interface

// == TYPE ====================================================================

// Store defines a setting store
type store struct {
	settingsTableName  string
	db                 *sql.DB
	dbDriverName       string
	timeoutSeconds     int64
	automigrateEnabled bool
	debugEnabled       bool
	// sqlLogger          *slog.Logger // slog is not defined
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

func (store *store) SettingCount() (int64, error) {
	q, _, err := store.settingSelectQuery(SettingQuery())

	if err != nil {
		return -1, err
	}

	sqlStr, params, errSql := q.Prepared(true).
		Limit(1).
		Select(goqu.COUNT(goqu.Star()).As("count")).
		ToSQL()

	if errSql != nil {
		return -1, nil
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

	data := setting.Data()

	sqlStr, sqlParams, sqlErr := goqu.Dialect(st.dbDriverName).
		Insert(st.settingsTableName).
		Prepared(true).
		Rows(data).
		ToSQL()

	if sqlErr != nil {
		return sqlErr
	}

	st.logSql("create", sqlStr, sqlParams...)

	_, err := st.db.Exec(sqlStr, sqlParams...)

	if err != nil {
		return err
	}

	// setting.MarkAsNotDirty() // MarkAsNotDirty is not defined

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
		Delete(store.settingsTableName).
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

// SettingFindByID finds a setting by id
func (store *store) SettingFindByID(settingID string) (SettingInterface, error) {
	if settingID == "" {
		return nil, errors.New("setting store > find by id: setting id is required")
	}

	query := SettingQuery()
	query.SetID(settingID)
	query.SetLimit(1)

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

	q, columns, err := store.settingSelectQuery(query)

	if err != nil {
		return []SettingInterface{}, err
	}

	sqlStr, sqlParams, errSql := q.Prepared(true).Select(columns...).ToSQL()

	if errSql != nil {
		return []SettingInterface{}, nil
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

	sqlStr, sqlParams, sqlErr := goqu.Dialect(store.dbDriverName).
		Update(store.settingsTableName).
		Prepared(true).
		Where(goqu.C(COLUMN_KEY).Eq(setting.GetKey())).
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

func (store *store) settingSelectQuery(options SettingQueryInterface) (selectDataset *goqu.SelectDataset, columns []any, err error) {
	if options == nil {
		return nil, []any{}, errors.New("setting query: cannot be nil")
	}

	q := goqu.Dialect(store.dbDriverName).From(store.settingsTableName)

	columns = []any{}

	for _, column := range options.Columns() {
		columns = append(columns, column)
	}

	return q, columns, nil
}

func (store *store) logSql(sqlOperationType string, sql string, params ...interface{}) {
	if !store.debugEnabled {
		return
	}

	log.Println("sql: "+sqlOperationType, "sql", sql, "params", params)
	// if store.sqlLogger != nil {
	// 	store.sqlLogger.Debug("sql: "+sqlOperationType, slog.String("sql", sql), slog.Any("params", params))
	// } // slog is not defined
}
