package middleware

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/go-chi/render"
)

// Pagination is ...
func (m *Middleware) Pagination(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		q := r.URL
		params, err := url.ParseQuery(q.RawQuery)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var pnumb int
		if params.Get("p") != "" {
			pnumb, err = strconv.Atoi(params.Get("p"))
			if err != nil {
				render.Status(r, http.StatusBadRequest)
			}
		}
		ctx := context.WithValue(r.Context(), "pnumb", pnumb)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
