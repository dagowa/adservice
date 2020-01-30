package httplogger

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/render"

	"github.com/dagowa/adservice/internal/chi_utils"
	"github.com/dagowa/adservice/pkg/logger"
)

func RequestBody(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		l := logger.Ctx(r.Context())

		body := r.Body
		var buf []byte
		if body != nil {

			b, err := ioutil.ReadAll(r.Body)
			if err != nil {
				l.Error().Err(err).Msg("Cannot read out a request body")
				if err := render.Render(w, r, chi_utils.InvalidRequest(err)); err != nil {
					logger.Ctx(r.Context()).Error().Msg("Can't render an error")
				}
				return
			}
			buf = b
			r.Body = ioutil.NopCloser(bytes.NewBuffer(buf))
		}

		l.Debug().Str("method", r.Method).Str("url", r.URL.String()).Interface("headers", r.Header).Bytes("request_body", buf).Msg("The incoming request")

		next.ServeHTTP(w, r)

		r.Body = body
	}
	return http.HandlerFunc(fn)
}
