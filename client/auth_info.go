package client

// Authentication data structure
type Authentication struct {
	AccesToken   string
	Scope        string
	RefreshToken string
}

var authInfo = Authentication{}
