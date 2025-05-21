package utils

import (
	"net/url"
	"strings"
)

// StripOrganisationFromPackageName removes the organization from package name
// For example, @scope/package-name -> package-name
func StripOrganisationFromPackageName(packageName string) string {
	parts := strings.Split(packageName, "/")
	if len(parts) <= 1 {
		return packageName
	}
	return parts[len(parts)-1]
}

// GetEncodedPackageName encodes package name for use in URLs
func GetEncodedPackageName(packageName string) string {
	return url.QueryEscape(packageName)
} 