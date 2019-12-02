package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Token struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ClientID string             `json:"client_id,omitempty" bson:"client_id,omitempty"`
	Type     string             `json:"type,omitempty" bson:"type,omitempty"`
	Token    string             `json:"token,omitempty" bson:"token,omitempty"`
	Scope    string             `json:"scope,omitempty" bson:"scope,omitempty"`

	Subject         string `json:"subject,omitempty" bson:"subject,omitempty"`
	ExpiredAt       int64  `json:"expired_at,omitempty" bson:"expired_at,omitempty"`
	Issuer          string `json:"duration,omitempty" bson:"duration,omitempty"`
	IssuedAt        int64  `json:"issued_at,omitempty" bson:"issued_at,omitempty"`
	AuthorizedParty string `json:"authorized_party,omitempty" bson:"authorized_party,omitempty"`
}
