package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ppai-plivo/mnlookup/api"
)

type LookupService interface {
	Lookup(number string) (*api.Record, error)
}

type Service struct {
	l LookupService
}

func (s *Service) Handler(w http.ResponseWriter, r *http.Request) {

	var number string
	if v, ok := r.URL.Query()["number"]; ok {
		number = v[0]
	}

	if number == "" {
		http.NotFound(w, r)
	}

	record, err := s.l.Lookup(number)
	if err != nil {
		log.Println(err)
		http.NotFound(w, r)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	enc := json.NewEncoder(w)
	if err := enc.Encode(record); err != nil {
		log.Println(err)
	}
}

func NewService(l LookupService) *Service {
	return &Service{
		l: l,
	}
}
