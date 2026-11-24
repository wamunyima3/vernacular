package httpx

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"strings"
)

// ETag gives successful GET responses a weak ETag over their body, and
// answers 304 Not Modified, with no body, when If-None-Match names it
// (NFR-3). Weak, because compression changes the bytes on the wire but not
// the meaning. cacheControl sets Cache-Control unless the handler already
// did. The body names each item's release in trust.release, so the tag
// changes with the release a read is pinned to (PRV-4, US-11).
func ETag(cacheControl func(*http.Request) string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet && r.Method != http.MethodHead {
				next.ServeHTTP(w, r)
				return
			}
			buf := &buffered{ResponseWriter: w, status: http.StatusOK}
			next.ServeHTTP(buf, r)
			if buf.status != http.StatusOK {
				buf.flush()
				return
			}
			sum := sha256.Sum256(buf.body.Bytes())
			tag := `W/"` + base64.RawURLEncoding.EncodeToString(sum[:16]) + `"`
			h := w.Header()
			h.Set("ETag", tag)
			if h.Get("Cache-Control") == "" {
				h.Set("Cache-Control", cacheControl(r))
			}
			if matches(r.Header.Get("If-None-Match"), tag) {
				h.Del("Content-Length")
				h.Del("Content-Type")
				w.WriteHeader(http.StatusNotModified)
				return
			}
			buf.flush()
		})
	}
}

// matches compares If-None-Match to a tag weakly, as RFC 9110 asks for GET.
func matches(ifNoneMatch, tag string) bool {
	opaque := strings.TrimPrefix(tag, "W/")
	for t := range strings.SplitSeq(ifNoneMatch, ",") {
		t = strings.TrimSpace(t)
		if t == "*" || strings.TrimPrefix(t, "W/") == opaque {
			return true
		}
	}
	return false
}

// buffered holds a response until its ETag is known.
type buffered struct {
	http.ResponseWriter
	status int
	body   bytes.Buffer
}

func (b *buffered) WriteHeader(status int) { b.status = status }

func (b *buffered) Write(p []byte) (int, error) { return b.body.Write(p) }

func (b *buffered) flush() {
	b.ResponseWriter.WriteHeader(b.status)
	_, _ = b.ResponseWriter.Write(b.body.Bytes()) // the client went away; nothing to do
}
