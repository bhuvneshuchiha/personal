package test3

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	name string
	a int
	b int
	expected int
}

func add(a, b int) int {
	return a + b
}

func TestAdd(t *testing.T) {
	tests := []testStruct{
		{
			name: "test_1",
			a: 10,
			b: 20,
			expected: 30,
		},
		{
			name: "test_2",
			a: 11,
			b: 26,
			expected: 37,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := add(tc.a, tc.b)
			assert.Equal(t, tc.expected, result)
		})
	}
}









