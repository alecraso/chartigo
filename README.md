# ChartiGo
ChartiGo is a golang API client for interacting with the [Chartio API](https://api.chartio.com/v1/docs).

## Installation
This is a client library so there is nothing to install.

## Usage
Download the library into your `$GOPATH`:
```bash
go get github.com/aaronbiller/chartigo/chartigo
```

Import the library into your tool:

```go
import "github.com/aaronbiller/chartigo/chartigo"
```

Put your API credentials in environment variables:
```bash
export CHARTIO_API_KEY='myapikey'
export CHARTIO_API_PASSWORD='muchsecrethowsecuritywow'
```

#### Usage Example
```go
package main

import (
        "fmt"
        "log"

        "github.com/aaronbiller/chartigo/chartigo"
)

func main() {
        client := chartigo.NewClient("my_org")

        teams, err := client.ListTeams()
        if err != nil {
                log.Fatal(err)
        }
        for _, t := range teams.Teams {
                fmt.Printf("%d, %s\n", t.ID, t.Name)
        }

        teamInput := chartigo.TeamInput{TeamID: "owners"}
        team, err := client.GetTeam(teamInput)
        if err != nil {
                log.Fatal(err)
        }
        fmt.Printf("NAME:  %s\n", team.Name)
        fmt.Printf("USERS: %d\n", len(team.Users.Users))
        for _, u := range team.Users.Users {
                fmt.Printf("    %s\n", u.Email)
        }
        
}

```
