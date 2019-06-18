package chartigo

import (
        "bytes"
        "encoding/json"
        "fmt"
        "log"
        "net/http"
        "net/url"
        "os"
        "runtime"
        "strings"
        "time"
)

const (
        // APIKeyEnvVar is the name of the environment variable where the
        // Chartio API key should be read from.
        APIKeyEnvVar = "CHARTIO_API_KEY"

        // APIPasswordEnvVar is the name of the environment variable where
        // the Chartio API password should be read from.
        APIPasswordEnvVar = "CHARTIO_API_PASSWORD"

        // ChartioOrgEnvVar is the name of the environment variable where
        // the desired chartio oganization name should be read from.
        ChartioOrgEnvVar = "CHARTIO_ORG"

        // DefaultURL is the base url for the Chartio API.
        DefaultURL = "https://api.chartio.com"

        // DefaultAPIVersion is the default Chartio API version supported
        // by this library.
        DefaultAPIVersion = "v1"

        // TimeParseFormat correctly parses timestamps returned by the
        // Chartio API.
        TimeParseFormat = "2006-01-02T15:04:05.000000"
)

var (
        // ProjectURL is the url for this library.
        ProjectURL = "github.com/aaronbiller/chartigo"

        // ProjectVersion is the version of this library.
        ProjectVersion = "0.0.1"

        // UserAgent is the user agent for this particular client.
        UserAgent = fmt.Sprintf("ChartiGo/%s (+%s; %s)",
                ProjectVersion, ProjectURL, runtime.Version())
)

type (
        // Client is the main entrypoint to the Chartio golang API library.
        Client struct {
                BaseURL     *url.URL
                HTTPClient  *http.Client
                UserAgent   string
                APIKey      string
                APIPassword string
                ChartioOrg  string
        }

        // Links represents the standard link map returned by the Chartio API.
        Links map[string]map[string]string

        // ChartioTS represents the timestamp format returned by
        // the Chartio API.
        ChartioTS struct{ time.Time }

        // Option is a self-referential function applied to a new Client.
        Option func(*Client)
)

func SetHTTPClient(httpClient *http.Client) Option {
        return func(c *Client) {
                c.HTTPClient = httpClient
        }
}

// NewClient creates a client object for querying the Chartio API
func NewClient(org string, options ...Option) *Client {
        if org == "" {
                org = os.Getenv(ChartioOrgEnvVar)
        }

        apiKey, apiPass, err := getAPIEnvVariables()
        if err != nil {
                log.Fatal(err)
        }

        URL := fmt.Sprintf("%s/%s/%s", DefaultURL, DefaultAPIVersion, org)
        baseURL, _ := url.Parse(URL)

        client := &Client{
                BaseURL:     baseURL,
                HTTPClient:  http.DefaultClient,
                UserAgent:   UserAgent,
                APIKey:      apiKey,
                APIPassword: apiPass,
                ChartioOrg:  org,
        }

        for _, option := range options {
                option(client)
        }

        return client
}

// Get issues a GET request on the given endpoint.
func (c *Client) Get(p string, t interface{}) (*http.Response, error) {
        return c.request("GET", p, nil, t)
}

// Head issues a HEAD request on the given endpoint.
func (c *Client) Head(p string, t interface{}) (*http.Response, error) {
        return c.request("HEAD", p, nil, t)
}

// Patch issues a PATCH request on the given endpoint.
func (c *Client) Patch(p string, i, t interface{}) (*http.Response, error) {
        return c.request("PATCH", p, i, t)
}

// Post issues a POST request on the given endpoint.
func (c *Client) Post(p string, i, t interface{}) (*http.Response, error) {
        return c.request("POST", p, i, t)
}

// Put issues a PUT request on the given endpoint.
func (c *Client) Put(p string, i, t interface{}) (*http.Response, error) {
        return c.request("PUT", p, i, t)
}

// Delete issues a DELETE request on the given endpoint.
func (c *Client) Delete(p string) (*http.Response, error) {
        return c.request("DELETE", p, nil, nil)
}

func (c *Client) request(m, p string, d, t interface{}) (*http.Response, error) {
        u := c.buildURL(p)
        body, err := json.Marshal(d)
        if err != nil {
                return nil, err
        }

        req, err := http.NewRequest(m, u, bytes.NewReader(body))
        if err != nil {
                return nil, err
        }
        req.Header.Set("Content-Type", "application/json")
        req.Header.Set("Accept", "application/json")
        req.Header.Set("User-Agent", c.UserAgent)
        req.SetBasicAuth(c.APIKey, c.APIPassword)

        resp, err := wrapHTTPResp(c.HTTPClient.Do(req))
        if err != nil {
                return resp, err
        }

        if t != nil {
                defer resp.Body.Close()
                if err := json.NewDecoder(resp.Body).Decode(t); err != nil {
                        return resp, err
                }
        }
        return resp, err
}

func (c *Client) buildURL(p string) string {
        return c.BaseURL.String() + p
}

func getAPIEnvVariables() (k, p string, err error) {
        k, ok := os.LookupEnv(APIKeyEnvVar)
        if !ok {
                return k, p, ErrMissingAPIKey
        }
        p, ok = os.LookupEnv(APIPasswordEnvVar)
        if !ok {
                return k, p, ErrMissingAPIPassword
        }
        return k, p, nil
}

func (cd *ChartioTS) UnmarshalJSON(input []byte) error {
        strInput := string(input)
        strInput = strings.Trim(strInput, `"`)
        newTime, err := time.Parse(TimeParseFormat, strInput)
        if err != nil {
                return err
        }

        cd.Time = newTime
        return nil
}
