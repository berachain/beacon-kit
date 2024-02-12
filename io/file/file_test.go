// SPDX-License-Identifier: MIT
//
// Copyright (c) 2023 Berachain Foundation
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

package file_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/itsdevbear/bolaris/io/file"
)

func TestHasDir(t *testing.T) {
	dir, err := os.MkdirTemp("", "testdir")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir) // clean up

	exists, err := file.HasDir(dir)
	if err != nil || !exists {
		t.Errorf("Expected directory to exist, got error: %v", err)
	}

	notExists, err := file.HasDir(filepath.Join(dir, "nonexistent"))
	if err != nil || notExists {
		t.Errorf("Expected directory to not exist, got error: %v", err)
	}
}

func TestMkdirAll(t *testing.T) {
	dir, err := os.MkdirTemp("", "testmkdir")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir) // clean up

	newDir := filepath.Join(dir, "newdir", "subdir")
	err = file.MkdirAll(newDir)
	if err != nil {
		t.Errorf("Failed to create directory: %v", err)
	}

	if _, err = os.Stat(newDir); os.IsNotExist(err) {
		t.Errorf("Directory was not created")
	}
}

func TestCreateFile(t *testing.T) {
	dir, err := os.MkdirTemp("", "testcreatefile")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir) // clean up

	filePath := filepath.Join(dir, "testfile.txt")
	err = file.CreateFile(filePath)
	if err != nil {
		t.Errorf("Failed to create file: %v", err)
	}

	if _, err = os.Stat(filePath); os.IsNotExist(err) {
		t.Errorf("File was not created")
	}
}

func TestExists(t *testing.T) {
	dir, err := os.MkdirTemp("", "testexists")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir) // clean up

	filePath := filepath.Join(dir, "testfile.txt")
	_, err = os.Create(filePath)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	if !file.Exists(filePath) {
		t.Errorf("File should exist")
	}

	if file.Exists(filepath.Join(dir, "nonexistent.txt")) {
		t.Errorf("File should not exist")
	}
}

func TestWriteFile(t *testing.T) {
	dir, err := os.MkdirTemp("", "testwritefile")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir) // clean up

	filePath := filepath.Join(dir, "testfile.txt")
	data := []byte("Hello, World!")
	err = file.WriteFile(filePath, data)
	if err != nil {
		t.Errorf("Failed to write to file: %v", err)
	}

	readData, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	if !bytes.Equal(readData, data) {
		t.Errorf("Data written and read does not match")
	}
}
