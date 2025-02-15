package settingstore

import (
	"context"
)

// StoreInterface defines the interface for a setting store.
type StoreInterface interface {
	// AutoMigrate creates the settings table if it does not exist
	//
	// Parameters:
	// - ctx: the context
	//
	// Returns:
	// - error - nil if no error, error otherwise
	AutoMigrate(ctx context.Context) error

	// EnableDebug - enables the debug option
	//
	// # If enabled will log the SQL statements to the provided logger
	//
	// Parameters:
	//   - debug: true to enable, false otherwise
	//
	// Returns:
	//   - void
	EnableDebug(debug bool)

	// SettingCount counts the settings based on the provided query.
	//
	// Parameters:
	// - ctx: the context
	// - query: the query
	//
	// Returns:
	// - int64: the count
	// - error: nil if no error, error otherwise
	SettingCount(ctx context.Context, query SettingQueryInterface) (int64, error)

	// SettingCreate creates a new setting
	//
	// Parameters:
	// - ctx: the context
	// - setting: the setting
	//
	// Returns:
	// - error: nil if no error, error otherwise
	SettingCreate(ctx context.Context, setting SettingInterface) error

	// SettingDelete deletes a setting
	//
	// Parameters:
	// - ctx: the context
	// - setting: the setting
	//
	// Returns:
	// - error: nil if no error, error otherwise
	SettingDelete(ctx context.Context, setting SettingInterface) error

	// SettingDeleteByID deletes a setting by id
	//
	// Parameters:
	// - ctx: the context
	// - id: the id of the setting to delete
	//
	// Returns:
	// - error: nil if no error, error otherwise
	SettingDeleteByID(ctx context.Context, settingID string) error

	// SettingFindByID finds a setting by id
	//
	// Parameters:
	// - ctx: the context
	// - id: the id of the setting to find
	//
	// Returns:
	// - SettingInterface: the setting
	// - error: nil if no error, error otherwise
	SettingFindByID(ctx context.Context, settingID string) (SettingInterface, error)

	// SettingFindByKey finds a setting by key
	//
	// Parameters:
	// - ctx: the context
	// - key: the key of the setting to find
	//
	// Returns:
	// - SettingInterface: the setting
	// - error: nil if no error, error otherwise
	SettingFindByKey(ctx context.Context, settingKey string) (SettingInterface, error)

	// SettingList retrieves a list of settings
	//
	// Parameters:
	// - ctx: the context
	// - query: the query
	//
	// Returns:
	// - []SettingInterface: the list of settings
	// - error: nil if no error, error otherwise
	SettingList(ctx context.Context, query SettingQueryInterface) ([]SettingInterface, error)

	// SettingSoftDelete soft deletes a setting
	//
	// Parameters:
	// - ctx: the context
	// - setting: the setting
	//
	// Returns:
	// - error: nil if no error, error otherwise
	SettingSoftDelete(ctx context.Context, setting SettingInterface) error

	// SettingSoftDeleteByID soft deletes a setting by id
	//
	// Parameters:
	// - ctx: the context
	// - id: the id of the setting to soft delete
	//
	// Returns:
	// - error: nil if no error, error otherwise
	SettingSoftDeleteByID(ctx context.Context, settingID string) error

	// SettingUpdate updates a setting
	//
	// Parameters:
	// - ctx: the context
	// - setting: the setting
	//
	// Returns:
	// - error: nil if no error, error otherwise
	SettingUpdate(ctx context.Context, setting SettingInterface) error

	// Delete is a shortcut method to delete a value by key
	//
	// Parameters:
	// - ctx: the context
	// - settingKey: the key of the setting to delete
	//
	// Returns:
	// - error - nil if no error, error otherwise
	Delete(ctx context.Context, settingKey string) error

	// Get is a shortcut method to get a value by key, or a default, if not found
	//
	// Parameters:
	// - ctx: the context
	// - settingKey: the key of the setting to get
	// - valueDefault: the default value to return if the setting is not found
	//
	// Returns:
	// - string - the value of the setting, or the default value if not found
	// - error - nil if no error, error otherwise
	Get(ctx context.Context, settingKey string, valueDefault string) (string, error)

	// GetAny is a shortcut method to get a value by key as an interface, or a default if not found
	//
	// Parameters:
	// - ctx: the context
	// - settingKey: the key of the setting to get
	// - valueDefault: the default value to return if the setting is not found
	//
	// Returns:
	// - interface{}, error
	GetAny(ctx context.Context, key string, valueDefault any) (any, error)

	// GetMap is a shortcut method to get a value by key as a map, or a default if not found
	//
	// Parameters:
	// - ctx: the context
	// - settingKey: the key of the setting to get
	// - valueDefault: the default value to return if the setting is not found
	//
	// Returns:
	// - map[string]any - the value of the setting, or the default value if not found
	// - error - nil if no error, error otherwise
	GetMap(ctx context.Context, key string, valueDefault map[string]any) (map[string]any, error)

	// Has is a shortcut method to check if a setting exists by key
	//
	// Parameters:
	// - ctx: the context
	// - settingKey: the key of the setting to check
	//
	// Returns:
	// - bool - true if the setting exists, false otherwise
	// - error - nil if no error, error otherwise
	Has(ctx context.Context, settingKey string) (bool, error)

	// MergeMap is a shortcut method to merge a map with an existing map
	//
	// Parameters:
	// - ctx: the context
	// - key: the key of the setting to merge
	// - mergeMap: the map to merge with the existing map
	//
	// Returns:
	// - error - nil if no error, error otherwise
	MergeMap(ctx context.Context, key string, mergeMap map[string]any) error

	// Set is a shortcut method to save a value by key, use Get to extract
	//
	// Parameters:
	// - ctx: the context
	// - settingKey: the key of the setting to save
	// - value: the value to save
	//
	// Returns:
	// - error - nil if no error, error otherwise
	Set(ctx context.Context, settingKey string, value string) error

	// SetAny is a shortcut method to save any value by key, use GetAny to extract
	//
	// Parameters:
	// - ctx: the context
	// - settingKey: the key of the setting to save
	// - value: the value to save
	//
	// Returns:
	// - error - nil if no error, error otherwise
	SetAny(ctx context.Context, key string, value interface{}, seconds int64) error

	// SetMap is a shortcut method to save a map by key, use GetMap to extract
	//
	// Parameters:
	// - ctx: the context
	// - settingKey: the key of the setting to save
	// - value: the value to save
	//
	// Returns:
	// - error - nil if no error, error otherwise
	SetMap(ctx context.Context, key string, value map[string]any) error

	// SettingDeleteByKey deletes a setting by id
	//
	// Parameters:
	// - ctx: the context
	// - settingKey: the key of the setting to delete
	//
	// Returns:
	// - error - nil if no error, error otherwise
	SettingDeleteByKey(ctx context.Context, settingKey string) error
}
