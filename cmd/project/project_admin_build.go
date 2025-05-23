package project

import (
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/allincart-org/allincart-cli/extension"
	"github.com/allincart-org/allincart-cli/internal/phpexec"
	"github.com/allincart-org/allincart-cli/logging"
	"github.com/allincart-org/allincart-cli/shop"
)

var projectAdminBuildCmd = &cobra.Command{
	Use:   "admin-build [project-dir]",
	Short: "Builds the Administration",
	RunE: func(cmd *cobra.Command, args []string) error {
		var projectRoot string
		var err error

		if len(args) == 1 {
			// We need an absolute path for webpack
			projectRoot, err = filepath.Abs(args[0])
			if err != nil {
				return err
			}
		} else if projectRoot, err = findClosestAllincartProject(); err != nil {
			return err
		}

		shopCfg, err := shop.ReadConfig(projectConfigPath, true)
		if err != nil {
			return err
		}

		logging.FromContext(cmd.Context()).Infof("Looking for extensions to build assets in project")

		if err := runTransparentCommand(commandWithRoot(phpexec.ConsoleCommand(phpexec.AllowBinCI(cmd.Context()), "feature:dump"), projectRoot)); err != nil {
			return err
		}

		sources, err := filterAndGetSources(cmd, projectRoot, shopCfg)
		if err != nil {
			return err
		}

		forceInstall, _ := cmd.PersistentFlags().GetBool("force-install-dependencies")

		allincartConstraint, err := extension.GetAllincartProjectConstraint(projectRoot)
		if err != nil {
			return err
		}

		assetCfg := extension.AssetBuildConfig{
			DisableStorefrontBuild: true,
			AllincartRoot:          projectRoot,
			AllincartVersion:       allincartConstraint,
			NPMForceInstall:        forceInstall,
			ContributeProject:      extension.IsContributeProject(projectRoot),
		}

		if err := extension.BuildAssetsForExtensions(cmd.Context(), sources, assetCfg); err != nil {
			return err
		}

		skipAssetsInstall, _ := cmd.PersistentFlags().GetBool("skip-assets-install")
		if skipAssetsInstall {
			return nil
		}

		return runTransparentCommand(commandWithRoot(phpexec.ConsoleCommand(cmd.Context(), "assets:install"), projectRoot))
	},
}

func init() {
	projectRootCmd.AddCommand(projectAdminBuildCmd)
	projectAdminBuildCmd.PersistentFlags().Bool("skip-assets-install", false, "Skips the assets installation")
	projectAdminBuildCmd.PersistentFlags().Bool("force-install-dependencies", false, "Force install NPM dependencies")
	projectAdminBuildCmd.PersistentFlags().String("only-extensions", "", "Only watch the given extensions (comma separated)")
	projectAdminBuildCmd.PersistentFlags().String("skip-extensions", "", "Skips the given extensions (comma separated)")
	projectAdminBuildCmd.PersistentFlags().Bool("only-custom-static-extensions", false, "Only build extensions from custom/static-plugins directory")
}
