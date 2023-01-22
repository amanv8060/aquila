/*
Copyright © 2022 Aman Verma. All rights reserved.
Use of this source code is governed by a BSD-style
license that can be found in the LICENSE file.
*/

package cmd

import (
	"aquila/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate code regions from your codebase.",
	Long: `Generate command generates code regions from your codebase.

It will generate code regions from a given codebase and store them in a directory called code_regions within a path 
specified by the user through config file. The code regions are stored in a JSON file with the name of the file as 
the key and the value being an array of strings containing the code region.
`,
	Run: func(cmd *cobra.Command, args []string) {
		codePath := viper.GetString("codePath")
		// recursively get all files path in the codePaths
		err := filepath.Walk(codePath,
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if !info.IsDir() {
					// print path
					log.Println(path)
					// get file extension of the path
					extension := filepath.Ext(path)

					// check if extension is present inside the supported extensions in the array
					var result = false
					for _, x := range supported_extension {
						if x == extension {
							result = true
							break
						}
					}

					if result {
						utils.ReadFile(path, codePath)
					} else {
						log.Println("File extension not supported")
					}
					return nil
				}
				return nil
			})
		if err != nil {
			log.Println(err)
		}
	},
}

// list of file extension
var supported_extension = []string{".go", ".py", ".js", ".java", ".c", ".cpp", ".cs", ".sh", ".html", ".css", ".md"}

func init() {
	rootCmd.AddCommand(generateCmd)
}
