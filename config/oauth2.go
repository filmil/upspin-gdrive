// Copyright 2017 The Upspin Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The config package holds OAuth2 configuration data shared by the drive storage
// and the setupstorage-drive command.
package config // import "github.com/filmil/upspin-gdrive/config"

import (
	"os"

	"golang.org/x/oauth2"
	"google.golang.org/api/drive/v3"
)

// OAuth2 holds OAuth configuration used by the Upspin Google Drive package. It is used by both
// the storage and the setup process.
var OAuth2 = &oauth2.Config{
	ClientID:     getEnvOrDefault("GOOGLE_CLIENT_ID", ""),
	ClientSecret: getEnvOrDefault("GOOGLE_CLIENT_SECRET", ""),
	Endpoint: oauth2.Endpoint{
		AuthURL:  "https://accounts.google.com/o/oauth2/auth",
		TokenURL: "https://oauth2.googleapis.com/token",
	},
	Scopes: []string{drive.DriveAppdataScope},
}

func getEnvOrDefault(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
