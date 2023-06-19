package base

import "os"

func FileExists(filename string) bool {
	// Use os.Stat to get file information
	_, err := os.Stat(filename)

	// Check if the error indicates that the file doesn't exist
	if os.IsNotExist(err) {
		return false
	}

	// Return true if no error occurred (file exists or other error occurred)
	return err == nil
}
