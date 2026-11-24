package server_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/wamunyima3/vernacular/apps/api/internal/accounts"
	"github.com/wamunyima3/vernacular/apps/api/internal/gen"
	"github.com/wamunyima3/vernacular/apps/api/internal/phrases"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/database"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/database/dbtest"
	"github.com/wamunyima3/vernacular/apps/api/internal/releases"
	"github.com/wamunyima3/vernacular/apps/api/internal/server"
)

// TestPinnedRead: `release` serves a read from the release's frozen rows,
// names it in trust.release and caches it for a day (ADR-0012, US-11).
func TestPinnedRead(t *testing.T) {
	tx := dbtest.Tx(t)
	svc := releases.NewService(tx)
	if _, err := tx.Exec(t.Context(), `UPDATE items SET status = 'reviewed' WHERE status = 'verified';
		WITH i AS (INSERT INTO items (kind, language_code, status) VALUES ('phrase', 'bem', 'verified') RETURNING id)
		INSERT INTO phrases (id, language_code, pack, text, source_en) SELECT id, 'bem', 'zz-pin', 'Mwapoleni', 'Hello' FROM i`); err != nil {
		t.Fatal(err)
	}
	built, err := svc.Build(accounts.WithPrincipal(t.Context(), accounts.Principal{ID: "a", Role: gen.RoleAdmin}), "")
	if err != nil {
		t.Fatal(err)
	}
	// Changed after the release: the pinned read must not see it.
	if _, err := tx.Exec(t.Context(), `UPDATE phrases SET text = 'Muli shani' WHERE pack = 'zz-pin'`); err != nil {
		t.Fatal(err)
	}

	h := newServer(server.Deps{
		Accounts: fakeAuth{}, Releases: svc, Phrases: phrases.NewStore(database.Routed{Default: tx}),
	})
	for _, c := range []struct{ release, text, cache string }{
		{"", "Muli shani", "private, max-age=300"},
		{built.Name, "Mwapoleni", "private, max-age=86400"},
	} {
		url := "/v1/phrases?pack=zz-pin&lang=bem"
		if c.release != "" {
			url += "&release=" + c.release
		}
		req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, url, nil)
		req.Header.Set("X-API-Key", "vk_good")
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		var body gen.PhrasePack
		if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
			t.Fatalf("%s: %v", rec.Body, err)
		}
		if rec.Code != http.StatusOK || len(body.Data) != 1 || body.Data[0].Text != c.text ||
			body.Data[0].Trust.Release == nil || *body.Data[0].Trust.Release != built.Name ||
			rec.Header().Get("Cache-Control") != c.cache {
			t.Errorf("release %q: %d %s %+v, want %q in %s cached %q", c.release, rec.Code, rec.Header().Get("Cache-Control"), body, c.text, built.Name, c.cache)
		}
	}
}
