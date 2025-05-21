package models

import (
	"time"
)

// Stats stores download statistics
type Stats struct {
	// Successful downloads
	SuccessfulDownloads int

	// Failed downloads
	FailedDownloads int

	// Start time of downloads
	StartTime time.Time

	// Total number of downloads to perform
	TotalDownloads int
}

// NewStats creates a new statistics instance
func NewStats(totalDownloads int) *Stats {
	return &Stats{
		SuccessfulDownloads: 0,
		FailedDownloads:     0,
		StartTime:           time.Now(),
		TotalDownloads:      totalDownloads,
	}
}

// GetDownloadSpeed returns download rate in downloads/second
func (s *Stats) GetDownloadSpeed() float64 {
	elapsed := time.Since(s.StartTime).Seconds()
	if elapsed <= 0 || s.SuccessfulDownloads <= 0 {
		return 0
	}
	
	return float64(s.SuccessfulDownloads) / elapsed
}

// GetTimeRemaining returns estimated remaining time in seconds
func (s *Stats) GetTimeRemaining() (float64, bool) {
	speed := s.GetDownloadSpeed()
	if speed <= 0 {
		return 0, false
	}
	
	downloadsRemaining := s.TotalDownloads - s.SuccessfulDownloads
	if downloadsRemaining <= 0 {
		return 0, true
	}
	
	return float64(downloadsRemaining) / speed, true
}

// GetProgress returns completion percentage
func (s *Stats) GetProgress() float64 {
	if s.TotalDownloads <= 0 {
		return 0
	}
	return float64(s.SuccessfulDownloads) / float64(s.TotalDownloads) * 100
} 