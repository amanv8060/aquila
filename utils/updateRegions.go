package utils

import (
	"aquila/models"
	"bufio"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var docsPath = "./docs/"

func UpdateRegions() {

	// recursively walk through docs directory
	// for each file, read it and update the code regions

	// walk through the docs directory
	// for each file, read it and update the code regions
	err := filepath.Walk(docsPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			// fmt.Println("reading file: ", path)
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

	// loop over the lines
	for _, line := range lines {
		// update codeRegionInserted
		wasUpdated := processWriteLine(line, &updatedLines)

		if !codeRegionInserted {
			codeRegionInserted = wasUpdated
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
func processWriteLine(line string, updatedLines *[]string) bool {
	// match line with regex

	codeRegion := models.GetCodeRegion(line)

	// insert the current line into the updatedLines array
	*updatedLines = append(*updatedLines, line)

	if codeRegion != nil {
		// fetch code region from yaml file
		splitPath := strings.Split(codeRegion.Path, ".")
		modifiedCodeRegionPath, extension := splitPath[0], splitPath[1]

		jsonPath := ExcerptsPath + modifiedCodeRegionPath + ".json"

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
