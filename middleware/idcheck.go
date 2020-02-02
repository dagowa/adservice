package middleware

import (
	"net/http"
	"strconv"

	"github.com/dagowa/adservice/internal/controllers/advertmanager"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/jackc/pgx"
)

// AdvertID designed to check advert existance
func (m *Middleware) AdvertID(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			render.Status(r, http.StatusBadRequest)
		}

		am := advertmanager.AdvertManager{
			ConnPool: m.ConnPool,
		}
		if err := am.IsExist(id); err != nil {
			if err == pgx.ErrNoRows {
				render.Status(r, http.StatusNotFound)
			}
			render.Status(r, http.StatusInternalServerError)
		}

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
