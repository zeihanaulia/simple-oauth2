package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger *zap.Logger

	clientPort        int
	authorizationPort int
	protectedPort     int
)

// RootCmd root for all command
var RootCmd = &cobra.Command{
	Use:   "simple-oauth2",
	Short: "DELEGATE. - A oauth2 demo application.",
	Long:  `DELEGATE. - A oauth2 demo application.`,
}

// Execute command for calling by main.go
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		logger.Fatal("What are you doing!", zap.Error(err))
		os.Exit(-1)
	}
}

func init() {
	RootCmd.PersistentFlags().IntVarP(&clientPort, "client-service-port", "c", 8081, "Port for client service")
	RootCmd.PersistentFlags().IntVarP(&authorizationPort, "authorization-service-port", "a", 8082, "Port for authorization service")
	RootCmd.PersistentFlags().IntVarP(&protectedPort, "protected-service-port", "p", 8083, "Port for protected service")

	logger, _ = zap.NewDevelopment(zap.AddStacktrace(zapcore.FatalLevel))
	cobra.OnInitialize(beforeInitialize)
}

func beforeInitialize() {
	if clientPort != 8081 {
		logger.Info("client service port changed", zap.Int("old", 8081), zap.Int("new", clientPort))
	}

	if authorizationPort != 8082 {
		logger.Info("authorization service port changed", zap.Int("old", 8082), zap.Int("new", authorizationPort))
	}

	if protectedPort != 8083 {
		logger.Info("protected service port changed", zap.Int("old", 8083), zap.Int("new", protectedPort))
	}
}

func logError(logger *zap.Logger, err error) error {
	if err != nil {
		logger.Error("Error running command", zap.Error(err))
	}
	return err
}
