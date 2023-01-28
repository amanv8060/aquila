/*
Copyright © 2022 Aman Verma. All rights reserved.
Use of this source code is governed by a BSD-style
license that can be found in the LICENSE file.
*/

package models

import (
	"reflect"
	"testing"
)

func TestGetDirective(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *Directive
		wantNil bool
	}{
		{
			name:  "start directive",
			input: "// #aqstart main",
			want: &Directive{
				Regions: []string{"main"},
				Kind:    StartDirective,
			},
		},
		{
			name:  "end directive",
			input: "// #aqend main",
			want: &Directive{
				Regions: []string{"main"},
				Kind:    EndDirective,
			},
		},
		{
			name:  "multiple regions",
			input: "// #aqstart region1,region2",
			want: &Directive{
				Regions: []string{"region1", "region2"},
				Kind:    StartDirective,
			},
		},
		{
			name:  "with whitespace",
			input: "  // #aqstart  main  ",
			want: &Directive{
				Regions: []string{"main"},
				Kind:    StartDirective,
			},
		},
		{
			name:    "invalid directive",
			input:   "// not a directive",
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetDirective(tt.input)
			if tt.wantNil {
				if got != nil {
					t.Errorf("GetDirective() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("GetDirective() returned nil, want non-nil")
			}
			if got.Kind != tt.want.Kind {
				t.Errorf("Kind = %v, want %v", got.Kind, tt.want.Kind)
			}
			if !reflect.DeepEqual(got.Regions, tt.want.Regions) {
				t.Errorf("Regions = %v, want %v", got.Regions, tt.want.Regions)
			}
		})
	}
}
