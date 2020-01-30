package server

import (
	"net/http"

	"github.com/dagowa/adservice/internal/chi_utils"
	"github.com/go-chi/render"
)

type server struct {
}


func (s *server) ListAdverts(w http.ResponseWriter, r *http.Request) {
	render.Respond(w, r, chi_utils.NotImplementedError())
}

func (s *server) AddAdvert(w http.ResponseWriter, r *http.Request) {
	render.Respond(w, r, chi_utils.NotImplementedError())
}

func (s *server) GetAdvert(w http.ResponseWriter, r *http.Request) {
	render.Respond(w, r, chi_utils.NotImplementedError())
}
