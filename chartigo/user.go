package chartigo

import (
        "fmt"
        "time"
)

type (
        // User represents a user object in Chartio
        User struct {
                ID          int       `mapstructure:id`
                DisplayName string    `mapstructure:display_name,omitempty`
                Email       string    `mapstructure:email`
                CreatedAt   time.Time `mapstructure:created_at`
                UpdatedAt   time.Time `mapstructure:updated_at`
                Links       Links     `mapstructure:_links`
                Embedded    Teams     `mapstructure:_embedded`
        }

        // Users represents an array of users as returned by the Chartio API
        Users struct {
                Users []User `mapstructure:"users,omitempty"`
                Count int    `mapstructure:"count,omitempty"`
                Links Links  `mapstructure:"_links,omitempty"`
        }
)

// ListUsers returns the full list of users for the current organization.
func (c *Client) ListUsers() (*Users, error) {
        users := new(Users)
        _, err := c.Get("/users", users)
        return users, err
}

// CreateUserInput is used as input to the CreateUser function.
type CreateUserInput struct {
        Email string `json:"email"`
        Team  struct {
                ID string `json:"id"`
        } `json:"team"`
}

// CreateUser creates a new user with the given name
func (c *Client) CreateUser(i CreateUserInput) (*User, error) {
        user := new(User)
        _, err := c.Post("/users", i, user)
        return user, err
}

// GetUserInput is used as input to the GetUser function.
type GetUserInput struct {
        ID    string `json:"id"`
        Email string `json:"email"`
}

// GetUser returns a particular user looked up by ID
func (c *Client) GetUser(i GetUserInput) (*User, error) {
        user := new(User)
        if i.ID == "" {
                i.ID = i.Email
        }
        path := fmt.Sprintf("/users/%s", i.ID)
        _, err := c.Get(path, user)
        return user, err
}

// DeleteUserInput is used as input to the DeleteUser function.
type DeleteUserInput struct {
        ID    string `json:"id"`
        Email string `json:"email"`
}

// DeleteUser deletes a particular team looked up by ID
func (c *Client) DeleteUser(i DeleteUserInput) error {
        if i.ID == "" {
                i.ID = i.Email
        }
        path := fmt.Sprintf("/users/%s", i.ID)
        _, err := c.Delete(path)
        return err
}
