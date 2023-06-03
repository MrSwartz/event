package utils

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSplitOsAndVersion(t *testing.T) {
	tests := []struct {
		name          string
		payload       string
		os            string
		version       string
		expectedError error
	}{
		{
			name:          "payload from request example",
			payload:       "IOS 13.5.1",
			os:            "ios",
			version:       "13.5.1",
			expectedError: nil,
		}, {
			name:          "payload with invalid version format",
			payload:       "IOS 13.5.1.1",
			os:            "",
			version:       "",
			expectedError: errors.New(`invalid os version format: "IOS 13.5.1.1"`),
		}, {
			name:          "payload with invalid version format",
			payload:       "IOS 13.5",
			os:            "",
			version:       "",
			expectedError: errors.New(`can't extract data from string: "IOS 13.5"`),
		}, {
			name:          "empty payload",
			payload:       "",
			os:            "",
			version:       "",
			expectedError: errors.New(`can't extract data from string: ""`),
		}, {
			name:          "lowercase os type with valid version",
			payload:       "ios 13.0.1",
			os:            "ios",
			version:       "13.0.1",
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os, ver, err := SplitOsAndVersion(tt.payload)

			if tt.expectedError != nil {
				require.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				require.Nil(t, err)
			}

			require.Equal(t, tt.os, os)
			require.Equal(t, tt.version, ver)
		})
	}
}
