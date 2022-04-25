package utils

import (
	directive "aquila/models"
	"bufio"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/rs/zerolog/log"

	"os"
	"path/filepath"
)

type AQSnippet struct {
	Key   string
	Value []string
}

var entireFile = ""

var openRegions = mapset.NewSet[string]()
var regions = make(map[string][]string)

func ReadFile() {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exePath := filepath.Dir(ex)
	f, err := os.OpenFile(exePath+"/main.go", os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Fatal().Msgf("open file error: %v", err)
		return
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Text()
		processLine(line)
	}
	if err := sc.Err(); err != nil {
		log.Fatal().Msgf("scan file error: %v", err)
		return
	}

	print()
}

func processLine(line string) {
	var _directive = directive.GetDirective(line)
	if _directive == nil {
		openRegions.Each(func(ele string) bool {
			regions[ele] = append(regions[ele], line)
			return false
		})
		return
	}

	if _directive.Kind == directive.StartDirective {
		regionStart(*_directive)
	} else {
		regionEnd(*_directive)
	}
}

func regionStart(directive directive.Directive) {

	var _regions = directive.Regions

	if len(_regions) == 0 {
		_regions = append(_regions, entireFile)
	}
	for _, region := range _regions {
		if openRegions.Contains(region) {
			log.Debug().Msg(region + " is opened again,  ignoring reopen.")
		} else {
			openRegions.Add(region)
		}
	}
}

func regionEnd(directive directive.Directive) {

	var _regions = directive.Regions

	if len(_regions) == 0 {
		_regions = append(_regions, entireFile)
	}

	for _, region := range _regions {
		if openRegions.Contains(region) {
			regions[region] = append(regions[region], "...")
		} else {
			log.Debug().Msg(region + "doesn't have a start")
		}
	}
}
