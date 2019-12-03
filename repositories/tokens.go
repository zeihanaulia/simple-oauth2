package repositories

import (
	"context"

	"github.com/zeihanaulia/simple-oauth2/models"
)

type Tokens interface {
	Save(ctx context.Context, token models.Token) (models.Token, error)
	FindByToken(ctx context.Context, tokenStr string, tokenType string) (resp models.Token, err error)
}
