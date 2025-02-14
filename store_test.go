package settingstore

import (
	"database/sql"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/gouniverse/sb"
	"github.com/gouniverse/utils"
	_ "github.com/mattn/go-sqlite3"
)

func initDB(filepath string) (*sql.DB, error) {
	if filepath != ":memory:" && utils.FileExists(filepath) {
		err := os.Remove(filepath) // remove database

		if err != nil {
			return nil, err
		}
	}

	dsn := filepath + "?parseTime=true"
	db, err := sql.Open("sqlite3", dsn)

	if err != nil {
		return nil, err
	}

	return db, nil
}

func initStore(filepath string) (StoreInterface, error) {
	db, err := initDB(filepath)

	if err != nil {
		return nil, err
	}

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		SettingTableName:   "setting",
		AutomigrateEnabled: true,
	})

	if err != nil {
		return nil, err
	}

	if store == nil {
		return nil, errors.New("unexpected nil store")
	}

	return store, nil
}

func TestStore_Create(t *testing.T) {
	store, err := initStore(":memory:")

	if err != nil {
		t.Fatal("Store could not be created: ", err.Error())
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}
}

func TestStore_Automigrate(t *testing.T) {
	store, err := initStore(":memory:")

	if err != nil {
		t.Fatal("Store could not be created: ", err.Error())
	}

	err = store.AutoMigrate()

	if err != nil {
		t.Fatal("Automigrate failed: " + err.Error())
	}
}

func TestStore_EnableDebug(t *testing.T) {
	store, err := initStore(":memory:")

	if err != nil {
		t.Fatal("Store could not be created: ", err.Error())
	}

	store.EnableDebug(true)

	err = store.AutoMigrate()

	if err != nil {
		t.Fatal("Automigrate failed: " + err.Error())
	}
}

// func TestSetGetMap(t *testing.T) {
// 	store, err := initStore(":memory:")

// 	if err != nil {
// 		t.Fatal("Store could not be created: ", err.Error())
// 	}

// 	value := map[string]any{
// 		"key1": "value1",
// 		"key2": "value2",
// 		"key3": "value3",
// 	}
// 	err = store.SetMap("mykey", value, 5, SettingOptions{})

// 	if err != nil {
// 		t.Fatalf("Set Map failed: " + err.Error())
// 	}

// 	result, err := store.GetMap("mykey", nil, SettingOptions{})

// 	if err != nil {
// 		t.Fatalf("Get JSON failed: " + err.Error())
// 	}

// 	if result == nil {
// 		t.Fatalf("GetMap failed: nil returned")
// 	}

// 	if result["key1"].(string) != value["key1"] {
// 		t.Fatalf("Key1 not correct: " + result["key1"].(string))
// 	}

// 	if result["key2"] != value["key2"] {
// 		t.Fatalf("Key2 not correct: " + result["key2"].(string))
// 	}

// 	if result["key3"] != value["key3"] {
// 		t.Fatalf("Key3 not correct: " + result["key3"].(string))
// 	}
// }

// func TestMergeMap(t *testing.T) {
// 	store, err := initStore(":memory:")

// 	if err != nil {
// 		t.Fatal("Store could not be created: ", err.Error())
// 	}

// 	value := map[string]any{
// 		"key1": "value1",
// 		"key2": "value2",
// 		"key3": "value3",
// 	}
// 	err = store.SetMap("mykey", value, 600, SettingOptions{})

// 	if err != nil {
// 		t.Fatalf("Set Map failed: " + err.Error())
// 	}

// 	valueMerge := map[string]any{
// 		"key2": "value22",
// 		"key3": "value33",
// 	}

// 	err = store.MergeMap("mykey", valueMerge, 600, SettingOptions{})

// 	if err != nil {
// 		t.Fatalf("Merge Map failed: " + err.Error())
// 	}

// 	result, err := store.GetMap("mykey", nil, SettingOptions{})

// 	if err != nil {
// 		t.Fatalf("Get JSON failed: " + err.Error())
// 	}

// 	if result == nil {
// 		t.Fatalf("GetMap failed: nil returned")
// 	}

// 	if result["key1"].(string) != value["key1"] {
// 		t.Fatalf("Key1 not correct: " + result["key1"].(string))
// 	}

// 	if result["key2"].(string) != valueMerge["key2"] {
// 		t.Fatalf("Key2 not correct: " + result["key2"].(string))
// 	}

// 	if result["key3"].(string) != valueMerge["key3"] {
// 		t.Fatalf("Key3 not correct: " + result["key3"].(string))
// 	}
// }

