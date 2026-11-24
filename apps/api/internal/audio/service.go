// Package audio holds pronunciation clips and the consent records of the
// speakers who recorded them (LEX-8, NFR-6, US-33). Voice is personal data
// under the Data Protection Act, 2021: nothing here runs until the server
// sets AUDIO_ENABLED, which waits on legal advice (tasks.md P2-08d).
package audio

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/wamunyima3/vernacular/apps/api/internal/accounts"
	"github.com/wamunyima3/vernacular/apps/api/internal/audio/dbgen"
	"github.com/wamunyima3/vernacular/apps/api/internal/gen"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/blob"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/database"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/httpx"
)

// MaxConsentSize caps a scanned consent form.
const MaxConsentSize = 10 << 20

// ErrDisabled answers every audio call while AUDIO_ENABLED is off.
var ErrDisabled = &httpx.Problem{Status: http.StatusNotFound, Code: gen.ErrorCodeNotFound, Message: "Audio is not switched on for this server"}

// Speaker is a speaker's consent record.
type Speaker struct {
	UserID      string
	Languages   []string
	ConsentedAt time.Time
	WithdrawnAt *time.Time
}

func fromRow(r dbgen.Speaker) Speaker {
	return Speaker{UserID: r.UserID, Languages: r.Languages, ConsentedAt: r.ConsentedAt, WithdrawnAt: httpx.OptionalTime(r.WithdrawnAt)}
}

// DB reads from the pool and writes in audited transactions.
type DB interface {
	database.Beginner
	dbgen.DBTX
}

// Service records consent and, from P2-08b, clips.
type Service struct {
	db      DB
	q       *dbgen.Queries // reads
	bucket  *blob.Bucket
	enabled bool
}

// NewService returns the audio service. With enabled false every method
// returns ErrDisabled.
func NewService(db DB, bucket *blob.Bucket, enabled bool) *Service {
	return &Service{db: db, q: dbgen.New(db), bucket: bucket, enabled: enabled}
}

// check is the gate in front of every method: audio is on, and require passes.
func (s *Service) check(ctx context.Context, require func(context.Context) (accounts.Principal, error)) (accounts.Principal, error) {
	if !s.enabled {
		return accounts.Principal{}, ErrDisabled
	}
	return require(ctx)
}

func signedIn(ctx context.Context) (accounts.Principal, error) { return accounts.Require(ctx, nil) }

func anyone(context.Context) (accounts.Principal, error) { return accounts.Principal{}, nil }

// A speaker's files sit under one prefix: their signed form (a new form
// replaces it) and their clips, so withdrawal can clear the clips at once.
func speakerPrefix(userID string) string { return "speakers/" + userID + "/" }

func consentKey(userID string) string { return speakerPrefix(userID) + "consent" }

// RecordConsent stores a speaker's signed consent form and records their
// consent to record in languages. Admins only.
func (s *Service) RecordConsent(ctx context.Context, userID string, languages []string, form io.Reader) (Speaker, error) {
	admin, err := s.check(ctx, accounts.RequireAdmin)
	if err != nil {
		return Speaker{}, err
	}
	data, contentType, err := httpx.ReadUpload(form, MaxConsentSize, httpx.DocumentTypes, "Consent forms", httpx.DocumentAccept)
	if err != nil {
		return Speaker{}, err
	}
	switch exists, err := s.q.UserExists(ctx, userID); {
	case err != nil:
		return Speaker{}, err
	case !exists:
		return Speaker{}, httpx.ErrNotFound
	}
	key := consentKey(userID)
	if err := s.bucket.WriteAll(ctx, key, data, &blob.WriterOptions{ContentType: contentType}); err != nil {
		return Speaker{}, fmt.Errorf("store consent form: %w", err)
	}
	var row dbgen.Speaker
	err = database.InTx(ctx, s.db, admin.ID, func(tx pgx.Tx) error {
		row, err = dbgen.New(tx).RecordConsent(ctx, dbgen.RecordConsentParams{UserID: userID, ConsentKey: key, RecordedBy: admin.ID, Languages: languages})
		return err
	})
	if err != nil {
		return Speaker{}, err
	}
	return fromRow(row), nil
}

// ConsentForm opens a speaker's consent form. Admins only. The caller closes
// the reader.
func (s *Service) ConsentForm(ctx context.Context, userID string) (io.ReadCloser, int64, error) {
	if _, err := s.check(ctx, accounts.RequireAdmin); err != nil {
		return nil, 0, err
	}
	row, err := s.speaker(ctx, userID)
	if err != nil {
		return nil, 0, err
	}
	rd, err := s.bucket.NewReader(ctx, row.ConsentKey, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("open consent form: %w", err)
	}
	return rd, rd.Size(), nil
}

// Me is the signed-in user's consent record.
func (s *Service) Me(ctx context.Context) (Speaker, error) {
	me, err := s.check(ctx, signedIn)
	if err != nil {
		return Speaker{}, err
	}
	row, err := s.speaker(ctx, me.ID)
	if err != nil {
		return Speaker{}, err
	}
	return fromRow(row), nil
}

func (s *Service) speaker(ctx context.Context, userID string) (dbgen.Speaker, error) {
	row, err := s.q.GetSpeaker(ctx, userID)
	if errors.Is(err, pgx.ErrNoRows) {
		return row, httpx.ErrNotFound
	}
	return row, err
}

// Withdraw ends the signed-in user's consent and deletes their clips and
// files (NFR-6).
func (s *Service) Withdraw(ctx context.Context) (Speaker, error) {
	me, err := s.check(ctx, signedIn)
	if err != nil {
		return Speaker{}, err
	}
	var row dbgen.Speaker
	err = database.InTx(ctx, s.db, me.ID, func(tx pgx.Tx) error {
		q := dbgen.New(tx)
		if err := q.DeleteSpeakerClips(ctx, me.ID); err != nil {
			return err
		}
		row, err = q.WithdrawConsent(ctx, me.ID)
		if errors.Is(err, pgx.ErrNoRows) {
			if _, err := q.GetSpeaker(ctx, me.ID); errors.Is(err, pgx.ErrNoRows) {
				return httpx.ErrNotFound
			}
			return errWithdrawn
		}
		return err
	})
	if err != nil && !errors.Is(err, errWithdrawn) {
		return Speaker{}, err
	}
	// The rows are gone and the clips hidden; now the files. If this fails,
	// withdrawing again retries it.
	if perr := s.purgeFiles(ctx, me.ID); perr != nil {
		return Speaker{}, perr
	}
	if err != nil {
		return Speaker{}, err
	}
	return fromRow(row), nil
}

var errWithdrawn = httpx.Conflict("You have already withdrawn your consent")
