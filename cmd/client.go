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
		port := net.JoinHostPort("0.0.0.0", strconv.Itoa(clientPort))
		logger.Info(fmt.Sprintf("client service listening to http://%s", port))
		zapLogger := logger.With(zap.String("service", "client"))
		server := client.NewServer(port)
		return logError(zapLogger, server.Run())
	},
}

func init() {
	RootCmd.AddCommand(clientCmd)
}
