package settingstore

import (
	"os"
	"testing"

	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB(filepath string) *sql.DB {
	os.Remove(filepath) // remove database
	dsn := filepath + "?parseTime=true"
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		panic(err)
	}

	return db
}

func InitStore() *Store {
	db := InitDB("test_settingstore.db")
	return &Store{
		settingsTableName:  "test_settingsTableName.db",
		dbDriverName:       "sql",
		db:                 db,
		debug:              false,
		automigrateEnabled: false,
	}
}

func TestWithAutoMigrate(t *testing.T) {
	db := InitDB("test_settingsTableName.db")

	s := Store{
		settingsTableName:  "test_settingsTableName.db",
		dbDriverName:       "sql",
		db:                 db,
		debug:              false,
		automigrateEnabled: false,
	}

	f := WithAutoMigrate(true)
	f(&s)

	if s.automigrateEnabled != true {
		t.Fatalf("automigrateEnabled: Expected [true] received [%v]", s.automigrateEnabled)
	}

	s = Store{
		settingsTableName:  "test_settingsTableName.db",
		dbDriverName:       "sql",
		db:                 db,
		debug:              false,
		automigrateEnabled: true,
	}

	f = WithAutoMigrate(false)
	f(&s)

	if s.automigrateEnabled == true {
		t.Fatalf("automigrateEnabled: Expected [true] received [%v]", s.automigrateEnabled)
	}
}

func TestWithDb(t *testing.T) {
	db := InitDB("test")
	s := Store{
		settingsTableName:  "test_settingsTableName.db",
		dbDriverName:       "sql",
		db:                 db,
		debug:              false,
		automigrateEnabled: true,
	}

	f := WithDb(db)
	f(&s)

	if s.db == nil {
		t.Fatalf("DB: Expected Initialized DB, received [%v]", s.db)
	}

}

func TestWithTableName(t *testing.T) {
	s := Store{
		settingsTableName:  "test_settingsTableName.db",
		dbDriverName:       "sql",
		db:                 nil,
		debug:              false,
		automigrateEnabled: true,
	}
	table_name := "Table1"
	f := WithTableName(table_name)
	f(&s)
	if s.settingsTableName != table_name {
		t.Fatalf("Expected logTableName [%v], received [%v]", table_name, s.settingsTableName)
	}
	table_name = "Table2"
	f = WithTableName(table_name)
	f(&s)
	if s.settingsTableName != table_name {
		t.Fatalf("Expected logTableName [%v], received [%v]", table_name, s.settingsTableName)
	}
}

func Test_Store_AutoMigrate(t *testing.T) {
	db := InitDB("test_settingsTableName.db")

	s, _ := NewStore(WithDb(db), WithTableName("log_with_automigrate"), WithAutoMigrate(true))

	s.AutoMigrate()

	if s.settingsTableName != "log_with_automigrate" {
		t.Fatalf("Expected logTableName [log_with_automigrate] received [%v]", s.settingsTableName)
	}
	if s.db == nil {
		t.Fatalf("DB Init Failure")
	}
	if s.automigrateEnabled != true {
		t.Fatalf("Failure:  WithAutoMigrate")
	}
}

func Test_Store_Set(t *testing.T) {
	db := InitDB("test_settingsTableName.db")
	s, _ := NewStore(WithDb(db), WithTableName("log_with_automigrate"), WithAutoMigrate(true))

	key := "1234z"
	val := "123zx"
	ok, _ := s.Set(key, val)

	if !ok {
		t.Fatalf("Failure: Set")
	}
}

func Test_Store_SetJSON(t *testing.T) {
	db := InitDB("test_settingsTableName.db")
	s, _ := NewStore(WithDb(db), WithTableName("log_with_automigrate"), WithAutoMigrate(true))

	key := "1234z"
	val := `{"a" : "b", "c" : "d"}`
	ok, _ := s.SetJSON(key, val)

	if !ok {
		t.Fatalf("Failure: Set")
	}
}

func Test_Store_Remove(t *testing.T) {
	db := InitDB("test_settingsTableName.db")
	s, _ := NewStore(WithDb(db), WithTableName("settings_test_autoremove"), WithAutoMigrate(true))

	key := "1234z"
	val := "123zx"
	ok, _ := s.Set(key, val)

	if !ok {
		t.Fatalf("Failure: Set")
	}

	s.Remove(key)
	ret, err := s.Get(key, "default")
	if err != nil {
		t.Fatalf("No errors are expected but %s", err.Error())
	}
	if ret != "default" {
		t.Fatalf("Unable to delete!!! Entry Persists")
	}
}

func Test_Store_Get(t *testing.T) {
	db := InitDB("test_settingsTableName.db")
	s, _ := NewStore(WithDb(db), WithTableName("log_with_automigrate"), WithAutoMigrate(true))

	key := "1234z"
	val := "123zx"
	ok, _ := s.Set(key, val)

	if !ok {
		t.Fatalf("Failure: Set")
	}

	ret, err := s.Get(key, "default")

	if err != nil {
		t.Fatalf("No errors are expected but %s", err.Error())
	}

	if ret != val {
		t.Fatalf("Unable to Get: Expected [%v] Received [%v]", val, ret)
	}
}

func Test_Store_FindByKey(t *testing.T) {
	db := InitDB("test_settingsTableName.db")
	s, _ := NewStore(WithDb(db), WithTableName("log_with_automigrate"), WithAutoMigrate(true))

	key := "1234z"
	val := "123zx"
	ok, _ := s.Set(key, val)

	if !ok {
		t.Fatalf("Failure: Set")
	}

	meta, err := s.FindByKey(key)
	if err != nil {
		t.Fatalf("No errors are expected but %s", err.Error())
	}

	if meta == nil {
		t.Fatalf("NIL Record Received")
	}
}

func Test_Store_GetJSON(t *testing.T) {
	db := InitDB("test_GetJSON.db")
	s, _ := NewStore(WithDb(db), WithTableName("setting_get_json"), WithAutoMigrate(true))

	key := "1234z"
	val := `{"a" : "b", "c" : "d"}`
	ok, _ := s.SetJSON(key, val)

	if !ok {
		t.Fatalf("Failure: Set")
	}

	ret, err := s.GetJSON(key, nil)

	if err != nil {
		t.Fatalf("No errors are expected but error thrown: %s", err.Error())
	}

	if ret == nil {
		t.Fatalf("Failure getting JSON value")
	}

	if ret != val {
		t.Fatalf("Retrieved value not the same as set value")
	}
}
