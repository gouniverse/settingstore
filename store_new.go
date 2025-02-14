package settingstore

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/gouniverse/sb"
)

// NewStoreOptions define the options for creating a new setting store
type NewStoreOptions struct {
	SettingTableName   string
	DB                 *sql.DB
	DbDriverName       string
	AutomigrateEnabled bool
	DebugEnabled       bool
	SqlLogger          *slog.Logger
}

// NewStore creates a new setting store
func NewStore(opts NewStoreOptions) (*store, error) {
	store := &store{
		settingTableName:   opts.SettingTableName,
		automigrateEnabled: opts.AutomigrateEnabled,
		db:                 opts.DB,
		dbDriverName:       opts.DbDriverName,
		debugEnabled:       opts.DebugEnabled,
		sqlLogger:          opts.SqlLogger,
	}

	if store.settingTableName == "" {
		return nil, errors.New("setting store: settingTableName is required")
	}

	if store.db == nil {
		return nil, errors.New("setting store: DB is required")
	}

	if store.dbDriverName == "" {
		store.dbDriverName = sb.DatabaseDriverName(store.db)
	}

	if store.sqlLogger == nil {
		store.sqlLogger = slog.Default()
	}

	if store.automigrateEnabled {
		store.AutoMigrate(context.Background())
	}

	return store, nil
}
