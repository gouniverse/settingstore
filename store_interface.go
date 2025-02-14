package settingstore

import (
	"context"
)

type StoreInterface interface {
	AutoMigrate(ctx context.Context) error
	EnableDebug(debug bool)

	SettingCount(ctx context.Context, query SettingQueryInterface) (int64, error)
	SettingCreate(ctx context.Context, setting SettingInterface) error
	SettingDelete(ctx context.Context, setting SettingInterface) error
	SettingDeleteByID(ctx context.Context, settingID string) error
	SettingFindByID(ctx context.Context, settingID string) (SettingInterface, error)
	SettingFindByKey(ctx context.Context, settingKey string) (SettingInterface, error)
	SettingList(ctx context.Context, query SettingQueryInterface) ([]SettingInterface, error)
	SettingSoftDelete(ctx context.Context, setting SettingInterface) error
	SettingSoftDeleteByID(ctx context.Context, settingID string) error
	SettingUpdate(ctx context.Context, setting SettingInterface) error
}
