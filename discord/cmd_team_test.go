package main

import (
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
)

func TestTeam_Match(t *testing.T) {
	c := &TeamCmd{}
	assert.Equal(t, true, c.Match("team 8410"))
	assert.Equal(t, false, c.Match("not a team message"))
}

func TestTeam_HandleSubCmd(t *testing.T) {
	th := newTestHarness()

	c := &TeamCmd{
		SubCmds: []Cmd{&CatchAllCmd{}},
	}
	md := map[string]string{}
	dChan, m, err := c.Handle(md, th.ts, nil, &discordgo.MessageCreate{
		Message: &discordgo.Message{
			ChannelID: "abcd",
			Content:   "<@1234> team 8410",
		}},
		"team 8410")
	assert.NoError(t, err)
	assert.Equal(t, "abcd", dChan)
	assert.Equal(t, &discordgo.MessageSend{
		Content: "catch all",
	}, m)
}

func TestTeam_HandleHelp(t *testing.T) {
	th := newTestHarness()

	c := &TeamCmd{
		SubCmds: []Cmd{&CatchAllCmd{}},
	}
	md := map[string]string{
		"path": "@hud",
	}
	dChan, m, err := c.Handle(md, th.ts, nil, &discordgo.MessageCreate{
		Message: &discordgo.Message{
			ChannelID: "abcd",
			Content:   "<@1234> team help",
		}},
		"team help")
	assert.NoError(t, err)
	assert.Equal(t, "abcd", dChan)
	assert.Equal(t, &discordgo.MessageSend{
		Content: "team help:\n@hud team [anything] - catch all",
	}, m)
}

func TestTeamID_Match(t *testing.T) {
	c := &TeamIDCmd{}
	assert.Equal(t, true, c.Match("8410"))
	assert.Equal(t, false, c.Match("notanumber"))
}

func TestTeamID_Handle(t *testing.T) {
	th := newTestHarness()
	th.baResponse = `{
	"city": "Durham",
	"nickname": "Oyster River Robotics",
	"state_prov": "New Hampshire"
}`

	c := &TeamIDCmd{}
	md := map[string]string{}
	dChan, m, err := c.Handle(md, th.ts, nil, &discordgo.MessageCreate{
		Message: &discordgo.Message{
			ChannelID: "abcd",
			Content:   "<@1234> team 8410",
		}},
		"8410")
	assert.NoError(t, err)
	assert.Equal(t, "abcd", dChan)
	assert.Equal(t, &discordgo.MessageSend{
		Content: "Oyster River Robotics from Durham, New Hampshire",
	}, m)
	assert.Equal(t, "/team/frc8410", th.baRequest.URL.Path)
}
