package random

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRandomString(t *testing.T) {
	tests := []struct {
		name string
		size uint8
	}{
		{
			name: "size = 1",
			size: 1,
		},
		{
			name: "size = 5",
			size: 5,
		},
		{
			name: "size = 10",
			size: 10,
		},
		{
			name: "size = 20",
			size: 20,
		},
		{
			name: "size = 30",
			size: 30,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			str1, err1 := NewRandomString(tt.size)
			str2, err2 := NewRandomString(tt.size)

			assert.Len(t, str1, int(tt.size))
			assert.Len(t, str2, int(tt.size))

			assert.NoError(t, err1)
			assert.NoError(t, err2)

			// Check that two generated strings are different
			// This is not an absolute guarantee that the function works correctly,
			// but this is a good heuristic for a simple random generator.
			assert.NotEqual(t, str1, str2)
		})
	}

	testsNilSize := []struct {
		name string
		size uint8
	}{
		{
			name: "size = 0",
			size: 0,
		},
	}

	for _, tt := range testsNilSize {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewRandomString(tt.size)

			assert.Error(t, err)
		})
	}
}
