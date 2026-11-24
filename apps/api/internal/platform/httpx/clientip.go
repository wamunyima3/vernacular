package httpx

import (
	"net"
	"net/http"
	"net/netip"
	"strings"
)

// ClientIP finds the address a request came from. X-Forwarded-For is
// believed only from trusted proxies: walking it right to left, the first
// address that is not a trusted proxy is the client. Without trusted proxies
// it is the connection's address, so a caller cannot pick its own IP by
// sending the header (unlike chi's RealIP).
type ClientIP struct {
	Trusted []netip.Prefix
}

// Of returns the client address of r.
func (c ClientIP) Of(r *http.Request) netip.Addr {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		host = r.RemoteAddr
	}
	addr, err := netip.ParseAddr(host)
	if err != nil {
		return netip.Addr{}
	}
	addr = addr.Unmap()
	hops := strings.Split(strings.Join(r.Header.Values("X-Forwarded-For"), ","), ",")
	for i := len(hops) - 1; i >= 0 && c.trusts(addr); i-- {
		hop, err := netip.ParseAddr(strings.TrimSpace(hops[i]))
		if err != nil {
			break // a malformed hop: keep the last address we could vouch for
		}
		addr = hop.Unmap()
	}
	return addr
}

func (c ClientIP) trusts(addr netip.Addr) bool {
	for _, p := range c.Trusted {
		if p.Contains(addr) {
			return true
		}
	}
	return false
}
