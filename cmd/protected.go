package cmd

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/zeihanaulia/simple-oauth2/repositories"
	"github.com/zeihanaulia/simple-oauth2/repositories/mongodb"
	"github.com/zeihanaulia/simple-oauth2/services/protected"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var protectedCmd = &cobra.Command{
	Use:   "protected",
	Short: "Start protected service",
	Long:  `Start protected servive`,
	RunE: func(cmd *cobra.Command, args []string) error {
		port := net.JoinHostPort("0.0.0.0", strconv.Itoa(protectedPort))
		logger.Info(fmt.Sprintf("protected service listening to http://%s", port))
		zapLogger := logger.With(zap.String("service", "protected"))

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		db, _ := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))

		var tokens repositories.Tokens
		{
			tokens = mongodb.NewTokens(db)
		}

		server := protected.NewServer(port, tokens)
		return logError(zapLogger, server.Run())
	},
}

func init() {
	RootCmd.AddCommand(protectedCmd)
}
