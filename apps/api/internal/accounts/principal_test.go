package accounts_test

import (
	"errors"
	"testing"

	"github.com/wamunyima3/vernacular/apps/api/internal/accounts"
	"github.com/wamunyima3/vernacular/apps/api/internal/gen"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/httpx"
)

var (
	developer = accounts.Principal{ID: "d", Role: gen.Developer}
	admin     = accounts.Principal{ID: "a", Role: gen.Admin}
	reviewer  = accounts.Principal{ID: "r", Role: gen.Contributor, Approvals: []accounts.Approval{
		{Language: "bem", Level: gen.Reviewer},
		{Language: "nya", Level: gen.Senior},
	}}
)

func TestRequire(t *testing.T) {
	contributors := func(p accounts.Principal) bool { return p.HasRole(gen.Contributor) }
	reviewsBemba := func(p accounts.Principal) bool { return p.CanReview("bem") }
	seniorBemba := func(p accounts.Principal) bool { return p.IsSenior("bem") }

	tests := []struct {
		name    string
		caller  *accounts.Principal
		allow   func(accounts.Principal) bool
		wantErr error
	}{
		{"anonymous, sign-in only", nil, nil, httpx.ErrUnauthorized},
		{"anonymous, role guard", nil, contributors, httpx.ErrUnauthorized},
		{"signed in, sign-in only", &developer, nil, nil},
		{"developer is not a contributor", &developer, contributors, httpx.ErrForbidden},
		{"contributor role", &reviewer, contributors, nil},
		{"admin has every role", &admin, contributors, nil},
		{"approved language", &reviewer, reviewsBemba, nil},
		{"admin is not a reviewer by role alone", &admin, reviewsBemba, httpx.ErrForbidden},
		{"reviewer is not senior", &reviewer, seniorBemba, httpx.ErrForbidden},
		{"senior in nya only", &reviewer, func(p accounts.Principal) bool { return p.IsSenior("nya") }, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()
			if tt.caller != nil {
				ctx = accounts.WithPrincipal(ctx, *tt.caller)
			}
			p, err := accounts.Require(ctx, tt.allow)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err = %v, want %v", err, tt.wantErr)
			}
			if err == nil && p.ID != tt.caller.ID {
				t.Errorf("principal = %q, want %q", p.ID, tt.caller.ID)
			}
		})
	}
}
