// Package accounts resolves who is calling: sessions from the web app
// (ADR-0007) now, API keys later (P1-API-11). It also holds the role guards
// every other module uses.
package accounts

import (
	"context"
	"slices"

	"github.com/wamunyima3/vernacular/apps/api/internal/gen"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/httpx"
)

// Roles and review levels are the contract's enums (gen.Role, gen.ReviewerLevel),
// so Go, the web app and the database CHECK share one list.

// Approval lets a user review one language, unless it is suspended because
// their agreement fell below the threshold (DQ-7).
type Approval struct {
	Language  string            `json:"language"`
	Level     gen.ReviewerLevel `json:"level"`
	Suspended bool              `json:"suspended"`
}

// Principal is the authenticated caller.
type Principal struct {
	ID        string
	Name      string
	Email     string
	Role      gen.Role
	Approvals []Approval
}

// HasRole reports whether the caller has any of the roles. Admins have every role.
func (p Principal) HasRole(roles ...gen.Role) bool {
	return p.Role == gen.RoleAdmin || slices.Contains(roles, p.Role)
}

// active returns the approvals that are not suspended.
func (p Principal) active() []Approval {
	return slices.DeleteFunc(slices.Clone(p.Approvals), func(a Approval) bool { return a.Suspended })
}

// CanReview reports whether the caller may review items in a language.
func (p Principal) CanReview(language string) bool {
	return slices.ContainsFunc(p.active(), func(a Approval) bool { return a.Language == language })
}

// IsSenior reports whether the caller is an active senior reviewer for a language.
func (p Principal) IsSenior(language string) bool {
	return slices.Contains(p.SeniorLanguages(), language)
}

// SeniorLanguages lists the languages the caller is an active senior reviewer for.
func (p Principal) SeniorLanguages() []string {
	var out []string
	for _, a := range p.active() {
		if a.Level == gen.ReviewerLevelSenior {
			out = append(out, a.Language)
		}
	}
	return out
}

type ctxKey struct{}

// WithPrincipal returns a context carrying the caller.
func WithPrincipal(ctx context.Context, p Principal) context.Context {
	return context.WithValue(ctx, ctxKey{}, p)
}

// FromContext returns the caller, if the request was authenticated.
func FromContext(ctx context.Context) (Principal, bool) {
	p, ok := ctx.Value(ctxKey{}).(Principal)
	return p, ok
}

// RequireAdmin returns the caller if they are an admin.
func RequireAdmin(ctx context.Context) (Principal, error) {
	return Require(ctx, func(p Principal) bool { return p.Role == gen.RoleAdmin })
}

// Require returns the caller if they pass allow, or an error the server turns
// into 401 (not signed in) or 403 (signed in, not allowed). A nil allow only
// requires sign-in.
func Require(ctx context.Context, allow func(Principal) bool) (Principal, error) {
	p, ok := FromContext(ctx)
	if !ok {
		return Principal{}, httpx.ErrUnauthorized
	}
	if allow != nil && !allow(p) {
		return Principal{}, httpx.ErrForbidden
	}
	return p, nil
}
