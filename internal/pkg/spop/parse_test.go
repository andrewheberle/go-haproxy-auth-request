package spop

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseBinaryHeader(t *testing.T) {
	// test cases

	// one header
	one := make([]byte, 0)
	one = append(one, 8)
	one = append(one, []byte("X-Header")...)
	one = append(one, 5)
	one = append(one, []byte("Value")...)
	one = append(one, 0, 0)

	// two headers
	two := make([]byte, 0)
	two = append(two, 8)
	two = append(two, []byte("X-Header")...)
	two = append(two, 5)
	two = append(two, []byte("Value")...)
	two = append(two, 16)
	two = append(two, []byte("X-Another-Header")...)
	two = append(two, 13)
	two = append(two, []byte("Another Value")...)
	two = append(two, 0, 0)

	// invalid (incorrect length)
	invalid := make([]byte, 0)
	invalid = append(invalid, 7)
	invalid = append(invalid, []byte("X-Header")...)
	invalid = append(invalid, 5)
	invalid = append(invalid, []byte("Value")...)
	invalid = append(invalid, 0, 0)

	// truncated (no empty values at end)
	truncated := make([]byte, 0)
	truncated = append(truncated, 7)
	truncated = append(truncated, []byte("X-Header")...)
	truncated = append(truncated, 5)
	truncated = append(truncated, []byte("Value")...)

	tests := []struct {
		name    string
		b       []byte
		want    http.Header
		wantErr bool
	}{
		{"empty", []byte{0, 0}, make(http.Header), false},
		{"one", one, http.Header{"X-Header": []string{"Value"}}, false},
		{"two", two, http.Header{"X-Header": []string{"Value"}, "X-Another-Header": []string{"Another Value"}}, false},
		{"truncated", truncated, nil, true},
		{"invalid", invalid, nil, true},
	}

	for _, tt := range tests {
		headers, err := ParseBinaryHeader(tt.b)
		if tt.wantErr {
			assert.NotNil(t, err, tt.name, "should not be nil")
			continue
		}

		assert.Nil(t, err, tt.name, "should not nil")
		assert.Equal(t, tt.want, headers, tt.name, "should equal")
	}
}
