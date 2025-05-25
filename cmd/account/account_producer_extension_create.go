package account

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	accountApi "github.com/allincart-org/allincart-cli/internal/account-api"
	"github.com/allincart-org/allincart-cli/logging"
)

var accountCompanyProducerExtensionCreateCmd = &cobra.Command{
	Use:   "create [name] [plugin|theme|app]",
	Short: "Creates a new extension",
	Args:  cobra.ExactArgs(2),
	ValidArgsFunction: func(_ *cobra.Command, args []string, _ string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 1 {
			return []string{accountApi.GenerationApp, accountApi.GenerationTheme, accountApi.GenerationPlugin}, cobra.ShellCompDirectiveNoFileComp
		}

		return []string{}, cobra.ShellCompDirectiveNoFileComp
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		p, err := services.AccountClient.Producer(cmd.Context())
		if err != nil {
			return fmt.Errorf("cannot get producer endpoint: %w", err)
		}

		profile, err := p.Profile(cmd.Context())
		if err != nil {
			return fmt.Errorf("cannot get producer profile: %w", err)
		}

		if args[1] != accountApi.GenerationApp && args[1] != accountApi.GenerationPlugin && args[1] != accountApi.GenerationTheme {
			return fmt.Errorf("generation must be one of these options: %s %s %s", accountApi.GenerationPlugin, accountApi.GenerationTheme, accountApi.GenerationApp)
		}

		if !strings.HasPrefix(args[0], profile.Prefix) {
			return fmt.Errorf("extension name must start with the prefix %s", profile.Prefix)
		}

		extension, err := p.CreateExtension(cmd.Context(), accountApi.CreateExtensionRequest{
			Name:       args[0],
			SubType:    args[1],
			ProducerID: p.GetId(),
		})
		if err != nil {
			return fmt.Errorf("cannot create extension: %w", err)
		}

		logging.FromContext(cmd.Context()).Infof("Extension with name %s has been successfully created", extension.Name)

		return nil
	},
}

func init() {
	accountCompanyProducerExtensionCmd.AddCommand(accountCompanyProducerExtensionCreateCmd)
}
