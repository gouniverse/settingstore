package settingstore

import (
	"github.com/dromara/carbon/v2"
	"github.com/gouniverse/dataobject"
	"github.com/gouniverse/sb"
	"github.com/gouniverse/uid"
)

var _ SettingInterface = (*Setting)(nil)

// Setting type
type Setting struct {
	dataobject.DataObject
}

// == CONSTRUCTORS ============================================================

func NewSetting() SettingInterface {
	createdAt := carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)
	updatedAt := carbon.Now(carbon.UTC).ToDateTimeString(carbon.UTC)
	deletedAt := sb.MAX_DATETIME

	o := (&Setting{})

	o.SetID(uid.HumanUid()).
		SetKey("").
		SetValue("").
		SetCreatedAt(createdAt).
		SetUpdatedAt(updatedAt).
		SetSoftDeletedAt(deletedAt)

	return o
}

func NewSettingFromExistingData(data map[string]string) SettingInterface {
	o := &Setting{}
	o.Hydrate(data)
	return o
}

// == METHODS =================================================================

func (o *Setting) IsSoftDeleted() bool {
	return o.GetSoftDeletedAtCarbon().Compare("<", carbon.Now(carbon.UTC))
}

// == SETTERS AND GETTERS =====================================================

func (setting *Setting) GetID() string {
	return setting.Get(COLUMN_ID)
}

func (setting *Setting) SetID(id string) SettingInterface {
	setting.Set(COLUMN_ID, id)
	return setting
}

func (setting *Setting) GetKey() string {
	return setting.Get(COLUMN_SETTING_KEY)
}

func (setting *Setting) SetKey(key string) SettingInterface {
	setting.Set(COLUMN_SETTING_KEY, key)
	return setting
}

func (setting *Setting) GetValue() string {
	return setting.Get(COLUMN_SETTING_VALUE)
}

func (setting *Setting) SetValue(value string) SettingInterface {
	setting.Set(COLUMN_SETTING_VALUE, value)
	return setting
}

func (setting *Setting) GetCreatedAt() string {
	return setting.Get(COLUMN_CREATED_AT)
}

func (setting *Setting) GetCreatedAtCarbon() carbon.Carbon {
	return carbon.Parse(setting.GetCreatedAt(), carbon.UTC)
}

func (setting *Setting) SetCreatedAt(createdAt string) SettingInterface {
	setting.Set(COLUMN_CREATED_AT, createdAt)
	return setting
}

func (setting *Setting) GetUpdatedAt() string {
	return setting.Get(COLUMN_UPDATED_AT)
}

func (setting *Setting) GetUpdatedAtCarbon() carbon.Carbon {
	return carbon.Parse(setting.GetUpdatedAt(), carbon.UTC)
}

func (setting *Setting) SetUpdatedAt(updatedAt string) SettingInterface {
	setting.Set(COLUMN_UPDATED_AT, updatedAt)
	return setting
}

func (setting *Setting) GetSoftDeletedAt() string {
	return setting.Get(COLUMN_SOFT_DELETED_AT)
}

func (setting *Setting) GetSoftDeletedAtCarbon() carbon.Carbon {
	return carbon.Parse(setting.GetSoftDeletedAt(), carbon.UTC)
}

func (setting *Setting) SetSoftDeletedAt(deletedAt string) SettingInterface {
	setting.Set(COLUMN_SOFT_DELETED_AT, deletedAt)
	return setting
}