// func TestExtend(t *testing.T) {
// 	store, err := initStore(":memory:")

// 	if err != nil {
// 		t.Fatal("Store could not be created: ", err.Error())
// 	}

// 	err = store.Set("mykey", "test", 5, SettingOptions{})

// 	if err != nil {
// 		t.Fatal("Set failed: " + err.Error())
// 	}

// 	err = store.Extend("mykey", 100, SettingOptions{})

// 	if err != nil {
// 		t.Fatal("Extend failed: " + err.Error())
// 	}

// 	settingExtended, err := store.FindByKey("mykey", SettingOptions{})

// 	if err != nil {
// 		t.Fatal("Extend failed: " + err.Error())
// 	}

// 	if settingExtended == nil {
// 		t.Fatal("Extend failed. Setting is NIL")
// 	}

// 	if settingExtended.GetValue() != "test" {
// 		t.Fatal("Extend failed. Value is wrong", settingExtended.GetValue())
// 	}

// 	diff := settingExtended.GetExpiresAtCarbon().DiffAbsInSeconds(carbon.Now(carbon.UTC))

// 	if diff < 90 {
// 		t.Fatal("Extend failed. ExpiresAt must be more than 90 seconds", settingExtended.GetExpiresAt(), diff)
// 	}

// 	if diff > 110 {
// 		t.Fatal("Extend failed. ExpiresAt must be less than 110 seconds", settingExtended.GetExpiresAt(), diff)
// 	}

// }

func TestStore_SettingCreate(t *testing.T) {
	store, err := initStore(":memory:")

	if err != nil {
		t.Fatal("Store could not be created: ", err.Error())
	}

	setting := NewSetting().
		SetKey("1").
		SetValue("one two three")

	if setting == nil {
		t.Fatal("unexpected nil setting")
	}

	if setting.GetID() == "" {
		t.Fatal("unexpected empty id:", setting.GetID())
	}

	if len(setting.GetID()) != 32 {
		t.Fatal("unexpected id length:", len(setting.GetID()))
	}

	if setting.GetKey() == "" {
		t.Fatal("unexpected empty key:", setting.GetKey())
	}

	if len(setting.GetKey()) != 1 {
		t.Fatal("unexpected key length:", len(setting.GetKey()))
	}

	err = store.SettingCreate(setting)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}
}

func TestStore_SettingDelete(t *testing.T) {
	store, err := initStore(":memory:")

	if err != nil {
		t.Fatal("Store could not be created: ", err.Error())
	}

	setting := NewSetting().
		SetKey("1").
		SetValue("one two three")

	err = store.SettingCreate(setting)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.SettingDeleteByID(setting.GetID())

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	settingFindWithDeleted, err := store.SettingList(SettingQuery().
		SetID(setting.GetID()).
		SetLimit(1).
		SetSoftDeletedIncluded(true))

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if len(settingFindWithDeleted) != 0 {
		t.Fatal("Setting MUST be deleted, but it is not")
	}
}

func TestStore_SettingDeleteByID(t *testing.T) {
	store, err := initStore(":memory:")

	if err != nil {
		t.Fatal("Store could not be created: ", err.Error())
	}

	setting := NewSetting().
		SetKey("1").
		SetValue("one two three")

	if setting == nil {
		t.Fatal("unexpected nil setting")
	}

	if setting.GetID() == "" {
		t.Fatal("unexpected empty id:", setting.GetID())
	}

	err = store.SettingCreate(setting)

	if err != nil {
		t.Error("unexpected error:", err)
	}

	err = store.SettingDeleteByID(setting.GetID())

	if err != nil {
		t.Error("unexpected error:", err)
	}

	settingFindWithDeleted, err := store.SettingList(SettingQuery().
		SetID(setting.GetID()).
		SetLimit(1).
		SetSoftDeletedIncluded(true))

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if len(settingFindWithDeleted) != 0 {
		t.Fatal("Setting MUST be deleted, but it is not")
	}
}

