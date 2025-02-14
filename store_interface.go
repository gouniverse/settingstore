package settingstore

type StoreInterface interface {
	AutoMigrate() error
	EnableDebug(debug bool)
	SessionExpiryGoroutine() error

	SettingCount(query SettingQueryInterface) (int64, error)
	SettingCreate(setting SettingInterface) error
	SettingDelete(setting SettingInterface) error
	SettingDeleteByID(settingID string) error
	SettingExtend(setting SettingInterface, seconds int64) error
	SettingFindByID(settingID string) (SettingInterface, error)
	SettingFindByKey(settingKey string) (SettingInterface, error)
	SettingList(query SettingQueryInterface) ([]SettingInterface, error)
	SettingSoftDelete(setting SettingInterface) error
	SettingSoftDeleteByID(settingID string) error
	SettingUpdate(setting SettingInterface) error
}
