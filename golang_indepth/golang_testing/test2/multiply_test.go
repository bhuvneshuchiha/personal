package test2

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func test_multiply(t *testing.T) {
	result := multiply(2,6)
	expected := 12
	assert.Equal(t, expected, result, "Both these should be equal")
}
