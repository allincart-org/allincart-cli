package cmd

import (
	"context"
	"os"
	"slices"

	"github.com/spf13/cobra"

	"github.com/allincart-org/allincart-cli/cmd/account"
	"github.com/allincart-org/allincart-cli/cmd/extension"
	"github.com/allincart-org/allincart-cli/cmd/project"
	accountApi "github.com/allincart-org/allincart-cli/internal/account-api"
	"github.com/allincart-org/allincart-cli/internal/config"
	"github.com/allincart-org/allincart-cli/logging"
)

var (
	cfgFile string
	version = "dev"
)

var rootCmd = &cobra.Command{
	Use:     "allincart-cli",
	Short:   "A cli for common Allincart tasks",
	Long:    `This application contains some utilities like extension management`,
	Version: version,
}

func Execute(ctx context.Context) {
	ctx = logging.WithLogger(ctx, logging.NewLogger(slices.Contains(os.Args, "--verbose")))
	accountApi.SetUserAgent("allincart-cli/" + version)

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		logging.FromContext(ctx).Fatalln(err)
	}
}

func init() {
	rootCmd.SilenceErrors = true

	cobra.OnInitialize(func() {
		_ = config.InitConfig(cfgFile)
	})

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.allincart-cli.yaml)")
	rootCmd.PersistentFlags().Bool("verbose", false, "show debug output")

	project.Register(rootCmd)
	extension.Register(rootCmd)
	account.Register(rootCmd, func(commandName string) (*account.ServiceContainer, error) {
		err := config.InitConfig(cfgFile)
		if err != nil {
			return nil, err
		}
		conf := config.Config{}
		if commandName == "login" || commandName == "logout" {
			return &account.ServiceContainer{
				Conf:          conf,
				AccountClient: nil,
			}, nil
		}
		client, err := accountApi.NewApi(rootCmd.Context(), conf)
		if err != nil {
			return nil, err
		}
		return &account.ServiceContainer{
			Conf:          conf,
			AccountClient: client,
		}, nil
	})
}
