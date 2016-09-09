package xmlx

import (
	"encoding/xml"
	"fmt"
	"io"
	"math/rand"
)

// Chunk returns the start and stop position of the first encountered token.
func Chunk(reader io.ReadSeeker, token string, offset int64) ([2]int64, error) {

	// Return the error whatever it is: an EOF means there is no node corresponding
	// to the token.
	start, err := nextStartOffset(reader, token, offset)
	if err != nil {
		return [2]int64{}, err
	}

	stop, err := nextStopOffset(reader, token, start)
	if err != nil {
		return [2]int64{}, err
	}

	return [2]int64{start, stop}, nil
}

// ChunkAll reads the reader and defines a list of segments that correspond to valid chunks.
// Each segment, except for the last one, contains a number of chunks superior or equal to the
// provided bulkLen value.
func ChunkAll(reader io.ReadSeeker, token string, bulkLen int) ([][2]int64, error) {
	var segments [][2]int64
	var start, stop int64

	size, err := guessTokenSize(reader, token, 10)
	if err != nil {
		return nil, err
	}
	size *= int64(bulkLen)

	EOF, err := reader.Seek(0, 2)
	if err != nil {
		return nil, err
	}

	// Calculate segments until the end of the file.
	for err != io.EOF {

		start, err = nextStartOffset(reader, token, stop)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		// Find next sto
		jump := start + size
		stop, err = nextStopOffset(reader, token, jump)
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			stop, err = lastStopOffset(reader, token, start, EOF)
			if err != nil {
				return nil, err
			}
		}

		segment := [2]int64{start, stop}
		segments = append(segments, segment)
	}

	return segments, nil
}

// guessTokenSize returns the size of one token found by the parser. If the parser could
// not find the expected token, this function returns an error. The size is calculed as
// the average size of random token found after n iterations.
func guessTokenSize(reader io.ReadSeeker, token string, iteration int) (int64, error) {

	var avgs []int64

	// Get the len of the reader.
	EOF, err := reader.Seek(0, 2)
	if err != nil {
		return 0, err
	}

	// Calculate an average size based on the number of iterations.
	// Start at position 0 to quickly find the node if it is unique.
	var pos int64
	for i := 0; i < iteration; i++ {

		_, err = reader.Seek(pos, 0)
		if err != nil {
			return 0, err
		}

		start, err := nextStartOffset(reader, token, pos)
		if err != nil {

			if err != io.EOF {
				return 0, nil
			}

			// This condition ensures the token exists: if the encountered error is an io.EOF,
			// the loop should continue. However, if an EOF is encountered when the position
			// is still 0, it means that the reader could not find any token matching the parsers
			// token.
			if pos == 0 {
				return 0, fmt.Errorf("cannot find token \"%s\" in reader", token)
			}

			continue
		}

		stop, err := nextStopOffset(reader, token, start)
		if err != nil {
			continue
		}

		avgs = append(avgs, stop-start)

		pos = rand.Int63n(EOF)
	}

	var sum int64
	for _, a := range avgs {
		sum += a
	}

	return int64(sum / int64(len(avgs))), nil
}

// nextStartOffset returns the offset of the byte before the next start token.
func nextStartOffset(reader io.ReadSeeker, token string, offset int64) (int64, error) {

	_, err := reader.Seek(offset, 0)
	if err != nil {
		return 0, err
	}

	var last int64
	decoder := xml.NewDecoder(reader)
	for {
		t, err := decoder.RawToken()
		if err != nil {
			return 0, err
		}

		// Break the loop as soon as the token is found.
		elt, ok := t.(xml.StartElement)
		if ok && elt.Name.Local == token {
			return last + offset, nil
		}

		last = decoder.InputOffset()
	}
}

// nextStopOffset returns the offset of the byte after the next stop token.
func nextStopOffset(reader io.ReadSeeker, token string, offset int64) (int64, error) {

	offset, err := reader.Seek(offset, 0)
	if err != nil {
		return 0, err
	}

	decoder := xml.NewDecoder(reader)
	for {
		t, err := decoder.RawToken()
		if err != nil {
			return 0, err
		}

		// Break the loop as soon as the token is found.
		elt, ok := t.(xml.EndElement)
		if ok && elt.Name.Local == token {
			return offset + decoder.InputOffset(), nil
		}
	}
}

// lastStopOffset returns the offset of the byte after the last stop token.
// It processes a dichotomial research in the reader.
func lastStopOffset(reader io.ReadSeeker, token string, start, stop int64) (int64, error) {

	half := start + (stop-start)/2
	pos, err := nextStopOffset(reader, token, half)
	if err != nil {

		// Went to far
		if err == io.EOF {
			return lastStopOffset(reader, token, start, half)
		}

		return 0, err
	}

	_, err = nextStopOffset(reader, token, pos)
	if err != nil {

		// Last offset was the good one.
		if err == io.EOF {
			return pos, nil
		}

		return 0, err
	}

	// position was too short
	return lastStopOffset(reader, token, half, stop)
}
