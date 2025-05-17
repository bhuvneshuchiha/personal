package test1

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	result := Add(1,2)
	expected := 3
	assert.Equal(t, expected, result, "expected and actual should be equal")
}

