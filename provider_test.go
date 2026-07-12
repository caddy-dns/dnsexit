package dnsexit

import (
	"testing"

	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/libdns/dnsexit"
)

func TestProviderInstantiation(t *testing.T) {
	p := &Provider{Provider: &dnsexit.Provider{APIKey: "test"}}
	if p.Provider.APIKey != "test" {
		t.Errorf("Expected APIKey to be 'test', got %s", p.Provider.APIKey)
	}
}

func TestUnmarshalCaddyfileAcceptsAPIToken(t *testing.T) {
	input := `dnsexit {
		api_token test-token
	}`
	p := &Provider{}
	d := caddyfile.NewTestDispenser(input)
	if err := p.UnmarshalCaddyfile(d); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if p.Provider.APIKey != "test-token" {
		t.Fatalf("expected API token to be parsed, got %q", p.Provider.APIKey)
	}
}

func TestUnmarshalCaddyfileRejectsZoneSubdirective(t *testing.T) {
	input := `dnsexit {
		api_token test-token
		zone example.com.
	}`
	p := &Provider{}
	d := caddyfile.NewTestDispenser(input)
	if err := p.UnmarshalCaddyfile(d); err == nil {
		t.Fatal("expected error for unsupported zone subdirective")
	}
}

func TestUnmarshalCaddyfileRejectsHostSubdirective(t *testing.T) {
	input := `dnsexit {
		api_token test-token
		host example.com.
	}`
	p := &Provider{}
	d := caddyfile.NewTestDispenser(input)
	if err := p.UnmarshalCaddyfile(d); err == nil {
		t.Fatal("expected error for unsupported host subdirective")
	}
}
