package protected

import (
	"html/template"
	"log"
	"net/http"
)

type Server struct {
	hostPort string
}

type ConfigOptions struct {
	ProtectedHostPort string
}

func NewServer(hostPort string) *Server {
	return &Server{
		hostPort: hostPort,
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
	t, err := template.ParseFiles("protected/templates/index.html")
	if err != nil {
		log.Fatal(err)
	}
	_ = t.Execute(w, nil)
}
