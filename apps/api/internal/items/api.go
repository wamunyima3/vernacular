package items

import (
	"github.com/wamunyima3/vernacular/apps/api/internal/gen"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/httpx"
)

// API returns the trust fields in the contract's shape, for any handler that
// returns an item (PRV-1).
func (t Trust) API() gen.Trust {
	sources := make([]gen.SourceRef, len(t.Sources))
	for i, s := range t.Sources {
		sources[i] = gen.SourceRef{Name: s.Name, Licence: s.Licence, Url: httpx.Optional(s.URL), Reference: httpx.Optional(s.Reference), Model: httpx.Optional(s.Model)}
	}
	return gen.Trust{
		Status:  gen.TrustStatus(t.Status),
		Sources: sources,
		Reviews: gen.ReviewSummary{Accepts: t.Accepts, Corrections: t.Corrections},
	}
}

// Minimum is the least-checked status a read returns: the caller's
// min_status, or verified when none is given (ADR-0004, TRN-4).
func Minimum(minStatus *gen.MinStatus) string {
	if minStatus == nil {
		return string(gen.MinStatusVerified)
	}
	return string(*minStatus)
}
