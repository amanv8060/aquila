/*
Copyright © 2022 Aman Verma. All rights reserved.
Use of this source code is governed by a BSD-style
license that can be found in the LICENSE file.
*/

package utils

import (
	directive "aquila/models"
	"bufio"
	"encoding/json"
	"fmt"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/rs/zerolog/log"
	"os"
	"path/filepath"
	"regexp"
)

var entireFile = ""

var openRegions = mapset.NewSet[string]()
var regions = make(map[string][]string)

var ExcerptsPath = "./code_regions/"

func ReadFile(path string, codePath string) {
	f, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Fatal().Msgf("open file error: %v", err)
		return
	}

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatal().Msgf("close file error: %v", err)
		}
	}(f)

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Text()
		processReadLine(line)
	}
	if err := sc.Err(); err != nil {
		log.Fatal().Msgf("scan file error: %v", err)
		return
	}
	// iterate over regions and update them
	for region, lines := range regions {
		//call function to update regions[region]
		regions[region] = *removeTrailingLines(&lines)
		regions[region] = *removeTrailingPlaster(&lines)
		regions[region] = *removeTrailingLines(&lines)
	}

	// get relative path of file
	relPath, err := filepath.Rel(codePath, path)
	if err != nil {
		log.Warn().Msg("Couldn't get relative path")
		return
	}
	saveToYamlFile(relPath)
}

func removeTrailingPlaster(lines *[]string) *[]string {
	// remove last element if it matches ...
	if len(*lines) > 0 {
		if (*lines)[len(*lines)-1] == "..." {
			*lines = (*lines)[:len(*lines)-1]
		}
	}
	return lines
}

func removeTrailingLines(lines *[]string) *[]string {

	// remove all blank lines from end in lines
	for i := len(*lines) - 1; i >= 0; i-- {
		blank, err := regexp.MatchString(directive.BlankLine, (*lines)[i])
		if err == nil && blank {
			*lines = (*lines)[:i]
		} else {
			break
		}
	}
	return lines
}
func saveToYamlFile(path string) {
	// return if regions are empty
	if len(regions) == 0 {
		return
	}

	// replace extension in path with .json
	path = ExcerptsPath + "/" + path + ".json"
	err := os.MkdirAll(filepath.Dir(path), os.ModePerm)
	if err != nil {
		log.Fatal().Msg("Failed to create directories")
		return
	}

	data, err := json.Marshal(&regions)
	if err != nil {
		log.Fatal().Msg("Error yaml couldn't be parsed")
	}
	err2 := os.WriteFile(path, data, os.ModePerm)
	if err2 != nil {
		log.Fatal().Msg("Couldn't create file")
	}
	regions = make(map[string][]string)
	fmt.Println("data written")
}

func processReadLine(line string) {
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
			// remove trailing lines from regions[region] before adding plaster.
			lines := regions[region]
			regions[region] = *removeTrailingLines(&lines)
			regions[region] = append(regions[region], "\t...")
			openRegions.Remove(region)
		} else {
			log.Debug().Msg(region + "doesn't have a start")
		}
	}
}
