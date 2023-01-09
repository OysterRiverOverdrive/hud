package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/oysterriveroverdrive/hud"
	"github.com/sirupsen/logrus"
)

// RulesCmd handles @hud rules ... commands
type RulesCmd struct {
	SubCmds []Cmd
}

func (c *RulesCmd) Stub() string {
	return "rules"
}

func (c *RulesCmd) Match(msg string) bool {
	logrus.Debugf("RulesCmd.Match %q", msg)
	return strings.HasPrefix(strings.TrimSpace(msg), "rules")
}

func (c *RulesCmd) Help() string {
	return c.Stub() + " - working with game rules"
}

func (c *RulesCmd) Handle(md map[string]string, ts *hud.TriviaService, s *discordgo.Session, m *discordgo.MessageCreate, msg string) (string, *discordgo.MessageSend, error) {
	logrus.Debugf("RulesCmd.Handle %v %q", md, msg)
	suffix := strings.TrimSpace(strings.TrimPrefix(msg, "rules"))
	if suffix == "help" || suffix == "" {
		var help []string
		for _, subCmd := range c.SubCmds {
			help = append(help, md["path"]+" "+c.Stub()+" "+subCmd.Help())
		}
		return m.ChannelID, &discordgo.MessageSend{
			Content: "rules help:\n" + strings.Join(help, "\n"),
		}, nil
	}
	for _, subCmd := range c.SubCmds {
		if subCmd.Match(suffix) {
			md["path"] += " " + c.Stub()
			return subCmd.Handle(md, ts, s, m, suffix)
		}
	}
	return "", nil, nil
}

// RulesListCmd handles @hud rules list ... commands
type RulesListCmd struct{}

func (c *RulesListCmd) Stub() string {
	return "list"
}

func (c *RulesListCmd) Match(msg string) bool {
	logrus.Debugf("RulesListCmd.Match %q", msg)
	return msg == "list"
}

func (c *RulesListCmd) Help() string {
	return c.Stub() + " - list the summary of all the new rules for this year"
}

func (c *RulesListCmd) Handle(md map[string]string, ts *hud.TriviaService, s *discordgo.Session, m *discordgo.MessageCreate, msg string) (string, *discordgo.MessageSend, error) {
	logrus.Debugf("RulesListCmd.Handle %v %q", md, msg)

	var summaries []string
	for _, rule := range ChargedUpRules {
		summaries = append(summaries, fmt.Sprintf("%s: %s", rule.Number, rule.Title))
	}

	return m.ChannelID, &discordgo.MessageSend{
		Content: fmt.Sprintf("Rules specific to 2023 Charge Up:\n%s", strings.Join(summaries, "\n")),
	}, nil
}

// RulesNumberCmd handles @hud rules [RuleNumber] ... commands
type RulesNumberCmd struct{}

func (c *RulesNumberCmd) Stub() string {
	return "[RuleNumber]"
}

func (c *RulesNumberCmd) Match(msg string) bool {
	logrus.Debugf("RulesNumberCmd.Match %q", msg)
	return regexp.MustCompile(`^\s*[a-zA-Z]\d{3}\s*`).MatchString(msg)
}

func (c *RulesNumberCmd) Help() string {
	return c.Stub() + " - show the details for a specific rule id."
}

func (c *RulesNumberCmd) Handle(md map[string]string, ts *hud.TriviaService, s *discordgo.Session, m *discordgo.MessageCreate, msg string) (string, *discordgo.MessageSend, error) {
	logrus.Debugf("RulesNumberCmd.Handle %v %q", md, msg)
	match := regexp.MustCompile(`\s*([a-zA-Z]\d{3})\s*`).FindStringSubmatch(msg)
	ruleID := strings.ToUpper(match[1])

	for _, rule := range ChargedUpRules {
		if rule.Number == ruleID {
			return m.ChannelID, &discordgo.MessageSend{
				Content: fmt.Sprintf("Rule Number: %s\nTitle: %s\nDetails: %s", rule.Number, rule.Title, rule.Details),
			}, nil
		}
	}
	return m.ChannelID, &discordgo.MessageSend{
		Content: fmt.Sprintf("unable to locate rule number %q", msg),
	}, nil

}

