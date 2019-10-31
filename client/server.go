package client

import (
	"html/template"
	"log"
	"net/http"
)

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
	return mux
}

func (s *Server) index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("client/templates/index.html")
	if err != nil {
		log.Fatal(err)
	}
	_ = t.Execute(w, authInfo)
}
