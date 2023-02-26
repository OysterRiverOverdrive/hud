package hud

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/oysterriveroverdrive/hud/bluealliance"
	"github.com/oysterriveroverdrive/hud/bluealliance/model"
)

type Team struct {
	Data *model.Team
}

type Teams struct {
	Data []*model.Team
}

func TeamByNumber(ba *bluealliance.Service, teamNum int) (*Team, error) {
	resp, err := ba.Get(ba.URL+fmt.Sprintf("/team/frc%d", teamNum), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	r := &model.Team{}
	if err := json.NewDecoder(resp.Body).Decode(r); err != nil {
		return nil, fmt.Errorf("unable to parse team response: %w", err)
	}
	if r.ErrorResponse.Error != "" {
		return nil, fmt.Errorf(r.ErrorResponse.Error)
	}
	return &Team{Data: r}, nil
}

func TeamsInDistrict(ba *bluealliance.Service, district string) (*Teams, error) {
	resp, err := ba.Get(ba.URL+fmt.Sprintf("/district/%s/teams", district), nil)
	if err != nil {
		log.Println("error", err)
		return nil, err
	}
	defer resp.Body.Close()

	r := []*model.Team{}
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, fmt.Errorf("unable to parse district teams response: %w", err)
	}
	return &Teams{Data: r}, nil
}
