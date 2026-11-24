package httpx

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

// DocumentTypes are the scanned or signed documents accepted as uploads
// (permission documents, consent forms), and DocumentAccept says so in words.
var DocumentTypes = []string{"application/pdf", "image/png", "image/jpeg"}

// DocumentAccept names DocumentTypes for messages.
const DocumentAccept = "a PDF, PNG or JPEG file"

// ReadUpload reads an uploaded file of at most limit bytes whose content, as
// sniffed from its bytes rather than what the client claims, is one of types.
// It returns the bytes and their type, or a 413 or 400 for the caller: what
// names the files ("Permission documents") and accept the types in words
// ("a PDF, PNG or JPEG file").
func ReadUpload(body io.Reader, limit int64, types []string, what, accept string) ([]byte, string, error) {
	data, err := io.ReadAll(io.LimitReader(body, limit+1))
	if err != nil {
		return nil, "", fmt.Errorf("read upload: %w", err)
	}
	if int64(len(data)) > limit {
		return nil, "", TooLarge(fmt.Sprintf("%s are limited to %s", what, sizeText(limit)))
	}
	ct := http.DetectContentType(data)
	for _, t := range types {
		if strings.HasPrefix(ct, t) {
			return data, t, nil
		}
	}
	return nil, "", BadRequest("Upload " + accept)
}

func sizeText(n int64) string {
	if n >= 1<<20 && n%(1<<20) == 0 {
		return fmt.Sprintf("%d MB", n>>20)
	}
	return fmt.Sprintf("%d KB", n>>10)
}
