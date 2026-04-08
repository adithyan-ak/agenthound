package common

import (
	"net"
	"net/url"
	"strings"

	"github.com/adithyan-ak/agenthound/internal/model"
)

type HostInfo struct {
	Hostname  string
	IP        string
	IsLocal   bool
	IsPrivate bool
	IsPublic  bool
}

func ClassifyHost(hostOrURL string) HostInfo {
	host := hostOrURL

	if strings.Contains(host, "://") {
		if u, err := url.Parse(host); err == nil && u.Hostname() != "" {
			host = u.Hostname()
		}
	}

	host = strings.TrimSpace(host)
	if h, _, err := net.SplitHostPort(host); err == nil {
		host = h
	}

	info := HostInfo{}

	if ip := net.ParseIP(host); ip != nil {
		info.IP = ip.String()
	} else {
		info.Hostname = host
	}

	if isLocal(host) {
		info.IsLocal = true
		return info
	}
	if isPrivate(host) {
		info.IsPrivate = true
		return info
	}

	info.IsPublic = true
	return info
}

func HostNodeID(hostname string) string {
	return model.ComputeNodeID("Host", hostname)
}

func isLocal(host string) bool {
	lower := strings.ToLower(host)
	if lower == "localhost" {
		return true
	}

	ip := net.ParseIP(host)
	if ip == nil {
		return false
	}

	if ip.IsLoopback() {
		return true
	}

	return false
}

func isPrivate(host string) bool {
	if host == "169.254.169.254" {
		return true
	}

	ip := net.ParseIP(host)
	if ip == nil {
		return false
	}

	privateRanges := []struct {
		network *net.IPNet
	}{
		{mustParseCIDR("10.0.0.0/8")},
		{mustParseCIDR("172.16.0.0/12")},
		{mustParseCIDR("192.168.0.0/16")},
	}

	for _, r := range privateRanges {
		if r.network.Contains(ip) {
			return true
		}
	}

	return false
}

func mustParseCIDR(cidr string) *net.IPNet {
	_, network, err := net.ParseCIDR(cidr)
	if err != nil {
		panic("invalid CIDR: " + cidr)
	}
	return network
}
