<!-- omit in toc -->
# go-teams-notify

A package to send messages to Microsoft Teams (channels)

[![Latest release][githubtag-image]][githubtag-url]
[![GoDoc][godoc-image]][godoc-url]
[![License][license-image]][license-url]
[![Validate Codebase](https://github.com/atc0005/go-teams-notify/workflows/Validate%20Codebase/badge.svg)](https://github.com/atc0005/go-teams-notify/actions?query=workflow%3A%22Validate+Codebase%22)
[![Validate Docs](https://github.com/atc0005/go-teams-notify/workflows/Validate%20Docs/badge.svg)](https://github.com/atc0005/go-teams-notify/actions?query=workflow%3A%22Validate+Docs%22)
[![Lint and Build using Makefile](https://github.com/atc0005/go-teams-notify/workflows/Lint%20and%20Build%20using%20Makefile/badge.svg)](https://github.com/atc0005/go-teams-notify/actions?query=workflow%3A%22Lint+and+Build+using+Makefile%22)
[![Quick Validation](https://github.com/atc0005/go-teams-notify/workflows/Quick%20Validation/badge.svg)](https://github.com/atc0005/go-teams-notify/actions?query=workflow%3A%22Quick+Validation%22)

<!-- omit in toc -->
## Table of contents

- [Usage](#usage)
- [References](#references)

## Usage

To get the package, execute:

```console
go get https://github.com/atc0005/go-teams-notify/v2
```

To import this package, add the following line to your code:

```golang
import "github.com/atc0005/go-teams-notify/v2"
```

And this is an example of a simple implementation ...

```golang
import (
  "github.com/atc0005/go-teams-notify/v2"
)

func main() {
  _ = sendTheMessage()
}

func sendTheMessage() error {
  // init the client
  mstClient := goteamsnotify.NewClient()

  // setup webhook url
  webhookUrl := "https://outlook.office.com/webhook/YOUR_WEBHOOK_URL_OF_TEAMS_CHANNEL"

  // setup message card
  msgCard := goteamsnotify.NewMessageCard()
  msgCard.Title = "Hello world"
  msgCard.Text = "Here are some examples of formatted stuff like "+
      "<br> * this list itself  <br> * **bold** <br> * *italic* <br> * ***bolditalic***"
  msgCard.ThemeColor = "#DF813D"

  // send
  return mstClient.Send(webhookUrl, msgCard)
}
```

## References

- MS Teams - adaptive cards
  ([de-de](https://docs.microsoft.com/de-de/outlook/actionable-messages/adaptive-card),
  [en-us](https://docs.microsoft.com/en-us/outlook/actionable-messages/adaptive-card))
- MS Teams - send via connectors
  ([de-de](https://docs.microsoft.com/de-de/outlook/actionable-messages/send-via-connectors),
  [en-us](https://docs.microsoft.com/en-us/outlook/actionable-messages/send-via-connectors))
- [adaptivecards.io](https://adaptivecards.io/designer)

[githubtag-image]: https://img.shields.io/github/release/atc0005/go-teams-notify.svg?style=flat
[githubtag-url]: https://github.com/atc0005/go-teams-notify

[godoc-image]: https://godoc.org/github.com/atc0005/go-teams-notify?status.svg
[godoc-url]: https://godoc.org/github.com/atc0005/go-teams-notify

[license-image]: https://img.shields.io/github/license/atc0005/go-teams-notify.svg?style=flat
[license-url]: https://github.com/atc0005/go-teams-notify/blob/master/LICENSE
