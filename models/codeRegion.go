/*
Copyright © 2022 Aman Verma. All rights reserved.
Use of this source code is governed by a BSD-style
license that can be found in the LICENSE file.
*/

package models

import "regexp"

type CodeRegion struct {
	Path       string
	RegionName string
}

var codeRegionRegex = "^<\\?code-region \"([^\"]*)\" region=\"([^\"]*)\"\\?>"

var codeRegion = regexp.MustCompile(codeRegionRegex)

func GetCodeRegion(line string) *CodeRegion {
	match := codeRegion.FindStringSubmatch(line)
	if match != nil {
		return &CodeRegion{
			Path:       match[1],
			RegionName: match[2],
		}
	} else {
		return nil
	}
}
