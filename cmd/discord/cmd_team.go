package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/oysterriveroverdrive/hud"
	"github.com/oysterriveroverdrive/hud/bluealliance"
	"github.com/sirupsen/logrus"
)

// TeamCmd handles @hud team ... commands
type TeamCmd struct {
	SubCmds []Cmd
}

func (c *TeamCmd) Stub() string {
	return "team"
}

func (c *TeamCmd) Match(msg string) bool {
	logrus.Debugf("TeamCmd.Match %q", msg)
	return strings.HasPrefix(strings.TrimSpace(msg), "team")
}

func (c *TeamCmd) Help() string {
	return c.Stub() + " - working with frc teams"
}

func (c *TeamCmd) Handle(md map[string]string, ts *bluealliance.Service, s *discordgo.Session, m *discordgo.MessageCreate, msg string) (string, []*discordgo.MessageSend, error) {
	logrus.Debugf("TeamCmd.Handle %v %q", md, msg)
	suffix := strings.TrimSpace(strings.TrimPrefix(msg, "team"))
	if suffix == "help" || suffix == "" {
		var help []string
		for _, subCmd := range c.SubCmds {
			help = append(help, md["path"]+" "+c.Stub()+" "+subCmd.Help())
		}
		return m.ChannelID, []*discordgo.MessageSend{{
			Content: "team help:\n" + strings.Join(help, "\n"),
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

// TeamIDCmd handles @hud team [id] commands
type TeamIDCmd struct{}

func (c *TeamIDCmd) Stub() string {
	return "[number,number,...]"
}

func (c *TeamIDCmd) Match(msg string) bool {
	logrus.Debugf("TeamIDCmd.Match %q", msg)
	return len(c.parseTeamNumbers(msg)) > 0
}

func (c *TeamIDCmd) Help() string {
	return c.Stub() + " - request frc team data (request multiple teams with team numbers separated by commas 1234,5678)"
}

func (c *TeamIDCmd) Handle(md map[string]string, ts *bluealliance.Service, s *discordgo.Session, m *discordgo.MessageCreate, msg string) (string, []*discordgo.MessageSend, error) {
	logrus.Debugf("TeamIDCmd.Handle %v %q", md, msg)
	teamNums := c.parseTeamNumbers(msg)

	var summaries []string
	for _, teamNum := range teamNums {

		team, err := hud.TeamByNumber(ts, teamNum)
		if err != nil {
			log.Printf("[ERROR] %s", err)
			summaries = append(summaries, fmt.Sprintf("%d: ERROR %s", teamNum, err))
			continue
		}
		summaries = append(summaries, fmt.Sprintf("%d: %s from %s, %s", team.Data.Number, team.Data.Nickname, team.Data.City, team.Data.StateProv))
	}
	return m.ChannelID, []*discordgo.MessageSend{{
		Content: strings.Join(summaries, "\n"),
	}}, nil
}

func (c *TeamIDCmd) parseTeamNumbers(msg string) []int {
	var teamNums []int
	commaSplit := strings.Split(msg, ",")
	for _, split := range commaSplit {
		split := strings.TrimSpace(split)
		teamNum, err := strconv.Atoi(split)
		if err != nil {
			break
		}
		teamNums = append(teamNums, teamNum)
	}
	return teamNums
}
