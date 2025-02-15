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

	Delete(ctx context.Context, settingKey string) error
	Get(ctx context.Context, settingKey string, valueDefault string) (string, error)
	GetAny(ctx context.Context, key string, valueDefault any) (any, error)
	GetMap(ctx context.Context, key string, valueDefault map[string]any) (map[string]any, error)
	Has(ctx context.Context, settingKey string) (bool, error)
	MergeMap(ctx context.Context, key string, mergeMap map[string]any) error
	Set(ctx context.Context, settingKey string, value string) error
	SetAny(ctx context.Context, key string, value interface{}, seconds int64) error
	SetMap(ctx context.Context, key string, value map[string]any) error
	SettingDeleteByKey(ctx context.Context, settingKey string) error
}
