package authorization

type client struct {
	ClientID     string
	ClientSecret string
	RedirectUris []string
	Scope        string
}

type authServer struct {
	AuthorizationEndpoint string
	TokenEndpoint         string
}

type serverResponse struct {
	Clients    []client
	AuthServer authServer
}

var clientInfo = serverResponse{
	Clients: []client{
		client{
			ClientID:     "oauth-client-1",
			ClientSecret: "oauth-client-secret-1",
			RedirectUris: []string{"http://localhost:8081/callback"},
			Scope:        "foo bar",
		},
		client{
			ClientID:     "oauth-client-2",
			ClientSecret: "oauth-client-secret-2",
			RedirectUris: []string{"http://localhost:8081/callback"},
			Scope:        "foo bar",
		},
	},
	AuthServer: authServer{
		AuthorizationEndpoint: "http://localhost:8081/authorize",
		TokenEndpoint:         "http://localhost:8081/token",
	},
}
