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
		Action: func(c *cli.Context) error {
			authKey := os.Getenv("BLUE_ALLIANCE_AUTH_KEY")
			teamNum := os.Args[1]
			ctx := context.Background()
			bac := bamm.NewBAClient(&http.Client{}, authKey)
			team, err := bac.TeamSimple(ctx, "frc"+teamNum)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%s\n\n", team)
			socials, err := bac.TeamSocialMedia(ctx, "frc"+teamNum)
			if err != nil {
				log.Fatal(err)
			}
			for _, social := range socials {
				fmt.Println(social)
			}
			return nil
		},
	}

	app.Run(os.Args)
}
