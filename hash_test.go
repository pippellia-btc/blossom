package blossom

import (
	"encoding/json"
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

func TestHashJSON(t *testing.T) {
	tests := []struct {
		name  string
		input string
		hex   string
		err   bool
	}{
		// error cases
		{name: "too short", input: `"0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcde"`, err: true},
		{name: "too long", input: `"0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0"`, err: true},
		{name: "invalid hex", input: `"g123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"`, err: true},
		{name: "number", input: `12345`, err: true},
		{name: "object", input: `{}`, err: true},
		{name: "empty", input: `""`, err: true},

		// success cases
		{name: "lowercase", input: `"0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"`, hex: "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"},
		{name: "uppercase", input: `"0123456789ABCDEF0123456789ABCDEF0123456789ABCDEF0123456789ABCDEF"`, hex: "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"},
		{name: "zeros", input: `"0000000000000000000000000000000000000000000000000000000000000000"`, hex: "0000000000000000000000000000000000000000000000000000000000000000"},
		{name: "all f", input: `"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"`, hex: "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var h Hash
			err := json.Unmarshal([]byte(test.input), &h)

			if test.err {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if h.Hex() != test.hex {
				t.Errorf("got %s, want %s", h.Hex(), test.hex)
			}

			// roundtrip
			data, _ := json.Marshal(h)
			var h2 Hash
			json.Unmarshal(data, &h2)
			if h != h2 {
				t.Errorf("roundtrip mismatch: %x != %x", h, h2)
			}
		})
	}
}

func TestHashScan(t *testing.T) {
	tests := []struct {
		name  string
		input any
		hex   string
		err   bool
	}{
		// error cases
		{name: "nil", input: nil, err: true},
		{name: "int", input: 12345, err: true},
		{name: "float", input: 123.45, err: true},
		{name: "string too short", input: "abc", err: true},
		{name: "string too long", input: "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0", err: true},
		{name: "string invalid hex", input: "g123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef", err: true},
		{name: "bytes too short", input: []byte("abc"), err: true},
		{name: "bytes too long", input: []byte("0123456789abcdef0123456789abcdef0"), err: true},

		// success cases
		{name: "string", input: "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef", hex: "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"},
		{name: "string uppercase", input: "0123456789ABCDEF0123456789ABCDEF0123456789ABCDEF0123456789ABCDEF", hex: "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"},
		{name: "bytes", input: []byte("01onal<>?!@#$%^&*()_+-={}[]|:;<>"), hex: "30316f6e616c3c3e3f21402324255e262a28295f2b2d3d7b7d5b5d7c3a3b3c3e"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var h Hash
			err := h.Scan(test.input)

			if test.err {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if h.Hex() != test.hex {
				t.Errorf("got %s, want %s", h.Hex(), test.hex)
			}
		})
	}
}
