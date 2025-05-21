package logger

import (
	"fmt"
	"math"
	"time"

	"github.com/briandowns/spinner"
	"npm-download-spammer/pkg/models"
)

var (
	// Spinner for displaying progress
	Spinner *spinner.Spinner
)

// Initialize initializes the logger
func Initialize() {
	Spinner = spinner.New(spinner.CharSets[14], 100*time.Millisecond)
}

// FormatDuration formats time in HH:MM:SS format
func FormatDuration(seconds float64) string {
	if seconds < 0 || math.IsNaN(seconds) || math.IsInf(seconds, 0) {
		return "--:--:--"
	}

	hours := int(seconds) / 3600
	minutes := (int(seconds) % 3600) / 60
	secs := int(seconds) % 60

	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, secs)
}

// LogDownload logs the current download status
func LogDownload(stats *models.Stats) {
	if Spinner == nil {
		Initialize()
	}

	if !Spinner.Active() {
		Spinner.Start()
	}

	downloadSpeed := stats.GetDownloadSpeed()
	timeRemaining, valid := stats.GetTimeRemaining()
	
	var timeStr string
	if valid {
		timeStr = FormatDuration(timeRemaining)
	} else {
		timeStr = "--:--:--"
	}

	text := fmt.Sprintf("\nDownload count:           %d/%d\n", 
		stats.SuccessfulDownloads, stats.TotalDownloads)
	text += fmt.Sprintf("Download speed:           %.2f dl/s\n", downloadSpeed)
	text += fmt.Sprintf("Estimated time remaining: %s\n", timeStr)

	Spinner.Suffix = text
}

// LogComplete logs successful completion
func LogComplete(packageName string, downloads int) {
	if Spinner != nil && Spinner.Active() {
		Spinner.Stop()
	}
	fmt.Printf("✅ Successfully completed %d downloads for package %s\n", downloads, packageName)
}

// LogError logs an error
func LogError(err error) {
	if Spinner != nil && Spinner.Active() {
		Spinner.Stop()
	}
	fmt.Printf("❌ Error: %v\n", err)
} 