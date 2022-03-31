// Schema for objects defined at https://www.thebluealliance.com/apidocs/v3
package bamm

import "fmt"

type TeamSimple struct {
	// TBA team key with the format frcXXXX with XXXX representing the team number.
	Key string `json:"key"`
	// Official team number issued by FIRST.
	TeamNumber int `json:"team_number"`
	// Team nickname provided by FIRST.
	Nickname string `json:"nickname"`
	// Official long name registered with FIRST.
	Name string `json:"name"`
	// City of team derived from parsing the address registered with FIRST.
	City string `json:"city"`
	// State of team derived from parsing the address registered with FIRST.
	StateProv string `json:"state_prov"`
	// Country of team derived from parsing the address registered with FIRST.
	Country string `json:"country"`
}

func (t TeamSimple) String() string {
	return fmt.Sprintf("TeamNumber: %d\nNickname: %s", t.TeamNumber, t.Nickname)
}

// Media - contains a reference for most any media associated with a team or event on TBA.
type Media struct {
	// String type of the media element.
	Type string `json:"type"`
	// The key used to identify this media on the media site.
	ForeignKey string            `json:"foreign_key"`
	Details    map[string]string `json:"details"`
	// True if the media is of high quality.
	Preferred bool `json:"preferred"`
	// Direct URL to the media.
	DirectURL string `json:"direct_url"`
	// The URL that leads to the full web page for the media, if one exists.
	ViewURL string `json:"view_url"`
}

func (t Media) String() string {
	return fmt.Sprintf("Type: %s\nViewURL: %s\nDirectURL: %s\nDetails: %v", t.Type, t.ViewURL, t.DirectURL, t.Details)
}

type Team struct {
	// TBA team key with the format frcXXXX with XXXX representing the team number.
	Key string `json:"key"`
	// Official team number issued by FIRST.
	TeamNumber int `json:"team_number"`
	// Team nickname provided by FIRST.
	Nickname string `json:"nickname"`
	// Official long name registered with FIRST.
	Name string `json:"name"`
	// Name of team school or affilited group registered with FIRST.
	SchoolName string `json:"school_name"`
	// City of team derived from parsing the address registered with FIRST.
	City string `json:"city"`
	// State of team derived from parsing the address registered with FIRST.
	StateProv string `json:"state_prov"`
	// Country of team derived from parsing the address registered with FIRST.
	Country string `json:"country"`
	// Will be NULL, for future development.
	Address string `json:"address"`
	// Postal code from the team address.
	PostalCode string `json:"postal_code"`
	// Will be NULL, for future development.
	GmapsPlaceID string `json:"gmaps_place_id"`
	// Will be NULL, for future development.
	GmapsURL string `json:"gmaps_url"`
	// Will be NULL, for future development.
	Lat float64 `json:"lat"`
	// Will be NULL, for future development.
	Lng float64 `json:"lng"`
	// Will be NULL, for future development.
	LocationName string `json:"location_name"`
	// Official website associated with the team.
	Website string `json:"website"`
	// First year the team officially competed.
	RookieYear int `json:"rookie_year"`
	// Location of the team's home championship each year as a key-value pair. The year (as a string) is the key, and the city is the value.
	HomeChampionship map[string]string `json:"home_championship"`
}

type MatchSimple struct {
	// TBA match key with the format yyyy[EVENT_CODE]_[COMP_LEVEL]m[MATCH_NUMBER], where yyyy is the year, and EVENT_CODE is the event code of the event, COMP_LEVEL is (qm, ef, qf, sf, f), and MATCH_NUMBER is the match number in the competition level. A set number may append the competition level if more than one match in required per set.
	Key string `json:"key"`
	// 	The competition level the match was played at.
	// Enum:
	// [ qm, ef, qf, sf, f ]
	CompLevel string `json:"comp_level"`
	// The set number in a series of matches where more than one match is required in the match series.
	SetNumber int `json:"set_number"`
	// The match number of the match in the competition level.
	MatchNumber int `json:"match_number"`
	// A list of alliances, the teams on the alliances, and their score.
	Alliances struct {
		Red  MatchAlliance `json:"red"`
		Blue MatchAlliance `json:"blue"`
	} `json:"alliances"`
	// The color (red/blue) of the winning alliance. Will contain an empty string in the event of no winner, or a tie.
	// Enum:
	// Array [ red, blue, ""]
	WinningAlliance string `json:"winning_alliance"`
	// Event key of the event the match was played at.
	EventKey string `json:"event_key"`
	// UNIX timestamp (seconds since 1-Jan-1970 00:00:00) of the scheduled match time, as taken from the published schedule.
	Time int64 `json:"time"`
	// UNIX timestamp (seconds since 1-Jan-1970 00:00:00) of the TBA predicted match start time.
	PredictedTime int64 `json:"predicted_time"`
	// UNIX timestamp (seconds since 1-Jan-1970 00:00:00) of actual match start time.
	ActualTime int64 `json:"actual_time"`
}

type MatchAlliance struct {
	// Score for this alliance. Will be null or -1 for an unplayed match.
	Score int `json:"score"`
	// TBA Team keys (eg frc254) for teams on this alliance.
	TeamKeys []string `json:"team_keys"`
	// TBA team keys (eg frc254) of any teams playing as a surrogate.
	SurrogateTeamKeys []string `json:"surrogate_team_keys"`
	// TBA team keys (eg frc254) of any disqualified teams.
	DQTeamKeys []string `json:"dq_team_keys"`
}
