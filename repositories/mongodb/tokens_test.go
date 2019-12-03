package mongodb

import (
	"context"
	"testing"
	"time"

	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/zeihanaulia/simple-oauth2/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Test_tokens_Save(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	db, _ := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))

	type fields struct {
		db *mongo.Client
	}
	type args struct {
		ctx   context.Context
		token models.Token
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "create a new token",
			fields: fields{
				db: db,
			},
			args: args{
				ctx: ctx,
				token: models.Token{
					Type:      "access_token",
					ClientID:  "oauth-client-1",
					Token:     "CJMGZAveWBqWMAbNFx8s3Z46fG7MJkuf",
					Scope:     "[openid purchase]",
					IssuedAt:  time.Now().UnixNano() / int64(time.Millisecond),
					ExpiredAt: time.Now().Add(3600*time.Second).UnixNano() / int64(time.Millisecond),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tkn := &tokens{
				db: tt.fields.db,
			}
			gotResp, err := tkn.Save(tt.args.ctx, tt.args.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("tokens.Save() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			oid, _ := primitive.ObjectIDFromHex("000000000000000000000000")
			if gotResp.ID == oid {
				t.Errorf("tokens.Save() = %v", gotResp)
			}
		})
	}
}

func Test_tokens_FindByToken(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	db, _ := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))

	type fields struct {
		db *mongo.Client
	}
	type args struct {
		ctx       context.Context
		tokenStr  string
		tokenType string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "find by token",
			fields: fields{
				db: db,
			},
			args: args{
				ctx:       ctx,
				tokenStr:  "CJMGZAveWBqWMAbNFx8s3Z46fG7MJkuf",
				tokenType: "access_token",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tkn := &tokens{
				db: tt.fields.db,
			}
			_, err := tkn.FindByToken(tt.args.ctx, tt.args.tokenStr, tt.args.tokenType)
			if (err != nil) != tt.wantErr {
				t.Errorf("tokens.FindByToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
