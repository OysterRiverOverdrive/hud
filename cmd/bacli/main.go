package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/mathyourlife/bluealliance"
)

func main() {
	authKey := os.Getenv("BLUE_ALLIANCE_AUTH_KEY")
	teamNum := os.Args[1]
	ctx := context.Background()
	bac := bluealliance.NewBAClient(&http.Client{}, authKey)
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
}
