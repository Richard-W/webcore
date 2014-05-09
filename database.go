package webcore

import (
	sql	"database/sql"
	uuid	"code.google.com/p/go-uuid/uuid"
	strconv	"strconv"
)

const DATABASE_VERSION = 1

// Maintains the database and upgrades it if neccessary
type Database struct {
	realDatabase	*sql.DB	// Wrapped SQL connection
}

// Wraps a sql.DB pointer in a database struct
func newDatabase (db *sql.DB) (*Database, error) {
	database := new (Database)
	database.realDatabase = db
	err := database.checkVersion ()
	if err != nil {
		return nil, err
	}
	return database, nil
}

// Checks if the database is up-to-date and upgrades it if it is not.
func (db Database) checkVersion () error {
	res, err := db.Query ("SELECT `version` FROM `versions` WHERE `name` = 'webcore'")
	if err != nil {
		return db.upgradeDatabase (0)
	}
	if !res.Next () {
		panic ("Unexpected behaviour of sql-module: ErrNoRows not thrown but res.Next returns false")
	}
	var database_version int
	res.Scan (&database_version)
	if database_version < DATABASE_VERSION {
		return db.upgradeDatabase (database_version)
	}
	return nil
}

// See sql.DB.Query
func (db Database) Query (query string, v ...interface{}) (*sql.Rows, error) {
	return db.realDatabase.Query (query, v...)
}

// See sql.DB.Exec
func (db Database) Exec (query string, v ...interface{}) (sql.Result, error) {
	return db.realDatabase.Exec (query, v...)
}

// Upgrades the database to the current version
func (db Database) upgradeDatabase (oldVersion int) error {
	switch oldVersion {
	case 0:
		rootUUID := uuid.New ()
		queries := []string {
			"CREATE TABLE `versions` ("+
			"	`name` TEXT NOT NULL PRIMARY KEY,"+
			"	`version` INTEGER"+
			")",

			"INSERT INTO `versions` (`name`, `version`) VALUES "+
			"	('webcore', '"+strconv.Itoa (DATABASE_VERSION)+"')",

			"CREATE TABLE `nodes` ("+
			"	`uuid` TEXT NOT NULL PRIMARY KEY,"+
			"	`parent_id` TEXT,"+
			"	`name` TEXT,"+
			"	`display_name` TEXT,"+
			"	`fragment` TEXT,"+
			"	`fragment_options` TEXT,"+
			"	`deep` INTEGER"+
			")",

			"INSERT INTO `nodes` "+
			"	(`uuid`, `parent_id`, `name`, `display_name`, `fragment`, `fragment_options`) VALUES "+
			"	('"+rootUUID+"', '', 'root', 'Home', 'markdown', '')",

			"CREATE TABLE `fragment_markdown` ("+
			"	`uuid` TEXT NOT NULL PRIMARY KEY, "+
			"	`content` TEXT "+
			")",

			"INSERT INTO `fragment_markdown` "+
			"	(`uuid`, `content`) VALUES "+
			"	('"+rootUUID+"', '"+
					"Home\n"+
					"====\n\n"+
					"This is a sample page"+
			"	')",
		}
		for _, query := range queries {
			_, err := db.Exec (query)
			if err != nil {
				return err
			}
		}
		break
	}
	return nil
}
