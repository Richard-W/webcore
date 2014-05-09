package webcore

import (
	sql	"database/sql"
	errors	"errors"
)

// Maintains the database and upgrades it if neccessary
type Database struct {
	realDatabase	*sql.DB			// Wrapped SQL connection
	plugins		[]DatabasePlugin	// Slice of plugins registered with this database
}

var ErrDBPluginNameTaken = errors.New ("Plugin name is already registered.")

// Wraps a sql.DB pointer in a database struct
func NewDatabase (db *sql.DB) (*Database, error) {
	database := new (Database)
	database.realDatabase = db
	err := database.RegisterPlugin (coreDatabasePlugin)
	if err != nil {
		return nil, err
	}
	return database, nil
}

// Registers a database plugin with this database and upgrades the tables if neccessary
func (db *Database) RegisterPlugin (plugin DatabasePlugin) error {
	for _, rplugin := range db.plugins {
		if plugin.Name == rplugin.Name {
			return ErrDBPluginNameTaken
		}
	}
	db.plugins = append (db.plugins, plugin)
	res, err := db.Query ("SELECT `version` FROM `plugin_versions` WHERE `name` = ?", plugin.Name)
	if err != nil {
		return db.upgradePlugin (plugin, 0)
	}
	res.Next ()
	var plugin_version int
	res.Scan (&plugin_version)
	if plugin_version < plugin.Version {
		return db.upgradePlugin (plugin, plugin_version)
	}
	return nil;
}

// Upgrades a plugin and sets the version in the version-table
func (db *Database) upgradePlugin (plugin DatabasePlugin, oldVersion int) error {
	err := plugin.Upgrade (db, oldVersion)
	if err != nil {
		return err
	}
	var query string
	if oldVersion == 0 {
		query = "INSERT INTO `plugin_versions` (`version`, `name`) VALUES (?, ?)"
	} else {
		query = "UPDATE `plugin_versions` SET `version` = ? WHERE `name` = ?"
	}
	_, err = db.Exec (query, plugin.Version, plugin.Name)
	return err
}

// See sql.DB.Query
func (db Database) Query (query string, v ...interface{}) (*sql.Rows, error) {
	return db.realDatabase.Query (query, v...)
}

// See sql.DB.Exec
func (db Database) Exec (query string, v ...interface{}) (sql.Result, error) {
	return db.realDatabase.Exec (query, v...)
}
