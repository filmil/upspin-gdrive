// Copyright 2017 The Upspin Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The upspin-setupstorage-drive command is an external upspin subcommand that
// executes the second step in establishing an upspinserver for Google Drive.
// Run upspin setupstorage-drive -help for more information.
package main // import "github.com/filmil/upspin-gdrive/cmd/upspin-setupstorage-drive"

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"upspin.io/subcmd"

	"github.com/filmil/upspin-gdrive/config"
	"golang.org/x/oauth2"
)

const help = `
setupstorage-drive is the second step in establishing an upspinserver,
It sets up Google Drive storage for your Upspin installation. You may
skip this step if you wish to store Upspin data on your server's local
disk.

The first step is 'setupdomain' and the final step is 'setupserver'.
setupstorage-drive will add the Upspin Drive Storage application to your
Drive account. It then writes the obtained OAuth2 token information to
$where/$domain/serviceaccount.json and updates the server configuration
files in that directory to use the selected account.

Simply follow the on-screen instructions which will guide you through the
process.
`

type state struct{ *subcmd.State }

func main() {
	const name = "setupstorage-drive"

	log.SetFlags(0)
	log.SetPrefix("upspin setupstorage-drive: ")

	where := flag.String("where", filepath.Join(os.Getenv("HOME"), "upspin", "deploy"), "`directory` to store private configuration files")
	domain := flag.String("domain", "", "domain `name` for this Upspin installation")

	s := &state{State: subcmd.NewState(name)}
	s.ParseFlags(flag.CommandLine, os.Args[1:], help, "setupstorage-drive -domain=<name>")
	if *domain == "" {
		s.Exitf("%s\nthe -domain flag must be provided")
	}

	tok := s.tokenFromWeb()
	cfgPath := filepath.Join(*where, *domain)
	cfg := s.ReadServerConfig(cfgPath)
	cfg.StoreConfig = []string{
		"backend=Drive",
		"accessToken=" + tok.AccessToken,
		"tokenType=" + tok.TokenType,
		"refreshToken=" + tok.RefreshToken,
		"expiry=" + tok.Expiry.Format(time.RFC3339),
	}
	s.WriteServerConfig(cfgPath, cfg)

	fmt.Fprintf(os.Stderr, "You should now deploy the upspinserver binary and run 'upspin setupserver'.\n")
	s.ExitNow()
}

// tokenFromWeb attempts to obtain an OAuth2 token via a local loopback server.
func (s *state) tokenFromWeb() *oauth2.Token {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		s.Exitf("unable to start local web server: %v", err)
	}
	defer listener.Close()

	port := listener.Addr().(*net.TCPAddr).Port
	redirectURL := fmt.Sprintf("http://127.0.0.1:%d", port)

	conf := *config.OAuth2
	conf.RedirectURL = redirectURL

	authURL := conf.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	fmt.Printf("Please visit the following URL in your browser to authorize Upspin:\n\n\t%s\n\n", authURL)

	// Attempt to open the browser automatically
	switch runtime.GOOS {
	case "linux":
		_ = exec.Command("xdg-open", authURL).Start()
	case "windows":
		_ = exec.Command("rundll32", "url.dll,FileProtocolHandler", authURL).Start()
	case "darwin":
		_ = exec.Command("open", authURL).Start()
	}

	fmt.Println("Waiting for authorization... (If your browser didn't open automatically, please click the link above)")

	codeCh := make(chan string)
	errCh := make(chan error)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			errStr := r.URL.Query().Get("error")
			fmt.Fprintf(w, "Error: %s. You can close this window and check the terminal.", errStr)
			errCh <- fmt.Errorf("authorization failed: %s", errStr)
			return
		}
		fmt.Fprintf(w, "Authorization successful! You can safely close this window and return to the terminal.")
		codeCh <- code
	})

	server := &http.Server{Handler: mux}
	go func() {
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	var code string
	select {
	case code = <-codeCh:
	case err := <-errCh:
		s.Exitf("authorization error: %v", err)
	}

	// Gracefully shutdown the server since we got what we needed
	_ = server.Shutdown(context.Background())

	tok, err := conf.Exchange(context.Background(), code)
	if err != nil {
		s.Exitf("unable to exchange code for token: %v", err)
	}
	return tok
}
