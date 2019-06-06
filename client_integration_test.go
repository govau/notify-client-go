package notify_test

import (
	"os"
	"testing"

	notify "github.com/govau/notify-client-go"

	"github.com/stretchr/testify/assert"
)

var apiKey string
var smsTemplateId string
var emailTemplateId string
var phoneNumber string
var emailAddress string
var client *notify.Client

func setup(t *testing.T) {
	if client == nil {
		apiKey = os.Getenv("API_KEY")
		if apiKey == "" {
			t.Fatal("API_KEY environment variable not set")
		}

		smsTemplateId = os.Getenv("SMS_TEMPLATE_ID")
		if smsTemplateId == "" {
			t.Fatal("SMS_TEMPLATE_ID environment variable not set")
		}

		phoneNumber = os.Getenv("SMS_RECIPIENT_NUMBER")
		if phoneNumber == "" {
			t.Fatal("SMS_RECIPIENT_NUMBER environment variable not set")
		}

		emailTemplateId = os.Getenv("EMAIL_TEMPLATE_ID")
		if emailTemplateId == "" {
			t.Fatal("EMAIL_TEMPLATE_ID environment variable not set")
		}

		emailAddress = os.Getenv("EMAIL_RECIPIENT")
		if emailAddress == "" {
			t.Fatal("EMAIL_RECIPIENT environment variable not set")
		}

		var err error
		client, err = notify.NewClient(
			apiKey,
		)
		if err != nil {
			t.Fatal("Error creating client", err)
		}
	}
}

func assertEqual(t *testing.T, expected interface{}, actual interface{}) {
	if expected != actual {
		t.Fatalf("Expected (%v) actual (%v)", expected, actual)
	}
}

func TestSendSMS(t *testing.T) {
	setup(t)
	ref := "TestSendSMS"
	resp, err := client.SendSMS(
		smsTemplateId,
		phoneNumber,
		notify.Personalisation{
			{"name", "John"},
			{"day", "Friday"},
		},
		notify.Reference(ref),
	)
	if err != nil {
		t.Fatalf("Error sending SMS: %v", err)
	}

	assert.NotNil(t, resp.ID, resp.URI)
	assert.Equal(t, ref, *resp.Reference)
	assert.Equal(t, "Hello John,\n\nToday is Friday.", resp.Content.Body)
}

func TestSendEmail(t *testing.T) {
	setup(t)
	ref := "TestSendEmail"
	resp, err := client.SendEmail(
		emailTemplateId,
		emailAddress,
		notify.Personalisation{
			{"name", "John"},
			{"colour", "pink"},
		},
		notify.Reference(ref),
	)
	if err != nil {
		t.Fatalf("Error sending email: %v", err)
	}

	assert.NotNil(t, resp.ID, resp.URI)
	assert.Equal(t, ref, *resp.Reference)
	assert.Equal(t, "Hi John,\n\nMy favourite colour is pink.\n", resp.Content.Body)
}

func TestGetAllTemplates(t *testing.T) {
	setup(t)
	resp, err := client.Templates("")
	if err != nil {
		t.Fatalf("Error fetching all templates: %v", err)
	}

	assert.NotNil(t, resp)
	assertTemplateFound(t, resp, []string{emailTemplateId, smsTemplateId})
}

func TestGetSMSTemplates(t *testing.T) {
	setup(t)
	resp, err := client.Templates("sms")
	if err != nil {
		t.Fatalf("Error fetching sms templates: %v", err)
	}

	assert.NotNil(t, resp)
	assertTemplateFound(t, resp, []string{smsTemplateId})
}

func TestGetEmailTemplates(t *testing.T) {
	setup(t)
	resp, err := client.Templates("email")
	if err != nil {
		t.Fatalf("Error fetching email templates: %v", err)
	}

	assert.NotNil(t, resp)
	assertTemplateFound(t, resp, []string{emailTemplateId})
}

func TestGetTemplateById(t *testing.T) {
	setup(t)
	resp, err := client.TemplateByID(smsTemplateId)
	if err != nil {
		t.Fatalf("Error fetching template by id: %v", err)
	}

	assert.NotNil(t, resp)
	assert.Equal(t, "go-sdk-test-sms", resp.Name)
	assert.Equal(t, "sms", resp.Type)
}

func TestGetTemplateVersion(t *testing.T) {
	setup(t)
	resp, err := client.TemplateVersion(emailTemplateId, 1)
	if err != nil {
		t.Fatalf("Error fetching template version: %v", err)
	}

	assert.NotNil(t, resp)
	assert.Equal(t, "go-sdk-test-email", resp.Name)
	assert.Equal(t, "email", resp.Type)
	assert.Equal(t, 1, resp.Version)
	assert.Equal(t, "Hi ((name)),\r\n\r\nMy favourite colour is ((colour)).", resp.Body)
}

func TestTemplatePreview(t *testing.T) {
	setup(t)
	resp, err := client.TemplatePreview(emailTemplateId, notify.Personalisation{
		{"name", "KD"},
		{"colour", "yellow"},
	})
	if err != nil {
		t.Fatalf("Error fetching email templates: %v", err)
	}

	assert.NotNil(t, resp)
	assert.Equal(t, "Hi KD,\n\nMy favourite colour is yellow.\n", resp.Body)
	assert.Equal(t, "email", resp.Type)
}

func assertTemplateFound(t *testing.T, templates notify.Templates, templateIds []string) {

	for _, id := range templateIds {
		found := false
		for _, template := range templates {
			if template.ID == id {
				found = true
			}
		}
		if !found {
			assert.Failf(t, "Template id not found: %s", id)
		}
	}
}
