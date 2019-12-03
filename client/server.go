package client

import (
	"net/http"

	"github.com/zeihanaulia/simple-oauth2/pkg/simplehttp"
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
	simplehttp.HTMLRender(w, "client/templates/index.html", tokenInfo)
}
