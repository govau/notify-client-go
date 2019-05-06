package notify

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"

	"github.com/govau/notify-client-go/internal/base"
)

type Client struct {
	c base.Client
}

type ClientOption func(base.Client) (base.Client, error)

func WithBaseURL(target string) ClientOption {
	return func(c base.Client) (base.Client, error) {
		baseURL, err := url.Parse(target)
		c.BaseURL = baseURL
		return c, err
	}
}

func validateAPIKey(apiKey string) error {
	if apiKey == "" {
		return errors.New("api key is empty")
	}
	// 73 is the min length accounting for an API key where the name prefix has
	// not been provided.
	if len(apiKey) < 73 {
		return errors.New("api key is too short")
	}
	return nil
}

func NewClient(apiKey string, options ...ClientOption) (*Client, error) {
	if err := validateAPIKey(apiKey); err != nil {
		return nil, fmt.Errorf("notify: %v", err)
	}

	var err error

	slice := func(start, end int) string {
		n := len(apiKey)
		return apiKey[n+start : n+end]
	}

	client := base.Client{
		ServiceID:   slice(-73, -37),
		APIKey:      slice(-36, 0),
		RouteSecret: "",
	}

	client, err = WithBaseURL(base.NotifyBaseURL)(client)
	if err != nil {
		return nil, err
	}

	for _, option := range options {
		client, err = option(client)
		if err != nil {
			return nil, err
		}
	}

	return &Client{client}, nil
}

func (c Client) Templates() (Templates, error) {
	var templates Templates
	err := c.c.Get("./v2/templates").JSON(&templates, "templates").Error
	return templates, err
}

func (c Client) SendEmail(
	templateID string,
	emailAddress string,
	options ...SendEmailOption,
) (SentEmail, error) {
	var response SentEmail
	var buf bytes.Buffer
	var p = payload{
		{"template_id", templateID},
		{"email_address", emailAddress},
	}

	for _, option := range options {
		p = option.updateEmailPayload(p)
	}

	err := json.NewEncoder(&buf).Encode(p)
	if err != nil {
		return response, err
	}

	err = c.c.Post("./v2/notifications/email", &buf).JSON(&response).Error
	return response, err
}

func (c Client) SendSMS(
	templateID string,
	phoneNumber string,
	options ...SendSMSOption,
) (SentSMS, error) {
	var response SentSMS
	var buf bytes.Buffer
	var p = payload{
		{"template_id", templateID},
		{"phone_number", phoneNumber},
	}

	for _, option := range options {
		p = option.updateSMSPayload(p)
	}

	err := json.NewEncoder(&buf).Encode(p)
	if err != nil {
		return response, err
	}

	err = c.c.Post("./v2/notifications/sms", &buf).JSON(&response).Error
	return response, err
}

type Template struct {
	ID        string `json:"id,omitempty"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	CreatedBy string `json:"created_by"`
	Version   int    `json:"version"`
	Subject   string `json:"subject"`
	Body      string `json:"body"`
}

type Templates []Template

type SentSMS struct {
	ID           string  `json:"id"`
	URI          string  `json:"uri"`
	Reference    *string `json:"reference"`
	ScheduledFor *string `json:"scheduled_for"`

	Content struct {
		Body       string `json:"body"`
		FromNumber string `json:"from_number"`
	} `json:"content"`

	Template struct {
		ID      string `json:"id"`
		URI     string `json:"uri"`
		Version int    `json:"version"`
	} `json:"template"`
}

type SentEmail struct {
	ID           string  `json:"id"`
	URI          string  `json:"uri"`
	Reference    *string `json:"reference"`
	ScheduledFor *string `json:"scheduled_for"`

	Content struct {
		Subject   string `json:"subject"`
		Body      string `json:"body"`
		FromEmail string `json:"from_email"`
	} `json:"content"`

	Template struct {
		ID      string `json:"id"`
		URI     string `json:"uri"`
		Version int    `json:"version"`
	} `json:"template"`
}
