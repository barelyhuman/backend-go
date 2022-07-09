package server

import (
	"log"
	"net/http"

	"github.com/barelyhuman/tasks/storage"
	"github.com/barelyhuman/tasks/views"
	"github.com/gorilla/csrf"
	"github.com/julienschmidt/httprouter"
)

func bail(err error) {
	if err != nil {
		log.Panic(err)
	}
}

type Server struct {
	Views   *views.ViewEngine
	Router  *httprouter.Router
	Storage *storage.DBClient
}

func NewServer() *Server {
	s := &Server{}
	s.Views = &views.ViewEngine{}
	err := s.Views.GetTemplates()
	bail(err)

	s.Router = httprouter.New()
	s.Storage = storage.NewDB()

	return s
}

func (s *Server) CreateHandler(f func(*Server, http.ResponseWriter, *http.Request, httprouter.Params)) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		f(s, w, r, p)
	}
}

func (s *Server) StartServer() {
	CSRF := csrf.Protect([]byte("phuWGkHK1heMDiIK"), csrf.Secure(false))
	log.Println("Started listening on port :8080")
	log.Fatal(http.ListenAndServe(":8080", CSRF(s.Router)))
}
