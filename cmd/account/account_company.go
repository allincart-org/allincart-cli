package account

import (
	"github.com/spf13/cobra"
)

var accountCompanyRootCmd = &cobra.Command{
	Use:   "company",
	Short: "Manage your Allincart company",
}

func init() {
	accountRootCmd.AddCommand(accountCompanyRootCmd)
}
