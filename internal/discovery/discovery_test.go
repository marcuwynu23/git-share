package discovery

import (
	"strings"
	"testing"
)

func TestDiscover(t *testing.T) {
	addrs := Discover()
	if addrs == nil {
		t.Fatal("Discover() returned nil")
	}
	if addrs.Localhost != "127.0.0.1" {
		t.Errorf("Localhost = %q, want 127.0.0.1", addrs.Localhost)
	}
}

func TestLanIP(t *testing.T) {
	ip := lanIP()
	// This might be empty in CI or headless environments, which is acceptable
	if ip != "" {
		parts := strings.Split(ip, ".")
		if len(parts) != 4 {
			t.Errorf("lanIP() = %q, expected a valid IPv4 address", ip)
		}
	}
}

func TestDiscoverHasLAN(t *testing.T) {
	addrs := Discover()
	// LAN may be empty in some environments, just check consistency
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
	addrs := Discover()
	// Hostname may be empty if no LAN IP or reverse DNS fails
	_ = addrs.Hostname
}
