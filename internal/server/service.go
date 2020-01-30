package server

import (
	"context"
	"net/http"
	"strconv"

	"github.com/jackc/pgx"

	"github.com/dagowa/adservice/middleware/timings"
	"github.com/dagowa/adservice/pkg/logger"

	httplogger "github.com/dagowa/adservice/pkg/logger/http"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// NewServer is ...
func NewServer(ctx context.Context, cfg *Config, p *pgx.ConnPool) *http.Server {
	r := chi.NewRouter()
	httpSrv := &http.Server{Addr: cfg.Host + ":" + strconv.Itoa(cfg.Port), Handler: r}

	srv := &server{}

	{
		r.Use(httplogger.NewHandler(*logger.Ctx(ctx)))
		r.Use(timings.ElapsedTime)
		r.Use(httplogger.RequestIDHandler("id_request", "X-Request-ID"))
		r.Use(httplogger.Recoverer)
		if cfg.LogRequests {
			r.Use(httplogger.RequestBody)
		}
		r.Get("/adverts", srv.ListAdverts)
		r.Post("/adverts", srv.AddAdvert)
		r.Get("/advert/{id}", srv.GetAdvert)

		r.Route("/internal", func(r chi.Router) {
			r.Mount("/debug", middleware.Profiler())
		})
	}

	return httpSrv
}
