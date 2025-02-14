package settingstore

import "errors"

type SettingQueryInterface interface {
	Validate() error

	IsCountOnly() bool

	Columns() []string
	SetColumns(columns []string) SettingQueryInterface

	HasCreatedAtGte() bool
	CreatedAtGte() string
	SetCreatedAtGte(createdAtGte string) SettingQueryInterface

	HasCreatedAtLte() bool
	CreatedAtLte() string
	SetCreatedAtLte(createdAtLte string) SettingQueryInterface

	HasExpiresAtGte() bool
	ExpiresAtGte() string
	SetExpiresAtGte(expiresAtGte string) SettingQueryInterface

	HasExpiresAtLte() bool
	ExpiresAtLte() string
	SetExpiresAtLte(expiresAtLte string) SettingQueryInterface

	HasID() bool
	ID() string
	SetID(id string) SettingQueryInterface

	HasIDIn() bool
	IDIn() []string
	SetIDIn(idIn []string) SettingQueryInterface

	HasKey() bool
	Key() string
	SetKey(key string) SettingQueryInterface

	HasUserID() bool
	UserID() string
	SetUserID(userID string) SettingQueryInterface

	HasUserIpAddress() bool
	UserIpAddress() string
	SetUserIpAddress(userIpAddress string) SettingQueryInterface

	HasUserAgent() bool
	UserAgent() string
	SetUserAgent(userAgent string) SettingQueryInterface

	HasOffset() bool
	Offset() int
	SetOffset(offset int) SettingQueryInterface

	HasLimit() bool
	Limit() int
	SetLimit(limit int) SettingQueryInterface

	HasSortOrder() bool
	SortOrder() string
	SetSortOrder(sortOrder string) SettingQueryInterface

	HasOrderBy() bool
	OrderBy() string
	SetOrderBy(orderBy string) SettingQueryInterface

	HasCountOnly() bool
	SetCountOnly(countOnly bool) SettingQueryInterface

	HasSoftDeletedIncluded() bool
	SoftDeletedIncluded() bool
	SetSoftDeletedIncluded(withSoftDeleted bool) SettingQueryInterface
}

// SettingQuery is a shortcut version of NewSettingQuery to create a new query
func SettingQuery() SettingQueryInterface {
	return NewSettingQuery()
}

// NewSettingQuery creates a new setting query
func NewSettingQuery() SettingQueryInterface {
	return &settingQuery{
		properties: make(map[string]interface{}),
	}
}

var _ SettingQueryInterface = (*settingQuery)(nil)

type settingQuery struct {
	properties map[string]interface{}
}

func (q *settingQuery) Validate() error {
	if q.HasCreatedAtGte() && q.CreatedAtGte() == "" {
		return errors.New("Setting query. created_at_gte cannot be empty")
	}

	if q.HasCreatedAtLte() && q.CreatedAtLte() == "" {
		return errors.New("Setting query. created_at_lte cannot be empty")
	}

	if q.HasID() && q.ID() == "" {
		return errors.New("Setting query. id cannot be empty")
	}

	if q.HasIDIn() && len(q.IDIn()) < 1 {
		return errors.New("Setting query. id_in cannot be empty array")
	}

	if q.HasLimit() && q.Limit() < 0 {
		return errors.New("Setting query. limit cannot be negative")
	}

	if q.HasOffset() && q.Offset() < 0 {
		return errors.New("Setting query. offset cannot be negative")
	}

	return nil
}

func (q *settingQuery) Columns() []string {
	if !q.hasProperty("columns") {
		return []string{}
	}

	return q.properties["columns"].([]string)
}

func (q *settingQuery) SetColumns(columns []string) SettingQueryInterface {
	q.properties["columns"] = columns
	return q
}

func (q *settingQuery) HasCountOnly() bool {
	return q.hasProperty("count_only")
}

func (q *settingQuery) IsCountOnly() bool {
	return q.hasProperty("count_only") && q.properties["count_only"].(bool)
}

func (q *settingQuery) SetCountOnly(countOnly bool) SettingQueryInterface {
	q.properties["count_only"] = countOnly
	return q
}

func (q *settingQuery) HasCreatedAtGte() bool {
	return q.hasProperty("created_at_gte")
}

func (q *settingQuery) CreatedAtGte() string {
	return q.properties["created_at_gte"].(string)
}

func (q *settingQuery) SetCreatedAtGte(createdAtGte string) SettingQueryInterface {
	q.properties["created_at_gte"] = createdAtGte
	return q
}

func (q *settingQuery) HasCreatedAtLte() bool {
	return q.hasProperty("created_at_lte")
}

func (q *settingQuery) CreatedAtLte() string {
	return q.properties["created_at_lte"].(string)
}

