package pkg

import (
	"context"
	"net/http"

	"github.com/friends/configs"
	log "github.com/sirupsen/logrus"
)

func AccessLog(r *http.Request) {
	log.WithFields(log.Fields{
		configs.ReqID: r.Context().Value(configs.ReqID),
		"method":      r.Method,
		"remote_addr": r.RemoteAddr,
	}).Info(r.URL.Path)
}

func ErrorLogWithCtx(ctx context.Context, err error) {
	log.WithFields(log.Fields{
		configs.ReqID: ctx.Value(configs.ReqID),
	}).Error(err)
}

func ErrorMessage(msg string) {
	log.Error(msg)
}
