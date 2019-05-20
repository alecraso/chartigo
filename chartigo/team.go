package chartigo

import (
        "fmt"
)

type (
        // Team represents a team object in Chartio
        Team struct {
                ID        int      `json:"id,omitempty"`
                Name      string   `json:"name,omitempty"`
                IsOwner   bool     `json:"is_owner,omitempty"`
                CreatedAt JSONTime `json:"created_at,omitempty"`
                UpdatedAt JSONTime `json:"updated_at,omitempty"`
                Links     Links    `json:"_links,omitempty"`
                Embedded  Users    `json:"_embedded,omitempty"`
        }

        // Teams represents an array of teams as returned by the Chartio API
        Teams struct {
                Teams []Team `json:"teams,omitempty"`
                Count int    `json:"count,omitempty"`
        }
)

// ListTeams returns the full list of teams for the current organization.
func (c *Client) ListTeams() (*Teams, error) {
        teams := new(Teams)
        _, err := c.request("GET", "/teams", teams)
        return teams, err
}

// CreateTeamInput is used as input to the CreateTeam function.
type CreateTeamInput struct {
        Name string
}

// CreateTeam creates a new team with the given name
// func (c *Client) CreateTeam(i CreateTeamInput) (*Team, error) {
//         return
// }

// GetTeamInput is used as input to the GetTeam function.
type GetTeamInput struct {
        ID string
}

// GetTeam returns a particular team looked up by ID
func (c *Client) GetTeam(i GetTeamInput) (*Team, error) {
        team := new(Team)
        path := fmt.Sprintf("/teams/%s", i.ID)
        _, err := c.request("GET", path, team)
        return team, err
}
