# Settings Store

[![Tests Status](https://github.com/gouniverse/settingstore/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/gouniverse/settingstore/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/gouniverse/settingstore)](https://goreportcard.com/report/github.com/gouniverse/settingstore)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/gouniverse/settingstore)](https://pkg.go.dev/github.com/gouniverse/settingstore)

Saves settings in an SQL database. 

## 🌏  Open in the Cloud 
Click any of the buttons below to start a new development environment to demo or contribute to the codebase without having to install anything on your machine:

[![Open in VS Code](https://img.shields.io/badge/Open%20in-VS%20Code-blue?logo=visualstudiocode)](https://vscode.dev/github/gouniverse/settingstore)
[![Open in Glitch](https://img.shields.io/badge/Open%20in-Glitch-blue?logo=glitch)](https://glitch.com/edit/#!/import/github/gouniverse/settingstore)
[![Open in GitHub Codespaces](https://github.com/codespaces/badge.svg)](https://codespaces.new/gouniverse/settingstore)
[![Edit in Codesandbox](https://codesandbox.io/static/img/play-codesandbox.svg)](https://codesandbox.io/s/github/gouniverse/settingstore)
[![Open in StackBlitz](https://developer.stackblitz.com/img/open_in_stackblitz.svg)](https://stackblitz.com/github/gouniverse/settingstore)
[![Open in Repl.it](https://replit.com/badge/github/withastro/astro)](https://replit.com/github/gouniverse/settingstore)
[![Open in Codeanywhere](https://codeanywhere.com/img/open-in-codeanywhere-btn.svg)](https://app.codeanywhere.com/#https://github.com/gouniverse/settingstore)
[![Open in Gitpod](https://gitpod.io/button/open-in-gitpod.svg)](https://gitpod.io/#https://github.com/gouniverse/settingstore)

## Description

Every application needs to preserve settings between multiple restarts. This package helps save the setting represented as key-value pairs in an SQL database.

## License

This project is dual-licensed under the following terms:

- For non-commercial use, you may choose either the GNU Affero General Public License v3.0 (AGPLv3) *or* a separate commercial license (see below). You can find a copy of the AGPLv3 at: https://www.gnu.org/licenses/agpl-3.0.txt

- For commercial use, a separate commercial license is required. Commercial licenses are available for various use cases. Please contact me via my [contact page](https://lesichkov.co.uk/contact) to obtain a commercial license.

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
settingStore, err = settingstore.NewStore(settingstore.NewStoreOptions{
	DB: databaseInstance,
	SettingTableName: "settings",
	AutomigrateEnabled: true,
})

if err != nil {
	panic(err)
}

// as multiple lines
settingStore, err = settingstore.NewStore(settingstore.NewStoreOptions{
	DB: databaseInstance,
	SettingTableName: "settings",
})

if err != nil {
	panic(err)
}

settingStore.AutoMigrate()

```

## Usage

1. Create a new key value setting pair
```
settingsStore.Set("app.name", "My Web App")
settingsStore.Set("app.url", "http://localhost")
settingsStore.Set("server.ip", "127.0.0.1")
settingsStore.Set("server.port", "80")
```

2. Retrieve a setting value (or default value if not exists)
```
appName = settingsStore.Get("app.name", "Default Name")
appUrl = settingsStore.Get("app.url", "")
serverIp = settingsStore.Get("server.ip", "")
serverPort = settingsStore.Get("server.port", "")
```

3. Check if required setting is setup
```
if serverIp == "" {
    log.Panic("server ip not setup")
}
```

## Methods

These methods may be subject to change as still in development

### Store Methods
- NewStore(opts NewStoreOptions) (*store, error) - creates a new setting store
- AutoMigrate(ctx context.Context) error - auto migrate (create the tables in the database) the settings store tables
- DriverName(db *sql.DB) string - the name of the driver used for SQL strings (you may use this if you need to debug)
- SettingCount(ctx context.Context, query SettingQueryInterface) (int64, error) - counts the number of settings
- SettingCreate(ctx context.Context, setting SettingInterface) error - creates a new setting
- SettingDelete(ctx context.Context, setting SettingInterface) error - deletes a setting
- SettingDeleteByID(ctx context.Context, settingID string) error - deletes a setting by ID
- SettingDeleteByKey(ctx context.Context, settingKey string) error - deletes a setting by key
- SettingFindByID(ctx context.Context, settingID string) (SettingInterface, error) - finds a setting by ID
- SettingList(ctx context.Context, query SettingQueryInterface) ([]SettingInterface, error) - lists settings
- SettingSoftDelete(ctx context.Context, setting SettingInterface) error - soft deletes a setting
- SettingSoftDeleteByID(ctx context.Context, settingID string) error - soft deletes a setting by ID
- SettingUpdate(ctx context.Context, setting SettingInterface) error - updates a setting


### Shortcut Methods

- Get(ctx context.Context, key string, valueDefault string) (string, error) - gets a value from key-value setting pair
- Set(ctx context.Context, key string, value string, seconds int64) error - sets new key value pair

- GetAny(ctx context.Context, key string, valueDefault interface{}) (interface{}, error) - gets a value from key-value setting pair

- GetJSON(key string, valueDefault interface{}) (interface{}, error) - gets a value as JSON from key-value setting pair
- SetJSON(ctx context.Context, key string, value interface{}, seconds int64) error - sets new key JSON value pair

- GetMap(ctx context.Context, key string, valueDefault map[string]any) (map[string]any, error) - gets a value as JSON from key-value setting pair
- MergeMap(ctx context.Context, key string, mergeMap map[string]any, seconds int64) error - merges a map with an existing map

- Has(ctx context.Context, settingKey string) (bool, error) - checks if a setting exists