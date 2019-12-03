package authorization

import "errors"

type clientinfo struct {
	clients map[string]*Client
}

type Client struct {
	ClientID     string
	ClientSecret string
	RedirectUris []string
	Scope        []string
}

func newClientInfo() *clientinfo {
	return &clientinfo{
		clients: map[string]*Client{
			"oauth-client-1": {
				ClientID:     "oauth-client-1",
				ClientSecret: "oauth-client-secret-1",
				RedirectUris: []string{"http://localhost:8081/callback"},
				Scope:        []string{"foo", "bar"},
			},
			"oauth-client-2": {
				ClientID:     "oauth-client-2",
				ClientSecret: "oauth-client-secret-2",
				RedirectUris: []string{"http://localhost:8081/callback"},
				Scope:        []string{"foo", "bar"},
			},
		},
	}
}

func (ci *clientinfo) FindAll() (clients []Client, err error) {
	for _, val := range ci.clients {
		cl := Client{
			ClientID:     val.ClientID,
			ClientSecret: val.ClientSecret,
			RedirectUris: val.RedirectUris,
			Scope:        val.Scope,
		}
		clients = append(clients, cl)
	}
	return
}

func (ci *clientinfo) FindByID(clientid string) (client Client, err error) {

	for _, val := range ci.clients {
		if val.ClientID == clientid {
			client = Client{
				ClientID:     val.ClientID,
				ClientSecret: val.ClientSecret,
				RedirectUris: val.RedirectUris,
				Scope:        val.Scope,
			}
			return
		}
	}

	err = errors.New("sql: no rows in result set")
	return
}
