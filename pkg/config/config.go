package config

import (
	"encoding/json"
	"os"
	"strconv"
	"strings"
)

// Config contains settings for program operation
type Config struct {
	// NPM package name to increase download counter (legacy)
	PackageName string `json:"packageName,omitempty"`

	// NPM package names to increase download counter (new field)
	PackageNames []string `json:"packageNames,omitempty"`

	// Number of downloads to add to the package
	NumDownloads int `json:"numDownloads"`

	// Number of concurrent downloads
	MaxConcurrentDownloads int `json:"maxConcurrentDownloads"`

	// Maximum download timeout (in milliseconds)
	DownloadTimeout int `json:"downloadTimeout"`
}

// GetPackageNames returns the list of package names from both old and new fields
func (c *Config) GetPackageNames() []string {
	if len(c.PackageNames) > 0 {
		return c.PackageNames
	}
	if c.PackageName != "" {
		return []string{c.PackageName}
	}
	return []string{}
}

// SetPackageNames sets package names from a comma-separated string
func (c *Config) SetPackageNames(input string) {
	c.PackageNames = []string{}
	c.PackageName = ""
	
	if input == "" {
		return
	}
	
	packages := strings.Split(input, ",")
	c.PackageNames = make([]string, 0, len(packages))
	
	for _, pkg := range packages {
		trimmed := strings.TrimSpace(pkg)
		if trimmed != "" {
			c.PackageNames = append(c.PackageNames, trimmed)
		}
	}
}

// DefaultConfig returns the default configuration
func DefaultConfig() Config {
	return Config{
		PackageName:            "",
		PackageNames:           []string{},
		NumDownloads:           1000,
		MaxConcurrentDownloads: 300,
		DownloadTimeout:        3000,
	}
}

// LoadConfig loads configuration from file and environment variables
func LoadConfig() (Config, error) {
	config := DefaultConfig()

	// Try to load from configuration file
	if _, err := os.Stat("npm-downloads-increaser.json"); err == nil {
		file, err := os.Open("npm-downloads-increaser.json")
		if err == nil {
			defer file.Close()
			decoder := json.NewDecoder(file)
			if err := decoder.Decode(&config); err != nil {
				return config, err
			}
		}
	}

	// Override from environment variables if they are set
	if packageName := os.Getenv("NPM_PACKAGE_NAME"); packageName != "" {
		config.SetPackageNames(packageName)
	}

	if numDownloads := os.Getenv("NPM_NUM_DOWNLOADS"); numDownloads != "" {
		if val, err := strconv.Atoi(numDownloads); err == nil {
			config.NumDownloads = val
		}
	}

	if maxConcurrent := os.Getenv("NPM_MAX_CONCURRENT_DOWNLOAD"); maxConcurrent != "" {
		if val, err := strconv.Atoi(maxConcurrent); err == nil {
			config.MaxConcurrentDownloads = val
		}
	}

	if timeout := os.Getenv("NPM_DOWNLOAD_TIMEOUT"); timeout != "" {
		if val, err := strconv.Atoi(timeout); err == nil {
			config.DownloadTimeout = val
		}
	}

	return config, nil
}
