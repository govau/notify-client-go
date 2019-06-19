package notify_test

import (
	"os"
	"testing"

	notify "github.com/govau/notify-client-go"
)

func setup(t *testing.T) (client *notify.Client, smsTemplateID string, emailTemplateID string) {
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		t.Fatal("API_KEY environment variable must be set")
	}

	smsTemplateID = os.Getenv("SMS_TEMPLATE_ID")
	if smsTemplateID == "" {
		t.Fatal("SMS_TEMPLATE_ID environment variable must be set")
	}

	emailTemplateID = os.Getenv("EMAIL_TEMPLATE_ID")
	if emailTemplateID == "" {
		t.Fatal("EMAIL_TEMPLATE_ID environment variable must be set")
	}

	var err error
	client, err = notify.NewClient(
		apiKey,
	)
	if err != nil {
		t.Fatal("could not create Notify client", err)
	}
	return
}

func TestSendSMS(t *testing.T) {
	client, smsTemplateID, _ := setup(t)
	phoneNumber := os.Getenv("SMS_RECIPIENT_NUMBER")
	if phoneNumber == "" {
		t.Fatal("SMS_RECIPIENT_NUMBER environment variable must be set")
	}

	ref := "TestSendSMS"
	resp, err := client.SendSMS(
		smsTemplateID,
		phoneNumber,
		notify.Personalisation{
			{"name", "John"},
			{"day", "Friday"},
		},
		notify.Reference(ref),
	)
	if err != nil {
		t.Fatalf("could not send SMS: %v", err)
	}
	if resp.ID == "" {
		t.Errorf("Response ID should not be empty")
	}
	if resp.URI == "" {
		t.Errorf("Response URI should not be empty")
	}
	if *resp.Reference != ref {
		t.Errorf("got %v, want %v", *resp.Reference, ref)
	}
	wantedBody := "Hello John,\n\nToday is Friday."
	if resp.Content.Body != wantedBody {
		t.Errorf("got %v, want %v", resp.Content.Body, wantedBody)
	}
}

func TestSendEmail(t *testing.T) {
	client, _, emailTemplateID := setup(t)

	emailAddress := os.Getenv("EMAIL_RECIPIENT")
	if emailAddress == "" {
		t.Fatal("EMAIL_RECIPIENT environment variable not set")
	}

	ref := "TestSendEmail"
	resp, err := client.SendEmail(
		emailTemplateID,
		emailAddress,
		notify.Personalisation{
			{"name", "John"},
			{"colour", "pink"},
		},
		notify.Reference(ref),
	)
	if err != nil {
		t.Fatalf("could not send email: %v", err)
	}
	if resp.ID == "" {
		t.Errorf("Response ID should not be empty")
	}
	if resp.URI == "" {
		t.Errorf("Response URI should not be empty")
	}
	wantedBody := "Hi John,\n\nMy favourite colour is pink.\n"
	if resp.Content.Body != wantedBody {
		t.Errorf("got %v, want %v", resp.Content.Body, wantedBody)
	}
}

func TestGetAllTemplates(t *testing.T) {
	client, smsTemplateID, emailTemplateID := setup(t)
	resp, err := client.GetAllTemplates("")
	if err != nil {
		t.Fatalf("could not fetch all templates: %v", err)
	}

	assertTemplateFound(t, resp, []string{emailTemplateID, smsTemplateID})
}

func TestGetSMSTemplates(t *testing.T) {
	client, smsTemplateID, _ := setup(t)
	resp, err := client.GetAllTemplates("sms")
	if err != nil {
		t.Fatalf("could not fetch sms templates: %v", err)
	}

	assertTemplateFound(t, resp, []string{smsTemplateID})
}

func TestGetEmailTemplates(t *testing.T) {
	client, _, emailTemplateID := setup(t)
	resp, err := client.GetAllTemplates("email")
	if err != nil {
		t.Fatalf("could not fetch email templates: %v", err)
	}

	assertTemplateFound(t, resp, []string{emailTemplateID})
}

func TestGetTemplateByID(t *testing.T) {
	client, smsTemplateID, _ := setup(t)
	resp, err := client.GetTemplateByID(smsTemplateID)
	if err != nil {
		t.Fatalf("could not fetch template by id: %v", err)
	}
	wantedName := "go-sdk-test-sms"
	if resp.Name != wantedName {
		t.Errorf("got %v, want %v", resp.Name, wantedName)
	}
	wantedType := "sms"
	if resp.Type != wantedType {
		t.Errorf("got %v, want %v", resp.Type, wantedType)
	}
}

func TestGetTemplateByIDAndVersion(t *testing.T) {
	client, _, emailTemplateID := setup(t)
	resp, err := client.GetTemplateByIDAndVersion(emailTemplateID, 1)
	if err != nil {
		t.Fatalf("could not fetch template version: %v", err)
	}

	wantedName := "go-sdk-test-email"
	if resp.Name != wantedName {
		t.Errorf("got %v, want %v", resp.Name, wantedName)
	}
	wantedType := "email"
	if resp.Type != wantedType {
		t.Errorf("got %v, want %v", resp.Type, wantedType)
	}
	wantedVersion := 1
	if resp.Version != wantedVersion {
		t.Errorf("got %v, want %v", resp.Version, wantedVersion)
	}
	wantedBody := "Hi ((name)),\r\n\r\nMy favourite colour is ((colour))."
	if resp.Body != wantedBody {
		t.Errorf("got %v, want %v", resp.Body, wantedBody)
	}
}

func TestGenerateTemplatePreview(t *testing.T) {
	client, _, emailTemplateID := setup(t)
	resp, err := client.GenerateTemplatePreview(emailTemplateID, notify.Personalisation{
		{"name", "KD"},
		{"colour", "yellow"},
	})
	if err != nil {
		t.Fatalf("could not fetch email templates: %v", err)
	}
	wantedBody := "Hi KD,\n\nMy favourite colour is yellow.\n"
	if resp.Body != wantedBody {
		t.Errorf("got %v, want %v", resp.Body, wantedBody)
	}
	wantedType := "email"
	if resp.Type != wantedType {
		t.Errorf("got %v, want %v", resp.Type, wantedType)
	}
}

func assertTemplateFound(t *testing.T, templates notify.Templates, templateIDs []string) {
	for _, id := range templateIDs {
		found := false
		for _, template := range templates {
			if template.ID == id {
				found = true
			}
		}
		if !found {
			t.Errorf("could not find template id: %s", id)
		}
	}
}
