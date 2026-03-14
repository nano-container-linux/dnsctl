package client

import "testing"

func TestACMEChallengeFQDN(t *testing.T) {
cases := map[string]string{
"example.com":                  "_acme-challenge.example.com.",
"example.com.":                 "_acme-challenge.example.com.",
"*.example.com.":               "_acme-challenge.example.com.",
"_acme-challenge.example.com.": "_acme-challenge.example.com.",
"_ACME-CHALLENGE.EXAMPLE.COM":  "_acme-challenge.example.com.",
}
for in, want := range cases {
got := acmeChallengeFQDN(in)
if got != want {
t.Fatalf("acmeChallengeFQDN(%q) = %q, want %q", in, got, want)
}
}
}
