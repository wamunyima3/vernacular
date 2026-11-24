package httpx_test

import (
	"net/http"
	"net/http/httptest"
	"net/netip"
	"testing"

	"github.com/wamunyima3/vernacular/apps/api/internal/platform/httpx"
)

func TestClientIP(t *testing.T) {
	proxy := httpx.ClientIP{Trusted: []netip.Prefix{netip.MustParsePrefix("10.0.0.0/8")}}
	tests := []struct {
		name       string
		ip         httpx.ClientIP
		remote, ff string
		want       string
	}{
		{"no proxies: the header is ignored", httpx.ClientIP{}, "203.0.113.5:4000", "198.51.100.1", "203.0.113.5"},
		{"from a trusted proxy", proxy, "10.0.0.2:4000", "198.51.100.1", "198.51.100.1"},
		{"spoofed left of the proxy's entry", proxy, "10.0.0.2:4000", "1.2.3.4, 198.51.100.1", "198.51.100.1"},
		{"through two proxies", proxy, "10.0.0.2:4000", "198.51.100.1, 10.0.0.9", "198.51.100.1"},
		{"untrusted sender", proxy, "203.0.113.5:4000", "198.51.100.1", "203.0.113.5"},
		{"malformed hop", proxy, "10.0.0.2:4000", "nonsense", "10.0.0.2"},
		{"IPv6", httpx.ClientIP{}, "[2001:db8::1]:4000", "", "2001:db8::1"},
	}
	for _, tt := range tests {
		r := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)
		r.RemoteAddr = tt.remote
		if tt.ff != "" {
			r.Header.Set("X-Forwarded-For", tt.ff)
		}
		if got := tt.ip.Of(r).String(); got != tt.want {
			t.Errorf("%s: %s, want %s", tt.name, got, tt.want)
		}
	}
}
