package client

import (
	"testing"

	"github.com/nano-container-linux/libdnsd"
)

func TestACMEChallengeFQDN(t *testing.T) {
	cases := map[string]string{
		"example.com":                  "_acme-challenge.example.com.",
		"example.com.":                 "_acme-challenge.example.com.",
		"*.example.com.":               "_acme-challenge.example.com.",
		"_acme-challenge.example.com.": "_acme-challenge.example.com.",
		"_ACME-CHALLENGE.EXAMPLE.COM":  "_acme-challenge.example.com.",
	}
	for in, want := range cases {
		got := libdnsd.AcmeChallengeFQDN(in)
		if got != want {
			t.Fatalf("AcmeChallengeFQDN(%q) = %q, want %q", in, got, want)
		}
	}
}
