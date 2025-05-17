package testing_challenge_2

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

type caseDump struct {
	name string
	a int
	expected int
	errExpected bool
}

func TestSqrt(t *testing.T) {
	tests := []caseDump {
		{
			name:"test_1",
			a: 144,
			expected: 12,
			errExpected: false,
		},
		{
			name:"test_2",
			a: -144,
			expected: 0,
			errExpected: true,
		},
		{
			name:"test_3",
			a: 625,
			expected: 25,
			errExpected: false,
		},
		{
			name:"test_4",
			a: 0,
			expected: 0,
			errExpected: false,
		},
		{
			name:"test_5",
			a: -1,
			expected: 0,
			errExpected: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := Sqrt(tc.a)
			if tc.errExpected {
				assert.Error(t, err)
			}else {
				assert.NoError(t, err)
				assert.Equal(t, result, tc.expected, "Result is not correct")
			}
		})
	}
}

