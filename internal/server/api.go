package server

import (
	"context"
	"net/http"
	"strconv"

	"github.com/rs/zerolog"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

func (s *server) ListAdverts(w http.ResponseWriter, r *http.Request) {
	l := zerolog.Ctx(r.Context())

	advertManager := s.controllers.AdvertManager()
	page := advertManager.ParsePage(r)

	adverts, err := advertManager.GetBatch(page)
	if err != nil {
		l.Error().Err(err).Msg("Cannot get batch of adverts")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx := context.WithValue(r.Context(), render.StatusCtxKey, advertManager.HTTPStatus)
	r = r.WithContext(ctx)

	render.JSON(w, r, adverts)
}

func (s *server) AddAdvert(w http.ResponseWriter, r *http.Request) {
	l := zerolog.Ctx(r.Context())

	advertManager := s.controllers.AdvertManager()
	advert := advertManager.ParseAdvert(r)

	_, err := advertManager.AddOne(advert)
	if err != nil {
		l.Error().Err(err).Msg("Cannot add advert to database")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx := context.WithValue(r.Context(), render.StatusCtxKey, advertManager.HTTPStatus)
	r = r.WithContext(ctx)

	render.JSON(w, r, advert)
}

func (s *server) GetAdvert(w http.ResponseWriter, r *http.Request) {
	l := zerolog.Ctx(r.Context())

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		return
	}

	advertManager := s.controllers.AdvertManager()
	af := advertManager.ParseAdditionalFileds(r)

	adverts, err := advertManager.GetOne(id, af)
	if err != nil {
		l.Error().Err(err).Msg("Cannot get advert")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx := context.WithValue(r.Context(), render.StatusCtxKey, advertManager.HTTPStatus)
	r = r.WithContext(ctx)

	render.JSON(w, r, adverts)
}
