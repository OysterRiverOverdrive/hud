package main

import (
	"fmt"
	"log"
	"regexp"
	"sort"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/oysterriveroverdrive/hud"
	"github.com/oysterriveroverdrive/hud/bluealliance"
	"github.com/sirupsen/logrus"
)

// DistrictCmd handles @hud district ... commands
type DistrictCmd struct {
	SubCmds []Cmd
}

func (c *DistrictCmd) Stub() string {
	return "district"
}

func (c *DistrictCmd) Match(msg string) bool {
	// TODO: match on "@hud district" but not "@hud districtsometing"
	logrus.Debugf("DistrictCmd.Match %q", msg)
	return strings.HasPrefix(msg, "district")
}

func (c *DistrictCmd) Help() string {
	return c.Stub() + " - working with frc districts"
}

func (c *DistrictCmd) Handle(md map[string]string, ts *bluealliance.Service, s *discordgo.Session, m *discordgo.MessageCreate, msg string) (string, []*discordgo.MessageSend, error) {
	logrus.Debugf("DistrictCmd.Handle %v %q", md, msg)
	suffix := strings.TrimSpace(strings.TrimPrefix(msg, "district"))

	if suffix == "help" || suffix == "" {
		var help []string
		for _, subCmd := range c.SubCmds {
			help = append(help, md["path"]+" "+c.Stub()+" "+subCmd.Help())
		}
		return m.ChannelID, []*discordgo.MessageSend{{
			Content: "district help:\n" + strings.Join(help, "\n"),
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

// DistrictIDCmd handles @hud district [id] ... commands
type DistrictIDCmd struct {
	SubCmds []Cmd
}

func (c *DistrictIDCmd) Stub() string {
	return "[id]"
}

func (c *DistrictIDCmd) Match(msg string) bool {
	logrus.Debugf("DistrictIDCmd.Match %q", msg)
	return regexp.MustCompile(`\s*[a-z\d]+\s*`).MatchString(msg)
}

func (c *DistrictIDCmd) Help() string {
	return c.Stub() + " - working with a frc district"
}

func (c *DistrictIDCmd) Handle(md map[string]string, ts *bluealliance.Service, s *discordgo.Session, m *discordgo.MessageCreate, msg string) (string, []*discordgo.MessageSend, error) {
	logrus.Debugf("DistrictIDCmd.Handle %v %q", md, msg)
	match := regexp.MustCompile(`\s*([a-z\d]+)\s*(.*)`).FindStringSubmatch(msg)
	var suffix string
	if len(match) > 1 {
		suffix = strings.TrimSpace(match[2])
	}

	if suffix == "help" || suffix == "" {
		var help []string
		for _, subCmd := range c.SubCmds {
			help = append(help, md["path"]+" "+c.Stub()+" "+subCmd.Help())
		}
		return m.ChannelID, []*discordgo.MessageSend{{
			Content: "district [id] help:\n" + strings.Join(help, "\n"),
		}}, nil
	}
	for _, subCmd := range c.SubCmds {
		if subCmd.Match(suffix) {
			md["path"] += " " + c.Stub()
			md["district_id"] = match[1]
			return subCmd.Handle(md, ts, s, m, suffix)
		}
	}
	return "", nil, nil
}

// DistrictIDTeamsCmd handles @hud district [id] teams commands
type DistrictIDTeamsCmd struct{}

func (c *DistrictIDTeamsCmd) Stub() string {
	return "teams"
}

func (c *DistrictIDTeamsCmd) Match(msg string) bool {
	logrus.Debugf("DistrictIDTeamsCmd.Match %q", msg)
	return msg == "teams"
}

func (c *DistrictIDTeamsCmd) Help() string {
	return c.Stub() + " - list district teams"
}

func (c *DistrictIDTeamsCmd) Handle(md map[string]string, ts *bluealliance.Service, s *discordgo.Session, m *discordgo.MessageCreate, msg string) (string, []*discordgo.MessageSend, error) {
	logrus.Debugf("DistrictIDTeamsCmd.Handle %v %q", md, msg)
	teams, err := hud.TeamsInDistrict(ts, md["district_id"])
	if err != nil {
		log.Printf("[ERROR] %s", err)
		return "", nil, nil
	}
	teamNums := make([]int, 0, len(teams.Data))
	for _, team := range teams.Data {
		teamNums = append(teamNums, team.Number)
	}
	sort.Ints(teamNums)
	var resp []string
	for _, teamNum := range teamNums {
		resp = append(resp, fmt.Sprintf("%d", teamNum))
	}
	return m.ChannelID, []*discordgo.MessageSend{{
		Content: strings.Join(resp, ","),
	}}, nil
}
