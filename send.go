package goteamsnotify

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
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

	fmt.Printf("%+v", string(webhookMessageByte))

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

// MessageCardSectionFact represents a section fact entry usually displayed in
// a two-column key/value format.
type MessageCardSectionFact struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// MessageCardPotentialAction represents an action that a user may take for a
// received Microsoft Teams message.
type MessageCardPotentialAction struct {
	Target          []string    `json:"target"`
	Context         string      `json:"@context"`
	Type            string      `json:"@type"`
	ID              interface{} `json:"@id"`
	Name            string      `json:"name"`
	IsPrimaryAction bool        `json:"isPrimaryAction"`
}

// MessageCardSection represents a section to include in a message card.
type MessageCardSection struct {
	Title    string                   `json:"title"`
	Text     string                   `json:"text"`
	Markdown bool                     `json:"markdown"`
	Facts    []MessageCardSectionFact `json:"facts,omitempty"`
}

// https://docs.microsoft.com/en-us/outlook/actionable-messages/send-via-connectors
// https://docs.microsoft.com/en-us/outlook/actionable-messages/message-card-reference
// https://docs.microsoft.com/en-us/microsoftteams/platform/webhooks-and-connectors/how-to/connectors-using
// https://mholt.github.io/json-to-go/
// https://messagecardplayground.azurewebsites.net/
// https://connectplayground.azurewebsites.net/
// https://github.com/atc0005/bounce/issues/21

// MessageCard represents a legacy actionable message card used via Office 365
// or Microsoft Teams connectors.
type MessageCard struct {
	// Required; must be set to "MessageCard"
	Type string `json:"@type"`

	// Required; must be set to "https://schema.org/extensions"
	Context string `json:"@context"`

	// Summary is required if the card does not contain a text property,
	// otherwise optional. The summary property is typically displayed in the
	// list view in Outlook, as a way to quickly determine what the card is
	// all about. Summary appears to only be used when there are sections defined
	Summary string `json:"summary,omitempty"`

	// Title is the title property of a card. is meant to be rendered in a
	// prominent way, at the very top of the card. Use it to introduce the
	// content of the card in such a way users will immediately know what to
	// expect.
	Title string `json:"title"`

	// Text is required if the card does not contain a summary property,
	// otherwise optional. The text property is meant to be displayed in a
	// normal font below the card's title. Use it to display content, such as
	// the description of the entity being referenced, or an abstract of a
	// news article.
	Text string `json:"text"`

	// Specifies a custom brand color for the card. The color will be
	// displayed in a non-obtrusive manner.
	ThemeColor string `json:"themeColor,omitempty"`

	// Sections is a collection of sections to include in the card.
	Sections []MessageCardSection `json:"sections,omitempty"`

	// PotentialAction is a collection of actions that can be invoked on this card.
	PotentialAction []MessageCardPotentialAction `json:"potentialAction,omitempty"`
}

// NewMessageCard - create new empty message card
func NewMessageCard() MessageCard {

	// define expected values to meet Office 365 Connector card requirements
	// https://docs.microsoft.com/en-us/outlook/actionable-messages/message-card-reference#card-fields
	// TODO: Move string values to constants list
	msgCard := MessageCard{
		Type:    "MessageCard",
		Context: "https://schema.org/extensions",
	}

	return msgCard
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
