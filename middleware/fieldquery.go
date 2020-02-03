package middleware

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/go-chi/render"

	"github.com/dagowa/adservice/internal/models/advert"

	"github.com/dagowa/adservice/internal/controllers/advertmanager"
)

// AdditionalFileds gathers all requrement fileds shown in advert response
func (m *Middleware) AdditionalFileds(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		af := &advertmanager.AdditionalFileds{}

		q := r.URL
		params, err := url.ParseQuery(q.RawQuery)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			return
		}
		if params.Get("date") != "" {
			af.Date = true
		}
		if params.Get("descr") != "" {
			af.Description = true
		}
		if params.Get("gal") != "" {
			af.Gallery = true
		}
		ctx := context.WithValue(r.Context(), "additional_fields", af)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

func (m *Middleware) FiledsValidation(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		a := advert.Advert{}

		buf, err := ioutil.ReadAll(r.Body)
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			return
		}
		if err := json.Unmarshal([]byte(buf), &a); err != nil {
			render.Status(r, http.StatusBadRequest)
			return
		}

		if len(*(a.Gallery)) > 3 || len(*(a.Description)) > 1000 || len(a.Title) > 200 {
			render.Status(r, http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), "advert", a)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
