package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	defaultLimit = 500
	period       = 1
	autoRenew    = false

	pathDomainRecords       = "%s/v1/domains/%s/records?limit=%d&offset=%d"
	pathDomainRecordsByType = "%s/v1/domains/%s/records/%s"
	pathDomains             = "%s/v1/domains/%s"
	pathAvailable           = "%s/v1/domains/available"
	pathAgreements          = "%s/v1/agreements?tlds=%s&privacy=%t"
	pathDomainPurchase      = "%s/v1/domains/purchase"
	pathDomainRenew         = "%s/v1/domains/%s/renew"
)

// GetDomains fetches the details for the provided domain
func (c *Client) GetDomains(customerID string) ([]Domain, error) {
	domainURL := fmt.Sprintf(pathDomains, c.baseURL, "")
	req, err := http.NewRequest(http.MethodGet, domainURL, nil)

	if err != nil {
		return nil, err
	}

	var d []Domain
	if err := c.execute(customerID, req, &d); err != nil {
		return nil, err
	}

	return d, nil
}

// GetDomain fetches the details for the provided domain
func (c *Client) GetDomain(customerID, domain string) (*Domain, error) {
	domainURL := fmt.Sprintf(pathDomains, c.baseURL, domain)
	req, err := http.NewRequest(http.MethodGet, domainURL, nil)

	if err != nil {
		return nil, err
	}

	d := new(Domain)
	if err := c.execute(customerID, req, &d); err != nil {
		return nil, err
	}

	return d, nil
}

// GetDomainRecords fetches all existing records for the provided domain
func (c *Client) GetDomainRecords(customerID, domain string) ([]*DomainRecord, error) {
	offset := 1
	records := make([]*DomainRecord, 0)
	for {
		page := make([]*DomainRecord, 0)
		domainURL := fmt.Sprintf(pathDomainRecords, c.baseURL, domain, defaultLimit, offset)
		req, err := http.NewRequest(http.MethodGet, domainURL, nil)

		if err != nil {
			return nil, err
		}

		if err := c.execute(customerID, req, &page); err != nil {
			return nil, err
		}
		if len(page) == 0 {
			break
		}
		offset += 1
		records = append(records, page...)
	}

	return records, nil
}

// UpdateDomainRecords adds records or replaces all existing records for the provided domain
func (c *Client) UpdateDomainRecords(customerID, domain string, records []*DomainRecord) error {
	for t := range supportedTypes {
		typeRecords := c.domainRecordsOfType(t, records)
		if IsDisallowed(t, typeRecords) {
			continue
		}

		msg, err := json.Marshal(typeRecords)
		if err != nil {
			return err
		}

		domainURL := fmt.Sprintf(pathDomainRecordsByType, c.baseURL, domain, t)
		buffer := bytes.NewBuffer(msg)

		log.Println(domainURL)
		log.Println(buffer)

		req, err := http.NewRequest(http.MethodPut, domainURL, buffer)
		if err != nil {
			return err
		}

		if err := c.execute(customerID, req, nil); err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) domainRecordsOfType(t string, records []*DomainRecord) []*DomainRecord {
	typeRecords := make([]*DomainRecord, 0)

	for _, record := range records {
		if strings.EqualFold(record.Type, t) {
			typeRecords = append(typeRecords, record)
		}
	}

	return typeRecords
}

func (c *Client) DomainAvailable(domainNames []string) (bool, error) {

	msg, err := json.Marshal(domainNames)
	if err != nil {
		return false, err
	}

	domainURL := fmt.Sprintf(pathAvailable, c.baseURL)
	buffer := bytes.NewBuffer(msg)

	req, err := http.NewRequest(http.MethodPost, domainURL, buffer)
	if err != nil {
		return false, err
	}
	var resp AvailableResp
	if err := c.execute("", req, &resp); err != nil {
		return false, err
	}

	return resp.DomainAvailable[0].Available, nil
}

// Retrieve the legal agreement(s) required to purchase the specified TLD and add-ons
func (c *Client) GetAgreement(tld string, privacy bool) ([]*AgreementsResp, error) {

	domainURL := fmt.Sprintf(pathAgreements, c.baseURL, tld, privacy)
	//buffer := bytes.NewBuffer(msg)
	req, err := http.NewRequest(http.MethodGet, domainURL, nil)
	if err != nil {
		return nil, err
	}

	resp := make([]*AgreementsResp, 0)
	if err := c.execute("", req, &resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) Purchase(domainName string, keys []string, customer string, info RegisterDomainInfo) error {
	//url
	domainURL := fmt.Sprintf(pathDomainPurchase, c.baseURL)

	//request
	info.Domain = domainName
	info.Consent.AgreementKeys = keys
	info.Consent.AgreedAt = time.Now().String()
	info.Consent.AgreedBy = info.ContactAdmin.NameFirst + " " + info.ContactAdmin.NameLast
	info.Period = period
	info.RenewAuto = autoRenew
	info.Privacy = false
	msg, err := json.Marshal(info)
	if err != nil {
		return err
	}
	buffer := bytes.NewBuffer(msg)
	req, err := http.NewRequest(http.MethodPost, domainURL, buffer)
	if err != nil {
		return err
	}
	//response
	var resp DomainPurchaseResponse
	if err := c.execute(customer, req, &resp); err != nil {
		return err
	}

	return nil
}

func (c *Client) DomainRenew(customerID, domain string, period int) error {
	//url
	domainURL := fmt.Sprintf(pathDomainRenew, c.baseURL, domain)
	//request
	var info DomainRenew
	info.Period = period
	msg, err := json.Marshal(info)
	if err != nil {
		return err
	}
	buffer := bytes.NewBuffer(msg)
	req, err := http.NewRequest(http.MethodPost, domainURL, buffer)
	if err != nil {
		return err
	}
	//response
	var resp DomainPurchaseResponse
	//do request
	if err := c.execute(customerID, req, &resp); err != nil {
		return err
	}

	return nil
}

func (c *Client) DomainCancel(customerID, domain string) error {
	//url
	domainURL := fmt.Sprintf(pathDomains, c.baseURL, domain)
	//request
	req, err := http.NewRequest(http.MethodDelete, domainURL, nil)
	if err != nil {
		return err
	}
	//do request
	if err := c.execute(customerID, req, nil); err != nil {
		return err
	}

	return nil
}

// Retrieve the schema to be submitted when registering a Domain for the specified TLD
func (c *Client) Schema(domainNames []string) (bool, error) {
	return false, nil
}
