package repositories

import (
	"context"

	"github.com/zeihanaulia/simple-oauth2/models"
)

type ClientTokens interface {
	Save(ctx context.Context, token models.Token) (models.Token, error)
}
