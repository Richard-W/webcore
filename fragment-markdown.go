package webcore

import (
	markdown        "github.com/russross/blackfriday"
)

var markdownDBPlugin = DatabasePlugin {
	Name:		"webcore_markdown",
	Version:	1,
	Upgrade:	upgradeMarkdownDatabase,
}

func init () {
	fragment := Fragment {"markdown", fragmentMarkdownGetHtml}
	RegisterFragment (fragment)
}

func fragmentMarkdownGetHtml (iface FragmentInterface) string {
	db := iface.Instance.GetDatabase ()
	err := db.RegisterPlugin (markdownDBPlugin)
	if err != nil && err != ErrDBPluginNameTaken {
		panic (err.Error ())
	}
	res, err := db.Query ("SELECT `content` FROM `fragment_markdown` WHERE `uuid` = ?", iface.UUID)
	if err != nil {
		panic (err.Error ())
	}
	if !res.Next () {
		return ""
	}
	var content string
	res.Scan (&content)
	return string (markdown.MarkdownCommon ([]byte (content)))
}

func upgradeMarkdownDatabase (db *Database, oldVersion int) error {
	switch oldVersion {
	case 0:
		node, err := getTopNode (db)
		if err != nil {
			return err
		}
		queries := []string {
			"CREATE TABLE `fragment_markdown` ("+
			"	`uuid` TEXT NOT NULL PRIMARY KEY, "+
			"	`content` TEXT "+
			")",

			"INSERT INTO `fragment_markdown` "+
			"	(`uuid`, `content`) VALUES "+
			"	('"+node.uuid+"', '"+
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
