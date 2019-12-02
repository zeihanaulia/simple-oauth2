package mongodb

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/zeihanaulia/simple-oauth2/models"
	"github.com/zeihanaulia/simple-oauth2/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type tokens struct {
	db *mongo.Client
}

func NewTokens(db *mongo.Client) repositories.Tokens {
	return &tokens{db: db}
}

func (t *tokens) Save(ctx context.Context, token models.Token) (resp models.Token, err error) {
	collection := t.db.Database("oauth").Collection("token")
	result, _ := collection.InsertOne(ctx, token)
	oid, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		err = errors.New("invalid get id")
		return
	}
	token.ID = oid
	resp = token
	return
}

func (t *tokens) FindByToken(ctx context.Context, tokenStr string, tokenType string) (resp models.Token, err error) {
	collection := t.db.Database("oauth").Collection("token")

	var token models.Token
	if err = collection.FindOne(ctx, bson.M{"token": tokenStr, "type": tokenType}).Decode(&token); err != nil {
		return
	}

	resp = token
	return
}
