/*
Copyright © 2022 Aman Verma. All rights reserved.
Use of this source code is governed by a BSD-style
license that can be found in the LICENSE file.
*/

package models

import (
	"regexp"
	"strconv"
	"strings"
)

type CodeRegion struct {
	Path       string
	RegionName string
	StartLine  int
	EndLine    int
}

var codeRegionRegex = `^<\?code-region\s+"([^"]+)"\s+region="([^"]+)"(?:\s+lines="(\d+)(?:-(\d+))?")?`

var codeRegion = regexp.MustCompile(codeRegionRegex)

func GetCodeRegion(line string) *CodeRegion {
	match := codeRegion.FindStringSubmatch(line)
	if match != nil {
		region := &CodeRegion{
			Path:       strings.TrimSpace(match[1]),
			RegionName: strings.TrimSpace(match[2]),
			StartLine:  0,
			EndLine:    0,
		}
		
		// Parse line number(s)
		if len(match) > 3 && match[3] != "" {
			startLine, err := strconv.Atoi(match[3])
			if err == nil {
				region.StartLine = startLine
				region.EndLine = startLine // Default to single line
				
				// If range is specified (match[4] exists and not empty)
				if len(match) > 4 && match[4] != "" {
					if endLine, err := strconv.Atoi(match[4]); err == nil {
						region.EndLine = endLine
					}
				}
			}
		}
		
		return region
	}
	return nil
}
