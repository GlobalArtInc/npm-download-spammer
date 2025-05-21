package config

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/manifoldco/promptui"
)

// GetConfigFromCLI requests configuration through interactive CLI
func GetConfigFromCLI() (Config, error) {
	config := DefaultConfig()

	// Request package name
	packageName, err := promptPackageName()
	if err != nil {
		return config, err
	}
	config.PackageName = packageName

	// Request number of downloads
	numDownloads, err := promptNumericValue("Number of downloads", 1000)
	if err != nil {
		return config, err
	}
	config.NumDownloads = numDownloads

	// Request maximum concurrent downloads
	maxConcurrent, err := promptNumericValue("Number of concurrent downloads", 300)
	if err != nil {
		return config, err
	}
	config.MaxConcurrentDownloads = maxConcurrent

	// Request download timeout
	timeout, err := promptNumericValue("Download timeout (in ms)", 3000)
	if err != nil {
		return config, err
	}
	config.DownloadTimeout = timeout

	return config, nil
}

// promptPackageName requests package name with validation
func promptPackageName() (string, error) {
	validate := func(input string) error {
		if input == "" {
			return errors.New("package name cannot be empty")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Package name",
		Validate: validate,
	}

	result, err := prompt.Run()
	if err != nil {
		return "", fmt.Errorf("prompt error: %v", err)
	}

	return result, nil
}

// promptNumericValue requests numeric value with validation
func promptNumericValue(label string, defaultValue int) (int, error) {
	validate := func(input string) error {
		_, err := strconv.ParseInt(input, 10, 64)
		if err != nil {
			return errors.New("please enter a valid number")
		}
		
		num, _ := strconv.Atoi(input)
		if num <= 0 {
			return errors.New("value must be greater than 0")
		}
		
		return nil
	}

	prompt := promptui.Prompt{
		Label:    label,
		Default:  strconv.Itoa(defaultValue),
		Validate: validate,
	}

	result, err := prompt.Run()
	if err != nil {
		return 0, fmt.Errorf("prompt error: %v", err)
	}

	value, _ := strconv.Atoi(result)
	return value, nil
} 