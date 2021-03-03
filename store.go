package settingstore

import (
	"encoding/json"
	"errors"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Store defines a session store
type Store struct {
	settingsTableName string
	db                *gorm.DB
}

// StoreOption options for the vault store
type StoreOption func(*Store)

// WithDriverAndDNS sets the driver and the DNS for the database for the cache store
func WithDriverAndDNS(driverName string, dsn string) StoreOption {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	return func(s *Store) {
		s.db = db
	}
}

// WithGormDb sets the GORM database for the cache store
func WithGormDb(db *gorm.DB) StoreOption {
	return func(s *Store) {
		s.db = db
	}
}

// WithTableName sets the table name for the cache store
func WithTableName(settingsTableName string) StoreOption {
	return func(s *Store) {
		s.settingsTableName = settingsTableName
	}
}

// NewStore creates a new entity store
func NewStore(opts ...StoreOption) *Store {
	store := &Store{}
	for _, opt := range opts {
		opt(store)
	}

	if store.settingsTableName == "" {
		log.Panic("Vault store: vaultTableName is required")
	}

	store.db.Table(store.settingsTableName).AutoMigrate(&Setting{})

	return store
}

// FindByKey finds a session by key
func (st *Store) FindByKey(key string) *Setting {
	setting := &Setting{}
	result := st.db.Table(st.settingsTableName).Where("`setting_key` = ?", key).First(&setting)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil
	}

	return setting
}

// Get gets a key from settings
func (st *Store) Get(key string, valueDefault string) string {
	cache := st.FindByKey(key)

	if cache != nil {
		return cache.Value
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

// Remove gets a JSON key from cache
func (st *Store) Remove(key string) {
	st.db.Table(st.settingsTableName).Where("`setting_key` = ?", key).Delete(Setting{})
}

// Set sets new key value pair
func (st *Store) Set(key string, value string) bool {
	setting := st.FindByKey(key)

	if setting != nil {
		setting.Value = value
		dbResult := st.db.Table(st.settingsTableName).Save(&setting)
		if dbResult != nil {
			return false
		}
		return true
	}

	var newSetting = Setting{Key: key, Value: value}

	dbResult := st.db.Table(st.settingsTableName).Create(&newSetting)

	if dbResult.Error != nil {
		return false
	}

	return true
}

// SetJSON sets new key JSON value pair
func (st *Store) SetJSON(key string, value interface{}) bool {
	jsonValue, jsonError := json.Marshal(value)
	if jsonError != nil {
		return false
	}

	return st.Set(key, string(jsonValue))
}
