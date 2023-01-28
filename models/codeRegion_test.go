/*
Copyright © 2022 Aman Verma. All rights reserved.
Use of this source code is governed by a BSD-style
license that can be found in the LICENSE file.
*/

package models

import (
	"testing"
)

func TestGetCodeRegion(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *CodeRegion
		wantNil bool
	}{
		{
			name:  "basic region without lines",
			input: `<?code-region "file.go" region="main"?>`,
			want: &CodeRegion{
				Path:       "file.go",
				RegionName: "main",
				StartLine:  0,
				EndLine:    0,
			},
		},
		{
			name:  "region with single line",
			input: `<?code-region "file.go" region="main" lines="5"?>`,
			want: &CodeRegion{
				Path:       "file.go",
				RegionName: "main",
				StartLine:  5,
				EndLine:    5,
			},
		},
		{
			name:  "region with line range",
			input: `<?code-region "file.go" region="main" lines="5-10"?>`,
			want: &CodeRegion{
				Path:       "file.go",
				RegionName: "main",
				StartLine:  5,
				EndLine:    10,
			},
		},
		{
			name:  "region with extra whitespace",
			input: `<?code-region  "file.go"  region="main"  lines="5-10" ?>`,
			want: &CodeRegion{
				Path:       "file.go",
				RegionName: "main",
				StartLine:  5,
				EndLine:    10,
			},
		},
		{
			name:    "invalid format",
			input:   `<?code-region invalid?>`,
			wantNil: true,
		},
		{
			name:  "invalid line number (should keep 0)",
			input: `<?code-region "file.go" region="main" lines="abc"?>`,
			want: &CodeRegion{
				Path:       "file.go",
				RegionName: "main",
				StartLine:  0,
				EndLine:    0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetCodeRegion(tt.input)
			if tt.wantNil {
				if got != nil {
					t.Errorf("GetCodeRegion() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("GetCodeRegion() returned nil, want non-nil")
			}
			if got.Path != tt.want.Path {
				t.Errorf("Path = %v, want %v", got.Path, tt.want.Path)
			}
			if got.RegionName != tt.want.RegionName {
				t.Errorf("RegionName = %v, want %v", got.RegionName, tt.want.RegionName)
			}
			if got.StartLine != tt.want.StartLine {
				t.Errorf("StartLine = %v, want %v", got.StartLine, tt.want.StartLine)
			}
			if got.EndLine != tt.want.EndLine {
				t.Errorf("EndLine = %v, want %v", got.EndLine, tt.want.EndLine)
			}
		})
	}
}
