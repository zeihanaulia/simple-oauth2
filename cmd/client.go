package cmd

import (
	"fmt"
	"net"
	"strconv"

	"github.com/zeihanaulia/simple-oauth2/client"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Start client service",
	Long:  `Start client servive`,
	RunE: func(cmd *cobra.Command, args []string) error {
		options.ClientHostPort = net.JoinHostPort("0.0.0.0", strconv.Itoa(clientPort))
		logger.Info(fmt.Sprintf("client service listening to :%s", options.ClientHostPort))
		zapLogger := logger.With(zap.String("service", "client"))
		server := client.NewServer(options)
		return logError(zapLogger, server.Run())
	},
}

var options client.ConfigOptions

func init() {
	RootCmd.AddCommand(clientCmd)
}
