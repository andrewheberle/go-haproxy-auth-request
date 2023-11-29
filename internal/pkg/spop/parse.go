package spop

import (
	"fmt"
	"net/http"
)

// ParseBinaryHeader takes binary headers from encoded as per
// http://docs.haproxy.org/2.8/configuration.html#7.3.6-req.hdrs_bin and parses
// it into a http.Header type
func ParseBinaryHeader(b []byte) (http.Header, error) {
	var blankFound bool

	// start with new http.Header struct
	h := make(http.Header)

	p := 0
	for p < len(b) {
		// start from p in slice
		buf := b[p:]

		// decode header name
		name, n, err := decodeString(buf)
		if err != nil {
			return nil, fmt.Errorf("parse binary header name: %w", err)
		}

		// advance n bytes if possible
		if p+n >= len(b) {
			return nil, fmt.Errorf("parse binary header name: went past end of slice")
		}
		p = p + n

		// start from new p
		buf = b[p:]
		value, n, err := decodeString(buf)
		if err != nil {
			return nil, fmt.Errorf("parse binary header value: %w", err)
		}

		// advance n bytes if possible
		if p+n > len(b) {
			return nil, fmt.Errorf("parse binary header value: went past end of slice")
		}
		p = p + n

		// final is blank header and value
		if name == "" && value == "" {
			blankFound = true
			break
		}

		// add decoded data
		h.Add(name, value)
	}

	// check that blanks were found
	if !blankFound {
		return nil, fmt.Errorf("headers possibly trucated as final blanks found")
	}

	return h, nil
}

func decodeVarint(b []byte) (int, int, error) {
	if len(b) == 0 {
		return 0, 0, fmt.Errorf("decode varint: empty slice")
	}

	val := int(b[0])
	off := 1

	if val < 240 {
		return val, 1, nil
	}

	r := uint(4)
	for {
		if off > len(b)-1 {
			return 0, 0, fmt.Errorf("decode varint: unterminated sequence")
		}

		v := int(b[off])
		val += v << r
		off++
		r += 7

		if v < 128 {
			break
		}
	}

	return val, off, nil
}

func decodeBytes(b []byte) ([]byte, int, error) {
	l, off, err := decodeVarint(b)
	if err != nil {
		return nil, 0, fmt.Errorf("decode bytes: %w", err)
	}

	if len(b) < l+off {
		return nil, 0, fmt.Errorf("decode bytes: unterminated sequence")
	}

	return b[off : off+l], off + l, nil
}

func decodeString(b []byte) (string, int, error) {
	b, n, err := decodeBytes(b)
	return string(b), n, err
}
