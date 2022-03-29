package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/mathyourlife/bamm"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "bamm",
		Usage: "A cli tool for the blue alliance api.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "api-key",
				EnvVars:  []string{"BLUE_ALLIANCE_AUTH_KEY"},
				Required: true,
			},
			&cli.IntFlag{
				Name:     "team",
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			authKey := c.String("api-key")
			teamNum := c.Int("team")
			ctx := context.Background()
			bac := bamm.NewBAClient(&http.Client{}, authKey)
			team, err := bac.TeamSimple(ctx, fmt.Sprintf("frc%d", teamNum))
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%s\n", team)
			socials, err := bac.TeamSocialMedia(ctx, fmt.Sprintf("frc%d", teamNum))
			if err != nil {
				log.Fatal(err)
			}
			for _, social := range socials {
				if social.Type == "github-profile" && social.ForeignKey != "" {
					fmt.Printf("Github: https://github.com/%s\n", social.ForeignKey)
				}
			}
			return nil
		},
	}

	app.Run(os.Args)
}
