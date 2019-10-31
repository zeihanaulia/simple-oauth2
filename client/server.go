package client

import (
	"html/template"
	"log"
	"net/http"
)

type Server struct {
	hostPort string
}

type ConfigOptions struct {
	ClientHostPort string
}

func NewServer(options ConfigOptions) *Server {
	return &Server{
		hostPort: options.ClientHostPort,
	}
}

func (s *Server) Run() error {
	mux := s.createServerMux()
	return http.ListenAndServe(s.hostPort, mux)
}

func (s *Server) createServerMux() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.index)
	return mux
}

func (s *Server) index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("client/templates/index.html")
	if err != nil {
		log.Fatal(err)
	}
	_ = t.Execute(w, authInfo)
}
