package chartigo

import (
        "fmt"
        "time"
)

type (
        // Team represents a team object in Chartio
        Team struct {
                ID        int        `mapstructure:"id"`
                Name      string     `mapstructure:"name"`
                IsOwner   bool       `mapstructure:"is_owner"`
                CreatedAt *time.Time `mapstructure:"created_at"`
                UpdatedAt *time.Time `mapstructure:"updated_at"`
                Links     Links      `mapstructure:"_links"`
                Embedded  Users      `mapstructure:"_embedded"`
        }

        // Teams represents an array of teams as returned by the Chartio API
        Teams struct {
                Teams []Team `mapstructure:"teams"`
                Count int    `mapstructure:"count"`
        }
)

// ListTeams returns the full list of teams for the current organization.
func (c *Client) ListTeams() (*Teams, error) {
        teams := new(Teams)
        _, err := c.Get("/teams", teams)
        return teams, err
}

// CreateTeamInput is used as input to the CreateTeam function.
type CreateTeamInput struct {
        Name string `json:"name"`
}

// CreateTeam creates a new team with the given name
func (c *Client) CreateTeam(i CreateTeamInput) (*Team, error) {
        team := new(Team)
        _, err := c.Post("/teams", i, team)
        return team, err
}

// GetTeamInput is used as input to the GetTeam function.
type GetTeamInput struct {
        ID string `json:"id"`
}

// GetTeam returns a particular team looked up by ID
func (c *Client) GetTeam(i GetTeamInput) (*Team, error) {
        team := new(Team)
        path := fmt.Sprintf("/teams/%s", i.ID)
        _, err := c.Get(path, team)
        return team, err
}

// DeleteTeamInput is used as input to the DeleteTeam function.
type DeleteTeamInput struct {
        ID string `json:"id"`
}

// DeleteTeam deletes a particular team looked up by ID
func (c *Client) DeleteTeam(i DeleteTeamInput) error {
        path := fmt.Sprintf("/teams/%s", i.ID)
        _, err := c.Delete(path)
        return err
}
