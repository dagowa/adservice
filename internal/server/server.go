package server

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/rs/zerolog/hlog"

	"github.com/dagowa/adservice/internal/chi_utils"
	"github.com/dagowa/adservice/internal/controllers/advertmanager"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type server struct {
	service *service
}

func (s *server) ListAdverts(w http.ResponseWriter, r *http.Request) {
	logger := hlog.FromRequest(r)

	am := s.service.AdvManager

	sortCriteria := advertmanager.SortCriteria{}
	sortCriteria.SetDateSortOrder(advertmanager.SortOrderASC)

	ctx := r.Context()
	pnumb := ctx.Value("pnumb")

	adverts, err := am.GetBatch(&sortCriteria, pnumb.(int), 10)
	if err != nil {
		logger.Error().Err(err).Msg("Cannot get batch of adverts")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	render.JSON(w, r, adverts)
}

func (s *server) AddAdvert(w http.ResponseWriter, r *http.Request) {
	render.Respond(w, r, chi_utils.NotImplementedError())
}

func (s *server) GetAdvert(w http.ResponseWriter, r *http.Request) {
	logger := hlog.FromRequest(r)

	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	q := r.URL
	params, err := url.ParseQuery(q.RawQuery)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	rf := advertmanager.RequirementFileds{}
	if params.Get("date") != "" {
		rf.Date = true
	}
	if params.Get("descr") != "" {
		rf.Description = true
	}
	if params.Get("gallery") != "" {
		rf.Gallery = true
	}

	am := s.service.AdvManager
	adverts, err := am.GetOne(id, &rf)
	if err != nil {
		logger.Error().Err(err).Msg("Cannot get advert")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	render.JSON(w, r, adverts)
}
