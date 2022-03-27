# bluealliance

Starting an query lib for the FRC Blue Alliance API https://www.thebluealliance.com/apidocs/v3 . This is not an official TBA project.

## Setup

Create a read API key at https://www.thebluealliance.com/account

```
go get github.com/mathyourlife/bluealliance/cmd/bacli
export BLUE_ALLIANCE_AUTH_KEY=[YOUR KEY]
cd cmd/bacli
go run . [TEAM NUMBER]
```