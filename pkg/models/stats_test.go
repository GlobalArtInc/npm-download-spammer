package models

import (
	"testing"
	"time"
)

func TestNewStats(t *testing.T) {
	totalDownloads := 1000
	stats := NewStats(totalDownloads)

	if stats.TotalDownloads != totalDownloads {
		t.Errorf("NewStats(%d).TotalDownloads = %d; want %d", 
			totalDownloads, stats.TotalDownloads, totalDownloads)
	}

	if stats.SuccessfulDownloads != 0 {
		t.Errorf("NewStats(%d).SuccessfulDownloads = %d; want 0", 
			totalDownloads, stats.SuccessfulDownloads)
	}

	if stats.FailedDownloads != 0 {
		t.Errorf("NewStats(%d).FailedDownloads = %d; want 0", 
			totalDownloads, stats.FailedDownloads)
	}

	// Check that StartTime is close to current time
	timeDiff := time.Since(stats.StartTime)
	if timeDiff > time.Second {
		t.Errorf("Start time differs from current time by more than 1 second: %v", timeDiff)
	}
}

func TestGetDownloadSpeed(t *testing.T) {
	testCases := []struct {
		name          string
		stats         *Stats
		expectedSpeed float64
	}{
		{
			name: "Zero speed with zero downloads",
			stats: &Stats{
				SuccessfulDownloads: 0,
				StartTime:           time.Now().Add(-10 * time.Second),
				TotalDownloads:      1000,
			},
			expectedSpeed: 0,
		},
		{
			name: "Speed calculated correctly",
			stats: &Stats{
				SuccessfulDownloads: 50,
				StartTime:           time.Now().Add(-10 * time.Second),
				TotalDownloads:      1000,
			},
			expectedSpeed: 5.0, // 50 downloads in 10 seconds = 5 downloads/sec
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			speed := tc.stats.GetDownloadSpeed()
			// Use approximate comparison for floating point numbers
			if speed < tc.expectedSpeed*0.9 || speed > tc.expectedSpeed*1.1 {
				t.Errorf("GetDownloadSpeed() = %.2f; want approximately %.2f", 
					speed, tc.expectedSpeed)
			}
		})
	}
}

func TestGetTimeRemaining(t *testing.T) {
	testCases := []struct {
		name                 string
		stats                *Stats
		expectedTimeRemaining float64
		expectedValid        bool
	}{
		{
			name: "Invalid time with zero speed",
			stats: &Stats{
				SuccessfulDownloads: 0,
				StartTime:           time.Now().Add(-10 * time.Second),
				TotalDownloads:      1000,
			},
			expectedTimeRemaining: 0,
			expectedValid:        false,
		},
		{
			name: "Invalid time with completed downloads",
			stats: &Stats{
				SuccessfulDownloads: 1000,
				StartTime:           time.Now().Add(-10 * time.Second),
				TotalDownloads:      1000,
			},
			expectedTimeRemaining: 0,
			expectedValid:        true,
		},
		{
			name: "Calculate remaining time",
			stats: &Stats{
				SuccessfulDownloads: 500,
				StartTime:           time.Now().Add(-10 * time.Second),
				TotalDownloads:      1000,
			},
			expectedTimeRemaining: 10.0, // 500 more downloads at 50/sec = 10 sec
			expectedValid:        true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			timeRemaining, valid := tc.stats.GetTimeRemaining()
			
			if valid != tc.expectedValid {
				t.Errorf("GetTimeRemaining() valid = %v; want %v", 
					valid, tc.expectedValid)
			}
			
			if valid && tc.expectedValid {
				// Use approximate comparison for floating point numbers
				if timeRemaining < tc.expectedTimeRemaining*0.9 || 
				   timeRemaining > tc.expectedTimeRemaining*1.1 {
					t.Errorf("GetTimeRemaining() = %.2f; want approximately %.2f", 
						timeRemaining, tc.expectedTimeRemaining)
				}
			}
		})
	}
}

func TestGetProgress(t *testing.T) {
	testCases := []struct {
		name            string
		stats           *Stats
		expectedProgress float64
	}{
		{
			name: "Zero progress with zero TotalDownloads",
			stats: &Stats{
				SuccessfulDownloads: 50,
				TotalDownloads:      0,
			},
			expectedProgress: 0,
		},
		{
			name: "0% progress",
			stats: &Stats{
				SuccessfulDownloads: 0,
				TotalDownloads:      1000,
			},
			expectedProgress: 0,
		},
		{
			name: "50% progress",
			stats: &Stats{
				SuccessfulDownloads: 500,
				TotalDownloads:      1000,
			},
			expectedProgress: 50,
		},
		{
			name: "100% progress",
			stats: &Stats{
				SuccessfulDownloads: 1000,
				TotalDownloads:      1000,
			},
			expectedProgress: 100,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			progress := tc.stats.GetProgress()
			
			if progress != tc.expectedProgress {
				t.Errorf("GetProgress() = %.2f; want %.2f", 
					progress, tc.expectedProgress)
			}
		})
	}
} 