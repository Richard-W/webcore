package webcore

import (
	strings	"strings"
	errors	"errors"
)

var ErrTopNode = errors.New ("Not exactly one top node")

// A node represents a single entity in the webservice.
type node struct {
	uuid		string	// Identifier
	parent		*node	// Pointer to parent
	children	[]*node	// Pointers to children
	name		string	// Used to build paths
	displayName	string	// Used in menus
	fragment	string	// Name of the body-fragment
	fragmentOptions	string	// An arbitrary string the fragment can use
	deep		bool	// True when this node should handle all subpaths
}

// Get a node for a specific path
func (n *node) getNodeForPath (path string) *node {
	pathFields := strings.FieldsFunc (path, func (r rune) bool {
		if r == '/' {
			return true
		}
		return false
	})
	return n.getNodeForPathFields (pathFields)
}


// See getNodeForPath
func (n *node) getNodeForPathFields (pathFields []string) *node {
	if len (pathFields) == 0 {
		return n
	}
	for _, child := range n.children {
		if child.name == pathFields[0] {
			if child.deep {
				return child
			}
			return child.getNodeForPathFields (pathFields[1:])
		}
	}
	return nil
}

// Get the top node of the node tree
func getTopNode (db *Database) (*node, error) {
	topNodes, err := getNodesByParent (nil, db)
	if err != nil {
		return nil, err
	}
	if len (topNodes) != 1 {
		return nil, ErrTopNode
	}
	return topNodes[0], nil
}

// Recursively get children of nodes
func getNodesByParent (parent *node, db *Database) ([]*node, error) {
	var uuid string
	if parent == nil {
		uuid = ""
	} else {
		uuid = parent.uuid
	}
	result := []*node {}
	res, err := db.Query ("SELECT `uuid`, `name`, `display_name`, `fragment`, `fragment_options`, `deep` FROM `nodes` WHERE `parent_id` = ?", uuid)
	if err != nil {
		return nil, err
	}
	for res.Next () {
		newNode := new (node)
		var deep int
		res.Scan (&newNode.uuid, &newNode.name, &newNode.displayName, &newNode.fragment, &newNode.fragmentOptions, &deep)
		if deep == 0 {
			newNode.deep = false
		} else {
			newNode.deep = true
		}
		newNode.parent = parent
		newNode.children, err = getNodesByParent (newNode, db)
		if err != nil {
			return nil, err
		}
		result = append (result, newNode)
	}
	return result, nil
}
