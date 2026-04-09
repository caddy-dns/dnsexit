package dnsexit

import (
	"testing"

	"github.com/libdns/dnsexit"
)

func TestProviderInstantiation(t *testing.T) {
	p := &Provider{Provider: &dnsexit.Provider{APIKey: "test"}}
	if p.Provider.APIKey != "test" {
		t.Errorf("Expected APIKey to be 'test', got %s", p.Provider.APIKey)
	}
}
