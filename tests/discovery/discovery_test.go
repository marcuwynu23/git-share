package discovery_test

import (
	"strings"
	"testing"

	"github.com/marcuwynu23/git-share/internal/discovery"
)

func TestDiscover(t *testing.T) {
	addrs := discovery.Discover()
	if addrs == nil {
		t.Fatal("Discover() returned nil")
	}
	if addrs.Localhost != "127.0.0.1" {
		t.Errorf("Localhost = %q, want 127.0.0.1", addrs.Localhost)
	}
}

func TestDiscoverHasLAN(t *testing.T) {
	addrs := discovery.Discover()
	if addrs.LAN != "" {
		if addrs.LAN == addrs.Localhost {
			t.Error("LAN should not equal Localhost")
		}
		parts := strings.Split(addrs.LAN, ".")
		if len(parts) != 4 {
			t.Errorf("LAN = %q, expected valid IPv4", addrs.LAN)
		}
	}
}

func TestDiscoverHostname(t *testing.T) {
	addrs := discovery.Discover()
	_ = addrs.Hostname
}
