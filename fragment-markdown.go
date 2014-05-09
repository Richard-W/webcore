package webcore

import (
	markdown        "github.com/russross/blackfriday"
)

func init () {
	fragment := Fragment {"markdown", fragmentMarkdownGetHtml}
	RegisterFragment (fragment)
}

func fragmentMarkdownGetHtml (iface FragmentInterface) string {
	db := iface.Instance.GetDatabase ()
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
