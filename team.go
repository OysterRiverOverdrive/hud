package hud

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/oysterriveroverdrive/hud/model"
)

type Team struct {
	Data *model.Team
}

type Teams struct {
	Data []*model.Team
}

func TeamByNumber(ts *TriviaService, teamNum int) (*Team, error) {
	resp, err := ts.Get(ts.URL+fmt.Sprintf("/team/frc%d", teamNum), nil)
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

func TeamsInDistrict(ts *TriviaService, district string) (*Teams, error) {
	resp, err := ts.Get(ts.URL+fmt.Sprintf("/district/%s/teams", district), nil)
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
