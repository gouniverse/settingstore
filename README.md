# Setting Store

Saves settings in an SQL database. 

## Description

Every application needs to preserve settings between multiple restarts. This package helps save the setting represented as key-value pairs in an SQL database.

## Features

- Saves settings data as key-value pairs
- Supports SQLite, MySQL and Postgres
- Uses sql.DB directly
- Automigration

## Installation
```
go get -u github.com/gouniverse/settingstore
```

## Setup

```
// as one line
settingStore = settingstore.NewStore(settingstore.WithDb(databaseInstance), settingstore.WithTableName("settings"), entitystore.WithAutoMigrate(true))


// as multiple lines
settingStore = settingstore.NewStore(settingstore.WithDb(databaseInstance), settingstore.WithTableName("settings"))
settingStore.AutoMigrate()

```

## Usage

1. Create a new key value pair
```
settingsStore.Set("app.name", "My Web App")
settingsStore.Set("app.url", "http://localhost")
settingsStore.Set("server.ip", "127.0.0.1")
settingsStore.Set("server.port", "80")
```

2. Retrieve an entity (or default value if not exists)
```
appName = settingsStore.Set("app.name", "Default Name")
appUrl = settingsStore.Set("app.url", "")
serverIp = settingsStore.Set("server.ip", "")
serverPort = settingsStore.Set("server.port", "")

3. Check if required setting is setup
```
if serverIp == "" {
    log.Panic("server ip not setup")
}
```

## Methods

These methods may be subject to change

### Store Methods

- NewStore(opts ...StoreOption) *Store - creates a new setting store
- WithAutoMigrate(automigrateEnabled bool) StoreOption - a store option. sets whether database migration will run on setup
- WithDb(db *sql.DB) StoreOption - a store option (required). sets the DB to be used by the store
- WithTableName(settingsTableName string) StoreOption - a store option (required). sets the table name for the setting store in the DB

- AutoMigrate() error - auto migrate (create the tables in the database) the settings store tables
- DriverName(db *sql.DB) string - the name of the driver used for SQL strings (you may use this if you need to debug)
- SqlCreateTable() string - SQL string for creating the tables (you may use this string if you want to set your own migrations)

### Setting Methods

- Delete() bool - deletes the entity
- FindByKey(key string) *Setting - finds a Setting by key
- Get(key string, valueDefault string) string - gets a value from key-value setting pair
- GetJSON(key string, valueDefault interface{}) interface{} - gets a value as JSON from key-value setting pair
- Keys() ([]string, error) - gets all keys sorted alphabetically (useful if you want to list these in admin panel)
- Remove(key string) error - removes a setting from store
- Set(key string, value string) (bool, error) - sets new key value pair
- SetJSON(key string, value interface{}) (bool, error) - sets new key JSON value pair
