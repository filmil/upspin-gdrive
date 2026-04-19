package drive

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/filmil/upspin-gdrive/config"
	"upspin.io/cloud/storage"
)

func TestTokenRefreshUsesClientIdAndSecret(t *testing.T) {
	var (
		gotClientID     string
		gotClientSecret string
		tokenCallCount  int
		apiCallCount    int
	)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/token" {
			tokenCallCount++
			err := r.ParseForm()
			if err != nil {
				t.Errorf("failed to parse form: %v", err)
			}
			gotClientID = r.Form.Get("client_id")
			gotClientSecret = r.Form.Get("client_secret")

			// Return a valid fake token
			w.Header().Set("Content-Type", "application/json")
			io.WriteString(w, `{"access_token":"new-fake-token","expires_in":3600,"token_type":"Bearer"}`)
			return
		}

		if r.URL.Path == "/drive/v3/files" {
			apiCallCount++
			// Return a fake list of files
			w.Header().Set("Content-Type", "application/json")
			io.WriteString(w, `{"files":[]}`)
			return
		}

		t.Fatalf("unexpected request to %s", r.URL.Path)
	}))
	defer ts.Close()

	// Temporarily override the token endpoint
	origEndpoint := config.OAuth2.Endpoint
	defer func() { config.OAuth2.Endpoint = origEndpoint }()
	config.OAuth2.Endpoint.TokenURL = ts.URL + "/token"

	// Create storage with expired token and specific client credentials
	opts := &storage.Opts{
		Opts: map[string]string{
			"accessToken":  "expired-token",
			"tokenType":    "Bearer",
			"refreshToken": "fake-refresh-token",
			"expiry":       time.Now().Add(-1 * time.Hour).Format(time.RFC3339),
			"clientId":     "expected-client-id",
			"clientSecret": "expected-client-secret",
		},
	}

	st, err := New(opts)
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	dImpl, ok := st.(*driveImpl)
	if !ok {
		t.Fatalf("expected *driveImpl, got %T", st)
	}

	// Hack the base path to point to our test server
	dImpl.files.BasePath = ts.URL + "/drive/v3/"

	// Trigger an API call, which should trigger a token refresh
	lister, ok := st.(storage.Lister)
	if !ok {
		t.Fatalf("storage does not implement Lister")
	}

	_, _, err = lister.List("")
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if tokenCallCount != 1 {
		t.Errorf("expected 1 token refresh call, got %d", tokenCallCount)
	}
	if apiCallCount != 1 {
		t.Errorf("expected 1 API call, got %d", apiCallCount)
	}

	if gotClientID != "expected-client-id" {
		t.Errorf("expected client_id %q, got %q", "expected-client-id", gotClientID)
	}
	if gotClientSecret != "expected-client-secret" {
		t.Errorf("expected client_secret %q, got %q", "expected-client-secret", gotClientSecret)
	}
}