// RulesSearchCmd handles @hud rules search [keyword] ... commands
type RulesSearchCmd struct{}

func (c *RulesSearchCmd) Stub() string {
	return "search [keyword]"
}

func (c *RulesSearchCmd) Match(msg string) bool {
	logrus.Debugf("RulesSearchCmd.Match %q", msg)
	return strings.HasPrefix(strings.TrimSpace(msg), "search ")
}

func (c *RulesSearchCmd) Help() string {
	return c.Stub() + " - search for a keyword in the rules."
}

func (c *RulesSearchCmd) Handle(md map[string]string, ts *hud.TriviaService, s *discordgo.Session, m *discordgo.MessageCreate, msg string) (string, *discordgo.MessageSend, error) {
	logrus.Debugf("RulesSearchCmd.Handle %v %q", md, msg)

	keyword := strings.ToLower(strings.TrimPrefix(strings.TrimSpace(msg), "search "))

	var ruleMatches []FRCRule

	for _, rule := range ChargedUpRules {
		rule := rule
		if strings.Contains(strings.ToLower(rule.Title), keyword) ||
			strings.Contains(strings.ToLower(rule.Details), keyword) {
			ruleMatches = append(ruleMatches, rule)
		}
	}

	if len(ruleMatches) == 0 {
		return m.ChannelID, &discordgo.MessageSend{
			Content: fmt.Sprintf("search term %q not found", keyword),
		}, nil
	}

	// Attempt to format a return message with all the details.
	var msgDetails []string
	for _, rule := range ruleMatches {
		msgDetails = append(msgDetails, fmt.Sprintf("Rule Number: %s\nTitle: %s\nDetails: %s", rule.Number, rule.Title, rule.Details))
	}
	detailedMsg := strings.Join(msgDetails, "\n------------------\n")
	if len(detailedMsg) < 2000 {
		return m.ChannelID, &discordgo.MessageSend{
			Content: detailedMsg,
		}, nil
	}

	// The search keyword resulted in too many hits, try a smaller summarization.
	msgDetails = []string{}
	for _, rule := range ruleMatches {
		msgDetails = append(msgDetails, fmt.Sprintf("Rule Number: %s\nTitle: %s", rule.Number, rule.Title))
	}
	titleMsg := "Too many hits. Rule details removed. Use @hud rule [RuleNumber] for more information.\n" +
		strings.Join(msgDetails, "\n------------------\n")
	if len(titleMsg) < 2000 {
		return m.ChannelID, &discordgo.MessageSend{
			Content: titleMsg,
		}, nil
	}

	// Even without the details, the message is still too long. Return just the rule numbers.
	msgDetails = []string{}
	for _, rule := range ruleMatches {
		msgDetails = append(msgDetails, rule.Number)
	}
	numberMsg := "Too many hits. Rule details removed. Use @hud rule [RuleNumber] for more information.\n" +
		strings.Join(msgDetails, "\n")
	if len(numberMsg) < 2000 {
		return m.ChannelID, &discordgo.MessageSend{
			Content: numberMsg,
		}, nil
	}

	return m.ChannelID, &discordgo.MessageSend{
		Content: "Too many results to display in discord. Consult the manual https://firstfrc.blob.core.windows.net/frc2023/Manual/2023FRCGameManual.pdf",
	}, nil
}

type FRCRule struct {
	Number  string
	Title   string
	Details string
}

