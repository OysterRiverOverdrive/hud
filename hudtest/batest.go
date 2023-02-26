package hudtest

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"

	"github.com/oysterriveroverdrive/hud/model"
)

type BlueAlliance struct {
	Server    *httptest.Server
	Districts []*model.District
	District  []*model.District
}

func NewBlueAlliance() (*BlueAlliance, error) {

	blueAlliance := &BlueAlliance{}
	mux := http.NewServeMux()
	mux.HandleFunc("/district", blueAlliance.districtHandler)
	mux.HandleFunc("/districts", blueAlliance.districtsHandler)

	blueAlliance.Server = httptest.NewServer(mux)

	return blueAlliance, nil
}

// Close cleans up the hudtest server resources.
func (s *BlueAlliance) Close() {
	s.Server.Close()
}

// DistrictHandlers serves /districts/* endpoints. Such as
func (s *BlueAlliance) districtsHandler(w http.ResponseWriter, r *http.Request) {
	b, _ := httputil.DumpRequest(r, true)
	log.Println(string(b))
	json.NewEncoder(w).Encode(s.Districts)
}

// DistrictHandler serves /district/* endpoints. Such as
func (s *BlueAlliance) districtHandler(w http.ResponseWriter, r *http.Request) {
	b, _ := httputil.DumpRequest(r, true)
	log.Println(string(b))
	json.NewEncoder(w).Encode(s.District)
}
