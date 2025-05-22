package config

import (
	"os"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.PackageName != "" {
		t.Errorf("DefaultConfig().PackageName = %s; want empty string", cfg.PackageName)
	}

	if cfg.NumDownloads != 1000 {
		t.Errorf("DefaultConfig().NumDownloads = %d; want 1000", cfg.NumDownloads)
	}

	if cfg.MaxConcurrentDownloads != 300 {
		t.Errorf("DefaultConfig().MaxConcurrentDownloads = %d; want 300", cfg.MaxConcurrentDownloads)
	}

	if cfg.DownloadTimeout != 3000 {
		t.Errorf("DefaultConfig().DownloadTimeout = %d; want 3000", cfg.DownloadTimeout)
	}
}

func TestLoadConfigFromFile(t *testing.T) {
	// Create a temporary config file
	configContent := `{
		"packageName": "file-test-package",
		"numDownloads": 2500,
		"maxConcurrentDownloads": 400,
		"downloadTimeout": 4000
	}`
	
	err := os.WriteFile("npm-downloads-increaser.json", []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}
	
	// Make sure we clean up afterwards
	defer os.Remove("npm-downloads-increaser.json")
	
	// Save and reset environment variables to ensure they don't interfere
	origPackageName := os.Getenv("NPM_PACKAGE_NAME")
	origNumDownloads := os.Getenv("NPM_NUM_DOWNLOADS")
	origMaxConcurrent := os.Getenv("NPM_MAX_CONCURRENT_DOWNLOAD")
	origTimeout := os.Getenv("NPM_DOWNLOAD_TIMEOUT")
	
	os.Unsetenv("NPM_PACKAGE_NAME")
	os.Unsetenv("NPM_NUM_DOWNLOADS")
	os.Unsetenv("NPM_MAX_CONCURRENT_DOWNLOAD")
	os.Unsetenv("NPM_DOWNLOAD_TIMEOUT")
	
	defer func() {
		os.Setenv("NPM_PACKAGE_NAME", origPackageName)
		os.Setenv("NPM_NUM_DOWNLOADS", origNumDownloads)
		os.Setenv("NPM_MAX_CONCURRENT_DOWNLOAD", origMaxConcurrent)
		os.Setenv("NPM_DOWNLOAD_TIMEOUT", origTimeout)
	}()
	
	// Load configuration
	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() returned error: %v", err)
	}
	
	// Verify configuration was loaded from file
	if cfg.PackageName != "file-test-package" {
		t.Errorf("LoadConfig().PackageName = %s; want file-test-package", cfg.PackageName)
	}
	
	if cfg.NumDownloads != 2500 {
		t.Errorf("LoadConfig().NumDownloads = %d; want 2500", cfg.NumDownloads)
	}
	
	if cfg.MaxConcurrentDownloads != 400 {
		t.Errorf("LoadConfig().MaxConcurrentDownloads = %d; want 400", cfg.MaxConcurrentDownloads)
	}
	
	if cfg.DownloadTimeout != 4000 {
		t.Errorf("LoadConfig().DownloadTimeout = %d; want 4000", cfg.DownloadTimeout)
	}
	
	// Test environment variables override file settings
	os.Setenv("NPM_PACKAGE_NAME", "env-test-package")
	
	cfg, err = LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() returned error: %v", err)
	}
	
	if cfg.PackageName != "env-test-package" {
		t.Errorf("LoadConfig().PackageName = %s; want env-test-package", cfg.PackageName)
	}
}

func TestLoadConfigFromEnv(t *testing.T) {
	// First save current environment variable values
	origPackageName := os.Getenv("NPM_PACKAGE_NAME")
	origNumDownloads := os.Getenv("NPM_NUM_DOWNLOADS")
	origMaxConcurrent := os.Getenv("NPM_MAX_CONCURRENT_DOWNLOAD")
	origTimeout := os.Getenv("NPM_DOWNLOAD_TIMEOUT")

	// Restore environment variables at the end of the test
	defer func() {
		os.Setenv("NPM_PACKAGE_NAME", origPackageName)
		os.Setenv("NPM_NUM_DOWNLOADS", origNumDownloads)
		os.Setenv("NPM_MAX_CONCURRENT_DOWNLOAD", origMaxConcurrent)
		os.Setenv("NPM_DOWNLOAD_TIMEOUT", origTimeout)
	}()

	// Set test environment variables
	os.Setenv("NPM_PACKAGE_NAME", "test-package")
	os.Setenv("NPM_NUM_DOWNLOADS", "2000")
	os.Setenv("NPM_MAX_CONCURRENT_DOWNLOAD", "500")
	os.Setenv("NPM_DOWNLOAD_TIMEOUT", "5000")

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() returned error: %v", err)
	}

	if cfg.PackageName != "test-package" {
		t.Errorf("LoadConfig().PackageName = %s; want test-package", cfg.PackageName)
	}

	if cfg.NumDownloads != 2000 {
		t.Errorf("LoadConfig().NumDownloads = %d; want 2000", cfg.NumDownloads)
	}

	if cfg.MaxConcurrentDownloads != 500 {
		t.Errorf("LoadConfig().MaxConcurrentDownloads = %d; want 500", cfg.MaxConcurrentDownloads)
	}

	if cfg.DownloadTimeout != 5000 {
		t.Errorf("LoadConfig().DownloadTimeout = %d; want 5000", cfg.DownloadTimeout)
	}
}

func TestLoadConfigInvalidEnv(t *testing.T) {
	// First save current environment variable values
	origNumDownloads := os.Getenv("NPM_NUM_DOWNLOADS")
	origMaxConcurrent := os.Getenv("NPM_MAX_CONCURRENT_DOWNLOAD")
	origTimeout := os.Getenv("NPM_DOWNLOAD_TIMEOUT")

	// Restore environment variables at the end of the test
	defer func() {
		os.Setenv("NPM_NUM_DOWNLOADS", origNumDownloads)
		os.Setenv("NPM_MAX_CONCURRENT_DOWNLOAD", origMaxConcurrent)
		os.Setenv("NPM_DOWNLOAD_TIMEOUT", origTimeout)
	}()

	// Set invalid environment variables
	os.Setenv("NPM_NUM_DOWNLOADS", "not a number")
	os.Setenv("NPM_MAX_CONCURRENT_DOWNLOAD", "not a number")
	os.Setenv("NPM_DOWNLOAD_TIMEOUT", "not a number")

	// Should use default values for invalid inputs
	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() returned error: %v", err)
	}

	if cfg.NumDownloads != 1000 {
		t.Errorf("LoadConfig() with invalid env vars - NumDownloads = %d; want 1000", cfg.NumDownloads)
	}

	if cfg.MaxConcurrentDownloads != 300 {
		t.Errorf("LoadConfig() with invalid env vars - MaxConcurrentDownloads = %d; want 300", cfg.MaxConcurrentDownloads)
	}

	if cfg.DownloadTimeout != 3000 {
		t.Errorf("LoadConfig() with invalid env vars - DownloadTimeout = %d; want 3000", cfg.DownloadTimeout)
	}
} 