package cmd

import (
	"aquila/utils"
	"github.com/spf13/cobra"
)

// generateCmd represents the generate command
var updateMdCmd = &cobra.Command{
	Use:   "update-md",
	Short: "Update Annotated Regions in your markdown with code regions",
	Long:  `Updates the annotated regions in your markdown with the code regions.`,
	Run: func(cmd *cobra.Command, args []string) {
		utils.UpdateRegions()
	},
}

func init() {
	rootCmd.AddCommand(updateMdCmd)
}
