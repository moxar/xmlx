// Package xmlx extends the features of the encoding/xml package with a 
// generic unmarshallable xml structure.
//
// It provides a new structure, Node, which can unmarshal any xml data. This node has two useful
// methods: Map and Split.
//
// The Map method returns a map of name, data, attributes and subnodes to
// their values.
//
// The Split method returns an array of nodes having the same property as the parent,
// splitted after a subnode name.
package xmlx
