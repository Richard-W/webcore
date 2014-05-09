package webcore

import (
	uuid	"code.google.com/p/go-uuid/uuid"
)

type DatabasePlugin struct {
	Name	string
	Version	int
	Upgrade	func (database *Database, oldVersion int) error
}

var coreDatabasePlugin = DatabasePlugin {
	Name:		"core",
	Version:	1,
	Upgrade:	upgradeCoreDatabase,
}

// Upgrades the core database to the current version
func upgradeCoreDatabase (db *Database, oldVersion int) error {
	switch oldVersion {
	case 0:
		rootUUID := uuid.New ()
		queries := []string {
			"CREATE TABLE `plugin_versions` ("+
			"	`name` TEXT NOT NULL PRIMARY KEY,"+
			"	`version` INTEGER"+
			")",

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
