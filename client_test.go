package notify_test

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/govau/notify-client-go"
)

func TestNewClient(t *testing.T) {
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
	} {
		if !strings.Contains(request, term) {
			t.Errorf("request did not contain %s", term)
		}
	}
}
