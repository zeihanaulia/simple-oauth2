package client

type AuthServer struct {
	AuthorizationEndpoint string
	TokenEndpoint         string
	ClientID              string
	ClientSecret          string
}

var authServerInfo = AuthServer{
	AuthorizationEndpoint: "http://localhost:8082/authorize",
	TokenEndpoint:         "http://localhost:8082/token",
	ClientID:              "oauth-client-1",
	ClientSecret:          "oauth-client-secret-1",
}
