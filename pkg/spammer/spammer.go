package spammer

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"npm-download-spammer/pkg/config"
	"npm-download-spammer/pkg/logger"
	"npm-download-spammer/pkg/models"
	"npm-download-spammer/pkg/utils"
)

// QueryNpms requests package information from the NPM API
func QueryNpms(packageName string) (*models.NpmjsResponse, error) {
	url := fmt.Sprintf("https://registry.npmjs.com/-/v1/search?text=%s&size=1", 
		utils.GetEncodedPackageName(packageName))
	
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("HTTP request error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid response status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("response reading error: %w", err)
	}

	var npmResponse models.NpmjsResponse
	if err := json.Unmarshal(body, &npmResponse); err != nil {
		return nil, fmt.Errorf("JSON parsing error: %w", err)
	}

	if len(npmResponse.Objects) == 0 {
		return nil, fmt.Errorf("package not found: %s", packageName)
	}

	return &npmResponse, nil
}

// DownloadPackage downloads a package to increase the download counter
func DownloadPackage(packageName, version string, stats *models.Stats, timeout int) error {
	unscopedPackageName := utils.StripOrganisationFromPackageName(packageName)
	url := fmt.Sprintf("https://registry.yarnpkg.com/%s/-/%s-%s.tgz", 
		packageName, unscopedPackageName, version)

	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Millisecond,
	}

	resp, err := client.Get(url)
	if err != nil {
		stats.FailedDownloads++
		return err
	}
	defer resp.Body.Close()

	// Just read the response body and discard it
	_, err = io.Copy(io.Discard, resp.Body)
	if err != nil {
		stats.FailedDownloads++
		return err
	}

	stats.SuccessfulDownloads++
	return nil
}

// SpamDownloads initiates parallel downloads of the package
func SpamDownloads(cfg config.Config, version string, stats *models.Stats, wg *sync.WaitGroup) {
	defer wg.Done()

	var downloadWg sync.WaitGroup
	downloadWg.Add(cfg.MaxConcurrentDownloads)

	for i := 0; i < cfg.MaxConcurrentDownloads; i++ {
		go func() {
			defer downloadWg.Done()
			_ = DownloadPackage(cfg.PackageName, version, stats, cfg.DownloadTimeout)
		}()
	}

	downloadWg.Wait()

	// If we need to download more, start a new batch of downloads
	if stats.SuccessfulDownloads < cfg.NumDownloads {
		wg.Add(1)
		go SpamDownloads(cfg, version, stats, wg)
	}
}

// Run starts the download counter increment process
func Run(cfg config.Config) error {
	// Initialize the logger
	logger.Initialize()

	// Get package information
	npmResponse, err := QueryNpms(cfg.PackageName)
	if err != nil {
		return err
	}

	version := npmResponse.Objects[0].Package.Version
	stats := models.NewStats(cfg.NumDownloads)

	// Start logging in a separate goroutine
	ticker := time.NewTicker(1 * time.Second)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-ticker.C:
				logger.LogDownload(stats)
			case <-done:
				ticker.Stop()
				return
			}
		}
	}()

	// Start downloads
	var wg sync.WaitGroup
	wg.Add(1)
	go SpamDownloads(cfg, version, stats, &wg)

	// Wait for all downloads to complete
	wg.Wait()
	done <- true

	logger.LogComplete(cfg.PackageName, stats.SuccessfulDownloads)
	return nil
} 