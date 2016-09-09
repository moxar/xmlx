package xmlx

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"strings"
)

// Node is a generic XML node.
type Node struct {

	// The name of the node
	Name string

	// The attribute list of the node
	Attrs map[string]string

	// The data located within the node
	Data string

	// The subnodes within the node
	Nodes []Node

	prefix string
}

// UnmarshalXML takes the content of an XML node and puts it into the Node structure.
func (n *Node) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {

	if len(start.Attr) != 0 {
		n.Attrs = map[string]string{}
		for _, v := range start.Attr {
			n.Attrs[v.Name.Local] = v.Value
		}
	}
	n.Name = start.Name.Local

	balance := 1

	for balance != 0 {
		token, err := d.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		switch t := token.(type) {

		case xml.CharData:
			v := bytes.TrimSpace(t)
			if len(v) == 0 {
				continue
			}
			n.Data = string(t)

		case xml.StartElement:
			if t.Name.Local == start.Name.Local {
				balance++
				continue
			}

			node := Node{}
			err := node.UnmarshalXML(d, t)
			if err != nil {
				return err
			}

			n.Nodes = append(n.Nodes, node)

		case xml.EndElement:
			if t.Name.Local == start.Name.Local {
				balance--
			}
		}
	}

	return nil
}

// Split the node into many: each time the split label is encountered within a subnode of the node,
// a new node is created.
func (n Node) Split(label string) []Node {

	// Return the node itself if no label is specified: there is no split to do.
	if len(label) == 0 {
		return []Node{n}
	}

	// Create a leveled array of children.
	terms := strings.Split(label, ".")
	gen := len(terms)
	children := make([][]Node, gen+1, gen+1)
	children[0] = n.Nodes

	// Explore the leveled array. On each level, put the children of the matching node
	// to the upper level.
	for i, term := range terms {
		for _, node := range children[i] {
			if node.Name == term {
				children[i+1] = append(children[i+1], node.Nodes...)
			}
		}
	}

	// Create a node for each child. Also rename the child with label.
	var nodes []Node
	for _, child := range children[gen] {
		node := n.clone()
		var i int
		for i < len(node.Nodes) {
			if node.Nodes[i].Name == terms[0] {
				node.Nodes = append(node.Nodes[:i], node.Nodes[i+1:]...)
			}
			i++
		}
		child.Name = terms[gen-1]
		node.Nodes = append(node.Nodes, child)
		nodes = append(nodes, node)
	}

	return nodes
}

// Map returns a flatten representation of the node. If a node contains nodes
// having the same name, only the last node will exist in the map.
func (n Node) Map() map[string]string {

	out := n.flatten()

	var toProcess []Node
	for _, node := range n.Nodes {
		node.prefix = fmt.Sprintf("#nodes.%s", node.Name)
		toProcess = append(toProcess, node)
	}

	var i int
	for i < len(toProcess) {
		node := toProcess[i]
		for k, v := range node.flatten() {
			name := fmt.Sprintf("%s.%s", node.prefix, k)
			out[name] = v
		}

		for _, child := range node.Nodes {
			child.prefix = fmt.Sprintf("%s.#nodes.%s", node.prefix, child.Name)
			toProcess = append(toProcess, child)
		}
		i++
	}

	return out
}

// flatten takes a node a generates a map with it.
func (n Node) flatten() map[string]string {

	var t = map[string]string{}

	// put simple values into transcient.
	if len(n.Name) != 0 {
		t["#name"] = n.Name
	}
	if len(n.Data) != 0 {
		t["#data"] = n.Data
	}

	// put attributes into transcient.
	for k, v := range n.Attrs {
		name := fmt.Sprintf("#attr.%s", k)
		t[name] = v
	}

	return t
}

// clone creates a copy of the node and returns it.
func (n Node) clone() Node {

	node := Node{
		Name:   n.Name,
		Attrs:  n.Attrs,
		Data:   n.Data,
		prefix: n.prefix,
	}
	for i := range n.Nodes {
		node.Nodes = append(node.Nodes, n.Nodes[i].clone())
	}

	return node
}