var ChargedUpRules = []FRCRule{{
	Number:  "G106",
	Title:   "Tall ROBOTS not allowed.",
	Details: "ROBOT height, as measured when it's resting normally on a flat floor, may not exceed 6 ft. 6 in. (~198 cm)) above the carpet during the MATCH.",
}, {
	Number:  "G107",
	Title:   "Don't overextend yourself.",
	Details: "ROBOTS may not extend beyond their FRAME PERIMETER in more than 48 in. (~122 cm). MOMENTARY and inconsequential extensions beyond 48 in. (~122 cm) are an exception to this rule.",
}, {
	Number:  "G108",
	Title:   "Opponent's zone, no extension.",
	Details: "A ROBOT whose BUMPERS are intersecting the opponent's LOADING ZONE or COMMUNITY may not extend beyond its FRAME PERIMETER. Extensions which are both MOMENTARY and inconsequential are an exception to this rule.",
}, {
	Number:  "G109",
	Title:   "Don't extend in multiple directions.",
	Details: "ROBOTS may not extend beyond their FRAME PERIMETER in more than one direction (i.e. over 1 side of the ROBOT) at a time. For the purposes of this rule, a round or circular section of FRAME PERIMETER is considered to have an infinite number of sides. Exceptions to this rule are: A. MOMENTARY and inconsequential extensions in multiple directions B. A ROBOT fully contained within its LOADING ZONE or COMMUNITY.",
}, {
	Number:  "G207",
	Title:   "Right of way.",
	Details: "A ROBOT with any part of itself in their opponent's LOADING ZONE or COMMUNITY may not contact an opponent ROBOT, regardless of who initiates contact.",
}, {
	Number:  "G208",
	Title:   "Don't climb on each other unless in the COMMUNITY.",
	Details: "A ROBOT may not be fully supported by a partner ROBOT unless the partner's BUMPERS intersect its COMMUNITY.",
}, {
	Number:  "G209",
	Title:   "During the ENDGAME, don't touch ROBOTS touching their CHARGE STATION.",
	Details: "During the ENDGAME, a ROBOT may not contact, either directly or transitively through a GAME PIECE, an opponent ROBOT contacting its CHARGE STATION or supported by a partner contacting its CHARGE STATION, regardless of who initiates contact. A ROBOT in contact with its CHARGE STATION and partially in its opponent's LOADING ZONE is not protected by this rule.",
}, {
	Number:  "G301",
	Title:   "Be careful what you interact with.",
	Details: "ROBOTS and OPERATOR CONSOLES are prohibited from the following actions with regards to interaction with ARENA elements. Items A-D exclude GAME PIECES. grabbing, A. grasping, B. attaching to (including the use of a vacuum or hook fastener to anchor to the FIELD carpet C. and excluding use of the DRIVER STATION hook-and-loop tape, plugging in to the provided power outlet, and plugging the provided Ethernet cable into the OPERATOR CONSOLE), D. deforming, E. becoming entangled with, F. suspending from, and G. damaging.",
}, {
	Number:  "G302",
	Title:   "Stay on your side in AUTO.",
	Details: "During AUTO, a ROBOT may not intersect the infinite vertical volume created by the CENTERLINE of the FIELD.",
}, {
	Number:  "G303",
	Title:   "Do not interfere with opponent GAME PIECES in AUTO.",
	Details: "During AUTO, a ROBOT action may not cause GAME PIECES staged on the opposing side of the FIELD to move from their starting locations.",
}, {
	Number:  "G304",
	Title:   "Don't mess with the opponent's CHARGE STATION.",
	Details: "ROBOTS, either directly or transitively through a GAME PIECE, may not cause or prevent the movement of the opponent CHARGE STATION. The following are exceptions to this rule: A. movement, or prevention of movement, of an opponent CHARGE STATION because of a MOMENTARY ROBOT action resulting in minimal CHARGE STATION movement B. a ROBOT forced to contact an opponent's CHARGE STATION because of contact by an opponent ROBOT, either directly or transitively through a GAME PIECE or other ROBOT (e.g. a ROBOT wedged underneath the CHARGE STATION by the opposing ALLIANCE either intentionally or accidentally).",
}, {
	Number:  "G305",
	Title:   "Don't trick the sensors.",
	Details: "Teams may not interfere with automated scoring hardware.",
}, {
	Number:  "G403",
	Title:   "1 GAME PIECE at a time (except in LOADING ZONE and COMMUNITY).",
	Details: "ROBOTS completely outside their LOADING ZONE or COMMUNITY may not have greater-than-MOMENTARY CONTROL of more than 1 GAME PIECE, either directly or transitively through other objects. A ROBOT is in CONTROL of a GAME PIECE if: A. the GAME PIECE is fully supported by the ROBOT, or B. the ROBOT is intentionally moving a GAME PIECE to a desired location or in a preferred direction",
}, {
	Number:  "G404",
	Title:   "Launching GAME PIECES is only okay in the COMMUNITY.",
	Details: "A ROBOT may not launch GAME PIECES unless any part of the ROBOT is in its own COMMUNITY.",
}, {
	Number:  "G405",
	Title:   "Don't mess with the opponents' GRIDS.",
	Details: "A ROBOT may not move a scored GAME PIECE from an opponent's NODE.",
}, {
	Number:  "H110",
	Title:   "Don't mess with GAME PIECES.",
	Details: "Teams may not modify GAME PIECES in any way. Temporary deformation (e.g.to pre-load a ROBOT) is an exception to this rule.",
}, {
	Number:  "H309",
	Title:   "Know your ROBOT setup.",
	Details: "When placed on the FIELD for a MATCH, each ROBOT must be: A. in compliance with all ROBOT rules, i.e. has passed inspection (for exceptions regarding Practice MATCHES, see Section 10 Inspection & Eligibility Rules), B. the only team-provided item left on the FIELD by the DRIVE TEAM, C. confined to its STARTING CONFIGURATION (reference R102 and R104), D. positioned such that it is fully contained within its COMMUNITY E. not in contact with the CHARGE STATION F. fully supported by FIELD carpet, and G. fully and solely supporting not more than 1 GAME PIECE (as described in Section 6.1 Setup).",
}, {
	Number:  "H310",
	Title:   "Know your DRIVE TEAM positions.",
	Details: "Prior to the start of the MATCH, DRIVE TEAM members must be positioned as follows: A. DRIVERS: inside their ALLIANCE AREA and behind the STARTING LINE, B. COACHES: inside their ALLIANCE AREA and behind the STARTING LINE, and C. HUMAN PLAYERS: a. at least one HUMAN PLAYER in their SUBSTATION AREA, b. any remaining HUMAN PLAYERS: inside their ALLIANCE AREA and behind the STARTING LINE, and D. TECHNICIANS: in the event-designated area near the FIELD.",
}, {
	Number:  "H311",
	Title:   "Leave the GAME PIECES alone.",
	Details: "Prior to the start of the MATCH, HUMAN PLAYERS may not rearrange the GAME PIECES within the SUBSTATION AREA.",
}, {
	Number:  "H504",
	Title:   "GAME PIECES through PORTALS only.",
	Details: "GAME PIECES may only be introduced to the FIELD A. by a HUMAN PLAYER, B. through a PORTAL, and C. during TELEOP.",
}, {
	Number:  "H505",
	Title:   "DRIVE TEAMS, watch your reach.",
	Details: "DRIVE TEAMS may not extend any body part into the SINGLE SUBSTATION PORTAL for a greater-than-MOMENTARY period of time.",
}, {
	Number:  "R104",
	Title:   "STARTING CONFIGURATION – max size.",
	Details: "A ROBOT'S STARTING CONFIGURATION may not have a FRAME PERIMETER greater than 120 in. (~304 cm) and may not be more than 4 ft. 6 in. (~137 cm) tall.",
}, {
	Number:  "R105",
	Title:   "ROBOT extension limit.",
	Details: "ROBOTS may not extend more than 48 in. (~121 cm) beyond their FRAME PERIMETER.",
}}
