package settingstore

import "github.com/dromara/carbon/v2"

type SettingInterface interface {
	// From data object

	Data() map[string]string
	DataChanged() map[string]string
	MarkAsNotDirty()

	// Methods

	IsSoftDeleted() bool

	// Setters and Getters

	GetCreatedAt() string
	GetCreatedAtCarbon() *carbon.Carbon
	SetCreatedAt(createdAt string) SettingInterface

	GetID() string
	SetID(id string) SettingInterface

	GetKey() string
	SetKey(key string) SettingInterface

	GetSoftDeletedAt() string
	GetSoftDeletedAtCarbon() *carbon.Carbon
	SetSoftDeletedAt(deletedAt string) SettingInterface

	GetUpdatedAt() string
	GetUpdatedAtCarbon() *carbon.Carbon
	SetUpdatedAt(updatedAt string) SettingInterface

	GetValue() string
	SetValue(value string) SettingInterface
}
