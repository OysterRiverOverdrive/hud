package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

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
			// 2022macma - shrewsburrp
			// 2022nhsea - pease
			&cli.StringFlag{
				Name:  "event-key",
				Value: "2022nhsea",
			},
			&cli.StringFlag{
				Name:  "district-key",
				Value: "2022ne",
			},
			&cli.StringFlag{
				Name:  "match-key",
				Value: "2022nhsea_qm53",
			},
			&cli.IntFlag{
				Name:  "team",
				Value: 8410,
			},
			&cli.IntFlag{
				Name:  "year",
				Value: 2022,
			},
		},
		Action: func(c *cli.Context) error {
			authKey := c.String("api-key")
			teamNum := c.Int("team")
			year := c.Int("year")
			eventKey := c.String("event-key")
			districtKey := c.String("district-key")
			matchKey := c.String("match-key")
			ctx := context.Background()
			bc := bamm.NewBAClient(&http.Client{}, authKey)
			// team, err := bc.TeamSimple(ctx, fmt.Sprintf("frc%d", teamNum))
			// if err != nil {
			// 	log.Fatal(err)
			// }
			// fmt.Printf("%s\n", team)
			// socials, err := bc.TeamSocialMedia(ctx, fmt.Sprintf("frc%d", teamNum))
			// if err != nil {
			// 	log.Fatal(err)
			// }
			// for _, social := range socials {
			// 	if social.Type == "github-profile" && social.ForeignKey != "" {
			// 		fmt.Printf("Github: https://github.com/%s\n", social.ForeignKey)
			// 	}
			// }

			// teams, err := bc.EventTeams(ctx, eventKey)
			// if err != nil {
			// 	log.Fatal(err)
			// }
			// for _, team := range teams {
			// 	fmt.Printf("%d\t%s\n", team.TeamNumber, team.Nickname)
			// }
			bc.Dump(ctx, year, teamNum, eventKey, districtKey, matchKey)
			return nil
			matches, err := bc.EventMatchesSimple(ctx, eventKey)
			if err != nil {
				log.Fatal(err)
			}
			next, err := bamm.NextMatch(matches, 8410)
			if err != nil {
				log.Fatal(err)
			}
			summary := bamm.PrintNextMatchSummary(next, 8410)
			fmt.Println(strings.Join(summary, "\n"))
			return nil
		},
	}

	app.Run(os.Args)
}
