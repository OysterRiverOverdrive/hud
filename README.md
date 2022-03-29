# BAMM

B.A.M.M. - Blue Alliance Match Manager

Starting an query lib for the FRC Blue Alliance API 
https://www.thebluealliance.com/apidocs/v3.  Will evetually 
help out with providing data to the team in the pits.

Not offiliated .

## Setup

Create a read API key at https://www.thebluealliance.com/account

```
go get github.com/mathyourlife/bamm/cmd/bamm
export BLUE_ALLIANCE_AUTH_KEY=[YOUR KEY]
cd cmd/bamm
go run . [TEAM NUMBER]
```