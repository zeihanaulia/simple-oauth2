package client

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"

	"github.com/zeihanaulia/simple-oauth2/pkg/simplestr"

	"github.com/zeihanaulia/simple-oauth2/pkg/simplehttp"

	"github.com/zeihanaulia/simple-oauth2/pkg/randomstring"

	"github.com/zeihanaulia/simple-oauth2/pkg/simpleurl"
)

var states = make([]string, 0)
var clientsToken = make(map[string]interface{}, 0)

// Server client service
type Server struct {
	hostPort string
}

// NewServer like constructor for inject dependency
func NewServer(hostPort string) *Server {
	return &Server{
		hostPort: hostPort,
	}
}

// Run just wrap for running server
func (s *Server) Run() error {
	mux := s.createServerMux()
	return http.ListenAndServe(s.hostPort, mux)
}

func (s *Server) createServerMux() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.index)
	mux.HandleFunc("/authorize", s.authorize)
	mux.HandleFunc("/callback", s.callback)
	mux.HandleFunc("/fetch_resource", s.fetch_resource)
	return mux
}

func (s *Server) index(w http.ResponseWriter, r *http.Request) {
	simplehttp.HTMLRender(w, "client/templates/index.html", tokenInfo)
}

func (s *Server) authorize(w http.ResponseWriter, r *http.Request) {
	var state = randomstring.Generator(32)
	states = append(states, state) // save state for comparing in callback request

	// don't cache authorize request
	w.Header().Set("Cache-Control", "no-cache, private, max-age=0")
	w.Header().Set("Expires", time.Unix(0, 0).Format(http.TimeFormat))
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("X-Accel-Expires", "0")

	http.Redirect(w, r, simpleurl.Builder(authServerInfo.AuthorizationEndpoint, map[string]string{
		"response_type": "code",
		"client_id":     "oauth-client-1",
		"redirect_uri":  "http://localhost:8081/callback",
		"state":         state,
	}), 301)
}

type HTTPResponse struct {
	Error string
}

type CallbackResponse struct {
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Scope        string `json:"scope,omitempty"`
}

func (s *Server) callback(w http.ResponseWriter, r *http.Request) {

	state, ok := r.URL.Query()["state"]
	if !ok || len(state[0]) <= 1 {
		w.WriteHeader(http.StatusBadRequest)
		simplehttp.HTMLRender(w, "client/templates/error.html", HTTPResponse{Error: "Unable to fetch access token, server response: 400"})
		return
	}

	if !simplestr.Contains(states, state[0]) {
		w.WriteHeader(http.StatusBadRequest)
		simplehttp.HTMLRender(w, "client/templates/error.html", HTTPResponse{Error: "Unable to fetch access token, server response: 400"})
		return
	}
	deleteState(states, state[0]) // remove state

	code, ok := r.URL.Query()["code"]
	if !ok || len(code[0]) <= 1 {
		w.WriteHeader(http.StatusBadRequest)
		simplehttp.HTMLRender(w, "client/templates/error.html", HTTPResponse{Error: "Unable to fetch access token, server response: 400"})
		return
	}

	requestBody, _ := json.Marshal(map[string]string{
		"grant_type":   "authorization_code",
		"code":         code[0],
		"redirect_uri": "http://localhost:8081/callback",
	})

	basicAuthEnc := base64.URLEncoding.EncodeToString([]byte(authServerInfo.ClientID + ":" + authServerInfo.ClientSecret))
	headers := map[string]string{
		"Content-type":  "application/json",
		"Authorization": "Basic " + basicAuthEnc,
	}

	resp, err := simplehttp.Post(authServerInfo.TokenEndpoint, headers, requestBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		simplehttp.HTMLRender(w, "client/templates/error.html", HTTPResponse{Error: "Unable to fetch access token, server response: 400"})
		return
	}

	var cr CallbackResponse
	_ = json.Unmarshal(resp, &cr)

	clientsToken[authServerInfo.ClientID] = map[string]interface{}{
		"access_token":  cr.AccessToken,
		"refresh_token": cr.RefreshToken,
		"scope":         cr.Scope,
	}

	simplehttp.HTMLRender(w, "client/templates/index.html", Token{
		AccessToken:  cr.AccessToken,
		RefreshToken: cr.RefreshToken,
		Scope:        cr.Scope,
	})
}

func deleteState(states []string, accessToken string) []string {
	var key int
	for k, v := range states {
		if v == accessToken {
			key = k
		}
	}
	states[key] = states[len(states)-1]
	states[len(states)-1] = ""
	states = states[:len(states)-1]
	return states
}

type ResourceResponse struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

func (s *Server) fetch_resource(w http.ResponseWriter, r *http.Request) {

	token := clientsToken[authServerInfo.ClientID].(map[string]interface{})

	headers := map[string]string{
		"Content-type": "application/json",
		"Authorization": "Bearer " + token["access_token"].(string),
	}

	resp, err := simplehttp.Post("http://localhost:8083/resource", headers, []byte(``))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		simplehttp.HTMLRender(w, "client/templates/error.html", HTTPResponse{Error: "Unable to fetch access token, server response: 400"})
		return
	}

	var rr ResourceResponse
	_ = json.Unmarshal(resp, &rr)

	simplehttp.HTMLRender(w, "client/templates/data.html", rr)
}
