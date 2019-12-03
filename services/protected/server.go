package protected

import (
	"context"
	"html/template"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/zeihanaulia/simple-oauth2/repositories"

	"github.com/zeihanaulia/simple-oauth2/pkg/simplehttp"
)

// Server protected service
type Server struct {
	hostPort string
	tokens   repositories.Tokens
}

// NewServer like constructor for inject dependency
func NewServer(hostPort string, tokens repositories.Tokens) *Server {
	return &Server{
		hostPort: hostPort,
		tokens:   tokens,
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
	mux.Handle("/resource", simplehttp.Middleware(
		http.HandlerFunc(s.resource),
		s.AuthMiddleware,
	))
	return mux
}

func (s *Server) index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("protected/templates/index.html")
	if err != nil {
		log.Fatal(err)
	}
	_ = t.Execute(w, nil)
}

type ResourceResponse struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (s *Server) resource(w http.ResponseWriter, r *http.Request) {
	simplehttp.JSONRender(w, ResourceResponse{Name: "Protected Resource", Description: "This data has been protected by OAuth 2.0"})
}

const (
	BASIC_SCHEMA  string = "Basic "
	BEARER_SCHEMA string = "Bearer "
	ACCESS_TOKEN  string = "access_token"
)

func (s *Server) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")
		if len(authorization) == 0 {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = io.WriteString(w, `{"error":"invalid_key"}`)
			return
		}

		act := authorization[len(BEARER_SCHEMA):]

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		_, err := s.tokens.FindByToken(ctx, act, ACCESS_TOKEN)
		if err != nil {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = io.WriteString(w, `{"error":"unauthorize"}`)
			return
		}

		next.ServeHTTP(w, r)
	})
}
