package cmd

import (
	"github.com/spf13/cobra"
)

var allCmd = &cobra.Command{
	Use:   "all",
	Short: "Start all service",
	Long:  `Start all servive`,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.Info("starting all service")
		go protectedCmd.RunE(protectedCmd, args)
		go clientCmd.RunE(clientCmd, args)
		return authorizationCmd.RunE(authorizationCmd, args)
	},
}

func init() {
	RootCmd.AddCommand(allCmd)
}
