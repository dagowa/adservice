package server

import (
	"context"
	"net/http"
	"strconv"

	"github.com/dagowa/adservice/internal/controllers/advertmanager"
	"github.com/dagowa/adservice/internal/storage"
	"github.com/dagowa/adservice/middleware"
	"github.com/rs/zerolog"

	"github.com/dagowa/adservice/pkg/logger"

	httplogger "github.com/dagowa/adservice/pkg/logger/http"
	"github.com/go-chi/chi"
	chimiddleware "github.com/go-chi/chi/middleware"
)

type service struct {
	AdvManager *advertmanager.AdvertManager
}

// NewServer is ...
func NewServer(ctx context.Context, cfg *Config) *http.Server {
	r := chi.NewRouter()
	httpSrv := &http.Server{Addr: cfg.Host + ":" + strconv.Itoa(cfg.Port), Handler: r}

	l := zerolog.Ctx(ctx)

	psqlConn, err := storage.New().NewPostgreSQLConn(cfg.PSQLConfig)
	if err != nil {
		l.Fatal().Err(err).Msg("Cannot set up psql connection")
	}
	//TODO: ой блять, наверное надо отсюда его вытащить -_-
	defer psqlConn.Close()

	service := &service{
		AdvManager: &advertmanager.AdvertManager{ConnPool: psqlConn.Pool()},
	}
	srv := &server{
		service: service,
	}
	md := middleware.Middleware{ConnPool: psqlConn.Pool()}

	{
		r.Use(httplogger.NewHandler(*logger.Ctx(ctx)))
		r.Use(md.ElapsedTime)
		r.Use(httplogger.RequestIDHandler("id_request", "X-Request-ID"))
		r.Use(httplogger.Recoverer)
		if cfg.LogRequests {
			r.Use(httplogger.RequestBody)
		}
		r.Get("/adverts", srv.ListAdverts)
		r.Post("/adverts", srv.AddAdvert)
		r.With(md.AdvertID).Get("/advert/{id}", srv.GetAdvert)

		r.Route("/internal", func(r chi.Router) {
			r.Mount("/debug", chimiddleware.Profiler())
		})
	}

	return httpSrv
}
