package accounts

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/wamunyima3/vernacular/apps/api/internal/accounts/dbgen"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/database"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/httpx"
)

// Usage counts calls per API key per UTC day (DEV-5, US-13). Counts gather
// in memory and Flush adds them to key_usage, so a busy key costs one write
// per flush rather than one per call. Counts not yet flushed are lost if the
// process dies; the dashboard is a guide, not a bill.
type Usage struct {
	db     database.Beginner
	now    func() time.Time
	mu     sync.Mutex
	counts map[usageKey]usageCount
}

type usageKey struct {
	key uuid.UUID
	day time.Time // midnight UTC
}

type usageCount struct{ calls, errors, limited int32 }

// NewUsage returns a usage recorder writing to db.
func NewUsage(db database.Beginner) *Usage {
	return &Usage{db: db, now: time.Now, counts: map[usageKey]usageCount{}}
}

// Middleware counts each call made with an API key, by its response status.
// It goes after the accounts middleware, which names the key, and before the
// rate limiter, so refused calls count too.
func (u *Usage) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p, ok := FromContext(r.Context())
		if !ok || !p.ByKey() {
			next.ServeHTTP(w, r)
			return
		}
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		next.ServeHTTP(ww, r)
		u.record(p.KeyID, ww.Status())
	})
}

func (u *Usage) record(key uuid.UUID, status int) {
	k := usageKey{key: key, day: u.now().UTC().Truncate(24 * time.Hour)}
	u.mu.Lock()
	defer u.mu.Unlock()
	c := u.counts[k]
	c.calls++
	switch {
	case status == http.StatusTooManyRequests:
		c.limited++
	case status >= 400:
		c.errors++
	}
	u.counts[k] = c
}

// Flush adds the counts gathered so far to key_usage. On failure they are
// kept for the next flush.
func (u *Usage) Flush(ctx context.Context) error {
	u.mu.Lock()
	counts := u.counts
	u.counts = map[usageKey]usageCount{}
	u.mu.Unlock()
	if len(counts) == 0 {
		return nil
	}
	var p dbgen.AddKeyUsageParams
	for k, c := range counts {
		p.KeyIds = append(p.KeyIds, k.key)
		p.Days = append(p.Days, pgtype.Date{Time: k.day, Valid: true})
		p.Calls, p.Errors, p.Limited = append(p.Calls, c.calls), append(p.Errors, c.errors), append(p.Limited, c.limited)
	}
	err := database.InTx(ctx, u.db, "system:usage", func(tx pgx.Tx) error {
		return dbgen.New(tx).AddKeyUsage(ctx, p)
	})
	if err != nil {
		u.mu.Lock()
		for k, c := range counts {
			old := u.counts[k]
			u.counts[k] = usageCount{old.calls + c.calls, old.errors + c.errors, old.limited + c.limited}
		}
		u.mu.Unlock()
		return fmt.Errorf("flush key usage: %w", err)
	}
	return nil
}

// Run flushes every interval until ctx ends. The caller flushes once more
// after the server has stopped taking calls.
func (u *Usage) Run(ctx context.Context, every time.Duration, log *slog.Logger) {
	t := time.NewTicker(every)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			if err := u.Flush(ctx); err != nil {
				log.Error("key usage", "err", err)
			}
		case <-ctx.Done():
			return
		}
	}
}

// KeyUsage is the signed-in user's daily counts for one of their keys over
// the last days days, newest first. Days without calls are left out.
func (s *Service) KeyUsage(ctx context.Context, keyID uuid.UUID, days int) ([]dbgen.KeyUsageRow, error) {
	me, err := Require(ctx, nil)
	if err != nil {
		return nil, err
	}
	since := time.Now().UTC().AddDate(0, 0, 1-days)
	var rows []dbgen.KeyUsageRow
	err = database.InTx(ctx, s.db, me.ID, func(tx pgx.Tx) error {
		q := dbgen.New(tx)
		owns, err := q.OwnsKey(ctx, dbgen.OwnsKeyParams{ID: keyID, UserID: me.ID})
		if err != nil {
			return err
		}
		if !owns {
			return httpx.ErrNotFound
		}
		rows, err = q.KeyUsage(ctx, dbgen.KeyUsageParams{KeyID: keyID, Since: pgtype.Date{Time: since, Valid: true}})
		return err
	})
	return rows, err
}
