// Copyright 2020 Enrico Hoffmann
// Copyright 2020 Adam Chalkley
//
// https:#github.com/atc0005/go-teams-notify
//
// Licensed under the MIT License. See LICENSE file in the project root for
// full license information.

package goteamsnotify

import (
	"errors"
	"net/http"
	"net/url"
	"testing"
)

// http://hassansin.github.io/Unit-Testing-http-client-in-Go
func TestTeamsClientSend(t *testing.T) {
	simpleMsgCard := NewMessageCard()
	simpleMsgCard.Text = "Hello World"
	var tests = []struct {
		reqURL    string
		reqMsg    MessageCard
		resStatus int   // httpClient response status
		resError  error // httpClient error
		error     error // method error
	}{
		// invalid webhookURL - url.Parse error
		{
			reqURL:    "ht\ttp://",
			reqMsg:    simpleMsgCard,
			resStatus: 0,
			resError:  nil,
			error:     &url.Error{},
		},
		// invalid webhookURL - missing prefix in webhook URL
		{
			reqURL:    "",
			reqMsg:    simpleMsgCard,
			resStatus: 0,
			resError:  nil,
			error:     errors.New(""),
		},
		// invalid httpClient.Do call
		{
			reqURL:    "https://outlook.office.com/webhook/xxx",
			reqMsg:    simpleMsgCard,
			resStatus: 200,
			resError:  errors.New("pling"),
			error:     &url.Error{},
		},
		// invalid httpClient.Do call
		{
			reqURL:    "https://outlook.office365.com/webhook/xxx",
			reqMsg:    simpleMsgCard,
			resStatus: 200,
			resError:  errors.New("pling"),
			error:     &url.Error{},
		},
		// invalid response status code
		{
			reqURL:    "https://outlook.office.com/webhook/xxx",
			reqMsg:    simpleMsgCard,
			resStatus: 400,
			resError:  nil,
			error:     errors.New(""),
		},
	}
	for idx, test := range tests {
		// Create range scoped var for use within closure
		test := test

		client := NewTestClient(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: test.resStatus,

				// NOTE: Intentionally NOT setting the Body field as a test
				// case between Go 1.14.x and Go 15.
				// https://github.com/atc0005/go-teams-notify/pull/43
				//Body:       ioutil.NopCloser(bytes.NewBufferString(`OK`)),

				// Must be set to non-nil value or it panics
				Header: make(http.Header),
			}, test.resError
		})
		c := &teamsClient{httpClient: client}

		err := c.Send(test.reqURL, test.reqMsg)

		if err != nil {
			if !errors.As(err, &test.error) {
				t.Fatalf(
					"FAIL: test %d; got %T, want %T",
					idx,
					errors.Unwrap(err),
					test.error,
				)
			} else {
				t.Logf(
					"OK: test %d; test.error is of type %T, err is of type %T",
					idx,
					test.error,
					err,
				)
			}
		} else {
			t.Logf("OK: test %d; no error received", idx)
		}
	}
}

type RoundTripFunc func(req *http.Request) (*http.Response, error)

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

// NewTestClient returns *http.API with Transport replaced to avoid making real calls
func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: fn,
	}
}
