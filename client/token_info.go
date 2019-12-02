package client

// Token data structure
type Token struct {
	AccessToken  string
	Scope        string
	RefreshToken string
}

var tokenInfo = Token{}
