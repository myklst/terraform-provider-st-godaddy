package godaddy_provider

import (
	"fmt"
	"log"
	"terraform-provider-st-godaddy/godaddy/api"
)

// Config provides the provider's configuration
type Config struct {
	Key     string
	Secret  string
	BaseURL string
}

// Client returns a new client for accessing GoDaddy.
func (c *Config) Client() (*api.Client, error) {
	client, err := api.NewClient(c.BaseURL, c.Key, c.Secret)

	if err != nil {
		return nil, fmt.Errorf("error setting up client: %s", err)
	}

	log.Print("GoDaddy Client configured")

	return client, nil
}
