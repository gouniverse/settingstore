package settingstore

import (
	"time"
)

// Setting type
type Setting struct {
	ID        string     `db:"id";`
	Key       string     `db:"setting_key";`
	Value     string     `db:"setting_value";`
	CreatedAt time.Time  `db:"created_at";`
	UpdatedAt time.Time  `db:"updated_at";`
	DeletedAt *time.Time `db:"deleted_at";`
}
