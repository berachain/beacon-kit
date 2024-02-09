package config

import "os"

// HasDir checks if a directory exists.
func HasDir(dirPath string) (bool, error) {
	stat, err := os.Stat(dirPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil // Directory does not exist
		}
		return false, err // An error occurred (e.g., permissions)
	}
	return stat.IsDir(), nil // Return true if the path is a directory
}

// MkdirAll creates a directory along with any necessary parents.
func MkdirAll(dirPath string) error {
	return os.MkdirAll(dirPath, os.ModePerm) // Creates the directory with default permissions
}

// CreateFile creates or truncates a file at the specified path.
func CreateFile(filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err // Return the error if file creation failed
	}
	return file.Close() // Close the file and return the result of Close()
}

// Exists checks if a file exists.
func Exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true // File exists
	}
	if os.IsNotExist(err) {
		return false // File does not exist
	}
	return false // An error occurred (e.g., permissions)
}

// WriteFile writes data to a file at the specified path.
func WriteFile(filePath string, data []byte) error {
	if ok := Exists(filePath); !ok {
		err := CreateFile(filePath)
		if err != nil {
			return err // Return the error if creating file failed
		}
	}

	err := os.WriteFile(filePath, data, os.ModePerm)
	if err != nil {
		return err // Return the error if writing to file failed
	}
	return nil // Return nil if writing to file was successful
}
