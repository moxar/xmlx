package xmlx

import (
	"reflect"
	"strings"
	"testing"
)

func Test_ChunkAll(t *testing.T) {

	for label, c := range map[string]struct {
		in   string
		size int
		out  []string
	}{
		"three chunks": {
			size: 3,
			in: `
<music>
	<album name="Black Album">
		<meta>
			<band>Metallica</band>
			<year>1991</year>
		</meta>
	</album>
	<songs>
		<song>
			<name>Enter Sandman     </name>
			<number>1</number>
		</song>
		<song>
			<name>Sad but True       </name>
			<number>2</number>
		</song>
		<song>
			<name>Holier Than You    </name>
			<number>3</number>
		</song>
		<song>
			<name>The Unforgiven     </name>
			<number>4</number>
		</song>
		<song>
			<name>Wherever I May Roam</name>
			<number>5</number>
		</song>
		<song>
			<name>Don't Tread on Me  </name>
			<number>6</number>
		</song>
		<song>
			<name>Through the Never  </name>
			<number>7</number>
		</song>
	</songs>
</music>`,
			out: []string{
				`<song>
			<name>Enter Sandman     </name>
			<number>1</number>
		</song>
		<song>
			<name>Sad but True       </name>
			<number>2</number>
		</song>
		<song>
			<name>Holier Than You    </name>
			<number>3</number>
		</song>`,

				`<song>
			<name>The Unforgiven     </name>
			<number>4</number>
		</song>
		<song>
			<name>Wherever I May Roam</name>
			<number>5</number>
		</song>
		<song>
			<name>Don't Tread on Me  </name>
			<number>6</number>
		</song>`,

				`<song>
			<name>Through the Never  </name>
			<number>7</number>
		</song>`,
			},
		},
		"two chunks": {
			size: 5,
			in: `
<music>
	<album name="Black Album">
		<meta>
			<band>Metallica</band>
			<year>1991</year>
		</meta>
	</album>
	<songs>
		<song>
			<name>Enter Sandman     </name>
			<number>1</number>
		</song>
		<song>
			<name>Sad but True       </name>
			<number>2</number>
		</song>
		<song>
			<name>Holier Than You    </name>
			<number>3</number>
		</song>
		<song>
			<name>The Unforgiven     </name>
			<number>4</number>
		</song>
		<song>
			<name>Wherever I May Roam</name>
			<number>5</number>
		</song>
		<song>
			<name>Don't Tread on Me  </name>
			<number>6</number>
		</song>
		<song>
			<name>Through the Never  </name>
			<number>7</number>
		</song>
	</songs>
</music>`,
			out: []string{
				`<song>
			<name>Enter Sandman     </name>
			<number>1</number>
		</song>
		<song>
			<name>Sad but True       </name>
			<number>2</number>
		</song>
		<song>
			<name>Holier Than You    </name>
			<number>3</number>
		</song>
		<song>
			<name>The Unforgiven     </name>
			<number>4</number>
		</song>
		<song>
			<name>Wherever I May Roam</name>
			<number>5</number>
		</song>`,

				`<song>
			<name>Don't Tread on Me  </name>
			<number>6</number>
		</song>
		<song>
			<name>Through the Never  </name>
			<number>7</number>
		</song>`,
			},
		},
	} {
		segments, err := ChunkAll(strings.NewReader(c.in), "song", c.size)
		if err != nil {
			t.Log("on case", label)
			t.Log("unexpected error", err)
			t.Fail()
		}

		var out []string
		for _, s := range segments {
			out = append(out, c.in[s[0]:s[1]])
		}

		if !reflect.DeepEqual(out, c.out) {
			t.Log("on case", label)

			if len(out) != len(c.out) {
				t.Logf("expected: len %d", len(c.out))
				t.Logf("having: len %d", len(out))
				t.Fail()
				continue
			}

			for i := range out {
				if out[i] == c.out[i] {
					continue
				}
				t.Logf("expected:\n%v", c.out[i])
				t.Logf("having:\n%v", out[i])
				t.Fail()
			}
		}
	}
}

func Test_Chunk(t *testing.T) {

	for label, c := range map[string]struct {
		in  string
		out string
	}{
		"babar": {
			in: `
<music>
	<album name="Black Album">
		<meta>
			<band>Metallica</band>
			<year>1991</year>
		</meta>
	</album>
	<songs>
		<song>
			<name>Enter Sandman     </name>
			<number>1</number>
		</song>
		<song>
			<name>Sad but True       </name>
			<number>2</number>
		</song>
		<song>
			<name>Holier Than You    </name>
			<number>3</number>
		</song>
		<song>
			<name>The Unforgiven     </name>
			<number>4</number>
		</song>
		<song>
			<name>Wherever I May Roam</name>
			<number>5</number>
		</song>
		<song>
			<name>Don't Tread on Me  </name>
			<number>6</number>
		</song>
		<song>
			<name>Through the Never  </name>
			<number>7</number>
		</song>
	</songs>
</music>`,
			out: `<song>
			<name>Enter Sandman     </name>
			<number>1</number>
		</song>`,
		},
	} {
		segment, err := Chunk(strings.NewReader(c.in), "song", 0)
		if err != nil {
			t.Log("on case", label)
			t.Log("unexpected error", err)
			t.Fail()
		}

		out := c.in[segment[0]:segment[1]]

		if !reflect.DeepEqual(out, c.out) {
			t.Log("on case", label)

			if len(out) != len(c.out) {
				t.Logf("expected: len %d", len(c.out))
				t.Logf("having: len %d", len(out))
				t.Fail()
				continue
			}

			for i := range out {
				if out[i] == c.out[i] {
					continue
				}
				t.Logf("expected:\n%v", c.out[i])
				t.Logf("having:\n%v", out[i])
				t.Fail()
			}
		}
	}
}
