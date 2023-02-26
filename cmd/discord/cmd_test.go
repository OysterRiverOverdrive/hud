package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/bwmarrin/discordgo"
	"github.com/oysterriveroverdrive/hud/bluealliance"
)

type testHarness struct {
	ba         *httptest.Server
	baRequest  *http.Request
	ts         *bluealliance.Service
	baResponse string
}

func newTestHarness() *testHarness {
	th := &testHarness{}
	th.ba = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		th.baRequest = r
		fmt.Fprintf(w, th.baResponse)
		log.Println("test server handling")
	}))

	th.ts = bluealliance.NewService(th.ba.Client(), "test-token")
	th.ts.URL = th.ba.URL

	return th
}

func (th *testHarness) Close() {
	th.ba.Close()
}

type CatchAllCmd struct{}

func (c *CatchAllCmd) Stub() string {
	return "[anything]"
}

func (c *CatchAllCmd) Match(msg string) bool {
	return true
}

func (c *CatchAllCmd) Help() string {
	return c.Stub() + " - catch all"
}

func (c *CatchAllCmd) Handle(md map[string]string, ts *bluealliance.Service, s *discordgo.Session, m *discordgo.MessageCreate, msg string) (string, []*discordgo.MessageSend, error) {
	return m.ChannelID, []*discordgo.MessageSend{{
		Content: "catch all",
	}}, nil
}
