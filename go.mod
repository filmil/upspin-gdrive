module github.com/filmil/upspin-gdrive

go 1.26.2

toolchain go1.26.2

require (
	golang.org/x/oauth2 v0.36.0
	google.golang.org/api v0.276.0
	upspin.io v42.1.6+incompatible
)

replace upspin.io => github.com/filmil/upspin v0.1.0
