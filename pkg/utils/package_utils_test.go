package utils

import (
	"testing"
)

func TestStripOrganisationFromPackageName(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Package without organization",
			input:    "package-name",
			expected: "package-name",
		},
		{
			name:     "Package with organization",
			input:    "@scope/package-name",
			expected: "package-name",
		},
		{
			name:     "Package with multiple slashes",
			input:    "@scope/subgroup/package-name",
			expected: "package-name",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := StripOrganisationFromPackageName(tc.input)
			if result != tc.expected {
				t.Errorf("StripOrganisationFromPackageName(%s) = %s; want %s", 
					tc.input, result, tc.expected)
			}
		})
	}
}

func TestGetEncodedPackageName(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple package name",
			input:    "simple-package",
			expected: "simple-package",
		},
		{
			name:     "Package name with organization",
			input:    "@scope/package-name",
			expected: "%40scope%2Fpackage-name",
		},
		{
			name:     "Name with special characters",
			input:    "package name with spaces",
			expected: "package+name+with+spaces",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := GetEncodedPackageName(tc.input)
			if result != tc.expected {
				t.Errorf("GetEncodedPackageName(%s) = %s; want %s", 
					tc.input, result, tc.expected)
			}
		})
	}
} 