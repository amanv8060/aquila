package models

import (
	"regexp"
	"strings"
)

var directiveRegex = "^(#(aqstart|aqend))"

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
	if reg.MatchString(line) {
		var _line = reg.ReplaceAllString(line, "")
		_line = strings.TrimSpace(_line)
		var kind DirectiveType
		if reg.FindString(line) == "#aqstart" {
			kind = StartDirective
		} else {
			kind = EndDirective
		}
		return &Directive{
			Regions: strings.Split(_line, ","),
			Kind:    kind,
		}
	} else {
		return nil
	}
}
