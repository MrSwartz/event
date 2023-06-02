package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIPstringV4ToInt(t *testing.T) {
	// https://www.browserling.com/tools/ip-to-dec
	tests := []struct {
		name     string
		payload  string
		expected uint32
	}{
		{
			name:     "empty string",
			payload:  "",
			expected: 0,
		}, {
			name:     "not valid ipv4 format",
			payload:  "0.0.0",
			expected: 0,
		}, {
			name:     "this network",
			payload:  "0.0.0.0",
			expected: 0,
		}, {
			name:     "parse localhost",
			payload:  "127.0.0.1",
			expected: 2130706433,
		}, {
			name:     "try to overflow",
			payload:  "255.255.255.255",
			expected: 4294967295,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IPstringV4ToInt(tt.payload)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestIPIntV4ToString(t *testing.T) {
	// https://www.browserling.com/tools/ip-to-dec
	tests := []struct {
		name     string
		payload  uint32
		expected string
	}{
		{
			name:     "zero value",
			payload:  0,
			expected: "0.0.0.0",
		}, {
			name:     "localhost",
			payload:  2130706433,
			expected: "127.0.0.1",
		}, {
			name:     "try to overflow",
			payload:  4294967295,
			expected: "255.255.255.255",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IPIntV4ToString(tt.payload)
			require.Equal(t, tt.expected, result)
		})
	}
}
