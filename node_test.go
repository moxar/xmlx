package xmlx

import (
	"encoding/xml"
	"fmt"
	"reflect"
	"testing"
)

func Test_NodeUnmarshalXML(t *testing.T) {

	for i, c := range []struct {
		input    string
		expected interface{}
	}{
		{
			input: `
            <foo>
                <bar/>
            </foo>
			`,
			expected: Node{
				Name:  "foo",
				Attrs: nil,
				Nodes: []Node{
					{
						Name:  "bar",
						Attrs: nil,
					},
				},
			},
		},
		{
			input: `
<music>
	<album name="Black Album">
		<meta>
			<band>Metallica</band>
			<year>1991</year>
		</meta>
	</album>
</music>
			`,
			expected: Node{
				Name:  "music",
				Nodes: []Node{
					{
						Name: "album",
						Attrs: map[string]string{"name": "Black Album"},
						Nodes: []Node{
							{
								Name:  "meta",
								Nodes: []Node{
									{
										Name: "band",
										Data: "Metallica",
									},
									{
										Name: "year",
										Data: "1991",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			input: `
<?xml version="1.0"?>
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
</music>
			`,
			expected: Node{
				Name:  "music",
				Nodes: []Node{
					{
						Name: "album",
						Attrs: map[string]string{"name": "Black Album"},
						Nodes: []Node{
							{
								Name:  "meta",
								Nodes: []Node{
									{
										Name: "band",
										Data: "Metallica",
									},
									{
										Name: "year",
										Data: "1991",
									},
								},
							},
						},
					},
					{
						Name: "songs",
						Nodes: []Node{
							{
								Name: "song",
								Nodes: []Node{
									{
										Name: "name",
										Data: "Enter Sandman",
									},
									{
										Name: "number",
										Data: "1",
									},
								},
							},
							{
								Name: "song",
								Nodes: []Node{
									{
										Name: "name",
										Data: "Sad but True",
									},
									{
										Name: "number",
										Data: "2",
									},
								},
							},
							{
								Name: "song",
								Nodes: []Node{
									{
										Name: "name",
										Data: "Holier Than You",
									},
									{
										Name: "number",
										Data: "3",
									},
								},
							},
							{
								Name: "song",
								Nodes: []Node{
									{
										Name: "name",
										Data: "The Unforgiven",
									},
									{
										Name: "number",
										Data: "4",
									},
								},
							},
							{
								Name: "song",
								Nodes: []Node{
									{
										Name: "name",
										Data: "Wherever I May Roam",
									},
									{
										Name: "number",
										Data: "5",
									},
								},
							},
							{
								Name: "song",
								Nodes: []Node{
									{
										Name: "name",
										Data: "Don't Tread on Me",
									},
									{
										Name: "number",
										Data: "6",
									},
								},
							},
							{
								Name: "song",
								Nodes: []Node{
									{
										Name: "name",
										Data: "Through the Never",
									},
									{
										Name: "number",
										Data: "7",
									},
								},
							},
						},
					},
				},
			},
		},
	} {

		having := Node{}
		err := xml.Unmarshal([]byte(c.input), &having)
		if err != nil {
			t.Logf("failed case %d: %s", i+1, err)
			t.Fail()
		}

		if !reflect.DeepEqual(having, c.expected) {
			t.Logf("failed case %d", i+1)
			t.Logf("having:\n\n%s\n", having)
			t.Logf("expected:\n\n%s\n", c.expected)
			t.Fail()
		}
	}
}

func Test_NodeMap(t *testing.T) {

	for i, c := range []struct {
		in  Node
		out map[string]string
		err error
	}{
		{
			in: Node{
				Name: "foo",
				Attrs: map[string]string{
					"len":      "7",
					"priority": "0",
				},
				Nodes: []Node{
					{
						Name: "band",
						Data: "ACDC",
					},
					{
						Name: "size",
						Data: "4",
					},
				},
			},
			out: map[string]string{
				"#name":             "foo",
				"#attr.len":         "7",
				"#attr.priority":    "0",
				"#nodes.band.#name": "band",
				"#nodes.band.#data": "ACDC",
				"#nodes.size.#name": "size",
				"#nodes.size.#data": "4",
			},
			err: nil,
		},
		{
			in: Node{
				Name: "foo",
				Nodes: []Node{
					{
						Name: "bars",
						Nodes: []Node{
							{
								Name: "poo",
								Data: "0",
							},
						},
					},
				},
			},
			out: map[string]string{
				"#name":                        "foo",
				"#nodes.bars.#name":            "bars",
				"#nodes.bars.#nodes.poo.#name": "poo",
				"#nodes.bars.#nodes.poo.#data": "0",
			},
			err: nil,
		},
		{
			in: Node{
				Name: "Paris",
				Attrs: map[string]string{
					"type": "city",
				},
				Nodes: []Node{
					{
						Name: "foo",
						Data: "bar",
					},
					{
						Name: "geo",
						Attrs: map[string]string{
							"mode": "carthesian",
						},
						Nodes: []Node{
							{
								Name: "lat",
								Data: "-2.41",
							},
							{
								Name: "long",
								Data: "13.4",
							},
						},
					},
				},
			},
			out: map[string]string{
				"#name":               "Paris",
				"#attr.type": "city",

				"#nodes.foo.#name": "foo",
				"#nodes.foo.#data": "bar",

				"#nodes.geo.#name":      "geo",
				"#nodes.geo.#attr.mode": "carthesian",

				"#nodes.geo.#nodes.lat.#name": "lat",
				"#nodes.geo.#nodes.lat.#data": "-2.41",

				"#nodes.geo.#nodes.long.#name": "long",
				"#nodes.geo.#nodes.long.#data": "13.4",
			},
			err: nil,
		},
	} {
		out := c.in.Map()
		if !reflect.DeepEqual(out, c.out) {
			t.Logf("case %d:", i+1)
			t.Logf("expecting: %v", c.out)
			t.Logf("having: %v", out)
			t.Fail()
		}
	}
}

func ExampleNode() {

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

	var node Node
	err := xml.Unmarshal([]byte(input), &node)
	if err != nil {
		// do stuff...
	}
	
	fmt.Println(node)
	for _, n := range node.Split("album.songs") {
		fmt.Println(n)
	}
	
	// Output:
	// {music map[]  [{album map[]  [{songs map[]  [{song map[]  [{name map[] Don't Tread on Me [] } {number map[] 6 [] }] } {song map[]  [{name map[] Through the Never [] } {number map[] 7 [] }] }] }] }] }
	// {music map[]  [{songs map[]  [{name map[] Don't Tread on Me [] } {number map[] 6 [] }] }] }
	// {music map[]  [{songs map[]  [{name map[] Through the Never [] } {number map[] 7 [] }] }] }
}
