package authorization

type authzinfo struct {
	authzs Authz
}

type Authz struct {
	AuthorizationEndpoint string
	TokenEndpoint         string
}

func newAuthzInfo() *authzinfo {
	return &authzinfo{
		authzs: Authz{
			AuthorizationEndpoint: "http://localhost:8081/authorize",
			TokenEndpoint:         "http://localhost:8081/token",
		},
	}
}

func (ai *authzinfo) Find() (authz Authz, err error) {
	authz = ai.authzs
	return
}
