package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/ppai-plivo/mnlookup/api"
	"github.com/ppai-plivo/mnlookup/store"

	"github.com/nyaruka/phonenumbers"
)

type LookupService interface {
	Lookup(number string) (store.Value, error)
}

type Service struct {
	l LookupService
}

func (s *Service) santize(number string) (string, error) {

	replacer := strings.NewReplacer(" ", "", "-", "", "(", "", ")", "", "+", "", "/", "")
	number = replacer.Replace(strings.TrimLeft(number, "0"))

	pNum, err := phonenumbers.Parse("+"+number, "")
	if err != nil {
		return "", err
	}

	if !phonenumbers.IsValidNumber(pNum) {
		return "", fmt.Errorf("Invalid phone number")
	}

	fNum := phonenumbers.Format(pNum, phonenumbers.E164)
	return fNum[1:], nil
}

func (s *Service) Handler(w http.ResponseWriter, r *http.Request) {

	var number string
	if v, ok := r.URL.Query()["number"]; ok {
		number = v[0]
	}

	if number == "" {
		http.NotFound(w, r)
	}

	var err error
	number, err = s.santize(number)
	if err != nil {
		log.Println(err)
		http.NotFound(w, r)
	}

	value, err := s.l.Lookup(number)
	if err != nil {
		log.Println(err)
		http.NotFound(w, r)
	}

	resp := &api.Response{
		MCC: fmt.Sprintf("%03d", value.MCC),
		MNC: fmt.Sprintf("%03d", value.MNC),
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	enc := json.NewEncoder(w)
	if err := enc.Encode(resp); err != nil {
		log.Println(err)
	}
}

func NewService(l LookupService) *Service {
	return &Service{
		l: l,
	}
}
