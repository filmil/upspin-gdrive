// Copyright 2026 The Upspin Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package drive

import (
	"errors"
	"strings"
	"testing"
	"time"

	"golang.org/x/oauth2"
)

// fakeSource is a TokenSource with a fixed answer.
type fakeSource struct {
	tok *oauth2.Token
	err error
}

func (f fakeSource) Token() (*oauth2.Token, error) { return f.tok, f.err }

func TestTokenCheck(t *testing.T) {
	ok := tokenCheck(fakeSource{tok: &oauth2.Token{AccessToken: "a", Expiry: time.Now().Add(time.Hour)}})
	if err := ok(); err != nil {
		t.Errorf("a valid token: %v", err)
	}
	bad := tokenCheck(fakeSource{err: errors.New(`oauth2: "invalid_grant" "Token has been expired or revoked."`)})
	err := bad()
	if err == nil || !strings.Contains(err.Error(), "invalid_grant") || !strings.Contains(err.Error(), "tokenCheck") {
		t.Errorf("a revoked refresh token: %v, want an error naming tokenCheck and invalid_grant", err)
	}
}
