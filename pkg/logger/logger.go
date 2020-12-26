package logger

import (
	"context"
	"net/http"
	"time"

	"github.com/friends/configs"
	log "github.com/sirupsen/logrus"
)

func AccessLog(r *http.Request, start time.Time) {
	log.WithFields(log.Fields{
		configs.ReqID:  r.Context().Value(configs.ReqID),
		"method":       r.Method,
		"remote_addr":  r.RemoteAddr,
		"req_duration": time.Since(start),
	}).Info(r.URL.Path)
}

func DataLog(r *http.Request, data interface{}) {
	log.WithFields(log.Fields{
		configs.ReqID: r.Context().Value(configs.ReqID),
		"data":        data,
	}).Debug()
}

func ErrorLogWithCtx(ctx context.Context, err error) {
	log.WithFields(log.Fields{
		configs.ReqID: ctx.Value(configs.ReqID),
	}).Error(err)
}

func ErrorMessage(msg string) {
	log.Error(msg)
}
