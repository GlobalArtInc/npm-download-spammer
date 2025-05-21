package models

// Maintainer represents an NPM package maintainer
type Maintainer struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
}

// NpmPackage represents information about an NPM package
type NpmPackage struct {
	Name        string            `json:"name"`
	Scope       string            `json:"scope"`
	Version     string            `json:"version"`
	Description string            `json:"description"`
	Keywords    []string          `json:"keywords"`
	Date        string            `json:"date"`
	Links       map[string]string `json:"links"`
	Publisher   map[string]string `json:"publisher"`
	Maintainers []Maintainer      `json:"maintainers"`
}

// NpmObject represents an object in the NPM API response
type NpmObject struct {
	Package NpmPackage `json:"package"`
}

// NpmjsResponse represents a response from the NPM API
type NpmjsResponse struct {
	Objects []NpmObject `json:"objects"`
} 