package webcore

import (
	http	"net/http"
)

// Struct that is given to a fragment to handle requests.
type FragmentInterface struct {
	UUID		string		// UUID of the node that is requested
	Request		*http.Request	// The request object
	Instance	*Instance	// Instance on which this request happened
	Session		*Session	// Storage that is maintained for a unique user
	Options		[]string	// Node-specific options for the fragment.
}
