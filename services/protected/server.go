package protected

import (
	"context"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
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
	setupCORS(&w, r)

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
	setupCORS(&w, r)

	simplehttp.JSONRender(w, ResourceResponse{Name: "Protected Resource", Description: "This data has been protected by OAuth 2.0"})
}

const (
	BASIC_SCHEMA  string = "Basic "
	BEARER_SCHEMA string = "Bearer "
	ACCESS_TOKEN  string = "access_token"
)

func (s *Server) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		setupCORS(&w, r)

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		authorization := r.Header.Get("Authorization")
		if len(authorization) == 0 {
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = io.WriteString(w, `{"error":"invalid_key"}`)
			return
		}

		if len(strings.Split(authorization, " ")) < 2 {
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

func setupCORS(w *http.ResponseWriter, r *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}
