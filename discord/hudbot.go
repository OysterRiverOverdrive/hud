package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/oysterriveroverdrive/hud"
)

// HudBot contains the subcommand structure for handling message processing.
type HudBot struct {
	ts      *hud.TriviaService
	SubCmds []Cmd
}

func (hb *HudBot) Handle(s *discordgo.Session, m *discordgo.MessageCreate) {
	// HudBot will only respond if @mentioned directly.
	if !isBotMentioned(s.State.User, m.Mentions) {
		return
	}
	// Prevent talking to yourself.
	if isFromBot(s.State.User, m.Author) {
		return
	}

	// Since we know the bot has been mentioned, cut out the bot mention to
	// help process the rest of the command.
	msg := m.Content
	suffix := strings.TrimSpace(strings.Replace(msg, fmt.Sprintf("<@%s>", s.State.User.ID), "", -1))

	// Find out if we know how to respond to this message.

	// If help
	if suffix == "help" || suffix == "" {
		var help []string
		for _, subCmd := range hb.SubCmds {
			help = append(help, "@hud "+subCmd.Help())
		}
		s.ChannelMessageSend(m.ChannelID, "hud help:\n"+strings.Join(help, "\n"))
	}
	for _, subCmd := range hb.SubCmds {
		if subCmd.Match(suffix) {
			md := map[string]string{
				"path": "@hud",
			}
			ch, m, err := subCmd.Handle(md, hb.ts, s, m, suffix)
			if err != nil {
				log.Printf("[ERROR] %s", err)
			}
			if m == nil {
				continue
			}
			s.ChannelMessageSend(ch, m.Content)
		}
	}
}

// IsBotMentioned is a helper function for eeeing if the bot id is among
// the list of users mentioned in a message.
func isBotMentioned(bot *discordgo.User, mentions []*discordgo.User) bool {
	for _, mention := range mentions {
		if mention.ID == bot.ID {
			return true
		}
	}
	return false
}

// Ignore all messages created by the bot itself
// This isn't required in this specific example but it's a good practice.
func isFromBot(bot *discordgo.User, author *discordgo.User) bool {
	return bot.ID == author.ID
}