func (q *settingQuery) SetCreatedAtLte(createdAtLte string) SettingQueryInterface {
	q.properties["created_at_lte"] = createdAtLte
	return q
}

func (q *settingQuery) HasExpiresAtGte() bool {
	return q.hasProperty("expires_at_gte")
}

func (q *settingQuery) ExpiresAtGte() string {
	return q.properties["expires_at_gte"].(string)
}

func (q *settingQuery) SetExpiresAtGte(expiresAtGte string) SettingQueryInterface {
	q.properties["expires_at_gte"] = expiresAtGte
	return q
}

func (q *settingQuery) HasExpiresAtLte() bool {
	return q.hasProperty("expires_at_lte")
}

func (q *settingQuery) ExpiresAtLte() string {
	return q.properties["expires_at_lte"].(string)
}

func (q *settingQuery) SetExpiresAtLte(expiresAtLte string) SettingQueryInterface {
	q.properties["expires_at_lte"] = expiresAtLte
	return q
}

func (q *settingQuery) HasID() bool {
	return q.hasProperty("id")
}

func (q *settingQuery) ID() string {
	return q.properties["id"].(string)
}

func (q *settingQuery) SetID(id string) SettingQueryInterface {
	q.properties["id"] = id
	return q
}

func (q *settingQuery) HasIDIn() bool {
	return q.hasProperty("id_in")
}

func (q *settingQuery) IDIn() []string {
	return q.properties["id_in"].([]string)
}

func (q *settingQuery) SetIDIn(idIn []string) SettingQueryInterface {
	q.properties["id_in"] = idIn
	return q
}

func (q *settingQuery) HasKey() bool {
	return q.hasProperty("key")
}

func (q *settingQuery) Key() string {
	return q.properties["key"].(string)
}

func (q *settingQuery) SetKey(key string) SettingQueryInterface {
	q.properties["key"] = key
	return q
}

func (q *settingQuery) HasLimit() bool {
	return q.hasProperty("limit")
}

func (q *settingQuery) Limit() int {
	return q.properties["limit"].(int)
}

func (q *settingQuery) SetLimit(limit int) SettingQueryInterface {
	q.properties["limit"] = limit
	return q
}

func (q *settingQuery) HasOffset() bool {
	return q.hasProperty("offset")
}

func (q *settingQuery) Offset() int {
	return q.properties["offset"].(int)
}

func (q *settingQuery) SetOffset(offset int) SettingQueryInterface {
	q.properties["offset"] = offset
	return q
}

func (q *settingQuery) HasOrderBy() bool {
	return q.hasProperty("order_by")
}

func (q *settingQuery) OrderBy() string {
	return q.properties["order_by"].(string)
}

func (q *settingQuery) SetOrderBy(orderBy string) SettingQueryInterface {
	q.properties["order_by"] = orderBy
	return q
}

func (q *settingQuery) HasSoftDeletedIncluded() bool {
	return q.hasProperty("soft_deleted_included")
}

func (q *settingQuery) SoftDeletedIncluded() bool {
	if !q.HasSoftDeletedIncluded() {
		return false
	}

	return q.properties["soft_deleted_included"].(bool)
}

func (q *settingQuery) SetSoftDeletedIncluded(softDeletedIncluded bool) SettingQueryInterface {
	q.properties["soft_deleted_included"] = softDeletedIncluded
	return q
}

func (q *settingQuery) HasSortOrder() bool {
	return q.hasProperty("sort_order")
}

func (q *settingQuery) SortOrder() string {
	return q.properties["sort_order"].(string)
}

func (q *settingQuery) SetSortOrder(sortOrder string) SettingQueryInterface {
	q.properties["sort_order"] = sortOrder
	return q
}

func (q *settingQuery) HasUserAgent() bool {
	return q.hasProperty("user_agent")
}

func (q *settingQuery) UserAgent() string {
	return q.properties["user_agent"].(string)
}

func (q *settingQuery) SetUserAgent(userAgent string) SettingQueryInterface {
	q.properties["user_agent"] = userAgent
	return q
}

func (q *settingQuery) HasUserID() bool {
	return q.hasProperty("user_id")
}

func (q *settingQuery) UserID() string {
	return q.properties["user_id"].(string)
}

func (q *settingQuery) SetUserID(userID string) SettingQueryInterface {
	q.properties["user_id"] = userID
	return q
}

func (q *settingQuery) HasUserIpAddress() bool {
	return q.hasProperty("user_ip_address")
}

func (q *settingQuery) UserIpAddress() string {
	return q.properties["user_ip_address"].(string)
}

func (q *settingQuery) SetUserIpAddress(userIpAddress string) SettingQueryInterface {
	q.properties["user_ip_address"] = userIpAddress
	return q
}

func (q *settingQuery) hasProperty(key string) bool {
	_, ok := q.properties[key]
	return ok
}
