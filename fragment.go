package webcore

import (
	fmt	"fmt"
)

var registeredFragments = map[string]Fragment {}

// Defines a fragment that is typically displayed as the body of a page
type Fragment struct {
	Name	string			// Name of the fragment
	GetHtml	func (Interface) string	// Should return the html-code of the fragment
}

func RegisterFragment (fragment Fragment) {
	registeredFragments[fragment.Name] = fragment
}

func getFragment (name string) Fragment {
	fragment, ok := registeredFragments[name]
	if !ok {
		panic (fmt.Sprintf ("Fragment \"%s\" is not registered.", name))
	}
	return fragment
}
