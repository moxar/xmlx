# xmlx

## introduction

This package xmlx extends the features of the encoding/xml package with a 
generic unmarshallable xml structure.

It provides a new structure, Node, which can unmarshal any xml data. This node has two useful
methods: Map and Split.

The Map method returns a map of name, data, attributes and subnodes to
their values.

The Split method returns an array of nodes having the same property as the parent,
splitted after a subnode name.

## installation

```
	$ go get github.com/moxar/xmlx
```

## usage

```go

	// Lets assume you have an xml input that looks like this:
	input := `<music>
			<album>
				<songs>
					<song>
						<name>Don't Tread on Me</name>
						<number>6</number>
					</song>
					<song>
						<name>Through the Never</name>
						<number>7</number>
					</song>
				</songs>
			</album>
		</music>`

		// unmarshall it into the generic Node structure.
		var node Node
		err := xml.Unmarshal([]byte(input), &node)
		if err != nil {
			// do stuff...
		}
		
		fmt.Println(node)
		
		// If you need to split it into several nodes, just call node.Split
		for _, n := range node.Split("album.songs") {
			fmt.Println(n)
		}
	
		// Output:
		// {music map[]  [{album map[]  [{songs map[]  [{song map[]  [{name map[] Don't Tread on Me [] } {number map[] 6 [] }] } {song map[]  [{name map[] Through the Never [] } {number map[] 7 [] }] }] }] }] }
		// {music map[]  [{songs map[]  [{name map[] Don't Tread on Me [] } {number map[] 6 [] }] }] }
		// {music map[]  [{songs map[]  [{name map[] Through the Never [] } {number map[] 7 [] }] }] }
```

Of course, this assumes that you know the incomming structure, and when you know it, you can create a custom
structure with reflect xml tags. In such case, this package is useless.

However, it was created to be part of an ETL where the data sources format can be very different from
one source to another, and can often change. With the **one structure per source** solution, we have to
write the structure, ensure te unmarshal behaves well, recompile and redeploy each time a source changes, 
or a new one is added.
With the xmlx.Node solution, we can associate a configuration to a source, which means there is no need
to recompile each time a source changes.
