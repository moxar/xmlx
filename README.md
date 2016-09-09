# xmlx

## introduction

This package xmlx extends the features of the encoding/xml package with a 
generic unmarshallable xml structure.

It provides a new structure, Node, which can unmarshal any xml data. This node has two useful
methods: Map and Split.

* The Map method returns a map of name, data, attributes and subnodes to
their values.

* The Split method returns an array of nodes having the same property as the parent,
splitted after a subnode name.

## installation

```
	$ go get github.com/moxar/xmlx
```

## usage

### Node

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
one source to another, and can often change. 

With the _one structure per source_ solution, we have to
write the structure, ensure te unmarshal behaves well, recompile and redeploy each time a source changes, 
or a new one is added.

With the xmlx.Node solution, we can associate a configuration to a source, which means there is no need
to recompile each time a source changes.

### Chunk

Lets assume you have a big XML input located in a file. Reading the file and unmarshalling it would
take a lot of time, so you decide to parallelise the unmarshalling. 
The Chunk and ChunkAll methods allows you to retreive start and stop offsets of valid xml chunks.

For instance, if you have an XML document with 7 nodes, the ChunkAll method could return two segments:
the first defines the position of the segment 1 to 5, the second defines the position of the segment 6 to 7.

Once you have it, you can extract the bytes from the given segments and parallelize the unmarshalling
on each segment.

```go
	
	type Song struct{
		Name   string
		Number int
	}

	// lets take the list of the Metallica's Black Album 7 first soundtracks.
	input := `
	<music>
		<album name="Black Album">
			<meta>
				<band>Metallica</band>
				<year>1991</year>
			</meta>
		</album>
		<songs>
			<song>
				<name>Enter Sandman</name>
				<number>1</number>
			</song>
			<song>
				<name>Sad but True</name>
				<number>2</number>
			</song>
			<song>
				<name>Holier Than You</name>
				<number>3</number>
			</song>
			<song>
				<name>The Unforgiven</name>
				<number>4</number>
			</song>
			<song>
				<name>Wherever I May Roam</name>
				<number>5</number>
			</song>
			<song>
				<name>Don't Tread on Me</name>
				<number>6</number>
			</song>
			<song>
				<name>Through the Never</name>
				<number>7</number>
			</song>
		</songs>
	</music>`
	
	// Chunk the input. Each chunk must contain valid "song" tokens, and each chunk must contain
	// at least 2 tokens. The last chunk may contain less.
	segments, err := ChunkAll(strings.NewReader(input), "song", 2)
	if err != nil {
		// do stuff...
	}
	
	// The input has been segmented, parallelization of unmarshall can be done.
	for _, s := range segments {
		go func(s [2]int64) {
			dec := xml.NewDecoder(strings.NewReader(input[s[0]:s[1]]))
			for {
				// also works with: 
				// var node Node
				var song Song
				
				err := dec.Decode(&song)
				if err != nil {
					if err == io.EOF {
						break
					}
					// to stuff...
				}
				
				// do stuff with the new song...
			}
		}(s)
	}
```
