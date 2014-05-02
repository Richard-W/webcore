package webcore

import "testing"

func TestGetNodeForPath (t *testing.T) {
	ntop := new (node)
	ntop.name = "top"
	n1 := new (node)
	n1.name = "1"
	n2 := new (node)
	n2.name = "2"
	n3 := new (node)
	n3.name = "3"
	na := new (node)
	na.name = "a"
	nb := new (node)
	nb.name = "b"
	ntop.children = []*node {n1, na}
	n1.children = []*node {n2}
	n2.children = []*node {n3, nb}

	path := "/1/2/3"
	node := ntop.getNodeForPath (path)
	if node != n3 {
		t.Error ("%s does not lead to n3", path)
	}

	path = "/1/2/b"
	node = ntop.getNodeForPath (path)
	if node != nb {
		t.Error ("%s does not lead to nb", path)
	}

	path = "/1/2/a"
	node = ntop.getNodeForPath (path)
	if node != nil {
		t.Error ("%s should not lead somewhere", path)
	}
}
