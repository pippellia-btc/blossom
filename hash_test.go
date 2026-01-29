package blossom

import (
	"fmt"
	"testing"
)

func TestValidateHash(t *testing.T) {
	tests := []struct {
		input   string
		isValid bool
	}{
		{
			input:   "", // empty string: length must be exactly 64
			isValid: false,
		},
		{
			input:   "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcde", // too short: 63 characters instead of 64
			isValid: false,
		},
		{
			input:   "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0", // too long: 65 characters instead of 64
			isValid: false,
		},
		{
			input:   "g123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef", // invalid character 'g' at start (not a hex digit)
			isValid: false,
		},
		{
			input:   "0123456789abcdef0123456789abcdef0123456789zbcdef0123456789abcdef", // invalid character 'z' in middle (not a hex digit)
			isValid: false,
		},
		{
			input:   "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdeg", // invalid character 'g' at end (not a hex digit)
			isValid: false,
		},
		{
			input:   "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdeG", // invalid character 'G' at end (not a hex digit)
			isValid: false,
		},
		{
			input:   "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcde ", // invalid character space at end (not a hex digit)
			isValid: false,
		},
		{
			input:   "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcde-", // invalid character dash at end (not a hex digit)
			isValid: false,
		},
		{
			input:   "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcde_", // invalid character underscore at end (not a hex digit)
			isValid: false,
		},
		{
			input:   "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
			isValid: true,
		},
		{
			input:   "0123456789ABCDEF0123456789ABCDEF0123456789ABCDEF0123456789ABCDEF",
			isValid: true,
		},
		{
			input:   "0123456789aBcDeF0123456789aBcDeF0123456789aBcDeF0123456789aBcDeF",
			isValid: true,
		},
		{
			input:   "0000000000000000000000000000000000000000000000000000000000000000",
			isValid: true,
		},
		{
			input:   "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
			isValid: true,
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("case=%d", i), func(t *testing.T) {
			err := ValidateHash(test.input)
			isValid := err == nil

			if isValid != test.isValid {
				t.Errorf("ValidateHash(%q) = %v, want isValid %v", test.input, err, test.isValid)
			}
		})
	}
}
