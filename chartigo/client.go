package chartigo

// Links represents the standard link map returned by the Chartio API
type Links struct {
        self     string
        next     string
        previous string
}

// APIKeyEnvVar is the name of the environment variable where the Chartio API
// key should be read from.
const APIKeyEnvVar = "CHARTIO_API_KEY"

// APIPasswordEnvVar is the name of the environment variable where the Chartio API
// password should be read from.
const APIPasswordEnvVar = "CHARTIO_API_PASSWORD"

// DefaultEndpoint is the default endpoint for Chartio.
const DefaultEndpoint = "https://api.chartio.com"
