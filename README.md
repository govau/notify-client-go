# Notify Go client

This documentation is for developers interested in using a Go client to integrate with Notify.

## Table of Contents

- [Installation](#installation)
- [Getting started](#getting-started)
- [Send messages](#send-messages)
- [Errors](#errors)
- [Tests](#tests)

## Installation

```shell
go get -u github.com/govau/notify-client-go
```

## Getting started

```Go
client, err := notify.NewClient(apiKey)
```

Generate an API key by logging in to [Notify.gov.au](https://notify.gov.au) and going to the _API integration_ page.

## Send messages

### Text message

#### Method

<details>
<summary>
Click here to expand for more information.
</summary>

```go
resp, err := client.SendSMS(
  templateID,
  phoneNumber,
  notify.Reference("Sam's reminders"),
  notify.Personalisation{
    {"name", "Sam"},
  },
)
```

</details>

#### Response

If the request is successful, `response` will be a `struct`.

<details>
<summary>
Click here to expand for more information.
</summary>

```go
{
    ID: "bfb50d92-100d-4b8b-b559-14fa3b091cda",
    Reference: "Sam's reminders",
    Content: {
        Body: "Hi Sam, just a reminder to visit the post office today.",
        FromNumber: "0400000000"
    },
    URI: "https://rest-api.notify.gov.au/v2/notifications/ceb50d92-100d-4b8b-b559-14fa3b091cd",
    Template: {
        ID: "ceb50d92-100d-4b8b-b559-14fa3b091cda",
        Version: 1,
        URI: "https://rest-api.notify.gov.au/v2/templates/bfb50d92-100d-4b8b-b559-14fa3b091cda"
    },
}
```

</details>

#### Arguments

<details>
<summary>
Click here to expand for more information.
</summary>

##### `phoneNumber`

The phone number of the recipient, only required for sms notifications.

##### `templateID`

Find by clicking **API info** for the template you want to send.

##### `options`

###### `Reference`

An optional identifier you generate. The `Reference` can be used as a unique reference for the notification. Because Notify does not require this reference to be unique you could also use this reference to identify a batch or group of notifications.

You can omit this argument if you do not require a reference for the notification.

##### `Personalisation`

If a template has placeholders, you need to provide their values, for example:

```go
p := notify.Personalisation{
    {"name", "Daniel Smith"},
    {"age", "23"}
},
```

This does not need to be provided if your template does not contain placeholders.

##### `SMSSenderID`

Optional. Specifies the identifier of the sms sender to set for the notification. The identifiers are found in your service Settings, when you 'Manage' your 'Text message sender'.

If you omit this argument your default sms sender will be set for the notification.

Example usage with optional reference -

</details>

### Email

#### Method

<details>
<summary>
Click here to expand for more information.
</summary>

```go
sent, erra := client.SendEmail(
    "effc255a-d233-4f3f-949a-15915c45b6f0",
    "dan@email.com",
    notify.Personalisation{
        {"name", "Dan"},
    },
)
```

</details>

#### Response

If the request is successful, `response` will be a `struct`.

<details>
<summary>
Click here to expand for more information.
</summary>

```go
{
    ID: "bfb50d92-100d-4b8b-b559-14fa3b091cda",
    Reference: "Sam's reminders",
    Content: {
        Subject: "Physio",
        Body: "Hi Sam, you have a physio appointment at 2pm.",
        FromEmail: "reminders@email.com"
    },
    URI: "https://rest-api.notify.gov.au/v2/notifications/ceb50d92-100d-4b8b-b559-14fa3b091cd",
    Template: {
        ID: "ceb50d92-100d-4b8b-b559-14fa3b091cda",
        Version: 1,
        URI: "https://rest-api.notify.gov.au/v2/templates/bfb50d92-100d-4b8b-b559-14fa3b091cda"
    },
}
```

</details>

#### Arguments

<details>
<summary>
Click here to expand for more information.
</summary>

##### `emailAddress`

The email address of the recipient, only required for email notifications.

##### `templateID`

Find by clicking **API info** for the template you want to send.

##### `options`

###### `Reference`

An optional identifier you generate. The `reference` can be used as a unique reference for the notification. Because Notify does not require this reference to be unique you could also use this reference to identify a batch or group of notifications.

You can omit this argument if you do not require a reference for the notification.

###### `EmailReplyToID`

Optional. Specifies the identifier of the email reply-to address to set for the notification. The identifiers are found in your service Settings, when you 'Manage' your 'Email reply to addresses'.

If you omit this argument your default email reply-to address will be set for the notification.

###### `Personalisation`

If a template has placeholders, you need to provide their values, for example:

```go
p := notify.Personalisation{
    {"name", "Daniel Smith"},
    {"age", "23"}
},
```

</details>

## Get all templates

#### Method

Returns the latest version of each template

<details>
<summary>
Click here to expand for more information.
</summary>

```go
templates, err := client.Templates()
```

</details>

#### Response

If the request is successful, `response` will be a `slice` of templates.

<details>
<summary>
Click here to expand for more information.
</summary>

```go
[{
    ID:        "template ID",
    Name:      "template name",
    Type:      "email",
    CreatedAt: "2019-04-10T00:46:29.076570Z",
    UpdatedAt: "2019-04-10T00:46:29.076570Z",
    CreatedBy: "someone@email.com",
    Version:   1,
    Subject:   "Subject of an email or letter notification, or nil if an sms message",
    Body:      "Body of the notification",
}, {
    ...template
}]
```

</details>

## Errors

#### Basic error

```go
resp, err := client.SendSMS(...)
// err = "Template not found"
```

#### Detailed API error

The `APIErrorFromError` method is provided to retrieve more granular error information when a request to the Notify service returns an error response.

```go

type APIError struct {
    Errors []struct {
        Error   string
        Message string
    }
    StatusCode int64
}
```

```go
resp, err := client.SendSMS(...)
if err != nil {
    apiErr, ok := notify.APIErrorFromError(err)
    if ok {
        log.Println(apiErr.Errors)
        log.Println(apiErr.StatusCode)
    } else {
        log.Println(err)
    }
}

```

Outputs:

```shell
[{
    Error: "BadRequestError"
    Message: "Template not found"
}]
400
```

or

```shell
Some non-request error
```

## Tests

To run the unit tests:

```sh
go test ./...
```
