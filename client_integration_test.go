package notify_test

import (
	"os"
	"testing"

	notify "github.com/govau/notify-client-go"
)

func setup(t *testing.T) (*notify.Client, string, string) {
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		t.Fatal("API_KEY environment variable not set")
	}

	smsTemplateId := os.Getenv("SMS_TEMPLATE_ID")
	if smsTemplateId == "" {
		t.Fatal("SMS_TEMPLATE_ID environment variable not set")
	}

	emailTemplateId := os.Getenv("EMAIL_TEMPLATE_ID")
	if emailTemplateId == "" {
		t.Fatal("EMAIL_TEMPLATE_ID environment variable not set")
	}

	var err error
	client, err := notify.NewClient(
		apiKey,
	)
	if err != nil {
		t.Fatal("Error creating client", err)
	}

	return client, smsTemplateId, emailTemplateId
}

func TestSendSMS(t *testing.T) {
	client, smsTemplateId, _ := setup(t)
	phoneNumber := os.Getenv("SMS_RECIPIENT_NUMBER")
	if phoneNumber == "" {
		t.Fatal("SMS_RECIPIENT_NUMBER environment variable not set")
	}

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

	if resp.ID == "" {
		t.Errorf("Response ID should not be empty")
	}
	if resp.URI == "" {
		t.Errorf("Response URI should not be empty")
	}
	assertEqual(t, ref, *resp.Reference)
	assertEqual(t, "Hello John,\n\nToday is Friday.", resp.Content.Body)
}

func TestSendEmail(t *testing.T) {
	client, _, emailTemplateId := setup(t)

	emailAddress := os.Getenv("EMAIL_RECIPIENT")
	if emailAddress == "" {
		t.Fatal("EMAIL_RECIPIENT environment variable not set")
	}

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

	if resp.ID == "" {
		t.Errorf("Response ID should not be empty")
	}
	if resp.URI == "" {
		t.Errorf("Response URI should not be empty")
	}
	assertEqual(t, "Hi John,\n\nMy favourite colour is pink.\n", resp.Content.Body)
}

func TestGetAllTemplates(t *testing.T) {
	client, smsTemplateId, emailTemplateId := setup(t)
	resp, err := client.Templates("")
	if err != nil {
		t.Fatalf("Error fetching all templates: %v", err)
	}

	assertTemplateFound(t, resp, []string{emailTemplateId, smsTemplateId})
}

func TestGetSMSTemplates(t *testing.T) {
	client, smsTemplateId, _ := setup(t)
	resp, err := client.Templates("sms")
	if err != nil {
		t.Fatalf("Error fetching sms templates: %v", err)
	}

	assertTemplateFound(t, resp, []string{smsTemplateId})
}

func TestGetEmailTemplates(t *testing.T) {
	client, _, emailTemplateId := setup(t)
	resp, err := client.Templates("email")
	if err != nil {
		t.Fatalf("Error fetching email templates: %v", err)
	}

	assertTemplateFound(t, resp, []string{emailTemplateId})
}

func TestGetTemplateById(t *testing.T) {
	client, smsTemplateId, _ := setup(t)
	resp, err := client.TemplateByID(smsTemplateId)
	if err != nil {
		t.Fatalf("Error fetching template by id: %v", err)
	}
	assertEqual(t, "go-sdk-test-sms", resp.Name)
	assertEqual(t, "sms", resp.Type)
}

func TestGetTemplateVersion(t *testing.T) {
	client, _, emailTemplateId := setup(t)
	resp, err := client.TemplateVersion(emailTemplateId, 1)
	if err != nil {
		t.Fatalf("Error fetching template version: %v", err)
	}

	assertEqual(t, "go-sdk-test-email", resp.Name)
	assertEqual(t, "email", resp.Type)
	assertEqual(t, 1, resp.Version)
	assertEqual(t, "Hi ((name)),\r\n\r\nMy favourite colour is ((colour)).", resp.Body)
}

func TestTemplatePreview(t *testing.T) {
	client, _, emailTemplateId := setup(t)
	resp, err := client.TemplatePreview(emailTemplateId, notify.Personalisation{
		{"name", "KD"},
		{"colour", "yellow"},
	})
	if err != nil {
		t.Fatalf("Error fetching email templates: %v", err)
	}

	assertEqual(t, "Hi KD,\n\nMy favourite colour is yellow.\n", resp.Body)
	assertEqual(t, "email", resp.Type)
}

func assertEqual(t *testing.T, expected, actual interface{}) {
	if expected != actual {
		t.Errorf("Expected %v but received %v", expected, actual)
	}
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
			t.Errorf("Template id not found: %s", id)
		}
	}
}
