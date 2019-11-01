package cmd

import (
	"fmt"
	"net"
	"strconv"

	"github.com/zeihanaulia/simple-oauth2/authorization"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var authorizationCmd = &cobra.Command{
	Use:   "authorization",
	Short: "Start authorization service",
	Long:  `Start authorization servive`,
	RunE: func(cmd *cobra.Command, args []string) error {
		port := net.JoinHostPort("0.0.0.0", strconv.Itoa(authorizationPort))
		logger.Info(fmt.Sprintf("authorization service listening to :%s", port))
		zapLogger := logger.With(zap.String("service", "authorization"))
		server := authorization.NewServer(port)
		return logError(zapLogger, server.Run())
	},
}

func init() {
	RootCmd.AddCommand(authorizationCmd)
}
