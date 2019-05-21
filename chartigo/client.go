package chartigo

import (
        "bytes"
        "encoding/json"
        "fmt"
        "io"
        "log"
        "net/http"
        "net/url"
        "os"
        "reflect"
        "runtime"
        "time"

        "github.com/mitchellh/mapstructure"
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
        TimeParseFormat = "2006-01-02T15:04:05.999999"
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
                href string
        }

        // Links represents the standard link map returned by the Chartio API
        Links struct {
                self     *Link
                next     *Link
                previous *Link
        }
)

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

// Get issues a GET request on the given endpoint
func (c *Client) Get(p string, v interface{}) (*http.Response, error) {
        return c.request("GET", p, nil, v)
}

// Head issues a HEAD request on the given endpoint
func (c *Client) Head(p string, v interface{}) (*http.Response, error) {
        return c.request("HEAD", p, nil, v)
}

// Patch issues a PATCH request on the given endpoint
func (c *Client) Patch(p string, i, v interface{}) (*http.Response, error) {
        return c.request("PATCH", p, i, v)
}

// Post issues a POST request on the given endpoint
func (c *Client) Post(p string, i, v interface{}) (*http.Response, error) {
        return c.request("POST", p, i, v)
}

// Put issues a PUT request on the given endpoint
func (c *Client) Put(p string, i, v interface{}) (*http.Response, error) {
        return c.request("PUT", p, i, v)
}

// Delete issues a DELETE request on the given endpoint
func (c *Client) Delete(p string) (*http.Response, error) {
        return c.request("DELETE", p, nil, nil)
}

func (c *Client) request(m, p string, data interface{}, v interface{}) (*http.Response, error) {
        u := c.buildURL(p)

        var body io.Reader
        if data != nil {
                b, err := json.Marshal(data)
                if err != nil {
                        return nil, err
                }
                body = bytes.NewReader(b)
        }

        req, err := http.NewRequest(m, u, body)
        if err != nil {
                return nil, err
        }

        req.Header.Set("Content-Type", "application/json")
        req.Header.Set("Accept", "application/json")
        req.Header.Set("User-Agent", c.UserAgent)
        req.SetBasicAuth(c.APIKey, c.APIPassword)

        resp, err := c.HTTPClient.Do(req)
        if err != nil {
                return nil, err
        }

        if v != nil {
                if err := decodeJSON(resp.Body, v); err != nil {
                        return resp, err
                }
        }
        return resp, err
}

func (c *Client) buildURL(p string) string {
        return c.BaseURL.String() + p
}

func getAPIEnvVariables() (k, p string, err error) {
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

// decodeJSON is used to decode an HTTP response body into an interface as JSON.
func decodeJSON(body io.ReadCloser, out interface{}) error {
        defer body.Close()

        var parsed interface{}
        dec := json.NewDecoder(body)
        if err := dec.Decode(&parsed); err != nil {
                return err
        }

        decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
                DecodeHook: mapstructure.ComposeDecodeHookFunc(
                        stringToTimeHookFunc(),
                        mapToLinksHookFunc(),
                ),
                WeaklyTypedInput: true,
                Result:           out,
        })
        if err != nil {
                return err
        }
        return decoder.Decode(parsed)
}

func mapToLinksHookFunc() mapstructure.DecodeHookFunc {
        return func(
                f reflect.Type,
                t reflect.Type,
                data interface{}) (interface{}, error) {
                if f.Kind() != reflect.Map {
                        return data, nil
                }
                if t.Name() != "Links" {
                        return data, nil
                }
                var out *Links
                err := mapstructure.Decode(data, &out)
                return out, err
        }
}

// stringToTimeHookFunc returns a function that converts strings to a time.Time value.
func stringToTimeHookFunc() mapstructure.DecodeHookFunc {
        return func(
                f reflect.Type,
                t reflect.Type,
                data interface{}) (interface{}, error) {
                if f.Kind() != reflect.String {
                        return data, nil
                }
                if t != reflect.TypeOf(time.Now()) {
                        return data, nil
                }

                // Convert it by parsing
                return time.Parse(TimeParseFormat, data.(string))
        }
}
