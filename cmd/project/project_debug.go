package project

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/allincart-org/allincart-cli/extension"
	"github.com/allincart-org/allincart-cli/internal/table"
	"github.com/allincart-org/allincart-cli/logging"
	"github.com/allincart-org/allincart-cli/shop"
)

var projectDebug = &cobra.Command{
	Use:   "debug",
	Short: "Shows detected Allincart version and detected extensions for further debugging",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		args[0], err = filepath.Abs(args[0])
		if err != nil {
			return err
		}

		shopCfg, err := shop.ReadConfig(projectConfigPath, true)
		if err != nil {
			return err
		}

		allincartConstraint, err := extension.GetAllincartProjectConstraint(args[0])
		if err != nil {
			return err
		}

		if shopCfg.IsFallback() {
			fmt.Printf("Could not find a %s, using fallback config\n", projectConfigPath)
		} else {
			fmt.Printf("Found config: Yes\n")
		}
		fmt.Printf("Detected following Allincart version: %s\n", allincartConstraint.String())

		sources := extension.FindAssetSourcesOfProject(logging.DisableLogger(cmd.Context()), args[0], shopCfg)

		fmt.Println("Following extensions/bundles has been detected")
		table := table.NewWriter(os.Stdout)
		table.Header([]string{"Name", "Path"})

		for _, source := range sources {
			_ = table.Append([]string{source.Name, source.Path})
		}

		_ = table.Render()

		return nil
	},
}

func init() {
	projectRootCmd.AddCommand(projectDebug)
}
