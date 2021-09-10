package settingstore

import (
	"database/sql"
	"encoding/json"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/gouniverse/uid"
)

// Store defines a session store
type Store struct {
	settingsTableName  string
	dbDriverName       string
	db                 *sql.DB
	automigrateEnabled bool
}

// StoreOption options for the vault store
type StoreOption func(*Store)

// WithAutoMigrate sets the table name for the cache store
func WithAutoMigrate(automigrateEnabled bool) StoreOption {
	return func(s *Store) {
		s.automigrateEnabled = automigrateEnabled
	}
}

// WithDb sets the database for the setting store
func WithDb(db *sql.DB) StoreOption {
	return func(s *Store) {
		s.db = db
		s.dbDriverName = s.DriverName(s.db)
	}
}

// WithTableName sets the table name for the cache store
func WithTableName(settingsTableName string) StoreOption {
	return func(s *Store) {
		s.settingsTableName = settingsTableName
	}
}

// NewStore creates a new setting store
func NewStore(opts ...StoreOption) *Store {
	store := &Store{}
	for _, opt := range opts {
		opt(store)
	}

	if store.settingsTableName == "" {
		log.Panic("Setting store: settingTableName is required")
	}

	if store.automigrateEnabled {
		store.AutoMigrate()
	}

	return store
}

// AutoMigrate auto migrate
func (st *Store) AutoMigrate() error {

	sql := st.SqlCreateTable()

	_, err := st.db.Exec(sql)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// DriverName finds the driver name from database
func (st *Store) DriverName(db *sql.DB) string {
	dv := reflect.ValueOf(db.Driver())
	driverFullName := dv.Type().String()
	if strings.Contains(driverFullName, "mysql") {
		return "mysql"
	}
	if strings.Contains(driverFullName, "postgres") || strings.Contains(driverFullName, "pq") {
		return "postgres"
	}
	if strings.Contains(driverFullName, "sqlite") {
		return "sqlite"
	}
	if strings.Contains(driverFullName, "mssql") {
		return "mssql"
	}
	return driverFullName
}

// FindByKey finds a session by key
func (st *Store) FindByKey(key string) *Setting {
	setting := &Setting{}
	sqlStr, _, _ := goqu.From(st.settingsTableName).Where(goqu.C("setting_key").Eq(key), goqu.C("deleted_at").IsNull()).Select(Setting{}).ToSQL()

	log.Println(sqlStr)

	err := st.db.QueryRow(sqlStr).Scan(&setting.CreatedAt, &setting.DeletedAt, &setting.ID, &setting.Key, &setting.Value, &setting.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		log.Fatal("Failed to execute query: ", err)
		return nil
	}

	return setting
}

// Get gets a key from settings
func (st *Store) Get(key string, valueDefault string) string {
	setting := st.FindByKey(key)

	if setting != nil {
		return setting.Value
	}

	return valueDefault
}

// GetJSON gets a JSON key from setting
func (st *Store) GetJSON(key string, valueDefault interface{}) interface{} {
	setting := st.FindByKey(key)

	if setting != nil {
		jsonValue := setting.Value
		var e interface{}
		jsonError := json.Unmarshal([]byte(jsonValue), e)
		if jsonError != nil {
			return valueDefault
		}

		return e
	}

	return valueDefault
}

// Keys gets all keys sorted alphabetically
func (st *Store) Keys() ([]string, error) {
	keys := []string{}
	// settingList := []Setting{}
	// sqlStr, _, _ := goqu.From(st.settingsTableName).Order(goqu.I("setting_key").Asc()).Where(goqu.C("deleted_at").IsNull()).Select(Setting{}).ToSQL()
	sqlStr, _, _ := goqu.From(st.settingsTableName).Order(goqu.I("setting_key").Asc()).Where(goqu.C("deleted_at").IsNull()).Select("setting_value").ToSQL()

	log.Println(sqlStr)

	rows, err := st.db.Query(sqlStr)

	if err != nil {
		return keys, err
	}

	for rows.Next() {
		var value string
		err := rows.Scan(&value)
		if err != nil {
			return keys, err
		}
		// settingList = append(settingList, value)
		keys = append(keys, value)
	}

	// for _, setting := range settingList {
	// 	keys = append(keys, setting.Key)
	// }

	return keys, nil
}

// Remove gets a JSON key from cache
func (st *Store) Remove(key string) error {
	sqlStr, _, _ := goqu.From(st.settingsTableName).Where(goqu.C("setting_key").Eq(key), goqu.C("deleted_at").IsNull()).Delete().ToSQL()

	log.Println(sqlStr)

	_, err := st.db.Exec(sqlStr)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		log.Fatal("Failed to execute query: ", err)
		return nil
	}

	return nil
}

// Set sets new key value pair
func (st *Store) Set(key string, value string) (bool, error) {
	setting := st.FindByKey(key)

	// log.Println(setting)

	var sqlStr string
	if setting == nil {
		var newSetting = Setting{
			ID:        uid.MicroUid(),
			Key:       key,
			Value:     value,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		sqlStr, _, _ = goqu.Insert(st.settingsTableName).Rows(newSetting).ToSQL()
	} else {
		setting.Value = value
		setting.UpdatedAt = time.Now()
		sqlStr, _, _ = goqu.Update(st.settingsTableName).Set(setting).ToSQL()
	}

	// log.Println(sqlStr)

	_, err := st.db.Exec(sqlStr)

	if err != nil {
		log.Println(err)
		return false, err
	}

	return true, nil
}

// SetJSON sets new key JSON value pair
func (st *Store) SetJSON(key string, value interface{}) (bool, error) {
	jsonValue, jsonError := json.Marshal(value)
	if jsonError != nil {
		return false, jsonError
	}

	return st.Set(key, string(jsonValue))
}

// SqlCreateTable returns a SQL string for creating the setting table
func (st *Store) SqlCreateTable() string {
	sqlMysql := `
	CREATE TABLE IF NOT EXISTS ` + st.settingsTableName + ` (
	  id varchar(40) NOT NULL PRIMARY KEY,
	  setting_key varchar(40) NOT NULL,
	  setting_value text,
	  created_at datetime NOT NULL,
	  updated_at datetime NOT NULL,
	  deleted_at datetime
	);
	`

	sqlPostgres := `
	CREATE TABLE IF NOT EXISTS "` + st.settingsTableName + `" (
	  "id" varchar(40) NOT NULL PRIMARY KEY,
	  "setting_key" varchar(40) NOT NULL,
	  "setting_value" text,
	  "created_at" timestamptz(6) NOT NULL,
	  "updated_at" timestamptz(6) NOT NULL,
	  "deleted_at" timestamptz(6)
	)
	`

	sqlSqlite := `
	CREATE TABLE IF NOT EXISTS "` + st.settingsTableName + `" (
	  "id" varchar(40) NOT NULL PRIMARY KEY,
	  "setting_key" varchar(40) NOT NULL,
	  "setting_value" text,
	  "created_at" timestamptz(6) NOT NULL,
	  "updated_at" timestamptz(6) NOT NULL,
	  "deleted_at" timestamptz(6)
	)
	`

	sql := "unsupported driver " + st.dbDriverName

	if st.dbDriverName == "mysql" {
		sql = sqlMysql
	}
	if st.dbDriverName == "postgres" {
		sql = sqlPostgres
	}
	if st.dbDriverName == "sqlite" {
		sql = sqlSqlite
	}

	return sql
}
