# HUD

Heads-Up Display (HUD) is FRC team 8410's data management system
that primarily integrates into the FRC Blue Alliance API
https://www.thebluealliance.com/apidocs/v3.

Goals:

* Automate the collection and distribution of data.
* Facilitate scouting operations.

## Setup

Create a read API key at https://www.thebluealliance.com/account

```
go get github.com/oysterriveroverdrive/hud/cmd/hud
export BLUE_ALLIANCE_AUTH_KEY=[YOUR KEY]
cd cmd/hud
go run . [TEAM NUMBER]
```

## Use Case

```

hub.CompetitionWatcher(8410)
hub.WatchTeam(8410)
```

```
hub.TeamSummary(8410)
```
