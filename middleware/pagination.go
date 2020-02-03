package middleware

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/dagowa/adservice/internal/models/page"
	"github.com/go-chi/render"
)

func (m *Middleware) Pagination(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		q := r.URL
		params, err := url.ParseQuery(q.RawQuery)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		p := &page.Page{}
		if params.Get("p") == "" {
			render.Status(r, http.StatusBadRequest)
		}
		p.Numb, err = strconv.Atoi(params.Get("p"))
		if err != nil {
			render.Status(r, http.StatusBadRequest)
		}

		if params.Get("psize") == "" {
			render.Status(r, http.StatusBadRequest)
		}
		p.Size, err = strconv.Atoi(params.Get("psize"))
		if err != nil {
			render.Status(r, http.StatusBadRequest)
		}

		if params.Get("price") == "" {
			render.Status(r, http.StatusBadRequest)
		}
		priceSortType := params.Get("price")
		switch priceSortType {
		case page.SortOrderASC:
			p.PriceAsc = true
		case page.SortOrderDESC:
			p.PriceAsc = false
		default:
			render.Status(r, http.StatusBadRequest)
		}

		if params.Get("date") == "" {
			render.Status(r, http.StatusBadRequest)
		}
		dateSortType := params.Get("date")
		switch dateSortType {
		case page.SortOrderASC:
			p.DateAsc = true
		case page.SortOrderDESC:
			p.DateAsc = false
		default:
			render.Status(r, http.StatusBadRequest)
		}

		ctx := context.WithValue(r.Context(), "page", p)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
