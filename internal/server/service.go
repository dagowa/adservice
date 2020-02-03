package server

import (
	"context"
	"net/http"
	"strconv"

	"github.com/dagowa/adservice/internal/controllers"

	"github.com/dagowa/adservice/middleware"

	"github.com/dagowa/adservice/pkg/logger"

	httplogger "github.com/dagowa/adservice/pkg/logger/http"
	"github.com/go-chi/chi"
	chimiddleware "github.com/go-chi/chi/middleware"
)

type server struct {
	controllers *controllers.Controllers
}

// NewServer is ...
func NewServer(ctx context.Context, cfg *Config, c *controllers.Controllers) *http.Server {
	r := chi.NewRouter()
	httpSrv := &http.Server{Addr: cfg.Host + ":" + strconv.Itoa(cfg.Port), Handler: r}

	srv := &server{
		controllers: c,
	}

	md := middleware.Middleware{ConnPool: c.ConnPool}

	{
		r.Use(httplogger.NewHandler(*logger.Ctx(ctx)))
		r.Use(md.ElapsedTime)
		r.Use(httplogger.RequestIDHandler("id_request", "X-Request-ID"))
		r.Use(httplogger.Recoverer)
		if cfg.LogRequests {
			r.Use(httplogger.RequestBody)
		}
		r.With(md.Pagination).Get("/adverts", srv.ListAdverts)
		r.With(md.FiledsValidation).Post("/adverts", srv.AddAdvert)
		r.With(md.AdditionalFileds).Get("/advert/{id}", srv.GetAdvert)

		r.Route("/internal", func(r chi.Router) {
			r.Mount("/debug", chimiddleware.Profiler())
		})
	}

	return httpSrv
}
