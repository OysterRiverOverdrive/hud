package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/oysterriveroverdrive/hud/bluealliance"
)

// Cmd defines the methods required for all HudBot commands and subcommands.
type Cmd interface {
	// Stub returns the identifier for how the subcommand is triggered. Useful in referencing in
	// help command outputs.
	Stub() string
	// Match determines if the supplied message should be processed by the subcommand.
	Match(msg string) bool
	// Help returns directions on what the subcommand is used for.
	Help() string
	// Handle processes the message.
	Handle(md map[string]string, ts *bluealliance.Service, s *discordgo.Session, m *discordgo.MessageCreate, msg string) (string, []*discordgo.MessageSend, error)
}
