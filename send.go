package goteamsnotify

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// API - interface of MS Teams notify
type API interface {
	Send(webhookURL string, webhookMessage MessageCard) error
}

type teamsClient struct {
	httpClient *http.Client
}

// NewClient - create a brand new client for MS Teams notify
func NewClient() (API, error) {
	client := teamsClient{
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
	return &client, nil
}

// Send - will post a notification to MS Teams incomingWebhookURL
func (c teamsClient) Send(webhookURL string, webhookMessage MessageCard) error {
	// validate url
	// needs to look like: https://outlook.office.com/webhook/xxx
	valid, err := isValidWebhookURL(webhookURL)
	if !valid {
		return err
	}
	// prepare message
	webhookMessageByte, _ := json.Marshal(webhookMessage)
	webhookMessageBuffer := bytes.NewBuffer(webhookMessageByte)

	// prepare request (error not possible)
	req, _ := http.NewRequest("POST", webhookURL, webhookMessageBuffer)
	req.Header.Add("Content-Type", "application/json;charset=utf-8")

	// do the request
	res, err := c.httpClient.Do(req)
	if err != nil {
		log.Println(err)
		return err
	}
	if res.StatusCode >= 299 {
		err = errors.New("error on notification: " + res.Status)
		log.Println(err)
		return err
	}

	return nil
}

// MessageCard - struct of message card
// https://docs.microsoft.com/en-us/outlook/actionable-messages/send-via-connectors
// https://docs.microsoft.com/en-us/outlook/actionable-messages/message-card-reference
// https://mholt.github.io/json-to-go/
// https://messagecardplayground.azurewebsites.net/
// https://github.com/atc0005/bounce/issues/21
type MessageCard struct {
	Summary    string `json:"summary,omitempty"`
	Title      string `json:"title"`
	Text       string `json:"text"`
	ThemeColor string `json:"themeColor,omitempty"`
	Sections   []struct {
		Title    string `json:"title"`
		Text     string `json:"text"`
		Markdown bool   `json:"markdown"`
		Facts    []struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		} `json:"facts,omitempty"`
	} `json:"sections,omitempty"`
	PotentialAction []struct {
		Target          []string    `json:"target"`
		Context         string      `json:"@context"`
		Type            string      `json:"@type"`
		ID              interface{} `json:"@id"`
		Name            string      `json:"name"`
		IsPrimaryAction bool        `json:"isPrimaryAction"`
	} `json:"potentialAction,omitempty"`
}

// NewMessageCard - create new empty message card
func NewMessageCard() MessageCard {
	return MessageCard{}
}

// helper --------------------------------------------------------------------------------------------------------------

func isValidWebhookURL(webhookURL string) (bool, error) {
	// basic URL check
	_, err := url.Parse(webhookURL)
	if err != nil {
		return false, err
	}
	// only pass MS teams webhook URLs
	switch {
	case strings.HasPrefix(webhookURL, "https://outlook.office.com/webhook/"):
	case strings.HasPrefix(webhookURL, "https://outlook.office365.com/webhook/"):
	default:
		err = errors.New("invalid ms teams webhook url")
		return false, err
	}
	return true, nil
}