func TestStore_SettingFindByID(t *testing.T) {
	store, err := initStore(":memory:")

	if err != nil {
		t.Fatal("Store could not be created: ", err.Error())
	}

	setting := NewSetting().
		SetKey("1").
		SetValue("one two three four")

	if setting == nil {
		t.Fatal("unexpected nil setting")
	}

	if setting.GetID() == "" {
		t.Fatal("unexpected empty id:", setting.GetID())
	}

	err = store.SettingCreate(setting)
	if err != nil {
		t.Error("unexpected error:", err)
	}

	settingFound, errFind := store.SettingFindByID(setting.GetID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if settingFound == nil {
		t.Fatal("Setting MUST NOT be nil")
	}

	if settingFound.GetID() != setting.GetID() {
		t.Fatal("IDs do not match")
	}

	if settingFound.GetValue() != setting.GetValue() {
		t.Fatal("Values do not match")
	}

	if settingFound.GetValue() != "one two three four" {
		t.Fatal("Values do not match, expected: one two three four, got: ", settingFound.GetValue())
	}
}

func TestStore_SettingFindByKey(t *testing.T) {
	store, err := initStore(":memory:")

	if err != nil {
		t.Fatal("Store could not be created: ", err.Error())
	}

	setting := NewSetting().
		SetKey("1").
		SetValue("one two three four")

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if setting == nil {
		t.Fatal("unexpected nil setting")
	}

	if setting.GetKey() == "" {
		t.Fatal("unexpected empty key:", setting.GetKey())
	}

	err = store.SettingCreate(setting)
	if err != nil {
		t.Error("unexpected error:", err)
	}

	settingFound, errFind := store.SettingFindByKey(setting.GetKey())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if settingFound == nil {
		t.Fatal("Setting MUST NOT be nil")
	}

	if settingFound.GetID() != setting.GetID() {
		t.Fatal("IDs do not match")
	}

	if settingFound.GetValue() != setting.GetValue() {
		t.Fatal("Values do not match")
	}

	if settingFound.GetValue() != "one two three four" {
		t.Fatal("Values do not match, expected: one two three four, got: ", settingFound.GetValue())
	}
}

func TestStore_SettingList(t *testing.T) {
	store, err := initStore(":memory:")

	if err != nil {
		t.Fatal("Store could not be created: ", err.Error())
	}

	setting1 := NewSetting().
		SetKey("1").
		SetValue("one two three")

	setting2 := NewSetting().
		SetKey("2").
		SetValue("four five six")

	setting3 := NewSetting().
		SetKey("3").
		SetValue("seven eight nine")

	for _, setting := range []SettingInterface{setting1, setting2, setting3} {
		err = store.SettingCreate(setting)
		if err != nil {
			t.Error("unexpected error:", err)
		}
	}

	settingList, errList := store.SettingList(SettingQuery().
		SetKey("2").
		SetLimit(2))

	if errList != nil {
		t.Fatal("unexpected error:", errList)
	}

	if len(settingList) != 1 {
		t.Fatal("unexpected setting list length:", len(settingList))
	}
}

func TestStore_SettingSoftDelete(t *testing.T) {
	store, err := initStore(":memory:")

	if err != nil {
		t.Fatal("Store could not be created: ", err.Error())
	}

	setting := NewSetting().SetKey("1").SetValue("one two three")

	err = store.SettingCreate(setting)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.SettingSoftDeleteByID(setting.GetID())

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if setting.GetSoftDeletedAt() != sb.MAX_DATETIME {
		t.Fatal("Setting MUST NOT be soft deleted")
	}

	settingFound, errFind := store.SettingFindByID(setting.GetID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if settingFound != nil {
		t.Fatal("Setting MUST be nil")
	}

	settingFindWithSoftDeleted, err := store.SettingList(SettingQuery().
		SetID(setting.GetID()).
		SetSoftDeletedIncluded(true).
		SetLimit(1))

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if len(settingFindWithSoftDeleted) == 0 {
		t.Fatal("Exam MUST be soft deleted")
	}

	if strings.Contains(settingFindWithSoftDeleted[0].GetSoftDeletedAt(), sb.MAX_DATETIME) {
		t.Fatal("Setting MUST be soft deleted", setting.GetSoftDeletedAt())
	}

	if !settingFindWithSoftDeleted[0].IsSoftDeleted() {
		t.Fatal("Setting MUST be soft deleted")
	}
}

func TestStore_SettingUpdate(t *testing.T) {
	store, err := initStore(":memory:")

	if err != nil {
		t.Fatal("Store could not be created: ", err.Error())
	}

	setting := NewSetting().SetKey("1").SetValue("one two three")

	err = store.SettingCreate(setting)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	setting.SetValue("one two three")

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.SettingUpdate(setting)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	settingFound, errFind := store.SettingFindByID(setting.GetID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if settingFound == nil {
		t.Fatal("Setting MUST NOT be nil")
	}

	if settingFound.GetValue() != "one two three" {
		t.Fatal("Value MUST be 'one two three', found: ", settingFound.GetValue())
	}
}
