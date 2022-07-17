/*
Copyright © 2022 Aman Verma. All rights reserved.
Use of this source code is governed by a BSD-style
license that can be found in the LICENSE file.
*/

package models

import (
	"regexp"
	"strings"
)

var BlankLine = "^\\s*$"
var LeadingWhiteSpace = "^[ \\t]*"
var directiveRegex = "^(\\s*)(\\S.*?)?#(?P<DirectiveType>aqstart|aqend)\\b\\s*(?P<args>.*?)(?:\\s*(?:-->|\\*\\/))?\\s*$"

type Directive struct {
	Regions []string
	Kind    DirectiveType
}

// DirectiveType indicates whether the directive starts or ends a doc region.
type DirectiveType uint

const (
	StartDirective DirectiveType = iota
	EndDirective
)

var reg = regexp.MustCompile(directiveRegex)

func GetDirective(line string) *Directive {
	match := reg.FindStringSubmatch(line)
	if match != nil {
		var kind DirectiveType
		if match[3] == "aqstart" {
			kind = StartDirective
		} else {
			kind = EndDirective
		}
		args := strings.Split(match[4], ",")
		for i, arg := range args {
			args[i] = strings.TrimSpace(arg)
		}
		return &Directive{
			Regions: args,
			Kind:    kind,
		}
	} else {
		return nil
	}
}
