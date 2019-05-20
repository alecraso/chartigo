package chartigo

import (
        "fmt"
)

type (
        // User represents a user object in Chartio
        User struct {
                ID          int      `json:id,omitempty`
                Email       string   `json:email,omitempty`
                DisplayName string   `json:display_name,omitempty`
                CreatedAt   JSONTime `json:created_at,omitempty`
                UpdatedAt   JSONTime `json:updated_at,omitempty`
                Links       Links    `json:_links,omitempty`
                Embedded    Teams    `json:_embedded,omitempty`
        }

        // Users represents an array of users as returned by the Chartio API
        Users struct {
                Users []User `json:"users,omitempty"`
                Count int    `json:"count,omitempty"`
                Links Links  `json:"_links,omitempty"`
        }
)

// ListUsers returns the full list of users for the current organization.
func (c *Client) ListUsers() (*Users, error) {
        users := new(Users)
        _, err := c.request("GET", "/users", users)
        return users, err
}

// CreateUserInput is used as input to the CreateUser function.
type CreateUserInput struct {
        Name string
}

// CreateUser creates a new user with the given name
// func (c *Client) CreateUser(i CreateUserInput) (*User, error) {
//         return
// }

// GetUserInput is used as input to the GetUser function.
type GetUserInput struct {
        ID    string
        Email string
}

// GetUser returns a particular user looked up by ID
func (c *Client) GetUser(i GetUserInput) (*User, error) {
        user := new(User)
        if i.ID == "" {
                i.ID = i.Email
        }
        path := fmt.Sprintf("/users/%s", i.ID)
        _, err := c.request("GET", path, user)
        return user, err
}
