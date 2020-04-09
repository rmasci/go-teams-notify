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
	"log"
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
	if valid, err := IsValidInput(webhookMessage, webhookURL); !valid {
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

	// Make sure that we close the response body once we're done with it
	defer func() {
		if err := res.Body.Close(); err != nil {
			log.Printf("error closing response body: %v", err)
		}
	}()

	// Get the response body, then convert to string for use with extended
	// error messages
	responseData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode >= 299 {
		// 400 Bad Response is likely an indicator that we failed to provide a
		// required field in our JSON payload. For example, when leaving out
		// the top level MessageCard Summary or Text field, the remote API
		// returns "Summary or Text is required." as a text string. We include
		// that response text in the error message that we return to the
		// caller.

		err = fmt.Errorf("error on notification: %v, %q", res.Status, string(responseData))
		return err
	}

	return nil
}

// helper --------------------------------------------------------------------------------------------------------------

// IsValidInput is a validation "wrapper" function. This function is intended
// to run current validation checks and offer easy extensibility for future
// validation requirements.
func IsValidInput(webhookMessage MessageCard, webhookURL string) (bool, error) {
	// validate url
	if valid, err := IsValidWebhookURL(webhookURL); !valid {
		return false, err
	}

	// validate message
	if valid, err := IsValidMessageCard(webhookMessage); !valid {
		return false, err
	}

	return true, nil
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

// IsValidMessageCard performs validation/checks for known issues with
// MessardCard values.
func IsValidMessageCard(webhookMessage MessageCard) (bool, error) {
	if (webhookMessage.Text == "") && (webhookMessage.Summary == "") {
		// This scenario results in:
		// 400 Bad Request
		// Summary or Text is required.
		return false, fmt.Errorf("invalid message card: summary or text field is required")
	}

	return true, nil
}
