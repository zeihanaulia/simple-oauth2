package cmd

import (
	"fmt"
	"net"
	"strconv"

	"github.com/zeihanaulia/simple-oauth2/protected"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var protectedCmd = &cobra.Command{
	Use:   "protected",
	Short: "Start protected service",
	Long:  `Start protected servive`,
	RunE: func(cmd *cobra.Command, args []string) error {
		port := net.JoinHostPort("0.0.0.0", strconv.Itoa(protectedPort))
		logger.Info(fmt.Sprintf("protected service listening to :%s", port))
		zapLogger := logger.With(zap.String("service", "protected"))
		server := protected.NewServer(port)
		return logError(zapLogger, server.Run())
	},
}

func init() {
	RootCmd.AddCommand(protectedCmd)
}
