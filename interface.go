package webcore

import (
	http	"net/http"
)

// Interface used to request content from the application
// and for use in the base template
type Interface struct {
	Request		*http.Request
	MainMenu	[]MenuItem
	Instance	*Instance
	Session		*Session
	node		*node
}

// Used for constructing the menus in the base template
type MenuItem struct {
	Name	string
	Path	string
	SubMenu	[]MenuItem
}

func (iface Interface) UUID () string {
	return iface.node.uuid
}

func (iface Interface) Body () string {
	fragment := getFragment (iface.node.fragment)
	return fragment.GetHtml (iface)
}

func (iface Interface) FragmentOptions () string {
	return iface.node.fragmentOptions
}
