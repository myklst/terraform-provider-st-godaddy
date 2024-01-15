//go:build integration
// +build integration

package api

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var key, secret, baseURL string

func init() {
	key = os.Getenv("GODADDY_API_KEY")
	secret = os.Getenv("GODADDY_API_SECRET")
	baseURL = "https://api.ote-godaddy.com"  // Prod: "https://api.godaddy.com"
}

func TestInvalidUrl(t *testing.T) {
	_, err := NewClient("api.godaddy.com", key, secret)
	assert.NotNil(t, err)
}

func TestAuthFailure(t *testing.T) {
	client, err := NewClient(baseURL, "ABC", "123")
	assert.Nil(t, err)
	assert.NotNil(t, client)

	_, err = client.GetDomainRecords("", "bogus.com")
	assert.NotNil(t, err)
}

func TestAvailable(t *testing.T) {
	client, err := NewClient(baseURL, key, secret)
	assert.Nil(t, err)
	assert.NotNil(t, client)
	var domains []string
	domains = append(domains, "hohojiang.com")
	client.domainAvailable(domains)
}

func TestAgreements(t *testing.T) {
	client, err := NewClient(baseURL, key, secret)
	assert.Nil(t, err)
	assert.NotNil(t, client)
	client.getAgreement("pl", false)
}

func TestGetRecords(t *testing.T) {
	client, err := NewClient(baseURL, key, secret)
	assert.Nil(t, err)
	assert.NotNil(t, client)

	getRecords(t, client, "n3integration.com")
}

func TestGetTooManyRecords(t *testing.T) {
	client, err := NewClient(baseURL, key, secret)
	assert.Nil(t, err)
	assert.NotNil(t, client)

	for i := 0; i < 75; i++ {
		_, err := getRecords(t, client, "n3integration.com")
		if err != nil {
			fmt.Println("Requests failing at", i+1)
		}
	}
}

func getRecords(t *testing.T, client *Client, domain string) ([]*DomainRecord, error) {
	records, err := client.GetDomainRecords("", domain)
	assert.Nil(t, err)
	assert.NotNil(t, records)

	for _, rec := range records {
		assert.NotNil(t, rec.Name)
		assert.NotNil(t, rec.Type)
		assert.NotNil(t, rec.Data)

		fmt.Println("REC -->", rec)
	}

	return records, err
}
