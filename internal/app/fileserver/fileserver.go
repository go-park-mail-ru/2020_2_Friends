package fileserver

import (
	"net/http"

	"github.com/friends/configs"
	"github.com/friends/internal/pkg/middleware"
	"github.com/sirupsen/logrus"
)

func StartFileServer() {
	mux := http.NewServeMux()

	staticHandler := http.StripPrefix(
		"/data/",
		http.FileServer(http.Dir("./static")),
	)
	mux.Handle("/data/", staticHandler)

	corsHandler := middleware.CORS(mux)
	siteHandler := middleware.Panic(corsHandler)

	logrus.Info("starting fileserver at port ", configs.FileServerPort)
	logrus.Fatal(http.ListenAndServe(configs.FileServerPort, siteHandler))
}
