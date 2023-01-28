/*
 Copyright © 2022 Aman Verma. All rights reserved.
 Use of this source code is governed by a BSD-style
 license that can be found in the LICENSE file.
*/

package utils

import (
	"aquila/models"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestProcessWriteLine(t *testing.T) {
	// Setup test directory and files
	tmpDir := t.TempDir()
	ExcerptsPath = filepath.Join(tmpDir, "code_regions") + "/"
	err := os.MkdirAll(ExcerptsPath, 0755)
	if err != nil {
		t.Fatal(err)
	}

	// Create a test JSON file with proper indentation
	codeRegions := map[string][]string{
		"main": {
			"func main() {",
			"    fmt.Println(\"Hello\")",
			"    doWork()",
			"    cleanup()",
			"}",
		},
	}

	testFile := "test.go"
	jsonData, err := json.MarshalIndent(codeRegions, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	jsonPath := filepath.Join(ExcerptsPath, testFile+".json")
	err = os.WriteFile(jsonPath, jsonData, 0644)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(jsonPath); os.IsNotExist(err) {
		t.Fatalf("Failed to create test JSON file at %s", jsonPath)
	}

	tests := []struct {
		name          string
		codeRegion    *models.CodeRegion
		expectedLines int
		wantContent   []string
	}{
		{
			name: "full region",
			codeRegion: &models.CodeRegion{
				Path:       testFile,
				RegionName: "main",
			},
			expectedLines: 7, // ```go + 5 lines + ```
			wantContent: []string{
				"```go",
				"func main() {",
				"    fmt.Println(\"Hello\")",
				"    doWork()",
				"    cleanup()",
				"}",
				"```",
			},
		},
		{
			name: "single line",
			codeRegion: &models.CodeRegion{
				Path:       testFile,
				RegionName: "main",
				StartLine:  2,
				EndLine:    2,
			},
			expectedLines: 3, // ```go + 1 line + ```
			wantContent: []string{
				"```go",
				"    fmt.Println(\"Hello\")",
				"```",
			},
		},
		{
			name: "line range",
			codeRegion: &models.CodeRegion{
				Path:       testFile,
				RegionName: "main",
				StartLine:  2,
				EndLine:    3,
			},
			expectedLines: 4, // ```go + 2 lines + ```
			wantContent: []string{
				"```go",
				"    fmt.Println(\"Hello\")",
				"    doWork()",
				"```",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var lines []string
			updated := processWriteLine(tt.codeRegion, &lines)
			if !updated {
				t.Errorf("processWriteLine() returned false, want true\nJSON path: %s", jsonPath)
				return
			}
			if len(lines) != tt.expectedLines {
				t.Errorf("got %d lines, want %d\nGot lines: %q", len(lines), tt.expectedLines, lines)
			}
			if !reflect.DeepEqual(lines, tt.wantContent) {
				t.Errorf("content mismatch\ngot:  %q\nwant: %q", lines, tt.wantContent)
			}
		})
	}
}

func TestLinesFromReader(t *testing.T) {
	input := "line1\nline2\nline3"
	reader := strings.NewReader(input)

	lines, err := LinesFromReader(reader)
	if err != nil {
		t.Fatalf("LinesFromReader() error = %v", err)
	}

	expected := []string{"line1", "line2", "line3"}
	if len(lines) != len(expected) {
		t.Errorf("got %d lines, want %d", len(lines), len(expected))
	}

	for i := range expected {
		if lines[i] != expected[i] {
			t.Errorf("line[%d] = %v, want %v", i, lines[i], expected[i])
		}
	}
}
