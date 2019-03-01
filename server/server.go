package server

import (
	"net/http"
	"time"
)

func New(svc *Service) *http.Server {

	mux := http.NewServeMux()
	mux.HandleFunc("/lookup", svc.Handler)

	return &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
}
