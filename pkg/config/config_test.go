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

	packageNames := cfg.GetPackageNames()
	if len(packageNames) != 1 || packageNames[0] != "env-test-package" {
		t.Errorf("LoadConfig().GetPackageNames() = %v; want [env-test-package]", packageNames)
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

	packageNames := cfg.GetPackageNames()
	if len(packageNames) != 1 || packageNames[0] != "test-package" {
		t.Errorf("LoadConfig().GetPackageNames() = %v; want [test-package]", packageNames)
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

func TestSetPackageNames(t *testing.T) {
	cfg := DefaultConfig()

	// Test single package
	cfg.SetPackageNames("test-package")
	packages := cfg.GetPackageNames()
	if len(packages) != 1 || packages[0] != "test-package" {
		t.Errorf("SetPackageNames(single) failed: got %v, want [test-package]", packages)
	}

	// Test multiple packages
	cfg.SetPackageNames("package1, package2,  package3  ")
	packages = cfg.GetPackageNames()
	expected := []string{"package1", "package2", "package3"}
	if len(packages) != 3 {
		t.Errorf("SetPackageNames(multiple) length failed: got %d, want 3", len(packages))
	}
	for i, pkg := range expected {
		if packages[i] != pkg {
			t.Errorf("SetPackageNames(multiple) failed: got %v, want %v", packages, expected)
		}
	}

	// Test empty input
	cfg.SetPackageNames("")
	packages = cfg.GetPackageNames()
	if len(packages) != 0 {
		t.Errorf("SetPackageNames(empty) failed: got %v, want []", packages)
	}

	// Test with empty elements
	cfg.SetPackageNames("package1, , package3")
	packages = cfg.GetPackageNames()
	expected = []string{"package1", "package3"}
	if len(packages) != 2 {
		t.Errorf("SetPackageNames(with empty) length failed: got %d, want 2", len(packages))
	}
	for i, pkg := range expected {
		if packages[i] != pkg {
			t.Errorf("SetPackageNames(with empty) failed: got %v, want %v", packages, expected)
		}
	}
}

func TestGetPackageNames(t *testing.T) {
	cfg := DefaultConfig()

	// Test empty config
	packages := cfg.GetPackageNames()
	if len(packages) != 0 {
		t.Errorf("GetPackageNames(empty) failed: got %v, want []", packages)
	}

	// Test legacy PackageName
	cfg.PackageName = "legacy-package"
	packages = cfg.GetPackageNames()
	if len(packages) != 1 || packages[0] != "legacy-package" {
		t.Errorf("GetPackageNames(legacy) failed: got %v, want [legacy-package]", packages)
	}

	// Test new PackageNames (should override legacy)
	cfg.PackageNames = []string{"new-package1", "new-package2"}
	packages = cfg.GetPackageNames()
	expected := []string{"new-package1", "new-package2"}
	if len(packages) != 2 {
		t.Errorf("GetPackageNames(new) length failed: got %d, want 2", len(packages))
	}
	for i, pkg := range expected {
		if packages[i] != pkg {
			t.Errorf("GetPackageNames(new) failed: got %v, want %v", packages, expected)
		}
	}
}

func TestLoadConfigMultiplePackages(t *testing.T) {
	// Save current environment variable value
	originalValue := os.Getenv("NPM_PACKAGE_NAME")
	defer func() {
		if originalValue == "" {
			os.Unsetenv("NPM_PACKAGE_NAME")
		} else {
			os.Setenv("NPM_PACKAGE_NAME", originalValue)
		}
	}()

	// Set environment variable with multiple packages
	os.Setenv("NPM_PACKAGE_NAME", "package1, package2, package3")

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() returned error: %v", err)
	}

	packages := cfg.GetPackageNames()
	expected := []string{"package1", "package2", "package3"}
	if len(packages) != 3 {
		t.Errorf("LoadConfig(multiple packages) length failed: got %d, want 3", len(packages))
	}
	for i, pkg := range expected {
		if packages[i] != pkg {
			t.Errorf("LoadConfig(multiple packages) failed: got %v, want %v", packages, expected)
		}
	}
}
