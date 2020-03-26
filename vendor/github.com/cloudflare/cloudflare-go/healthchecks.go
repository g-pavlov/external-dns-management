package cloudflare

import (
	"encoding/json"
	"time"

	"github.com/pkg/errors"
)

// Healthcheck describes a Healthcheck object.
type Healthcheck struct {
	ID            string                  `json:"id,omitempty"`
	CreatedOn     *time.Time              `json:"created_on,omitempty"`
	ModifiedOn    *time.Time              `json:"modified_on,omitempty"`
	Name          string                  `json:"name"`
	Description   string                  `json:"description"`
	Suspended     bool                    `json:"suspended"`
	Address       string                  `json:"address"`
	Retries       int                     `json:"retries,omitempty"`
	Timeout       int                     `json:"timeout,omitempty"`
	Interval      int                     `json:"interval,omitempty"`
	Type          string                  `json:"type,omitempty"`
	CheckRegions  []string                `json:"check_regions"`
	HTTPConfig    *HealthcheckHTTPConfig  `json:"http_config,omitempty"`
	Notification  HealthcheckNotification `json:"notification,omitempty"`
	Status        string                  `json:"status"`
	FailureReason string                  `json:"failure_reason"`
}

// HealthcheckHTTPConfig describes configuration for a HTTP healthcheck.
type HealthcheckHTTPConfig struct {
	Method          string              `json:"method"`
	Port            uint16              `json:"port,omitempty"`
	Path            string              `json:"path"`
	ExpectedCodes   []string            `json:"expected_codes"`
	ExpectedBody    string              `json:"expected_body"`
	FollowRedirects bool                `json:"follow_redirects"`
	AllowInsecure   bool                `json:"allow_insecure"`
	Header          map[string][]string `json:"header"`
}

// HealthcheckNotification describes notification configuration for a healthcheck.
type HealthcheckNotification struct {
	Suspended      bool     `json:"suspended,omitempty"`
	EmailAddresses []string `json:"email_addresses,omitempty"`
}

// HealthcheckListResponse is the API response, containing an array of healthchecks.
type HealthcheckListResponse struct {
	Response
	Result     []Healthcheck `json:"result"`
	ResultInfo `json:"result_info"`
}

// HealthcheckResponse is the API response, containting a single healthcheck.
type HealthcheckResponse struct {
	Response
	Result Healthcheck `json:"result"`
}

// Healthchecks returns all healthchecks for a zone.
//
// API reference: https://api.cloudflare.com/#health-checks-list-health-checks
func (api *API) Healthchecks(zoneID string) ([]Healthcheck, error) {
	uri := "/zones/" + zoneID + "/healthchecks"
	res, err := api.makeRequest("GET", uri, nil)
	if err != nil {
		return []Healthcheck{}, errors.Wrap(err, errMakeRequestError)
	}
	var r HealthcheckListResponse
	err = json.Unmarshal(res, &r)
	if err != nil {
		return []Healthcheck{}, errors.Wrap(err, errUnmarshalError)
	}
	return r.Result, nil
}

// Healthcheck returns a single healthcheck by ID.
//
// API reference: https://api.cloudflare.com/#health-checks-health-check-details
func (api *API) Healthcheck(zoneID, healthcheckID string) (Healthcheck, error) {
	uri := "/zones/" + zoneID + "/healthchecks/" + healthcheckID
	res, err := api.makeRequest("GET", uri, nil)
	if err != nil {
		return Healthcheck{}, errors.Wrap(err, errMakeRequestError)
	}
	var r HealthcheckResponse
	err = json.Unmarshal(res, &r)
	if err != nil {
		return Healthcheck{}, errors.Wrap(err, errUnmarshalError)
	}
	return r.Result, nil
}

// CreateHealthcheck creates a new healthcheck in a zone.
//
// API reference: https://api.cloudflare.com/#health-checks-create-health-check
func (api *API) CreateHealthcheck(zoneID string, healthcheck Healthcheck) (Healthcheck, error) {
	uri := "/zones/" + zoneID + "/healthchecks"
	res, err := api.makeRequest("POST", uri, healthcheck)
	if err != nil {
		return Healthcheck{}, errors.Wrap(err, errMakeRequestError)
	}
	var r HealthcheckResponse
	err = json.Unmarshal(res, &r)
	if err != nil {
		return Healthcheck{}, errors.Wrap(err, errUnmarshalError)
	}
	return r.Result, nil
}

// UpdateHealthcheck updates an existing healthcheck.
//
// API reference: https://api.cloudflare.com/#health-checks-update-health-check
func (api *API) UpdateHealthcheck(zoneID string, healthcheckID string, healthcheck Healthcheck) (Healthcheck, error) {
	uri := "/zones/" + zoneID + "/healthchecks/" + healthcheckID
	res, err := api.makeRequest("PUT", uri, healthcheck)
	if err != nil {
		return Healthcheck{}, errors.Wrap(err, errMakeRequestError)
	}
	var r HealthcheckResponse
	err = json.Unmarshal(res, &r)
	if err != nil {
		return Healthcheck{}, errors.Wrap(err, errUnmarshalError)
	}
	return r.Result, nil
}

// DeleteHealthcheck deletes a healthcheck in a zone.
//
// API reference: https://api.cloudflare.com/#health-checks-delete-health-check
func (api *API) DeleteHealthcheck(zoneID string, healthcheckID string) error {
	uri := "/zones/" + zoneID + "/healthchecks/" + healthcheckID
	res, err := api.makeRequest("DELETE", uri, nil)
	if err != nil {
		return errors.Wrap(err, errMakeRequestError)
	}
	var r HealthcheckResponse
	err = json.Unmarshal(res, &r)
	if err != nil {
		return errors.Wrap(err, errUnmarshalError)
	}
	return nil
}
