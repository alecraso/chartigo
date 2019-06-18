package chartigo

import (
        "fmt"
)

type (
        // Datasource represents a datasource object in Chartio
        Datasource struct {
                ID              int        `json:"id"`
                Alias           string     `json:"alias,omitempty"`
                CreatedAt       *ChartioTS `json:"created_at,omitempty"`
                UpdatedAt       *ChartioTS `json:"updated_at,omitempty"`
                LastReflected   *ChartioTS `json:"last_reflected,omitempty"`
                LastRefreshedAt *ChartioTS `json:"last_refreshed_at,omitempty"`
                Links           Links      `json:"_links,omitempty"`
        }

        // Datasources represents an array of datasources as returned
        // by the Chartio API
        Datasources struct {
                Datasources []Datasource `json:"datasources,omitempty"`
                Count       int          `json:"count,omitempty"`
                Links       Links        `json:"_links,omitempty"`
        }

        // DatasourceInput is used as input on all datasource
        // object functions.
        DatasourceInput struct {
                DatasourceID string `json:"-"`
                Alias        string `json:"alias,omitempty"`
        }
)

// ListDatasources returns the full list of datasources for the
// current organization.
func (c *Client) ListDatasources() (*Datasources, error) {
        ds := new(Datasources)
        _, err := c.Get("/datasources", ds)
        return ds, err
}

// GetDatasource returns a particular datasource looked up by ID
func (c *Client) GetDatasource(i DatasourceInput) (*Datasource, error) {
        ds := new(Datasource)
        path := fmt.Sprintf("/datasources/%s", i.DatasourceID)
        _, err := c.Get(path, ds)
        return ds, err
}

// UpdateDatasource deletes a particular team looked up by ID
func (c *Client) UpdateDatasource(i DatasourceInput) (*Datasource, error) {
        ds := new(Datasource)
        path := fmt.Sprintf("/datasources/%s", i.DatasourceID)
        _, err := c.Patch(path, i, ds)
        return ds, err
}
