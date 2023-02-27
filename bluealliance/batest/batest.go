package batest

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"

	"github.com/oysterriveroverdrive/hud/bluealliance/model"
)

type Service struct {
	Server    *httptest.Server
	Districts []*model.District
	District  []*model.District
}

func NewService() (*Service, error) {

	s := &Service{}
	mux := http.NewServeMux()
	mux.HandleFunc("/district", s.districtHandler)
	mux.HandleFunc("/districts", s.districtsHandler)

	s.Server = httptest.NewServer(mux)

	return s, nil
}

// Close cleans up the batest server resources.
func (s *Service) Close() {
	s.Server.Close()
}

// DistrictHandlers serves /districts/* endpoints. Such as
func (s *Service) districtsHandler(w http.ResponseWriter, r *http.Request) {
	b, _ := httputil.DumpRequest(r, true)
	log.Println(string(b))
	json.NewEncoder(w).Encode(s.Districts)
}

// DistrictHandler serves /district/* endpoints. Such as
func (s *Service) districtHandler(w http.ResponseWriter, r *http.Request) {
	b, _ := httputil.DumpRequest(r, true)
	log.Println(string(b))
	json.NewEncoder(w).Encode(s.District)
}
