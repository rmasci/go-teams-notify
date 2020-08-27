// Copyright 2020 Enrico Hoffmann
// Copyright 2020 Adam Chalkley
//
// https:#github.com/atc0005/go-teams-notify
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package goteamsnotify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Known webhook URL prefixes for submitting messages to Microsoft Teams
const (
	WebhookURLOfficecomPrefix = "https://outlook.office.com"
	WebhookURLOffice365Prefix = "https://outlook.office365.com"
)

// API - interface of MS Teams notify
type API interface {
	Send(webhookURL string, webhookMessage MessageCard) error
}

type teamsClient struct {
	httpClient *http.Client
}

// NewClient - create a brand new client for MS Teams notify
func NewClient() API {
	client := teamsClient{
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
	return &client
}

// Send - will post a notification to MS Teams webhook URL
func (c teamsClient) Send(webhookURL string, webhookMessage MessageCard) error {
	// Validate input data
	if valid, err := IsValidWebhookURL(webhookURL); !valid {
		return err
	}

	// prepare message
	webhookMessageByte, _ := json.Marshal(webhookMessage)
	webhookMessageBuffer := bytes.NewBuffer(webhookMessageByte)

	// prepare request (error not possible)
	req, _ := http.NewRequest(http.MethodPost, webhookURL, webhookMessageBuffer)
	req.Header.Add("Content-Type", "application/json;charset=utf-8")

	// do the request
	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	// Get the response body, then convert to string for use with extended
	// error messages
	responseData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode >= 299 {
		err = fmt.Errorf("error on notification: %v, %q", res.Status, string(responseData))
		return err
	}

	return nil
}

// IsValidWebhookURL performs validation checks on the webhook URL used to
// submit messages to Microsoft Teams.
func IsValidWebhookURL(webhookURL string) (bool, error) {
	switch {
	case strings.HasPrefix(webhookURL, WebhookURLOfficecomPrefix):
	case strings.HasPrefix(webhookURL, WebhookURLOffice365Prefix):
	default:
		u, err := url.Parse(webhookURL)
		if err != nil {
			return false, fmt.Errorf(
				"unable to parse webhook URL %q: %w",
				webhookURL,
				err,
			)
		}
		userProvidedWebhookURLPrefix := u.Scheme + "://" + u.Host

		return false, fmt.Errorf(
			"webhook URL does not contain expected prefix; got %q, expected one of %q or %q",
			userProvidedWebhookURLPrefix,
			WebhookURLOfficecomPrefix,
			WebhookURLOffice365Prefix,
		)
	}

	return true, nil
}
