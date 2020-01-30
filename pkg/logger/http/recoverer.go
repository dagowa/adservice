package httplogger

import (
	"errors"
	"net/http"
	"runtime/debug"

	"github.com/go-chi/render"

	"github.com/dagowa/adservice/internal/chi_utils"
	"github.com/dagowa/adservice/pkg/logger"
)

var err = errors.New(http.StatusText(http.StatusInternalServerError))

func Recoverer(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if info := recover(); info != nil {
				logger.Ctx(r.Context()).Error().Interface("recover_info", info).Bytes("debug_stack", debug.Stack()).Msg("panic_on_request")
				if err := render.Render(w, r, chi_utils.InternalServerError(err)); err != nil {
					logger.Ctx(r.Context()).Error().Msg("Can't render an error")
				}
			}
		}()
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
