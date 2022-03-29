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
	return fmt.Sprintf("TeamNumber: %d\nNickname: %s\n", t.TeamNumber, t.Nickname)
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
