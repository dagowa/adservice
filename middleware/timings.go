package middleware

import (
	"net/http"
	"time"

	httplogger "github.com/dagowa/adservice/pkg/logger/http"
)

// ElapsedTime designed to log the request processing time
func (Middleware) ElapsedTime(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		stop := time.Now()
		httplogger.FromRequest(r).Debug().
			Str("start_time", start.Format(time.RFC3339Nano)).
			Str("stop_time", stop.Format(time.RFC3339Nano)).
			TimeDiff("duration", stop, start).
			Str("method", r.Method).
			Str("url", r.URL.String()).
			Msg("Elapsed time")
	}
	return http.HandlerFunc(fn)
}
