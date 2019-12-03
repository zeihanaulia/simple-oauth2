package authorization

import (
	"context"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/zeihanaulia/simple-oauth2/models"

	"github.com/zeihanaulia/simple-oauth2/repositories"

	"github.com/zeihanaulia/simple-oauth2/pkg/simpleurl"

	"github.com/zeihanaulia/simple-oauth2/pkg/simplestr"

	"github.com/zeihanaulia/simple-oauth2/pkg/randomstring"

	"github.com/zeihanaulia/simple-oauth2/pkg/simplehttp"
)

var requests = make(map[string]interface{}, 0)
var codes = make(map[string]interface{}, 0)

// Server authorization service
type Server struct {
	hostPort     string
	templatePath string
	clientinfo   *clientinfo
	authzinfo    *authzinfo
	tokens       repositories.Tokens
}

// NewServer like constructor for inject dependency
func NewServer(hostPort, templatePath string, tokens repositories.Tokens) *Server {
	return &Server{
		hostPort:     hostPort,
		templatePath: templatePath,
		clientinfo:   newClientInfo(),
		authzinfo:    newAuthzInfo(),
		tokens:       tokens,
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
	mux.HandleFunc("/approve", s.approve)
	mux.HandleFunc("/token", s.token)
	return mux
}

type indexResponse struct {
	Clients    []Client
	AuthServer Authz
}

func (s *Server) index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(s.templatePath + "index.html")
	if err != nil {
		log.Fatal(err)
	}

	clients, _ := s.clientinfo.FindAll()
	authzinfo, _ := s.authzinfo.Find()

	_ = t.Execute(w, indexResponse{Clients: clients, AuthServer: authzinfo})
}

type authzResponse struct {
	RequestID string
	Client
}

func (s *Server) authorize(w http.ResponseWriter, r *http.Request) {

	clientID, ok := r.URL.Query()["client_id"]
	if !ok || len(clientID[0]) <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		simplehttp.HTMLRender(w, s.templatePath+"error.html", nil)
		return
	}

	redirectURI, ok := r.URL.Query()["redirect_uri"]
	if !ok || len(redirectURI[0]) <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		simplehttp.HTMLRender(w, s.templatePath+"error.html", nil)
		return
	}

	// Pertama, kita perlu mencari tahu klien mana yang membuat permintaan.
	client, err := s.clientinfo.FindByID(clientID[0])
	if err != nil {
		// Unknown client
		w.WriteHeader(http.StatusBadRequest)
		simplehttp.HTMLRender(w, s.templatePath+"error.html", nil)
		return
	}

	// Setelah kita tahu klien mana yang meminta, kita perlu melakukan pemeriksaan atas permintaan itu
	// Pada titik ini, satu-satunya hal yang diteruskan melalui browser adalah client_id
	// dan karena dilewatkan melalui browser, ini dianggap informasi publik.
	// Pada titik ini, siapa pun bisa berpura-pura menjadi klien ini
	// tapi kita kan punya beberapa hal yang dapat membantu kami memastikan itu permintaan yang sah,
	// kita dapat melakukan pemeriksaan ke redirect uri yang telah didaftarkan
	if !simplestr.Contains(client.RedirectUris, redirectURI[0]) {
		w.WriteHeader(http.StatusBadRequest)
		simplehttp.HTMLRender(w, s.templatePath+"error.html", nil)
		return
	}

	// Akhirnya, jika klien kami berhasil, kami perlu membuat halaman untuk meminta izin kepada pengguna

	requestID := randomstring.Generator(8)

	// kami memberikan requestID agar nantinya kita bisa mengambil data dari r.URL.Query() lagi untuk dicocokan
	// dalam skala production kita bisa memasukan diserver side storage
	// requestID akan dimasukan sebagai hidden value, atau lebih dikenal sebagai CSRF protection
	requests[requestID] = r.URL.Query()

	// lalu kita tampikan halaman approval, agar resource owner yang memutuskan apakah ingin mengizinkan atau tidak
	simplehttp.HTMLRender(w, s.templatePath+"approve.html", authzResponse{RequestID: requestID, Client: client})
}

type Code struct {
	AuthorizationEndpointRequest interface{}
	Scope                        string
	User                         string
}

