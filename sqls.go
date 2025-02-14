package settingstore

import (
	"github.com/gouniverse/sb"
)

// SQLCreateTable returns a SQL string for creating the cache table
func (store *store) SQLCreateTable() string {
	sql := sb.NewBuilder(store.dbDriverName).
		Table(store.settingTableName).
		Column(sb.Column{
			Name:       COLUMN_ID,
			Type:       sb.COLUMN_TYPE_STRING,
			Length:     40,
			PrimaryKey: true,
		}).
		Column(sb.Column{
			Name:   COLUMN_SETTING_KEY,
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 255,
		}).
		Column(sb.Column{
			Name: COLUMN_SETTING_VALUE,
			Type: sb.COLUMN_TYPE_TEXT,
		}).
		Column(sb.Column{
			Name: COLUMN_CREATED_AT,
			Type: sb.COLUMN_TYPE_DATETIME,
		}).
		Column(sb.Column{
			Name: COLUMN_UPDATED_AT,
			Type: sb.COLUMN_TYPE_DATETIME,
		}).
		Column(sb.Column{
			Name: COLUMN_SOFT_DELETED_AT,
			Type: sb.COLUMN_TYPE_DATETIME,
		}).
		CreateIfNotExists()

	return sql
}
