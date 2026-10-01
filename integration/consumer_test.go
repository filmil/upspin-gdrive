// SPDX-License-Identifier: Apache-2.0

// Package integration uses upspin_drive from outside its module, as a
// module that depends on it would. It needs no Google account: it checks
// that linking the drive package registers the "Drive" storage with upspin,
// and that the storage rejects a dial with no credentials.
package integration

import (
	"testing"

	_ "github.com/filmil/upspin-gdrive/cloud/storage/drive"
	"upspin.io/cloud/storage"
	"upspin.io/errors"
)

func TestDriveStorageIsRegistered(t *testing.T) {
	_, err := storage.Dial("Drive")
	if err == nil {
		t.Fatal("Dial with no credentials succeeded")
	}
	// An unregistered name fails as NotExist; the Drive constructor
	// rejects the missing accessToken as Invalid.
	if !errors.Is(errors.Invalid, err) {
		t.Fatalf("Dial(\"Drive\") = %v, want an Invalid error from the Drive constructor", err)
	}
}