func (s *Server) approve(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		simplehttp.HTMLRender(w, s.templatePath+"error.html", nil)
		return
	}

	var approve = r.FormValue("approve")
	var requestID = r.FormValue("reqid")
	var user = r.FormValue("user")

	// Jika kami tidak menemukan permintaan yang tertunda untuk kode ini,
	// kemungkinan itu adalah serangan pemalsuan lintas situs
	// dan kami dapat mengirim pengguna ke halaman kesalahan.
	var que = requests[requestID]
	if que == nil {
		w.WriteHeader(http.StatusBadRequest)
		simplehttp.HTMLRender(w, s.templatePath+"error.html", nil)
		return
	}
	delete(requests, requestID) // remove becouse it's has been used

	query := que.(url.Values)

	var redirectURI = query["redirect_uri"][0]
	var state = query["state"][0]

	if approve == "" {
		// deny
		w.WriteHeader(http.StatusBadRequest)
		simplehttp.HTMLRender(w, s.templatePath+"error.html", nil)
		return
	}

	if query["response_type"][0] != "code" {
		// unsupported_response_type
		w.WriteHeader(http.StatusBadRequest)
		simplehttp.HTMLRender(w, s.templatePath+"error.html", nil)
		return
	}

	code := randomstring.Generator(8)

	codes[code] = Code{
		AuthorizationEndpointRequest: query,
		Scope:                        "",
		User:                         user,
	}

	// don't cache approve redirect
	w.Header().Set("Cache-Control", "no-cache, private, max-age=0")
	w.Header().Set("Expires", time.Unix(0, 0).Format(http.TimeFormat))
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("X-Accel-Expires", "0")

	http.Redirect(w, r, simpleurl.Builder(redirectURI, map[string]string{
		"code":  code,
		"state": state,
	}), 301)
}

type HTTPResponse struct {
	Error string
}

type TokenRequest struct {
	GrantType   string `json:"grant_type,omitempty"`
	Code        string `json:"code,omitempty"`
	RedirectUri string `json:"redirect_uri,omitempty"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

func (s *Server) token(w http.ResponseWriter, r *http.Request) {
	ua := r.Header.Get("Authorization")
	if ua == "" {
		w.WriteHeader(http.StatusUnauthorized)
		simplehttp.HTMLRender(w, s.templatePath+"error.html", HTTPResponse{Error: "unauthorized"})
		return
	}

	clientID, clientSecret, authOK := r.BasicAuth()
	if !authOK {
		w.WriteHeader(http.StatusUnauthorized)
		simplehttp.HTMLRender(w, s.templatePath+"error.html", HTTPResponse{Error: "unauthorized"})
		return
	}

	client, err := s.clientinfo.FindByID(clientID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		simplehttp.HTMLRender(w, s.templatePath+"error.html", HTTPResponse{Error: "unauthorized"})
		return
	}

	if client.ClientSecret != clientSecret {
		w.WriteHeader(http.StatusUnauthorized)
		simplehttp.HTMLRender(w, s.templatePath+"error.html", HTTPResponse{Error: "unauthorized"})
		return
	}

	if r.Body != nil {
		defer r.Body.Close()
	}

	var tr TokenRequest
	if err := json.NewDecoder(r.Body).Decode(&tr); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		simplehttp.HTMLRender(w, s.templatePath+"error.html", HTTPResponse{Error: "unauthorized"})
		return
	}

	if tr.GrantType != "authorization_code" {
		w.WriteHeader(http.StatusUnauthorized)
		simplehttp.HTMLRender(w, s.templatePath+"error.html", HTTPResponse{Error: "unauthorized"})
		return
	}

	var code = codes[tr.Code]

	if code != nil {
		delete(codes, tr.Code) // remove becouse it's has been used

		c, ok := code.(Code)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			simplehttp.HTMLRender(w, s.templatePath+"error.html", HTTPResponse{Error: "unauthorized"})
			return
		}

		d, ok := c.AuthorizationEndpointRequest.(url.Values)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			simplehttp.HTMLRender(w, s.templatePath+"error.html", HTTPResponse{Error: "unauthorized"})
			return
		}

		if d.Get("client_id") == clientID {
			accessToken := randomstring.Generator(32)
			refreshToken := randomstring.Generator(32)

			// TODO: save access token and refresh token to persistence storage with client id
			if err := s.saveToken(clientID, accessToken, refreshToken); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				simplehttp.HTMLRender(w, s.templatePath+"error.html", HTTPResponse{Error: err.Error()})
				return
			}

			simplehttp.JSONRender(w, TokenResponse{AccessToken: accessToken, RefreshToken: refreshToken, Scope: ""})
			return
		}
	}
}

func (s *Server) saveToken(clientID, accessToken, refreshToken string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	issueAt := time.Now().UnixNano() / int64(time.Millisecond)
	accessTokenExp := time.Now().Add(3600*time.Second).UnixNano() / int64(time.Millisecond)
	refreshTokenExp := time.Now().Add(48*time.Hour).UnixNano() / int64(time.Millisecond)

	if _, err = s.tokens.Save(ctx, models.Token{
		ClientID:  clientID,
		Type:      "access_token",
		Token:     accessToken,
		IssuedAt:  issueAt,
		ExpiredAt: accessTokenExp,
	}); err != nil {
		return
	}

	if _, err = s.tokens.Save(ctx, models.Token{
		ClientID:  clientID,
		Type:      "refresh_token",
		Token:     refreshToken,
		IssuedAt:  issueAt,
		ExpiredAt: refreshTokenExp,
	}); err != nil {
		return
	}

	return
}
