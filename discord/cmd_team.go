package main

import (
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/oysterriveroverdrive/hud"
)

// TeamCmd handles @hud team ... commands
type TeamCmd struct {
	SubCmds []Cmd
}

func (c *TeamCmd) Stub() string {
	return "team"
}

func (c *TeamCmd) Match(msg string) bool {
	return strings.HasPrefix(strings.TrimSpace(msg), "team")
}

func (c *TeamCmd) Help() string {
	return c.Stub() + " - working with frc teams"
}

func (c *TeamCmd) Handle(md map[string]string, ts *hud.TriviaService, s *discordgo.Session, m *discordgo.MessageCreate, msg string) (string, *discordgo.MessageSend, error) {
	suffix := strings.TrimSpace(strings.TrimPrefix(msg, "team"))
	if suffix == "help" || suffix == "" {
		var help []string
		for _, subCmd := range c.SubCmds {
			help = append(help, md["path"]+" "+c.Stub()+" "+subCmd.Help())
		}
		return m.ChannelID, &discordgo.MessageSend{
			Content: "team help:\n" + strings.Join(help, "\n"),
		}, nil
	}
	for _, subCmd := range c.SubCmds {
		if subCmd.Match(msg) {
			md["path"] += " " + c.Stub()
			return subCmd.Handle(md, ts, s, m, suffix)
		}
	}
	return "", nil, nil
}

// TeamIDCmd handles @hud team [id] commands
type TeamIDCmd struct{}

func (c *TeamIDCmd) Stub() string {
	return "[number]"
}

func (c *TeamIDCmd) Match(msg string) bool {
	return regexp.MustCompile(`\s*\d+\s*`).MatchString(msg)
}

func (c *TeamIDCmd) Help() string {
	return c.Stub() + " - request frc team data"
}

func (c *TeamIDCmd) Handle(md map[string]string, ts *hud.TriviaService, s *discordgo.Session, m *discordgo.MessageCreate, msg string) (string, *discordgo.MessageSend, error) {
	teamNum, err := strconv.Atoi(strings.TrimSpace(msg))
	if err != nil {
		// Didn't get an expected team number.
		return "", nil, nil
	}

	team, err := hud.TeamByNumber(ts, teamNum)
	if err != nil {
		log.Printf("[ERROR] %s", err)
		return "", nil, nil
	}
	return m.ChannelID, &discordgo.MessageSend{
		Content: team.Data.Nickname + " from " + team.Data.City + ", " + team.Data.StateProv,
	}, nil
}
