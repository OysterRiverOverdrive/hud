package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/oysterriveroverdrive/hud"
	"github.com/sirupsen/logrus"
)

var (
	Token   string
	BAToken string
)

func parseFlags() {
	flag.StringVar(&Token, "discord-token", "", "Discord Bot Token")
	flag.StringVar(&BAToken, "blue-alliance-token", "", "Blue Alliance Token")
	flag.Parse()
}

func main() {
	fmt.Println("starting hud bot")

	logrus.SetLevel(logrus.DebugLevel)

	parseFlags()

	ts := hud.NewTriviaService(&http.Client{}, BAToken)

	disc, err := discordgo.New("Bot " + Token)
	if err != nil {
		log.Fatal(err)
	}
	disc.Client = &http.Client{Timeout: (20 * time.Second)}

	// This creates the structure of the commands hud will respond to.
	// A "@hud team 8410" command would be proccessed through the
	// TeamCmd - for the "team"
	// TeamIDCmd - for the "8410"
	hb := &HudBot{
		ts: ts,
		SubCmds: []Cmd{
			&TeamCmd{
				SubCmds: []Cmd{
					&TeamIDCmd{},
				},
			},
			&DistrictCmd{
				SubCmds: []Cmd{
					&DistrictIDCmd{
						SubCmds: []Cmd{
							&DistrictIDTeamsCmd{},
						},
					},
				},
			},
		},
	}
	disc.AddHandler(hb.Handle)

	// Let HudBot read messages.
	disc.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = disc.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	disc.Close()
}
