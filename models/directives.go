package models

import (
	"regexp"
	"strings"
)

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
		return &Directive{
			Regions: strings.Split(match[4], ","),
			Kind:    kind,
		}
	} else {
		return nil
	}
}
