package server

import (
	"net/http"
	_ "net/http/pprof"
	"time"
)

func New(svc *Service) *http.Server {

	http.HandleFunc("/lookup", svc.Handler)

	return &http.Server{
		Addr:         ":8080",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
}
