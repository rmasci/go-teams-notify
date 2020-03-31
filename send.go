package goteamsnotify

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// logger is a package logger that can be enabled from client code to allow
// logging output from this package when desired/needed for troubleshooting
var logger *log.Logger

// API - interface of MS Teams notify
type API interface {
	Send(webhookURL string, webhookMessage MessageCard) error
}

type teamsClient struct {
	httpClient *http.Client
}

func init() {

	// Disable logging output by default unless client code explicitly
	// requests it
	logger = log.New(os.Stderr, "[goteamsnotify] ", 0)
	logger.SetOutput(ioutil.Discard)

}

// EnableLogging enables logging output from this package. Output is muted by
// default unless explicitly requested (by calling this function).
func EnableLogging() {
	logger.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	logger.SetOutput(os.Stderr)
}

// DisableLogging reapplies default package-level logging settings of muting
// all logging output.
func DisableLogging() {
	logger.SetFlags(0)
	logger.SetOutput(ioutil.Discard)
}

// NewClient - create a brand new client for MS Teams notify
// FIXME: Upstream v2 removed error return code here
// TODO: Look at other changes made to v2 and merge
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
	if valid, err := IsValidWebhookURL(webhookURL); !valid {
		return err
	}

	// validate message
	if valid, err := IsValidMessageCard(webhookMessage); !valid {
		return err
	}

	// prepare message
	webhookMessageByte, _ := json.Marshal(webhookMessage)
	webhookMessageBuffer := bytes.NewBuffer(webhookMessageByte)

	// Basic, unformatted JSON
	//logger.Printf("DEBUG: %+v\n", string(webhookMessageByte))

	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, webhookMessageByte, "", "\t"); err != nil {
		return err
	}
	logger.Printf("DEBUG: %v\n", prettyJSON.String())

	// prepare request (error not possible)
	req, _ := http.NewRequest("POST", webhookURL, webhookMessageBuffer)
	req.Header.Add("Content-Type", "application/json;charset=utf-8")

	// do the request
	res, err := c.httpClient.Do(req)
	if err != nil {
		logger.Println(err)
		return err
	}

	// Make sure that we close the response body once we're done with it
	defer res.Body.Close()

	// Get the response body, then convert to string for use with extended
	// error messages
	responseData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger.Println(err)
		return err
	}
	responseString := string(responseData)

	if res.StatusCode >= 299 {

		// 400 Bad Response is likely an indicator that we failed to provide a
		// required field in our JSON payload, such as "Summary or Text is
		// required." when failing to supply such a field in the top level of
		// the MessageCard value that we send to the webhook URL.

		err = fmt.Errorf("error on notification: %v, %q", res.Status, responseString)
		logger.Println(err)
		return err
	}

	// log the response string
	logger.Printf("DEBUG: Response string: %v\n", responseString)

	return nil
}

// MessageCardSectionFact represents a section fact entry that is usually
// displayed in a two-column key/value format.
type MessageCardSectionFact struct {

	// Name is the key for an associated value in a key/value pair
	Name string `json:"name"`

	// Value is the value for an associated key in a key/value pair
	Value string `json:"value"`
}

// MessageCardPotentialAction represents an action that a user may take for a
// received Microsoft Teams message. Actions may be defined for the card
// as a whole and/or in specifc sections.
// FIXME: Add doc comments to each field.
type MessageCardPotentialAction struct {
	Target  []string `json:"target"`
	Context string   `json:"@context"`
	Type    string   `json:"@type"`

	// FIXME: The original conversion had `interface{}` as the type here. The
	// example JSON payload I found had `null` (no quotes) as the value for
	// this field. We may need to change this back to `interface{}` for
	// compatibility reasons.
	// https://stackoverflow.com/questions/31048557/assigning-null-to-json-fields-instead-of-empty-strings
	ID              string `json:"@id"`
	Name            string `json:"name"`
	IsPrimaryAction bool   `json:"isPrimaryAction"`
}

// https://golang.org/pkg/encoding/json/
//
// The "omitempty" option specifies that the field should be omitted from the
// encoding if the field has an empty value, defined as false, 0, a nil
// pointer, a nil interface value, and any empty array, slice, map, or string.

// MessageCardSection represents a section to include in a message card.
type MessageCardSection struct {

	// Title is the title property of a section. This property  is displayed
	// in a font that stands out, while not as prominent as the card's title.
	// It is meant to introduce the section and summarize its content,
	// similarly to how the card's title property is meant to summarize the
	// whole card.
	Title string `json:"title,omitempty"`

	// Text is the section's text property. This property is very similar to
	// the text property of the card. It can be used for the same purpose.
	Text string `json:"text,omitempty"`

	// Markdown represents a toggle to enable or disable Markdown formatting.
	// By default, all text fields in a card and its sections can be formatted
	// using basic Markdown.
	Markdown bool `json:"markdown,omitempty"`

	// Facts is a collection of MessageCardSectionFact values. A section entry
	// usually is displayed in a two-column key/value format.
	Facts []MessageCardSectionFact `json:"facts,omitempty"`

	// PotentialAction is a collection of actions that can be invoked on this
	// section.
	PotentialAction []MessageCardPotentialAction `json:"potentialAction,omitempty"`

	// StartGroup is the section's startGroup property. This property marks
	// the start of a logical group of information. Typically, sections with
	// startGroup set to true will be visually separated from previous card
	// elements.
	StartGroup bool `json:"startGroup,omitempty"`
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
	Title string `json:"title,omitempty"`

	// Text is required if the card does not contain a summary property,
	// otherwise optional. The text property is meant to be displayed in a
	// normal font below the card's title. Use it to display content, such as
	// the description of the entity being referenced, or an abstract of a
	// news article.
	Text string `json:"text,omitempty"`

	// Specifies a custom brand color for the card. The color will be
	// displayed in a non-obtrusive manner.
	ThemeColor string `json:"themeColor,omitempty"`

	// Sections is a collection of sections to include in the card.
	Sections []MessageCardSection `json:"sections,omitempty"`

	// PotentialAction is a collection of actions that can be invoked on this card.
	PotentialAction []MessageCardPotentialAction `json:"potentialAction,omitempty"`
}

// NewMessageCard creates a new message card with required fields required by
// the legacy message card format already predefined
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

// NewMessageCardSection creates an empty message card section
func NewMessageCardSection() MessageCardSection {

	msgCardSection := MessageCardSection{}

	return msgCardSection

}

// NewMessageCardPotentialAction creates an empty potential action value. This
// value may be added or attached to the main message card value OR to a
// specific message card section.
func NewMessageCardPotentialAction() MessageCardPotentialAction {

	// Example potentialAction portion of MS Teams JSON payload:
	//
	// "potentialAction": [
	// 	{
	// 		"target": [
	// 			"http://web.example.local:8000/app/search/@go?sid=scheduler_admin_search_W2_at_14232356_132"
	// 		],
	// 		"@context": "http://schema.org",
	// 		"@type": "ViewAction",
	// 		"@id": null,
	// 		"name": "View full Splunk report",
	// 		"isPrimaryAction": true
	// 	}
	// ]

	msgCardPotentialAction := MessageCardPotentialAction{}

	// prefill required fields
	msgCardPotentialAction.Context = "http://schema.org"

	return msgCardPotentialAction
}

// helper --------------------------------------------------------------------------------------------------------------

// IsValidWebhookURL performs validation checks on the webhook URL used to
// submit messages to Microsoft Teams.
func IsValidWebhookURL(webhookURL string) (bool, error) {
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
