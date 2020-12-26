package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/friends/configs"
	log "github.com/friends/pkg/logger"
	"github.com/lithammer/shortuuid"
)

func AccessLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := shortuuid.New()
		ctx := context.WithValue(r.Context(), configs.ReqID, reqID)
		r = r.WithContext(ctx)
		start := time.Now()
		next.ServeHTTP(w, r)
		log.AccessLog(r, start)
	})
}
