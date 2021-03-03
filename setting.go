package settingstore

import (
	"time"

	"github.com/gouniverse/uid"
	"gorm.io/gorm"
)

// import (
// 	"errors"
// 	"time"

// 	"github.com/gouniverse/uid"
// 	"gorm.io/gorm"
// )

// Setting type
type Setting struct {
	ID        string     `gorm:"type:varchar(40);column:id;primary_key;"`
	Key       string     `gorm:"type:varchar(40);column:setting_key;"`
	Value     string     `gorm:"type:longtext;column:setting_value;"`
	CreatedAt time.Time  `gorm:"type:datetime;column:created_at;DEFAULT NULL;"`
	UpdatedAt time.Time  `gorm:"type:datetime;column:updated_at;DEFAULT NULL;"`
	DeletedAt *time.Time `gorm:"type:datetime;olumn:deleted_at;DEFAULT NULL;"`
}

// BeforeCreate adds UID to model
func (c *Setting) BeforeCreate(tx *gorm.DB) (err error) {
	uuid := uid.NanoUid()
	c.ID = uuid
	return nil
}

// // SettingGet gets a key from cache
// func SettingGet(key string, valueDefault string) string {
// 	setting := SettingFindByKey(key)

// 	if setting != nil {
// 		return setting.Value
// 	}

// 	return valueDefault
// }

// // SettingSet sets a key in cache
// func SettingSet(key string, value string) bool {
// 	setting := SettingFindByKey(key)

// 	if setting != nil {
// 		setting.Value = value
// 		dbResult := GetDb().Save(&setting)
// 		if dbResult != nil {
// 			return false
// 		}
// 		return true
// 	}

// 	var newSetting = Setting{Key: key, Value: value}

// 	dbResult := GetDb().Create(&newSetting)

// 	if dbResult.Error != nil {
// 		return false
// 	}

// 	return true
// }
