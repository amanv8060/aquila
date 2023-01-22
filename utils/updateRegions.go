/*
Copyright © 2022 Aman Verma. All rights reserved.
Use of this source code is governed by a BSD-style
license that can be found in the LICENSE file.
*/

package utils

import (
	"aquila/models"
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func UpdateRegions() {

	// recursively walk through docs directory
	// for each file, read it and update the code regions
	var docsPath = viper.GetString("docsPath")

	// walk through the docs directory
	// for each file, read it and update the code regions
	err := filepath.Walk(docsPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			fmt.Println("reading file: ", path)
			// read the file and update the code regions
			processFile(path)
		}
		return nil
	})
	if err != nil {
		log.Fatal().Msgf("open file error: %v", err)
		return
	}
}

func processFile(path string) {
	// open the files and read them to modify the code regions
	f, err := os.OpenFile(path, os.O_RDWR, os.ModePerm)
	if err != nil {
		log.Fatal().Msgf("open file error: %v", err)
		return
	}

	sc := bufio.NewReadWriter(bufio.NewReader(f), bufio.NewWriter(f))

	defer func(sc *bufio.ReadWriter, f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatal().Msgf("close file error: %v", err)
		}
	}(sc, f)

	var updatedLines []string
	// read the file line by line
	lines, err := LinesFromReader(sc)

	// flag to check if the code region was inserted or not
	var codeRegionInserted bool

	// indexed loop over the lines
	var linesLen = len(lines)
	for i := 0; i < linesLen; i++ {
		line := lines[i]

		// add the line to the updatedLines array
		updatedLines = append(updatedLines, line)

		// find if the current line has regex in it.
		codeRegion := models.GetCodeRegion(line)

		if codeRegion != nil {
			var wasUpdated = processWriteLine(codeRegion, &updatedLines)

			// check if this is the last line
			if i <= linesLen-1 {
				nextLine := lines[i+1]
				// if this is not the last line, check if the next starts with code block.
				if strings.HasPrefix(nextLine, "```") {
					//loop over until we find the next line that ends the code block
					for j := i + 2; j < linesLen; j++ {
						if strings.HasPrefix(lines[j], "```") {
							// break the loop
							i = j
							break
						}
					}
				}
			}

			if !codeRegionInserted {
				codeRegionInserted = wasUpdated
			}
		}
	}

	// rewrite the file if code regions were inserted
	if codeRegionInserted {
		// delete the contents of the file
		err = f.Truncate(0)
		if err != nil {
			log.Fatal().Msgf("truncate file error: %v", err)
			return
		}
		// go back to the beginning of the file
		_, err = f.Seek(0, 0)
		if err != nil {
			log.Fatal().Msgf("seek file error: %v", err)
			return
		}
		// write the updated lines to the file
		for _, line := range updatedLines {
			_, err = f.WriteString(line + "\n")
		}
	}

	println("done reading file: ", path)
}

// function to read the lines, and update the newlines array, which is passed by reference.
// returns a bool indicating whether the code region was inserted or not
func processWriteLine(codeRegion *models.CodeRegion, updatedLines *[]string) bool {

	// fetch code region from yaml file
	splitPath := strings.Split(codeRegion.Path, ".")
	_, extension := splitPath[0], splitPath[1]

	jsonPath := ExcerptsPath + codeRegion.Path + ".json"

	// read the yaml file
	jsonFile, err := os.OpenFile(jsonPath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Fatal().Msgf("open file error: %v", err)
		return false
	}

	defer func(yamlFile *os.File) {
		err := yamlFile.Close()
		if err != nil {
			log.Fatal().Msgf("close file error: %v", err)
		}
	}(jsonFile)

	var regions = make(map[string][]string)

	decoder := json.NewDecoder(jsonFile)

	// read the json file
	err = decoder.Decode(&regions)
	if err != nil {
		return false
	}

	// get the code region from the yaml file
	codeRegionLines := regions[codeRegion.RegionName]

	if err != nil {
		log.Fatal().Msgf("write to file error: %v", err)
		return false
	}
	header := "```" + extension
	footer := "```"
	// if codeRegionLines are not empty or nil, insert the code region
	if codeRegionLines != nil && len(codeRegionLines) > 0 {
		// insert header
		*updatedLines = append(*updatedLines, header)

		// insert code region lines
		*updatedLines = append(*updatedLines, codeRegionLines...)

		// insert footer
		*updatedLines = append(*updatedLines, footer)

		return true
	}
	return false
}

func LinesFromReader(r io.Reader) ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}
