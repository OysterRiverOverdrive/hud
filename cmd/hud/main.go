package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/oysterriveroverdrive/hud"
	"github.com/urfave/cli/v2"
)

type Recipient struct {
	TeamKey string `json:"team_key"`
}

type Award struct {
	AwardType     int         `json:"award_type"`
	EventKey      string      `json:"event_key"`
	Name          string      `json:"name"`
	RecipientList []Recipient `json:"recipient_list"`
	Year          int         `json:"year"`
}

func main() {
	app := &cli.App{
		Name:  "hud",
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
			// teamNum := c.Int("team")
			// year := c.Int("year")
			eventKey := c.String("event-key")
			// districtKey := c.String("district-key")
			// matchKey := c.String("match-key")
			ctx := context.Background()
			client := hud.NewTriviaService(&http.Client{}, authKey)
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
			teams := []int{
				1027,
				1058,
				1071,
				1073,
				1099,
				1100,
				1124,
				1153,
				121,
				1247,
				125,
				126,
				1277,
				1307,
				131,
				133,
				1350,
				138,
				1474,
				151,
				1512,
				155,
				157,
				166,
				1699,
				172,
				1721,
				1729,
				173,
				1735,
				1740,
				175,
				1757,
				176,
				1761,
				1768,
				177,
				178,
				181,
				1831,
				190,
				1922,
				195,
				1965,
				1991,
				2064,
				2067,
				2079,
				2084,
				2168,
				2170,
				2262,
				228,
				230,
				2342,
				236,
				237,
				2370,
				238,
				2423,
				246,
				2523,
				2648,
				2712,
				2713,
				2785,
				2876,
				2877,
				3146,
				3182,
				319,
				3205,
				3323,
				3451,
				3461,
				3464,
				3467,
				348,
				3566,
				3623,
				3634,
				3654,
				3719,
				3958,
				4034,
				4041,
				4048,
				4097,
				4169,
				4176,
				4311,
				4473,
				4546,
				4564,
				4572,
				4628,
				467,
				4761,
				4796,
				4905,
				4906,
				4908,
				4909,
				4925,
				4987,
				5000,
				501,
				509,
				5112,
				5142,
				5265,
				5347,
				5422,
				5459,
				5491,
				5494,
				5556,
				5563,
				558,
				5687,
				571,
				5735,
				5752,
				58,
				5813,
				5846,
				5856,
				5902,
				5962,
				61,
				6153,
				6161,
				6201,
				6324,
				6328,
				6329,
				6333,
				6346,
				6367,
				6529,
				6620,
				663,
				6690,
				6723,
				6731,
				6762,
				6763,
				6895,
				69,
				6933,
				7127,
				7153,
				716,
				7314,
				7407,
				7462,
				7674,
				7694,
				7760,
				78,
				7822,
				7869,
				7907,
				7913,
				8013,
				8023,
				8046,
				8085,
				811,
				8167,
				839,
				8410,
				8544,
				8567,
				8604,
				8626,
				8708,
				8709,
				8724,
				88,
				8883,
				8889,
				95,
				97,
				999,
			}
			for _, team := range teams {
				for _, year := range []int{2022, 2021, 2020} {
					resp, err := client.Get(client.URL+fmt.Sprintf("/team/frc%d/awards/%d", team, year), nil)
					if err != nil {
						log.Fatal(err)
					}
					as := []Award{}
					if err := json.NewDecoder(resp.Body).Decode(&as); err != nil {
						log.Fatal(err)
					}
					for _, a := range as {
						fmt.Printf("%d\t%d\t%s\n", team, year, a.Name)
					}
					// data, err := ioutil.ReadAll(resp.Body)
					// if err != nil {
					// 	log.Fatal(err)
					// }
					// fname := fmt.Sprintf("awards-%d-%d.json", team, year)
					// err = os.WriteFile(fname, data, 0666)
					// if err != nil {
					// 	log.Fatal(err)
					// }
				}
			}

			// client.Dump(ctx, year, teamNum, eventKey, districtKey, matchKey)
			return nil
			matches, err := client.EventMatchesSimple(ctx, eventKey)
			if err != nil {
				log.Fatal(err)
			}
			next, err := hud.NextMatch(matches, 8410)
			if err != nil {
				log.Fatal(err)
			}
			summary := hud.PrintNextMatchSummary(nil, next, 8410)
			fmt.Println(strings.Join(summary, "\n"))
			return nil
		},
	}

	app.Run(os.Args)
}
