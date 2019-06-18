package chartigo

import (
        "fmt"
)

type (
        // Team represents a team object in Chartio.
        Team struct {
                ID        int        `json:"id"`
                Name      string     `json:"name"`
                IsOwner   bool       `json:"is_owner,omitempty"`
                CreatedAt *ChartioTS `json:"created_at,omitempty"`
                UpdatedAt *ChartioTS `json:"updated_at,omitempty"`
                Users     Users      `json:"_embedded,omitempty"`
                Links     Links      `json:"_links,omitempty"`
        }

        // Teams represents an array of teams as returned by the Chartio API.
        Teams struct {
                Teams []Team `json:"teams,omitempty"`
                Count *int   `json:"count,omitempty"`
                Links Links  `json:"_links,omitempty"`
        }

        // TeamInput is used as input on all team object functions.
        TeamInput struct {
                TeamID    string `json:"-"`
                Name      string `json:"name,omitempty"`
                UserID    string `json:"id,omitempty"`
                UserEmail string `json:"email,omitempty"`
        }
)

func (i *TeamInput) user() string {
        if i.UserID != "" {
                return i.UserID
        } else {
                return i.UserEmail
        }
}

// ListTeams returns the full list of teams for the current organization.
func (c *Client) ListTeams() (*Teams, error) {
        teams := new(Teams)
        _, err := c.Get("/teams", teams)
        return teams, err
}

// CreateTeam creates a new team with the given name
func (c *Client) CreateTeam(i TeamInput) (*Team, error) {
        team := new(Team)
        _, err := c.Post("/teams", i, team)
        return team, err
}

// GetTeam returns a particular team looked up by ID
func (c *Client) GetTeam(i TeamInput) (*Team, error) {
        team := new(Team)
        path := fmt.Sprintf("/teams/%s", i.TeamID)
        _, err := c.Get(path, team)
        return team, err
}

// DeleteTeam deletes a particular team looked up by ID
func (c *Client) DeleteTeam(i TeamInput) error {
        path := fmt.Sprintf("/teams/%s", i.TeamID)
        _, err := c.Delete(path)
        return err
}

// RenameTeam adds a new user to the given team.
func (c *Client) RenameTeam(i TeamInput) (*Team, error) {
        team := new(Team)
        path := fmt.Sprintf("/teams/%s", i.TeamID)
        _, err := c.Patch(path, i, team)
        return team, err
}

// AddTeamUser adds a new user to the given team.
func (c *Client) AddTeamUser(i TeamInput) (*Team, error) {
        team := new(Team)
        path := fmt.Sprintf("/teams/%s/users", i.TeamID)
        _, err := c.Patch(path, i, team)
        return team, err
}

// DeleteTeamUser deletes a user from the given team.
func (c *Client) DeleteTeamUser(i TeamInput) error {
        path := fmt.Sprintf("/teams/%s/users/%s", i.TeamID, i.user())
        _, err := c.Delete(path)
        return err
}
