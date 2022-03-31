package bamm

import (
	"fmt"
	"strings"
	"time"
)

const (
	DEFAULT_SERVER string = "https://www.thebluealliance.com/api/v3"
)

// NextMatch - Searches the list of matches provided for the next match that
// hasn't occured yet.  If no more matches are found, the function returns
// nil, nil.
func NextMatch(matches []*MatchSimple, teamNumber int) (*MatchSimple, error) {
	var nextMatch *MatchSimple
	matchTime := int64(99999999999) // default to a timestamp in the year 5138

	teamKey := fmt.Sprintf("frc%d", teamNumber)
	for _, match := range matches {
		match := match
		// If the match has an actual time, it has already happened.
		if match.ActualTime > 0 {
			continue
		}
		// Check red alliance
		for _, key := range match.Alliances.Red.TeamKeys {
			if key == teamKey && match.PredictedTime < matchTime {
				nextMatch = match
				matchTime = match.PredictedTime
			}
		}
		// Check blue alliance
		for _, key := range match.Alliances.Blue.TeamKeys {
			if key == teamKey {
				nextMatch = match
				matchTime = match.PredictedTime
			}
		}
	}

	return nextMatch, nil
}

func PrintNextMatchSummary(next *MatchSimple, teamNumber int) []string {
	if next == nil {
		return []string{
			"No more matches found.",
		}
	}
	data := []string{}
	data = append(data, fmt.Sprintf("Match Starts At: %s", time.Unix(next.PredictedTime, 0)))
	data = append(data, fmt.Sprintf("Starts In: %s", time.Unix(next.PredictedTime, 0).Sub(time.Now()).Truncate(time.Second)))
	teamKey := fmt.Sprintf("frc%d", teamNumber)
	isRed := false
	for _, key := range next.Alliances.Red.TeamKeys {
		if key == teamKey {
			isRed = true
		}
	}
	if isRed {
		data = append(data, fmt.Sprintf("Alliance: Red"))
		data = append(data, fmt.Sprintf("Allies: %s", strings.Join(next.Alliances.Red.TeamKeys, " ")))
		data = append(data, fmt.Sprintf("Opponents: %s", strings.Join(next.Alliances.Blue.TeamKeys, " ")))
	} else {
		data = append(data, fmt.Sprintf("Alliance: Blue"))
		data = append(data, fmt.Sprintf("Allies: %s", strings.Join(next.Alliances.Blue.TeamKeys, " ")))
		data = append(data, fmt.Sprintf("Opponents: %s", strings.Join(next.Alliances.Red.TeamKeys, " ")))
	}
	return data
}
