// Webcore implements a lightweight and extensible CMS.
package webcore

import (
	http		"net/http"
	sql		"database/sql"
	template	"text/template"
	regexp		"regexp"
	strings		"strings"
)

// The default instance
var DefaultInstance = NewInstance ()

// Regexp validate requests for static files
var staticRegexp = regexp.MustCompile ("^/static/([a-zA-Z0-9\\-_/]+(\\.[a-zA-Z0-9\\-_/]+)*)$")

// An Instance is the representation of an entire webservice.
type Instance struct {
	database	*Database		// SQL database used by this instance
	httpServer	http.Server		// http server serving requests
	httpHandler	*http.ServeMux		// associated http handler
	topNode		*node			// the root of the node tree
	baseTemplate	*template.Template	// the (trusted) template for the site
	staticPath	string			// the path of the static directory
}

// Creates a new instance
func NewInstance () *Instance {
	in := new (Instance)
	in.httpHandler = http.NewServeMux ()
	in.httpServer.Handler = in.httpHandler
	in.httpHandler.HandleFunc ("/static/", in.staticHandler)
	in.httpHandler.HandleFunc ("/", in.dynamicHandler)
	return in
}

// Set the database of the instance
func (in *Instance) SetDatabase (database *sql.DB) {
	db, err := newDatabase (database)
	if err != nil {
		panic (err.Error ())
	}
	in.database = db
}

// Set the database of the default instance
func SetDatabase (database *sql.DB) {
	DefaultInstance.SetDatabase (database)
}

// Get the database of the instance
func (in *Instance) GetDatabase () *Database {
	return in.database
}

// Get the database of the instance
func GetDatabase () *Database {
	return DefaultInstance.GetDatabase ()
}

// Set the address of the instance (e.g. "localhost:8080")
func (in *Instance) SetAddress (addr string) {
	in.httpServer.Addr = addr
}

// Set the address of the default instance (e.g. "localhost:8080")
func SetAddress (addr string) {
	DefaultInstance.SetAddress (addr)
}

// Sets the path of the base-template
func (in *Instance) SetBaseTemplate (path string) error {
	var err error
	in.baseTemplate, err = template.ParseFiles (path)
	return err
}

// Sets the path of the base-template for the default instance
func SetBaseTemplate (path string) error {
	return DefaultInstance.SetBaseTemplate (path)
}

// Sets the path were static files are stored
func (in *Instance) SetStaticPath (path string) {
	in.staticPath = path
}

// Sets the path were static files are stored for the default instance
func SetStaticPath (path string) {
	DefaultInstance.SetStaticPath (path)
}

// Start serving requests
func (in *Instance) Run () error {
	var err error
	in.topNode, err = getTopNode (in.database)
	if err != nil {
		return err
	}
	err = in.httpServer.ListenAndServe ()
	return err
}

// Start serving requests
func Run () error {
	return DefaultInstance.Run ()
}

// Handles dynamic http requests
func (in *Instance) dynamicHandler (w http.ResponseWriter, r *http.Request) {
	node := in.topNode.getNodeForPath (r.URL.Path)
	if node == nil {
		http.NotFound (w, r)
		return
	}
	fragmentOptions := strings.FieldsFunc (node.fragmentOptions, func (r rune) bool {
		if r == ',' {
			return true
		}
		return false
	})
	fIface := FragmentInterface {
		UUID:		node.uuid,
		Request:	r,
		Instance:	in,
		Session:	getSession (w, r),
		Options:	fragmentOptions,
	}
	fragment := getFragment (node.fragment)
	tIface := TemplateInterface {
		MainMenu:	in.getMainMenu (),
		Body:		fragment.GetHtml (fIface),
	}
	in.baseTemplate.Execute (w, tIface)
}

// Serves static files
func (in *Instance) staticHandler (w http.ResponseWriter, r *http.Request) {
	m := staticRegexp.FindStringSubmatch (r.URL.Path)
	if m == nil {
		http.NotFound (w, r)
		return
	}
	http.ServeFile (w, r, in.staticPath+"/"+m[1])
}

// Creates the main menu from the node tree
func (in *Instance) getMainMenu () []MenuItem {
	mainMenu := []MenuItem {}
	mainMenu = append (mainMenu, MenuItem {in.topNode.displayName, "/", nil})
	for _, child := range in.topNode.children {
		if child.displayName != "" {
			menuItem := MenuItem {child.displayName, "/"+child.name, []MenuItem {}}
			for _, subchild := range child.children {
				if subchild.displayName != "" {
					menuItem.SubMenu = append (menuItem.SubMenu, MenuItem {subchild.displayName, "/"+child.name+"/"+subchild.name, nil})
				}
			}
			mainMenu = append (mainMenu, menuItem)
		}
	}
	return mainMenu
}
