package utils

import (
	"path/filepath"
	"strings"
)

func LocalPathToFileURL(localPath string) (string, error) {

	// Get absolute path
	absPath, err := filepath.Abs(localPath)
	if err != nil {
		return "", err
	}

	// Convert backslashes to forward slashes (for Windows compatibility)
	urlPath := filepath.ToSlash(absPath)

	// Ensure path starts with /
	if !strings.HasPrefix(urlPath, "/") {
		urlPath = "/" + urlPath
	}

	// Construct the file URL
	return "file://" + urlPath, nil
}
