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

func (c *RulesCmd) Handle(md map[string]string, ts *hud.TriviaService, s *discordgo.Session, m *discordgo.MessageCreate, msg string) (string, []*discordgo.MessageSend, error) {
	logrus.Debugf("RulesCmd.Handle %v %q", md, msg)
	suffix := strings.TrimSpace(strings.TrimPrefix(msg, "rules"))
	if suffix == "help" || suffix == "" {
		var help []string
		for _, subCmd := range c.SubCmds {
			help = append(help, md["path"]+" "+c.Stub()+" "+subCmd.Help())
		}
		return m.ChannelID, []*discordgo.MessageSend{{
			Content: "rules help:\n" + strings.Join(help, "\n"),
		}}, nil
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

func (c *RulesListCmd) Handle(md map[string]string, ts *hud.TriviaService, s *discordgo.Session, m *discordgo.MessageCreate, msg string) (string, []*discordgo.MessageSend, error) {
	logrus.Debugf("RulesListCmd.Handle %v %q", md, msg)

	var msgs []*discordgo.MessageSend
	summary := "2023 Charge Up Rules:"
	for _, rule := range ChargedUpRules {
		ruleStr := fmt.Sprintf("\n%s: %s", rule.Number, rule.Title)
		if len(summary)+len(ruleStr) >= 2000 {
			msgs = append(msgs, &discordgo.MessageSend{
				Content: summary,
			})
			summary = ruleStr
		} else {
			summary += ruleStr
		}
	}
	msgs = append(msgs, &discordgo.MessageSend{
		Content: summary,
	})

	return m.ChannelID, msgs, nil
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

func (c *RulesNumberCmd) Handle(md map[string]string, ts *hud.TriviaService, s *discordgo.Session, m *discordgo.MessageCreate, msg string) (string, []*discordgo.MessageSend, error) {
	logrus.Debugf("RulesNumberCmd.Handle %v %q", md, msg)
	match := regexp.MustCompile(`\s*([a-zA-Z]\d{3})\s*`).FindStringSubmatch(msg)
	ruleID := strings.ToUpper(match[1])

	for _, rule := range ChargedUpRules {
		if rule.Number == ruleID {
			return m.ChannelID, []*discordgo.MessageSend{{
				Content: fmt.Sprintf("Rule Number: %s\nTitle: %s\nDetails: %s", rule.Number, rule.Title, rule.Details),
			}}, nil
		}
	}
	return m.ChannelID, []*discordgo.MessageSend{{
		Content: fmt.Sprintf("unable to locate rule number %q", msg),
	}}, nil

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

func (c *RulesSearchCmd) Handle(md map[string]string, ts *hud.TriviaService, s *discordgo.Session, m *discordgo.MessageCreate, msg string) (string, []*discordgo.MessageSend, error) {
	logrus.Debugf("RulesSearchCmd.Handle %v %q", md, msg)

	keyword := strings.ToLower(strings.TrimPrefix(strings.TrimSpace(msg), "search "))

	var ruleMatches []FRCRule

	for _, rule := range ChargedUpRules {
		rule := rule
		// See if the keyword is in the rule title or description, or if the
		// keyword is a rule number.
		if strings.Contains(strings.ToLower(rule.Title), keyword) ||
			strings.Contains(strings.ToLower(rule.Details), keyword) ||
			strings.ToLower(rule.Number) == keyword {
			ruleMatches = append(ruleMatches, rule)
		}
	}

	// No matches, exit here with a sorry, not found.
	if len(ruleMatches) == 0 {
		return m.ChannelID, []*discordgo.MessageSend{{
			Content: fmt.Sprintf("search term %q not found", keyword),
		}}, nil
	}

	// Attempt to format a return message with all the details.
	var msgDetails []string
	for _, rule := range ruleMatches {
		msgDetails = append(msgDetails, fmt.Sprintf("Rule Number: %s\nTitle: %s\nDetails: %s", rule.Number, rule.Title, rule.Details))
	}
	detailedMsg := strings.Join(msgDetails, "\n------------------\n")
	if len(detailedMsg) < 2000 {
		return m.ChannelID, []*discordgo.MessageSend{{
			Content: detailedMsg,
		}}, nil
	}

	// The search keyword resulted in too many hits, try a smaller summarization.
	msgDetails = []string{}
	for _, rule := range ruleMatches {
		msgDetails = append(msgDetails, fmt.Sprintf("Rule Number: %s\nTitle: %s", rule.Number, rule.Title))
	}
	titleMsg := "Too many hits. Rule details removed. Use @hud rules [RuleNumber] for more information.\n" +
		strings.Join(msgDetails, "\n------------------\n")
	if len(titleMsg) < 2000 {
		return m.ChannelID, []*discordgo.MessageSend{{
			Content: titleMsg,
		}}, nil
	}

	// Even without the details, the message is still too long. Return just the rule numbers.
	msgDetails = []string{}
	for _, rule := range ruleMatches {
		msgDetails = append(msgDetails, rule.Number)
	}
	numberMsg := "Too many hits. Rule details removed. Use @hud rules [RuleNumber] for more information.\n" +
		strings.Join(msgDetails, "\n")
	if len(numberMsg) < 2000 {
		return m.ChannelID, []*discordgo.MessageSend{{
			Content: numberMsg,
		}}, nil
	}

	// We tried, but the search term is just too generic to return data
	// via discord. Let them know the'll need to just RTM.
	return m.ChannelID, []*discordgo.MessageSend{{
		Content: "Too many results to display in discord. Consult the manual https://firstfrc.blob.core.windows.net/frc2023/Manual/2023FRCGameManual.pdf",
	}}, nil
}

// FRCRule contains the main text data for a rule.
type FRCRule struct {
	// Number contains the formatted rule number such as R101.
	Number string
	// Title contains a summary of the rule.
	Title string
	// Detils contains the full (or a significant portion) of the rule description.
	Details string
	// Evergreen rules are rules that are expected to "go relatively unchanged"
	// from year to year.
	Evergreen bool
}

// ChargedUpRules lists all the rules in the 2023 FRC "Charged Up" competition.
// Authoritative source for the rules is:
//
//	https://firstfrc.blob.core.windows.net/frc2023/Manual/2023FRCGameManual.pdf
var ChargedUpRules = []FRCRule{{
	Number:    "G101",
	Title:     "Dangerous ROBOTS: not allowed.",
	Details:   "ROBOTS whose operation or design is dangerous or unsafe are not permitted.",
	Evergreen: true,
}, {
	Number:    "G102",
	Title:     "ROBOTS, stay on the FIELD during the MATCH.",
	Details:   "ROBOTS and anything they control, e.g. GAME PIECES, may not contact anything outside the FIELD except for MOMENTARY incursions into PORTALS.",
	Evergreen: true,
}, {
	Number:    "G103",
	Title:     "Keep your BUMPERS low.",
	Details:   "BUMPERS must be in the BUMPER ZONE (see R402) during the MATCH.",
	Evergreen: true,
}, {
	Number:    "G104",
	Title:     "Keep your BUMPERS together.",
	Details:   "BUMPERS may not fail such that a segment completely detaches, any corner (as defined in R401) of a ROBOT'S FRAME PERIMETER is exposed, or the team number or ALLIANCE color are indeterminate.",
	Evergreen: true,
}, {
	Number:    "G105",
	Title:     "Keep it together.",
	Details:   "ROBOTS may not intentionally detach or leave parts on the FIELD.",
	Evergreen: true,
}, {
	Number:    "G106",
	Title:     "Tall ROBOTS not allowed.",
	Details:   "ROBOT height, as measured when it's resting normally on a flat floor, may not exceed 6 ft. 6 in. (~198 cm)) above the carpet during the MATCH.",
	Evergreen: false,
}, {
	Number:    "G107",
	Title:     "Don't overextend yourself.",
	Details:   "ROBOTS may not extend beyond their FRAME PERIMETER in more than 48 in. (~122 cm). MOMENTARY and inconsequential extensions beyond 48 in. (~122 cm) are an exception to this rule.",
	Evergreen: false,
}, {
	Number:    "G108",
	Title:     "Opponent's zone, no extension.",
	Details:   "A ROBOT whose BUMPERS are intersecting the opponent's LOADING ZONE or COMMUNITY may not extend beyond its FRAME PERIMETER. Extensions which are both MOMENTARY and inconsequential are an exception to this rule.",
	Evergreen: false,
}, {
	Number:    "G109",
	Title:     "Don't extend in multiple directions.",
	Details:   "ROBOTS may not extend beyond their FRAME PERIMETER in more than one direction (i.e. over 1 side of the ROBOT) at a time. For the purposes of this rule, a round or circular section of FRAME PERIMETER is considered to have an infinite number of sides. Exceptions to this rule are: A. MOMENTARY and inconsequential extensions in multiple directions B. A ROBOT fully contained within its LOADING ZONE or COMMUNITY.",
	Evergreen: false,
}, {
	Number:    "G201",
	Title:     "Don't expect to gain by doing others harm.",
	Details:   "Strategies clearly aimed at forcing the opponent ALLIANCE to violate a rule are not in the spirit of FIRST Robotics Competition and not allowed. Rule violations forced in this manner will not result in an assignment of a penalty to the targeted ALLIANCE.",
	Evergreen: true,
}, {
	Number:    "G202",
	Title:     "There's a 5-count on PINS.",
	Details:   "ROBOTS may not PIN an opponent's ROBOT for more than 5 seconds. A ROBOT is PINNING if it is preventing the movement of an opponent ROBOT by contact, either direct or transitive (such as against a FIELD element). A ROBOT is considered PINNED until the ROBOTS have separated by at least 6 ft. (~183 cm) from each other, either ROBOT has moved 6 ft. from where the PIN initiated, or the PINNING ROBOT gets PINNED, whichever comes first. The PINNING ROBOT(S) must then wait for at least 3 seconds before attempting to PIN the same ROBOT again.",
	Evergreen: true,
}, {
	Number:    "G203",
	Title:     "Don't collude with your partners to shut down major parts of game play.",
	Details:   "2 or more ROBOTS that appear to a REFEREE to be working together may neither isolate nor close off any major element of MATCH play.",
	Evergreen: true,
}, {
	Number:    "G204",
	Title:     "Stay out of other ROBOTS.",
	Details:   "A ROBOT may not use a COMPONENT outside its FRAME PERIMETER (except its BUMPERS) to initiate contact with an opponent ROBOT inside the vertical projection of that opponent ROBOT'S FRAME PERIMETER. Contact with an opponent in an opening of their BUMPERS or in the space above the BUMPER opening are exceptions to this rule.",
	Evergreen: true,
}, {
	Number:    "G205",
	Title:     "This isn't combat robotics.",
	Details:   "A ROBOT may not damage or functionally impair an opponent ROBOT in either of the following ways: A. deliberately, as perceived by a REFEREE. B. regardless of intent, by initiating contact inside the vertical projection of an opponent ROBOT'S FRAME PERIMETER. Contact between the ROBOT'S BUMPERS or COMPONENTS inside the ROBOT'S FRAME PERIMETER and COMPONENTS inside an opening of an opponent's BUMPERS is an exception to this rule.",
	Evergreen: true,
}, {
	Number:    "G206",
	Title:     "Don't tip or entangle.",
	Details:   "A ROBOT may not deliberately, as perceived by a REFEREE, attach to, tip, or entangle with an opponent ROBOT.",
	Evergreen: true,
}, {
	Number:    "G207",
	Title:     "Right of way.",
	Details:   "A ROBOT with any part of itself in their opponent's LOADING ZONE or COMMUNITY may not contact an opponent ROBOT, regardless of who initiates contact.",
	Evergreen: false,
}, {
	Number:    "G208",
	Title:     "Don't climb on each other unless in the COMMUNITY.",
	Details:   "A ROBOT may not be fully supported by a partner ROBOT unless the partner's BUMPERS intersect its COMMUNITY.",
	Evergreen: false,
}, {
	Number:    "G209",
	Title:     "During the ENDGAME, don't touch ROBOTS touching their CHARGE STATION.",
	Details:   "During the ENDGAME, a ROBOT may not contact, either directly or transitively through a GAME PIECE, an opponent ROBOT contacting its CHARGE STATION or supported by a partner contacting its CHARGE STATION, regardless of who initiates contact. A ROBOT in contact with its CHARGE STATION and partially in its opponent's LOADING ZONE is not protected by this rule.",
	Evergreen: false,
}, {
	Number:    "G301",
	Title:     "Be careful what you interact with.",
	Details:   "ROBOTS and OPERATOR CONSOLES are prohibited from the following actions with regards to interaction with ARENA elements. Items A-D exclude GAME PIECES. grabbing, A. grasping, B. attaching to (including the use of a vacuum or hook fastener to anchor to the FIELD carpet C. and excluding use of the DRIVER STATION hook-and-loop tape, plugging in to the provided power outlet, and plugging the provided Ethernet cable into the OPERATOR CONSOLE), D. deforming, E. becoming entangled with, F. suspending from, and G. damaging.",
	Evergreen: false,
}, {
	Number:    "G302",
	Title:     "Stay on your side in AUTO.",
	Details:   "During AUTO, a ROBOT may not intersect the infinite vertical volume created by the CENTERLINE of the FIELD.",
	Evergreen: false,
}, {
	Number:    "G303",
	Title:     "Do not interfere with opponent GAME PIECES in AUTO.",
	Details:   "During AUTO, a ROBOT action may not cause GAME PIECES staged on the opposing side of the FIELD to move from their starting locations.",
	Evergreen: false,
}, {
	Number:    "G304",
	Title:     "Don't mess with the opponent's CHARGE STATION.",
	Details:   "ROBOTS, either directly or transitively through a GAME PIECE, may not cause or prevent the movement of the opponent CHARGE STATION. The following are exceptions to this rule: A. movement, or prevention of movement, of an opponent CHARGE STATION because of a MOMENTARY ROBOT action resulting in minimal CHARGE STATION movement B. a ROBOT forced to contact an opponent's CHARGE STATION because of contact by an opponent ROBOT, either directly or transitively through a GAME PIECE or other ROBOT (e.g. a ROBOT wedged underneath the CHARGE STATION by the opposing ALLIANCE either intentionally or accidentally).",
	Evergreen: false,
}, {
	Number:    "G305",
	Title:     "Don't trick the sensors.",
	Details:   "Teams may not interfere with automated scoring hardware.",
	Evergreen: false,
}, {
	Number:    "G401",
	Title:     "Keep GAME PIECES in bounds.",
	Details:   "ROBOTS may not intentionally eject GAME PIECES from the FIELD (either directly or by bouncing off a FIELD element or other ROBOT).",
	Evergreen: true,
}, {
	Number:    "G402",
	Title:     "GAME PIECES: use as directed.",
	Details:   "ROBOTS may not deliberately use GAME PIECES in an attempt to ease or amplify challenges associated with FIELD elements.",
	Evergreen: true,
}, {
	Number:    "G403",
	Title:     "1 GAME PIECE at a time (except in LOADING ZONE and COMMUNITY).",
	Details:   "ROBOTS completely outside their LOADING ZONE or COMMUNITY may not have greater-than-MOMENTARY CONTROL of more than 1 GAME PIECE, either directly or transitively through other objects. A ROBOT is in CONTROL of a GAME PIECE if: A. the GAME PIECE is fully supported by the ROBOT, or B. the ROBOT is intentionally moving a GAME PIECE to a desired location or in a preferred direction",
	Evergreen: false,
}, {
	Number:    "G404",
	Title:     "Launching GAME PIECES is only okay in the COMMUNITY.",
	Details:   "A ROBOT may not launch GAME PIECES unless any part of the ROBOT is in its own COMMUNITY.",
	Evergreen: false,
}, {
	Number:    "G405",
	Title:     "Don't mess with the opponents' GRIDS.",
	Details:   "A ROBOT may not move a scored GAME PIECE from an opponent's NODE.",
	Evergreen: false,
}, {
	Number:    "H101",
	Title:     "Be a good person.",
	Details:   "All teams must be civil toward everyone and respectful of team and event equipment while at a FIRST Robotics Competition event.",
	Evergreen: true,
}, {
	Number:    "H102",
	Title:     "Enter only 1 ROBOT.",
	Details:   "Each registered FIRST Robotics Competition team may enter only 1 ROBOT (or “robot,” a ROBOT -like assembly equipped with most of its drive base, i.e. its MAJOR MECHANISM that enables it to move around a FIELD) into a 2023 FIRST Robotics Competition Event.",
	Evergreen: true,
}, {
	Number:    "H103",
	Title:     "Humans, stay off the FIELD until green.",
	Details:   "Team members may only enter the FIELD if the DRIVER STATION LED strings are green, unless explicitly instructed by a REFEREE or an FTA.",
	Evergreen: true,
}, {
	Number:    "H104",
	Title:     "Never step over the guardrail.",
	Details:   "Team members may only enter or exit the FIELD through open gates.",
	Evergreen: true,
}, {
	Number:    "H105",
	Title:     "Asking other teams to throw a MATCH - not cool.",
	Details:   "A team may not encourage an ALLIANCE, of which it is not a member, to play beneath its ability. NOTE: This rule is not intended to prevent an ALLIANCE from planning and/or executing its own strategy in a specific MATCH in which all the teams are members of the ALLIANCE.",
	Evergreen: true,
}, {
	Number:    "H106",
	Title:     "Letting someone coerce you in to throwing a MATCH - also not cool.",
	Details:   "A team, as the result of encouragement by a team not on their ALLIANCE, may not play beneath its ability. NOTE: This rule is not intended to prevent an ALLIANCE from planning and/or executing its own strategy in a specific MATCH in which all the ALLIANCE members are participants.",
	Evergreen: true,
}, {
	Number:    "H107",
	Title:     "Throwing your own MATCH is bad.",
	Details:   "A team may not intentionally lose a MATCH or sacrifice ranking points in an effort to lower their own ranking or manipulate the rankings of other teams.",
	Evergreen: true,
}, {
	Number:    "H108",
	Title:     "Don't abuse ARENA access.",
	Details:   "Team members (except DRIVERS, HUMAN PLAYERS, and COACHES) granted access to restricted areas in and around the ARENA (e.g. via TECHNICIAN button, event issued Media badges, etc.) may not assist or use signaling devices during the MATCH. Exceptions will be granted for inconsequential infractions and in cases concerning safety.",
	Evergreen: true,
}, {
	Number:    "H109",
	Title:     "Be careful what you interact with.",
	Details:   "Team members are prohibited from the following actions with regards to interaction with ARENA elements. Temporary deformation of a GAME PIECE (e.g.to pre-load a ROBOT) is an exception to this rule.",
	Evergreen: true,
}, {
	Number:    "H110",
	Title:     "Don't mess with GAME PIECES.",
	Details:   "Teams may not modify GAME PIECES in any way. Temporary deformation (e.g.to pre-load a ROBOT) is an exception to this rule.",
	Evergreen: false,
}, {
	Number:    "H201",
	Title:     "Egregious or exceptional violations.",
	Details:   "Egregious behavior beyond what is listed in the rules or subsequent violations of any rule or procedure during the event is prohibited. In addition to rule violations explicitly listed in this manual and witnessed by a REFEREE, the Head REFEREE may assign a YELLOW or RED CARD for egregious ROBOT actions or team member behavior at any time during the event. This includes violations of the event rules found on the FIRST® Robotics Competition District & Regional Events page. Please see Section 11.2.2 YELLOW and RED CARDS for additional detail.",
	Evergreen: true,
}, {
	Number:    "H202",
	Title:     "1 STUDENT, 1 Head REFEREE.",
	Details:   "A team may only address the Head REFEREE with 1 STUDENT. The STUDENT may not be accompanied by more than 1 silent observer.",
	Evergreen: true,
}, {
	Number:    "H301",
	Title:     "Be prompt.",
	Details:   "DRIVE TEAMS may not cause significant delays to the start of a MATCH. Causing a significant delay requires both of the following to be true: The expected MATCH start time has passed, and The DRIVE TEAM is neither MATCH ready nor making a good faith effort, as perceived by the Head REFEREE, to quickly become MATCH ready.",
	Evergreen: true,
}, {
	Number:    "H302",
	Title:     "Teams may not enable their ROBOTS on the FIELD.",
	Details:   "Teams may not tether to the ROBOT while on the FIELD except in special circumstances (e.g. after Opening Ceremonies, before an immediate MATCH replay, etc.) and with the express permission from the FTA or a REFEREE.",
	Evergreen: true,
}, {
	Number:    "H303",
	Title:     "You can't bring/use anything you want.",
	Details:   "The only equipment that may be brought to the ARENA and used by DRIVE TEAMS during a MATCH is listed below. Regardless of if equipment fits criteria below, it may not be employed in a way that breaks any other rules, introduces a safety hazard, blocks visibility for FIELD STAFF or audience members, or jams or interferes with the remote sensing capabilities of another team or the FIELD. A. the OPERATOR CONSOLE, B. non-powered signaling devices, C. reasonable decorative items, D. special clothing and/or equipment required due to a disability, E. devices used solely for planning or tracking strategy, F. devices used solely to record gameplay, and G. non-powered Personal Protective Equipment (examples include, but aren't limited to, gloves, eye protection, and hearing protection). Items brought to the ARENA under allowances B-G must meet all following conditions: I. do not connect or attach to the OPERATOR CONSOLE, FIELD, or ARENA, II. do not connect or attach to another ALLIANCE member (other than items in category G), III. do not communicate with anything or anyone outside of the ARENA, IV. do not communicate with the TECHNICIAN, V. do not include any form of enabled wireless electronic communication with the exception of medically required equipment, and VI. do not in any way affect the outcome of a MATCH, other than by allowing the drive team to a. plan or track strategy for the purposes of communication of that strategy to other ALLIANCE members or b. use items allowed per B to communicate with the ROBOT.",
	Evergreen: true,
}, {

	Number:    "H304",
	Title:     "By invitation only.",
	Details:   "Only DRIVE TEAMS for the current MATCH are allowed in their respective ALLIANCE AREAS and SUBSTATION AREAS.",
	Evergreen: true,
}, {

	Number:    "H305",
	Title:     "Show up to your MATCHES.",
	Details:   "Upon each team's ROBOT passing initial, complete inspection, the team must send at least 1 member of its DRIVE TEAM to the ARENA and participate in each of the team's assigned Qualification and Playoff MATCHES.",
	Evergreen: true,
}, {

	Number:    "H306",
	Title:     "Identify yourself.",
	Details:   "DRIVE TEAMS must wear proper identification while in the ARENA. Proper identification consists of: A. all DRIVE TEAM members wearing their designated buttons above the waist in a clear visible location at all times while in the ARENA B. the COACH wearing the “COACH” button C. the DRIVERS and HUMAN PLAYERS each wearing a “DRIVE TEAM” button D. the TECHNICIAN wearing the “TECHNICIAN” button E. during a Playoff MATCH, the ALLIANCE CAPTAIN clearly displaying the designated ALLIANCE CAPTAIN identifier (e.g. hat or armband)",
	Evergreen: true,
}, {

	Number:    "H307",
	Title:     "Plug in to/be in your DRIVER STATION.",
	Details:   "The OPERATOR CONSOLE must be used in the DRIVER STATION to which the team is assigned, as indicated on the team sign.",
	Evergreen: true,
}, {

	Number:    "H308",
	Title:     "Don't bang on the glass.",
	Details:   "Team members may never strike or hit the DRIVER STATION plastic windows.",
	Evergreen: true,
}, {

	Number:    "H309",
	Title:     "Know your ROBOT setup.",
	Details:   "When placed on the FIELD for a MATCH, each ROBOT must be: A. in compliance with all ROBOT rules, i.e. has passed inspection (for exceptions regarding Practice MATCHES, see Section 10 Inspection & Eligibility Rules), B. the only team-provided item left on the FIELD by the DRIVE TEAM, C. confined to its STARTING CONFIGURATION (reference R102 and R104), D. positioned such that it is fully contained within its COMMUNITY E. not in contact with the CHARGE STATION F. fully supported by FIELD carpet, and G. fully and solely supporting not more than 1 GAME PIECE (as described in Section 6.1 Setup).",
	Evergreen: false,
}, {

	Number:    "H310",
	Title:     "Know your DRIVE TEAM positions.",
	Details:   "Prior to the start of the MATCH, DRIVE TEAM members must be positioned as follows: A. DRIVERS: inside their ALLIANCE AREA and behind the STARTING LINE, B. COACHES: inside their ALLIANCE AREA and behind the STARTING LINE, and C. HUMAN PLAYERS: a. at least one HUMAN PLAYER in their SUBSTATION AREA, b. any remaining HUMAN PLAYERS: inside their ALLIANCE AREA and behind the STARTING LINE, and D. TECHNICIANS: in the event-designated area near the FIELD.",
	Evergreen: false,
}, {

	Number:    "H311",
	Title:     "Leave the GAME PIECES alone.",
	Details:   "Prior to the start of the MATCH, HUMAN PLAYERS may not rearrange the GAME PIECES within the SUBSTATION AREA.",
	Evergreen: false,
}, {

	Number:    "H401",
	Title:     "Behind the lines.",
	Details:   "During AUTO, DRIVE TEAM members in ALLIANCE AREAS and HUMAN PLAYERS in their SUBSTATION AREAS may not contact anything in front of the STARTING LINES, unless for personal or equipment safety or granted permission by a Head REFEREE or FTA.",
	Evergreen: true,
}, {

	Number:    "H402",
	Title:     "Disconnect or set down controllers.",
	Details:   "Prior to the start of the MATCH, any control devices worn or held by HUMAN PLAYERS and/or DRIVERS must be disconnected from the OPERATOR CONSOLE.",
	Evergreen: true,
}, {

	Number:    "H403",
	Title:     "Let the ROBOT do its thing.",
	Details:   "During AUTO, DRIVE TEAMS may not directly or indirectly interact with ROBOTS or OPERATOR CONSOLES unless for personal safety, OPERATOR CONSOLE safety, or pressing an E-Stop.",
	Evergreen: true,
}, {

	Number:    "H501",
	Title:     "COACHES and other teams: hands off the controls.",
	Details:   "A ROBOT shall be operated only by the DRIVERS and/or HUMAN PLAYERS of that team.",
	Evergreen: true,
}, {

	Number:    "H502",
	Title:     "No wandering.",
	Details:   "DRIVE TEAMS may not contact anything outside the area in which they started the MATCH (i.e. the ALLIANCE AREA, the SUBSTATION AREA, or the designated TECHNICIAN space). Exceptions are granted in cases concerning safety and for actions that are inadvertent, MOMENTARY, and inconsequential.",
	Evergreen: true,
}, {

	Number:    "H503",
	Title:     "COACHES, GAME PIECES are off limits.",
	Details:   "COACHES may not touch GAME PIECES, unless for safety purposes.",
	Evergreen: true,
}, {

	Number:    "H504",
	Title:     "GAME PIECES through PORTALS only.",
	Details:   "GAME PIECES may only be introduced to the FIELD A. by a HUMAN PLAYER, B. through a PORTAL, and C. during TELEOP.",
	Evergreen: false,
}, {

	Number:    "H505",
	Title:     "DRIVE TEAMS, watch your reach.",
	Details:   "DRIVE TEAMS may not extend any body part into the SINGLE SUBSTATION PORTAL for a greater-than-MOMENTARY period of time.",
	Evergreen: false,
}, {

	Number:    "R101",
	Title:     "FRAME PERIMETER must be fixed.",
	Details:   "The ROBOT (excluding BUMPERS) must have a FRAME PERIMETER, contained within the BUMPER ZONE and established while in the ROBOT'S STARTING CONFIGURATION, that is comprised of fixed, non-articulated structural elements of the ROBOT. Minor protrusions no greater than 1⁄4 in. (~6 mm) such as bolt heads, fastener ends, weld beads, and rivets are not considered part of the FRAME PERIMETER.",
	Evergreen: true,
}, {

	Number:    "R102",
	Title:     "STARTING CONFIGURATION - no overhang.",
	Details:   "In the STARTING CONFIGURATION (the physical configuration in which a ROBOT starts a MATCH), no part of the ROBOT shall extend outside the vertical projection of the FRAME PERIMETER, with the exception of its BUMPERS and minor protrusions such as bolt heads, fastener ends, rivets, cable ties, etc.",
	Evergreen: true,
}, {

	Number:    "R103",
	Title:     "ROBOT weight limit.",
	Details:   "The ROBOT weight must not exceed 125 lbs. (~56 kg). When determining weight, the basic ROBOT structure and all elements of all additional MECHANISMS that might be used in a single configuration of the ROBOT shall be weighed together (see I103). For the purposes of determining compliance with the weight limitations, the following items are excluded: A. ROBOT BUMPERS, B. ROBOT battery and its associated half of the Anderson cable quick connect/disconnect pair (including no more than 12 in. (~30 cm) of cable per leg, the associated cable lugs, connecting bolts, and insulation), and C. tags used for location detection systems if provided by the event.",
	Evergreen: true,
}, {

	Number:    "R104",
	Title:     "STARTING CONFIGURATION - max size.",
	Details:   "A ROBOT'S STARTING CONFIGURATION may not have a FRAME PERIMETER greater than 120 in. (~304 cm) and may not be more than 4 ft. 6 in. (~137 cm) tall.",
	Evergreen: false,
}, {

	Number:    "R105",
	Title:     "ROBOT extension limit.",
	Details:   "ROBOTS may not extend more than 48 in. (~121 cm) beyond their FRAME PERIMETER.",
	Evergreen: false,
}, {

	Number:    "R201",
	Title:     "No digging into carpet.",
	Details:   "Traction devices must not have surface features that could damage the ARENA (e.g. metal, sandpaper, hard plastic studs, cleats, hook-loop fasteners or similar attachments). Traction devices include all parts of the ROBOT that are designed to transmit any propulsive and/or braking forces between the ROBOT and FIELD carpet.",
	Evergreen: true,
}, {

	Number:    "R202",
	Title:     "No exposed sharp edges.",
	Details:   "Protrusions from the ROBOT and exposed surfaces on the ROBOT shall not pose hazards to the ARENA elements (including GAME PIECES) or people.",
	Evergreen: true,
}, {
	Number:    "R203",
	Title:     "General safety.",
	Details:   "ROBOT parts shall not be made from hazardous materials, be unsafe, cause an unsafe condition, or interfere with the operation of other ROBOTS.",
	Evergreen: true,
}, {

	Number:    "R204",
	Title:     "GAME PIECES stays with the FIELD.",
	Details:   "ROBOTS must allow removal of GAME PIECES from the ROBOT and the ROBOT from FIELD elements while DISABLED and powered off.",
	Evergreen: true,
}, {

	Number:    "R205",
	Title:     "Don't contaminate the FIELD.",
	Details:   "Lubricants may be used only to reduce friction within the ROBOT. Lubricants must not contaminate the FIELD or other ROBOTS.",
	Evergreen: true,
}, {

	Number:    "R206",
	Title:     "Don't damage GAME PIECES.",
	Details:   "ROBOT elements likely to come in contact with a GAME PIECE shall not pose a significant hazard to the GAME PIECE.",
	Evergreen: true,
}, {

	Number:    "R301",
	Title:     "Individual item cost limit.",
	Details:   "No individual, non-KOP item or software shall have a Fair Market Value (FMV) that exceeds $600 USD. The total cost of COMPONENTS purchased in bulk may exceed $600 USD as long as the cost of an individual COMPONENT does not exceed $600 USD.",
	Evergreen: true,
}, {

	Number:    "R302",
	Title:     "Custom parts, generally from this year only.",
	Details:   "FABRICATED ITEMS created before Kickoff are not permitted. Exceptions are: A. OPERATOR CONSOLE, B. BUMPERS, C. battery assemblies as described in R103-B, D. FABRICATED ITEMS consisting of 1 COTS electrical device (e.g. a motor or motor controller) and attached COMPONENTS associated with any of the following modifications:a. wires modified to facilitate connection to a ROBOT (including removal of existing connectors), b. connectors and any materials to secure and insulate those connectors added (note: passive PCBs such as those used to adapt motor terminals to connectors are considered connectors), c. motor shafts modified and/or gears, pulleys, or sprockets added, and d. motors modified with a filtering capacitor as described in the blue box below R625. E. COTS items, or functional equivalents, with any of the following modifications: a. non-functional decoration or labeling, b. assembly of COTS items per manufacturer specs, unless the result constitutes a MAJOR MECHANISM as defined in I101, and c. work that could be reasonably accomplished in fewer than 30 minutes with the use of handheld tools (e.g. drilling a small number of holes in a COTS part).",
	Evergreen: true,
}, {

	Number:    "R303",
	Title:     "Create new designs and software, unless they're public.",
	Details:   "ROBOT software and designs created before Kickoff are only permitted if the source files (complete information sufficient to produce the design) are available publicly prior to Kickoff.",
	Evergreen: true,
}, {

	Number:    "R304",
	Title:     "During an event, only work during pit hours.",
	Details:   "During an event a team is attending (regardless of whether the team is physically at the event location), the team may neither work on nor practice with their ROBOT or ROBOT elements outside of the hours that pits are open, with the following exceptions: A. exceptions listed in R302, other than R302-E-c, B. software development, and C. charging batteries.",
	Evergreen: true,
}}
