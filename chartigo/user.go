package chartigo

import (
        "fmt"
)

type (
        // User represents a user object in Chartio
        User struct {
                ID          int        `json:"id"`
                DisplayName string     `json:"display_name,omitempty"`
                Email       string     `json:"email"`
                CreatedAt   *ChartioTS `json:"created_at,omitempty"`
                UpdatedAt   *ChartioTS `json:"updated_at,omitempty"`
                Teams       *Teams     `json:"_embedded,omitempty"`
                Links       Links      `json:"_links,omitempty"`
        }

        // Users represents an array of users as returned by the Chartio API
        Users struct {
                Users []User `json:"users,omitempty"`
                Count *int   `json:"count,omitempty"`
                Links Links  `json:"_links,omitempty"`
        }

        // UserInput is used as input on all user object functions.
        UserInput struct {
                UserID    string `json:"id,omitempty"`
                UserEmail string `json:"email,omitempty"`
                Team      *struct {
                        ID string `json:"id"`
                } `json:"team,omitempty"`
        }
)

// User coalesces the optional ID and Email fields.
func (i *UserInput) User() string {
        if i.UserID != "" {
                return i.UserID
        } else {
                return i.UserEmail
        }
}

// ListUsers returns the full list of users for the current organization.
func (c *Client) ListUsers() (*Users, error) {
        users := new(Users)
        _, err := c.Get("/users", users)
        return users, err
}

// AddUser creates a new user with the given name
func (c *Client) AddUser(i UserInput) (*User, error) {
        user := new(User)
        _, err := c.Post("/users", i, user)
        return user, err
}

// GetUser returns a particular user looked up by ID
func (c *Client) GetUser(i UserInput) (*User, error) {
        user := new(User)
        path := fmt.Sprintf("/users/%s", i.User())
        _, err := c.Get(path, user)
        return user, err
}

// DeleteUser deletes a particular user looked up by ID
func (c *Client) DeleteUser(i UserInput) error {
        path := fmt.Sprintf("/users/%s", i.User())
        _, err := c.Delete(path)
        return err
}
