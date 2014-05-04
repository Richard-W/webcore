package webcore

import (
	markdown        "github.com/russross/blackfriday"
)

func init () {
	fragment := Fragment {"markdown", fragmentMarkdownGetHtml}
	RegisterFragment (fragment)
}

func fragmentMarkdownGetHtml (iface FragmentInterface) string {
	if len (iface.Options) != 3 {
		panic ("Invalid fragment options")
	}
	db := iface.Instance.GetDatabase ()
	res, err := db.Query ("SELECT `"+iface.Options[2]+"` FROM `"+iface.Options[0]+"` WHERE `"+iface.Options[1]+"` = ?", iface.UUID)
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
