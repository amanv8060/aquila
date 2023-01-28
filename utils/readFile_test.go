/*
Copyright © 2022 Aman Verma. All rights reserved.
Use of this source code is governed by a BSD-style
license that can be found in the LICENSE file.
*/

package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadFile(t *testing.T) {
	// Setup test directory and files
	tmpDir := t.TempDir()
	ExcerptsPath = filepath.Join(tmpDir, "code_regions")

	testFile := filepath.Join(tmpDir, "test.go")
	err := os.WriteFile(testFile, []byte(`
package main

// #aqstart init
func init() {
    setup()
    configure()
}
// #aqend init

// #aqstart main
func main() {
    // Some code
    doWork()
}
// #aqend main
`), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Run the test
	ReadFile(testFile, tmpDir)

	// Verify the output
	jsonFile := filepath.Join(ExcerptsPath, "test.go.json")
	if _, err := os.Stat(jsonFile); os.IsNotExist(err) {
		t.Error("Expected JSON file was not created")
	}

	// Clean up
	os.RemoveAll(tmpDir)
}

func TestRemoveTrailingLines(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  []string
	}{
		{
			name:  "empty slice",
			input: []string{},
			want:  []string{},
		},
		{
			name:  "no blank lines",
			input: []string{"code", "more code"},
			want:  []string{"code", "more code"},
		},
		{
			name:  "trailing blank lines",
			input: []string{"code", "", "  ", "\t"},
			want:  []string{"code"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := removeTrailingLines(&tt.input)
			if len(*got) != len(tt.want) {
				t.Errorf("removeTrailingLines() got len = %v, want %v", len(*got), len(tt.want))
			}
			for i := range tt.want {
				if i < len(*got) && (*got)[i] != tt.want[i] {
					t.Errorf("removeTrailingLines()[%d] = %v, want %v", i, (*got)[i], tt.want[i])
				}
			}
		})
	}
}
