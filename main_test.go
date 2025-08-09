package main

import (
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"npm-download-spammer/pkg/config"
)

// Test patch for http.DefaultClient
func patchHTTPClient(mockResp *http.Response, mockErr error) func() {
	// Save the original client
	originalClient := http.DefaultClient
	
	// Create a mock transport
	mockTransport := &mockTransport{
		mockResponse: mockResp,
		mockError:    mockErr,
	}
	
	// Create a mock client
	mockClient := &http.Client{
		Transport: mockTransport,
	}
	
	// Replace the default client with our mock
	http.DefaultClient = mockClient
	
	// Return function to restore the original state
	return func() {
		http.DefaultClient = originalClient
	}
}

// Mock implementation of http.RoundTripper for testing
type mockTransport struct {
	mockResponse *http.Response
	mockError    error
}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.mockResponse, m.mockError
}

func TestMainWithEnvVars(t *testing.T) {
	// This test only checks initialization with environment variables,
	// without actually executing downloads

	// Save current environment variable values
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
	os.Setenv("NPM_NUM_DOWNLOADS", "1")
	os.Setenv("NPM_MAX_CONCURRENT_DOWNLOAD", "1")
	os.Setenv("NPM_DOWNLOAD_TIMEOUT", "100")

	// Create mock HTTP responses
	npmResp := &http.Response{
		StatusCode: http.StatusOK,
		Body: io.NopCloser(strings.NewReader(`{
			"objects": [
				{
					"package": {
						"name": "test-package",
						"version": "1.0.0"
					}
				}
			]
		}`)),
	}
	
	// Patch HTTP client
	restore := patchHTTPClient(npmResp, nil)
	defer restore()
	
	// Since we can't safely mock time.Sleep and we can't directly test main(),
	// we'll just check that the environment variables are read correctly
	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() returned error: %v", err)
	}
	
	// Verify configuration values
	packageNames := cfg.GetPackageNames()
	if len(packageNames) != 1 || packageNames[0] != "test-package" {
		t.Errorf("cfg.GetPackageNames() = %v; want [test-package]", packageNames)
	}
	
	if cfg.NumDownloads != 1 {
		t.Errorf("cfg.NumDownloads = %d; want 1", cfg.NumDownloads)
	}
	
	if cfg.MaxConcurrentDownloads != 1 {
		t.Errorf("cfg.MaxConcurrentDownloads = %d; want 1", cfg.MaxConcurrentDownloads)
	}
	
	if cfg.DownloadTimeout != 100 {
		t.Errorf("cfg.DownloadTimeout = %d; want 100", cfg.DownloadTimeout)
	}
} 