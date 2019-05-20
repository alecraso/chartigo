package chartigo

import (
        "encoding/json"
        "fmt"
        "io/ioutil"
        "log"
        "net/http"
        "net/url"
        "os"
        "runtime"
        "strings"
        "time"
)

const (
        // APIKeyEnvVar is the name of the environment variable where the Chartio API
        // key should be read from.
        APIKeyEnvVar = "CHARTIO_API_KEY"

        // APIPasswordEnvVar is the name of the environment variable where the Chartio API
        // password should be read from.
        APIPasswordEnvVar = "CHARTIO_API_PASSWORD"

        // ChartioOrgEnvVar is the name of the environment variable where the desired
        // chartio oganization name should be read from.
        ChartioOrgEnvVar = "CHARTIO_ORG"

        // DefaultURL is the base url for the Chartio API.
        DefaultURL = "https://api.chartio.com"

        // DefaultAPIVersion is the default Chartio API version supported by this library.
        DefaultAPIVersion = "v1"

        // TimeParseFormat correctly parses timestamps returned by the Chartio API
        TimeParseFormat
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

        // NewClientInput contains the input parameters to create a new client
        NewClientInput struct {
                Version string
                Org     string
        }

        // Link represents a url value returned from the Chartio API
        Link struct {
                href string `json:"href"`
        }

        // Links represents the standard link map returned by the Chartio API
        Links struct {
                self     Link `json:"self,omitempty"`
                next     Link `json:"next,omitempty"`
                previous Link `json:"previous,omitempty"`
        }

        // Unmarshaler interface for types that can unmarshal JSON
        Unmarshaler interface {
                UnmarshalJSON([]byte) error
        }

        // JSONTime type for parsing timestamps returned by the Chartio API
        JSONTime struct {
                time.Time
        }
)

// UnmarshalJSON properly parses timestamps from the ChartioAPI
func (j *JSONTime) UnmarshalJSON(b []byte) (err error) {
        s := strings.Trim(string(b), "\"")
        j.Time, err = time.Parse("2006-01-02T15:04:05.999999", s)
        return
}

// NewClient creates a client object for querying the Chartio API
func NewClient(i NewClientInput) *Client {
        if i.Version == "" {
                i.Version = DefaultAPIVersion
        }
        if i.Org == "" {
                i.Org = os.Getenv(ChartioOrgEnvVar)
        }

        apiKey, apiPass, err := getAPIEnvVariables()
        if err != nil {
                log.Fatal(err)
        }

        URL := fmt.Sprintf("%s/%s/%s", DefaultURL, i.Version, i.Org)
        baseURL, _ := url.Parse(URL)

        client := &Client{
                BaseURL:     baseURL,
                HTTPClient:  http.DefaultClient,
                UserAgent:   UserAgent,
                APIKey:      apiKey,
                APIPassword: apiPass,
                ChartioOrg:  i.Org,
        }
        return client
}

func (c *Client) request(m, p string, v interface{}) (*http.Response, error) {
        u := c.buildURL(p)

        req, reqErr := http.NewRequest(m, u, nil)
        if reqErr != nil {
                return nil, reqErr
        }

        req.Header.Set("Content-Type", "application/json")
        req.Header.Set("Accept", "application/json")
        req.Header.Set("User-Agent", c.UserAgent)
        req.SetBasicAuth(c.APIKey, c.APIPassword)

        resp, respErr := c.HTTPClient.Do(req)
        if respErr != nil {
                return nil, respErr
        }
        defer resp.Body.Close()

        jsn, jsnErr := ioutil.ReadAll(resp.Body)
        if jsnErr != nil {
                log.Fatal("Error reading the body", jsnErr)
        }
        jsnErr = json.Unmarshal(jsn, &v)

        return resp, jsnErr
}

func (c *Client) buildURL(p string) string {
        return c.BaseURL.String() + p
}

func getAPIEnvVariables() (k string, p string, err error) {
        errMsg := "Error: environment variable %s not found, required for authentication."
        k, kOK := os.LookupEnv(APIKeyEnvVar)
        if !kOK {
                return "", "", fmt.Errorf(errMsg, APIKeyEnvVar)
        }
        p, pOK := os.LookupEnv(APIPasswordEnvVar)
        if !pOK {
                return "", "", fmt.Errorf(errMsg, APIPasswordEnvVar)
        }
        return k, p, nil
}
