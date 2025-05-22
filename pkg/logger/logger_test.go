package logger

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"os"
	"strings"
	"testing"
)

func TestFormatDuration(t *testing.T) {
	testCases := []struct {
		name     string
		seconds  float64
		expected string
	}{
		{
			name:     "Negative time",
			seconds:  -10.0,
			expected: "--:--:--",
		},
		{
			name:     "NaN (Not a Number)",
			seconds:  math.NaN(),
			expected: "--:--:--",
		},
		{
			name:     "Zero time",
			seconds:  0.0,
			expected: "00:00:00",
		},
		{
			name:     "Seconds only",
			seconds:  45.0,
			expected: "00:00:45",
		},
		{
			name:     "Minutes and seconds",
			seconds:  125.0, // 2 minutes 5 seconds
			expected: "00:02:05",
		},
		{
			name:     "Hours, minutes and seconds",
			seconds:  3725.0, // 1 hour 2 minutes 5 seconds
			expected: "01:02:05",
		},
		{
			name:     "Large number of hours",
			seconds:  90000.0, // 25 hours
			expected: "25:00:00",
		},
		{
			name:     "Infinity",
			seconds:  math.Inf(1),
			expected: "--:--:--",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := FormatDuration(tc.seconds)
			if result != tc.expected {
				t.Errorf("FormatDuration(%.2f) = %s; want %s", 
					tc.seconds, result, tc.expected)
			}
		})
	}
}

func TestInitialize(t *testing.T) {
	// First ensure that spinner is not initialized
	Spinner = nil
	
	// Initialize the logger
	Initialize()
	
	// Verify that spinner was created
	if Spinner == nil {
		t.Error("Initialize() did not create spinner, Spinner = nil")
	}
}

func TestLogDownload(t *testing.T) {
	// Skip this test as it's hard to reliably test spinner output
	t.Skip("Skipping TestLogDownload because it depends on spinner output which is hard to capture")
}

func TestLogComplete(t *testing.T) {
	// Ensure spinner is initialized
	Initialize()
	Spinner.Start()
	
	// Capture stdout to verify output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	
	// Call the function
	LogComplete("test-package", 100)
	
	// Restore stdout
	w.Close()
	os.Stdout = oldStdout
	
	// Read captured output
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()
	
	// Verify output contains important information
	expectedOutput := "Successfully completed 100 downloads for package test-package"
	if !strings.Contains(output, expectedOutput) {
		t.Errorf("LogComplete() output doesn't contain expected text. Got: %s", output)
	}
}

func TestLogError(t *testing.T) {
	// Ensure spinner is initialized
	Initialize()
	Spinner.Start()
	
	// Capture stdout to verify output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	
	// Call the function
	testError := fmt.Errorf("test error")
	LogError(testError)
	
	// Restore stdout
	w.Close()
	os.Stdout = oldStdout
	
	// Read captured output
	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()
	
	// Verify output contains important information
	if !strings.Contains(output, "Error: test error") {
		t.Errorf("LogError() output doesn't contain expected text. Got: %s", output)
	}
} 