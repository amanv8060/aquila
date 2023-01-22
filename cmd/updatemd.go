package cmd

import (
	"aquila/utils"
	"github.com/spf13/cobra"
)

// updateMdCmd represents the generate command
var updateMdCmd = &cobra.Command{
	Use:   "update-md",
	Short: "Update Annotated Regions in your markdown with code regions",
	Long: `Updates the annotated regions in your markdown with the code regions.

It will update the annotated regions in your markdown with the code regions from the code_regions directory.
It works on insert-or-replace basis. If the annotated region is not present in the markdown, it will be inserted.
If the annotated region has a code region present in the markdown, it will be replaced with the newly generated code
region.
`,
	Run: func(cmd *cobra.Command, args []string) {
		utils.UpdateRegions()
	},
}

func init() {
	rootCmd.AddCommand(updateMdCmd)
}
