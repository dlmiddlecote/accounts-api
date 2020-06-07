package transport

import (
	"net/http"
	"time"
)

func (s *server) log(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		v, ok := r.Context().Value(KeyValues).(*Values)
		if !ok {
			return
		}

		s.logger.Infow("request",
			"request_id", v.TraceID,
			"method", v.Method,
			"path", v.RequestPath,
			"status", v.StatusCode,
			"duration", time.Since(v.Now),
		)
	}
}
