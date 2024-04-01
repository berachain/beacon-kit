// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package file

import (
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strings"
)

// AbsPath expands the environment variables in a path and converts it
// to an absolute path.
func AbsPath(inputPath string) (string, error) {
	if strings.HasPrefix(inputPath, "~/") || inputPath == "~" {
		homeDir, err := HomeDir()
		if err != nil {
			return "", err
		}
		inputPath = filepath.Join(homeDir, inputPath[1:])
	}

	return filepath.Abs(path.Clean(os.ExpandEnv(inputPath)))
}

// HomeDir returns the home directory for the current user.
func HomeDir() (string, error) {
	if homeDir := os.Getenv("HOME"); homeDir != "" {
		return homeDir, nil
	}
	user, err := user.Current()
	if err != nil {
		return "", err
	}
	return user.HomeDir, nil
}

// HasDir checks if a directory exists.
func HasDir(dirPath string) (bool, error) {
	absDirPath, err := AbsPath(dirPath)
	if err != nil {
		return false, err
	}

	stat, err := os.Stat(absDirPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return stat.IsDir(), nil
}

// MkdirAll creates a directory along with any necessary parents.
// If the directory already exists, MkdirAll does nothing and returns nil.
func MkdirAll(dirPath string) error {
	absDirPath, err := AbsPath(dirPath)
	if err != nil {
		return err
	}
	return os.MkdirAll(absDirPath, os.ModePerm)
}

// Create creates or truncates a file at the specified path.
func Create(filePath string) error {
	absFilePath, err := AbsPath(filePath)
	if err != nil {
		return err
	}

	file, err := os.Create(absFilePath)
	if err != nil {
		return err
	}
	return file.Close()
}

// Exists checks if a file exists.
func Exists(path string) bool {
	absPath, err := AbsPath(path)
	if err != nil {
		return false
	}

	_, err = os.Stat(absPath)
	if err == nil {
		return true // File exists
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

// Read reads the file at the specified path and returns its content.
func Read(filePath string) ([]byte, error) {
	absFilePath, err := AbsPath(filePath)
	if err != nil {
		return nil, err
	}
	if ok := Exists(absFilePath); !ok {
		return nil, os.ErrNotExist
	}

	data, err := os.ReadFile(absFilePath)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Write writes data to a file at the specified path.
func Write(filePath string, data []byte) error {
	absFilePath, err := AbsPath(filePath)
	if err != nil {
		return err
	}

	if ok := Exists(absFilePath); !ok {
		err = Create(absFilePath)
		if err != nil {
			return err
		}
	}

	err = os.WriteFile(absFilePath, data, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}
