package main

import (
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
)

func TestRules_Match(t *testing.T) {
	c := &RulesCmd{}
	assert.Equal(t, true, c.Match("rules search something"))
	assert.Equal(t, false, c.Match("not a rules message"))
}

func TestRules_HandleSubCmd(t *testing.T) {
	th := newTestHarness()

	c := &RulesCmd{
		SubCmds: []Cmd{&CatchAllCmd{}},
	}
	md := map[string]string{}
	dChan, m, err := c.Handle(md, th.ts, nil, &discordgo.MessageCreate{
		Message: &discordgo.Message{
			ChannelID: "abcd",
			Content:   "<@1234> rules search something",
		}},
		"rules search something")
	assert.NoError(t, err)
	assert.Equal(t, "abcd", dChan)
	assert.Equal(t, []*discordgo.MessageSend{{
		Content: "catch all",
	}}, m)
}

func TestRules_HandleHelp(t *testing.T) {
	th := newTestHarness()

	c := &RulesCmd{
		SubCmds: []Cmd{&CatchAllCmd{}},
	}
	md := map[string]string{
		"path": "@hud",
	}
	dChan, m, err := c.Handle(md, th.ts, nil, &discordgo.MessageCreate{
		Message: &discordgo.Message{
			ChannelID: "abcd",
			Content:   "<@1234> rules help",
		}},
		"rules help")
	assert.NoError(t, err)
	assert.Equal(t, "abcd", dChan)
	assert.Equal(t, []*discordgo.MessageSend{{
		Content: "rules help:\n@hud rules [anything] - catch all",
	}}, m)
}

func TestRulesSearch_Match(t *testing.T) {
	c := &RulesSearchCmd{}
	assert.Equal(t, true, c.Match("search something"))
	assert.Equal(t, true, c.Match("search SOMETHING"))
	assert.Equal(t, false, c.Match("not a search command"))
}

func TestRulesSearch_Handle(t *testing.T) {
	th := newTestHarness()

	c := &RulesSearchCmd{}
	md := map[string]string{}
	dChan, m, err := c.Handle(md, th.ts, nil, &discordgo.MessageCreate{
		Message: &discordgo.Message{
			ChannelID: "abcd",
			Content:   "<@1234> rules search height",
		}},
		"search height")
	assert.NoError(t, err)
	assert.Equal(t, "abcd", dChan)
	assert.Equal(t, []*discordgo.MessageSend{{
		Content: "Rule Number: G106\nTitle: Tall ROBOTS not allowed.\nDetails: ROBOT height, as measured when it's resting normally on a flat floor, may not exceed 6 ft. 6 in. (~198 cm)) above the carpet during the MATCH.",
	}}, m)

	dChan, m, err = c.Handle(md, th.ts, nil, &discordgo.MessageCreate{
		Message: &discordgo.Message{
			ChannelID: "abcd",
			Content:   "<@1234> rules search community",
		}},
		"search community")
	assert.NoError(t, err)
	assert.Equal(t, "abcd", dChan)
	assert.Equal(t, []*discordgo.MessageSend{{
		Content: "Too many hits. Rule details removed. Use @hud rules [RuleNumber] for more information.\nRule Number: G108\nTitle: Opponent's zone, no extension.\n------------------\nRule Number: G109\nTitle: Don't extend in multiple directions.\n------------------\nRule Number: G207\nTitle: Right of way.\n------------------\nRule Number: G208\nTitle: Don't climb on each other unless in the COMMUNITY.\n------------------\nRule Number: G403\nTitle: 1 GAME PIECE at a time (except in LOADING ZONE and COMMUNITY).\n------------------\nRule Number: G404\nTitle: Launching GAME PIECES is only okay in the COMMUNITY.\n------------------\nRule Number: H309\nTitle: Know your ROBOT setup.",
	}}, m)

}
