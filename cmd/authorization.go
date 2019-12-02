package cmd

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/zeihanaulia/simple-oauth2/repositories"
	"github.com/zeihanaulia/simple-oauth2/repositories/mongodb"

	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/zeihanaulia/simple-oauth2/services/authorization"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var authorizationCmd = &cobra.Command{
	Use:   "authorization",
	Short: "Start authorization service",
	Long:  `Start authorization servive`,
	RunE: func(cmd *cobra.Command, args []string) error {
		port := net.JoinHostPort("0.0.0.0", strconv.Itoa(authorizationPort))
		logger.Info(fmt.Sprintf("authorization service listening to http://%s", port))
		zapLogger := logger.With(zap.String("service", "authorization"))

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		db, _ := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))

		var tokens repositories.Tokens
		{
			tokens = mongodb.NewTokens(db)
		}

		server := authorization.NewServer(port, "services/authorization/templates/", tokens)
		return logError(zapLogger, server.Run())
	},
}

func init() {
	RootCmd.AddCommand(authorizationCmd)
}
