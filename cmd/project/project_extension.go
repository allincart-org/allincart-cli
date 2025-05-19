package project

import "github.com/spf13/cobra"

var projectExtensionCmd = &cobra.Command{
	Use:   "extension",
	Short: "Manage the extensions of the Allincart shop",
}

func init() {
	projectRootCmd.AddCommand(projectExtensionCmd)
}
