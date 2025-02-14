package settingstore

import "database/sql"

// NewStore creates a new store
func NewStore(db *sql.DB, dbDriverName string, settingsTableName string, automigrateEnabled bool, debugEnabled bool) StoreInterface {
	s := &store{
		db:                 db,
		dbDriverName:       dbDriverName,
		settingsTableName:  settingsTableName,
		timeoutSeconds:     3600,
		automigrateEnabled: automigrateEnabled,
		debugEnabled:       debugEnabled,
	}

	return s
}
