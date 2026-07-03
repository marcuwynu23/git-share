package discovery

import (
	"net"
	"strings"
)

type Addresses struct {
	Localhost string
	LAN       string
	Hostname  string
}

func Discover() *Addresses {
	return &Addresses{
		Localhost: "127.0.0.1",
		LAN:       lanIP(),
		Hostname:  hostname(),
	}
}

func lanIP() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		return ""
	}

	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			ipnet, ok := addr.(*net.IPNet)
			if !ok {
				continue
			}
			if ipnet.IP.IsLoopback() {
				continue
			}
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}

	return ""
}

func hostname() string {
	h, err := net.LookupAddr(lanIP())
	if err == nil && len(h) > 0 {
		return strings.TrimSuffix(h[0], ".")
	}
	return ""
}
