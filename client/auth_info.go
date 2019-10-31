package client

type Authentication struct {
	AccesToken   string
	Scope        string
	RefreshToken string
}

var authInfo = Authentication{}
