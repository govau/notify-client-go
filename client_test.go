package notify_test

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	notify "github.com/govau/notify-client-go"
)

func TestNewClientAPIKey(t *testing.T) {
	type args struct {
	}
	tests := []struct {
		name    string
		apiKey  string
		wantErr bool
	}{
		{
			name:    "API key not provided",
			apiKey:  "",
			wantErr: true,
		},
		{
			name:    "API key too short",
			apiKey:  "1-2-3-4",
			wantErr: true,
		},
		{
			name:    "API key too short",
			apiKey:  "20bd0d9a-feda-4c75-97bd-a206ecb4019b",
			wantErr: true,
		},
		{
			name:    "API key without name prefix",
			apiKey:  "1af19ba3-1f4b-4014-af6f-bb917e0e14b3-20bd0d9a-feda-4c75-97bd-a206ecb4019b",
			wantErr: false,
		},
		{
			name:    "API key with name prefix",
			apiKey:  "key_name-1af19ba3-1f4b-4014-af6f-bb917e0e14b3-20bd0d9a-feda-4c75-97bd-a206ecb4019b",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := notify.NewClient(tt.apiKey); (err != nil) != tt.wantErr {
				t.Errorf("NewClient(%s) error = %v, wantErr %v", tt.apiKey, err, tt.wantErr)
			}
		})
	}
}

func TestNewClientAndSend(t *testing.T) {
	c := make(chan string, 1)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var buf bytes.Buffer
		io.Copy(&buf, r.Body)
		c <- buf.String()

		fmt.Fprintln(w, "{}")
	}))
	defer ts.Close()

	client, err := notify.NewClient(
		"key_name-95b3b534-bdd6-4f26-ad91-84b4e2301cca-e8a5f59a-b445-4dc0-9513-c5831615f937",
		notify.WithBaseURL(ts.URL),
	)

	if err != nil {
		t.Fatal(err)
	}

	_, err = client.SendEmail(
		"83f8a64f-74ec-4d90-ae48-394a8af3fe7c",
		"someone@example.com",
		notify.Reference("your-local-identifier"),
		notify.EmailReplyToID("ef661d9d-17bc-4fc2-ae12-1d3e00ae202e"),
		notify.Personalisation{
			{"user_name", "Sam"},
			{"amount_owing", "$205.20"},
		},
		notify.StatusCallback("https://localhost/callback", "1234567890"),
	)
	if err != nil {
		t.Fatal(err)
	}

	request := <-c
	for _, term := range []string{
		"email_address",
		"email_reply_to_id",
		"personalisation",
		"reference",
		"template_id",
		"your-local-identifier",
		"Sam",
		"status_callback_url",
		"status_callback_bearer_token",
		"https://localhost/callback",
		"1234567890",
	} {
		if !strings.Contains(request, term) {
			t.Errorf("request did not contain %s", term)
		}
	}
}
