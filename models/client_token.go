package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ClientToken struct {
	ID           primitive.ObjectID `json:"id,omitempty" bson:"id,omitempty"`
	AccessToken  string             `json:"access_token,omitempty" bson:"access_token,omitempty"`
	RefreshToken string             `json:"refresh_token,omitempty" bson:"refresh_token,omitempty"`
	Scope        string             `json:"scope,omitempty" bson:"scope,omitempty"`
}
