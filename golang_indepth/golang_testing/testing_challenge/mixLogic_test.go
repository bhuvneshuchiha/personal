package testing_challenge

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type caseDump struct {
	name string
	a int
	b int
	expected int
	errExpected bool
}

func Test_Add(t *testing.T) {
	tests := []caseDump{
		{
			name:"add_1",
			a: 10,
			b: 20,
			expected: 30,
			errExpected: false,
		},
		{
			name:"add_2",
			a: 21,
			b: 70,
			expected: 91,
			errExpected: false,
		},
		{
			name:"add_3",
			a: 0,
			b: -1,
			expected: -1,
			errExpected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := Add(tc.a, tc.b)
			assert.Equal(t, tc.expected, result, "Actual result != Expected result")
		})
	}

}

func TestSubtract(t *testing.T) {
	tests := []caseDump{
		{
			name:"subtract_1",
			a: 10,
			b: 20,
			expected: 10,
			errExpected: false,
		},
		{
			name:"subtract_2",
			a: 21,
			b: 70,
			expected: 49,
			errExpected: false,
		},
		{
			name:"subtract_3",
			a: 0,
			b: -1,
			expected: 1,
			errExpected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := Subtract(tc.a, tc.b)
			assert.Equal(t, tc.expected, result, "Actual result != Expected result")
		})
	}
}


func TestDivide(t *testing.T) {
	tests := []caseDump{
		{
			name:"divide_1",
			a: 20,
			b: 20,
			expected: 1,
			errExpected: false,
		},
		{
			name:"divide_2",
			a: 140,
			b: 70,
			expected: 2,
			errExpected: false,
		},
		{
			name:"divide_3",
			a: 1,
			b: 0,
			expected: 0,
			errExpected: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result , err := Divide(tc.a, tc.b)
			if tc.errExpected {
				assert.Error(t, err)
			}else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, result, "Actual result != Expected result")
			}
		})
	}
}


func TestMultiply(t *testing.T) {
	tests := []caseDump{
		{
			name:"multiply_1",
			a: 10,
			b: 20,
			expected: 200,
			errExpected: false,
		},
		{
			name:"multiply_2",
			a: 20,
			b: 70,
			expected: 1400,
			errExpected: false,
		},
		{
			name:"multiply_3",
			a: 0,
			b: -1,
			expected: 0,
			errExpected: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := Multiply(tc.a, tc.b)
			assert.Equal(t, tc.expected, result, "Actual result != Expected result")
		})
	}
}
