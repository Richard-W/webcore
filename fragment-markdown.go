package webcore

import (
	markdown        "github.com/russross/blackfriday"
	strings		"strings"
)

func init () {
	fragment := Fragment {"markdown", fragmentMarkdownGetHtml}
	RegisterFragment (fragment)
}

func fragmentMarkdownGetHtml (iface Interface) string {
	optionstring := iface.FragmentOptions ()
	options := strings.FieldsFunc (optionstring, func (r rune) bool {
		if r == '|' {
			return true
		}
		return false
	})
	if len (options) != 3 {
		panic ("Invalid optionstring: "+optionstring)
	}
	db := iface.Instance.GetDatabase ()
	res, err := db.Query ("SELECT `"+options[2]+"` FROM `"+options[0]+"` WHERE `"+options[1]+"` = ?", iface.UUID ())
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
