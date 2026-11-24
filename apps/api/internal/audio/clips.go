package audio

import (
	"context"
	"errors"
	"fmt"
	"io"
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"gocloud.dev/gcerrors"

	"github.com/wamunyima3/vernacular/apps/api/internal/audio/dbgen"
	"github.com/wamunyima3/vernacular/apps/api/internal/gen"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/blob"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/database"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/httpx"
)

// MaxClipSize caps a recording: a spoken word is a few seconds.
const MaxClipSize = 1 << 20

// clipTypes pairs what a recording sniffs as, in the order tried, with the
// type it is served with.
var clipTypes = [][2]string{
	{"application/ogg", "audio/ogg"},
	{"video/webm", "audio/webm"}, // browsers' MediaRecorder
	{"video/mp4", "audio/mp4"},   // Safari's MediaRecorder
	{"audio/mpeg", "audio/mpeg"},
	{"audio/wave", "audio/wav"},
}

var sniffedClipTypes = func() []string {
	out := make([]string, len(clipTypes))
	for i, t := range clipTypes {
		out[i] = t[0]
	}
	return out
}()

// Clip is a pronunciation of an entry by one speaker.
type Clip struct {
	ID, EntryID uuid.UUID
	SpeakerID   string
	ContentType string
	Size        int32
	RecordedAt  time.Time
}

func clipFromRow(r dbgen.AudioClip) Clip {
	return Clip{ID: r.ID, EntryID: r.EntryID, SpeakerID: r.SpeakerID, ContentType: r.ContentType, Size: r.Size, RecordedAt: r.RecordedAt}
}

// clipKey is where a clip's file lives; recording the entry again replaces it.
func clipKey(entryID uuid.UUID, speakerID string) string {
	return speakerPrefix(speakerID) + "audio/" + entryID.String()
}

// Record saves the signed-in speaker's pronunciation of an entry. They need
// consent in force that names the entry's language.
func (s *Service) Record(ctx context.Context, entryID uuid.UUID, body io.Reader) (Clip, error) {
	me, err := s.check(ctx, signedIn)
	if err != nil {
		return Clip{}, err
	}
	sp, err := s.speaker(ctx, me.ID)
	if errors.Is(err, httpx.ErrNotFound) {
		return Clip{}, httpx.Forbidden("Recording needs your signed consent form on file; ask an admin")
	}
	if err != nil {
		return Clip{}, err
	}
	if sp.WithdrawnAt.Valid {
		return Clip{}, httpx.Forbidden("You withdrew your consent to be recorded")
	}
	lang, err := s.q.EntryLanguage(ctx, entryID)
	if errors.Is(err, pgx.ErrNoRows) {
		return Clip{}, httpx.ErrNotFound
	}
	if err != nil {
		return Clip{}, err
	}
	if !slices.Contains(sp.Languages, lang) {
		return Clip{}, httpx.Forbidden("Your consent form does not cover " + lang)
	}
	data, t, err := httpx.ReadUpload(body, MaxClipSize, sniffedClipTypes, "Recordings", "an Ogg, WebM, MP4, MP3 or WAV recording")
	if err != nil {
		return Clip{}, err
	}
	contentType := clipTypes[slices.Index(sniffedClipTypes, t)][1]
	if err := s.bucket.WriteAll(ctx, clipKey(entryID, me.ID), data, &blob.WriterOptions{ContentType: contentType}); err != nil {
		return Clip{}, fmt.Errorf("store clip: %w", err)
	}
	var row dbgen.AudioClip
	err = database.InTx(ctx, s.db, me.ID, func(tx pgx.Tx) error {
		row, err = dbgen.New(tx).SaveClip(ctx, dbgen.SaveClipParams{
			EntryID: entryID, SpeakerID: me.ID, ContentType: contentType, Size: int32(len(data)),
		})
		return err
	})
	if err != nil {
		return Clip{}, err
	}
	return clipFromRow(row), nil
}

// Clips lists an entry's clips from speakers whose consent is in force.
func (s *Service) Clips(ctx context.Context, entryID uuid.UUID) ([]Clip, error) {
	if _, err := s.check(ctx, anyone); err != nil {
		return nil, err
	}
	rows, err := s.q.ListClips(ctx, entryID)
	if err != nil {
		return nil, err
	}
	out := make([]Clip, len(rows))
	for i, r := range rows {
		out[i] = clipFromRow(r)
	}
	return out, nil
}

// Open opens a clip's file. The caller closes the reader.
func (s *Service) Open(ctx context.Context, id uuid.UUID) (io.ReadCloser, Clip, error) {
	if _, err := s.check(ctx, anyone); err != nil {
		return nil, Clip{}, err
	}
	c, err := s.clip(ctx, id)
	if err != nil {
		return nil, Clip{}, err
	}
	rd, err := s.bucket.NewReader(ctx, clipKey(c.EntryID, c.SpeakerID), nil)
	if err != nil {
		return nil, Clip{}, fmt.Errorf("open clip: %w", err)
	}
	return rd, c, nil
}

// Delete removes a clip and its file. Its speaker or an admin may.
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	me, err := s.check(ctx, signedIn)
	if err != nil {
		return err
	}
	c, err := s.clip(ctx, id)
	if err != nil {
		return err
	}
	if c.SpeakerID != me.ID && me.Role != gen.RoleAdmin {
		return httpx.ErrForbidden
	}
	if err := s.deleteFile(ctx, clipKey(c.EntryID, c.SpeakerID)); err != nil {
		return err
	}
	return database.InTx(ctx, s.db, me.ID, func(tx pgx.Tx) error {
		return dbgen.New(tx).DeleteClip(ctx, id)
	})
}

func (s *Service) clip(ctx context.Context, id uuid.UUID) (Clip, error) {
	row, err := s.q.GetClip(ctx, id)
	if errors.Is(err, pgx.ErrNoRows) {
		return Clip{}, httpx.ErrNotFound
	}
	if err != nil {
		return Clip{}, err
	}
	return clipFromRow(row), nil
}

// deleteFile deletes a file; one already gone is fine.
func (s *Service) deleteFile(ctx context.Context, key string) error {
	err := s.bucket.Delete(ctx, key)
	if err != nil && gcerrors.Code(err) != gcerrors.NotFound {
		return fmt.Errorf("delete clip file: %w", err)
	}
	return nil
}

// purgeFiles deletes all of a speaker's clip files.
func (s *Service) purgeFiles(ctx context.Context, speakerID string) error {
	it := s.bucket.List(&blob.ListOptions{Prefix: speakerPrefix(speakerID) + "audio/"})
	for {
		obj, err := it.Next(ctx)
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return fmt.Errorf("list clip files: %w", err)
		}
		if err := s.deleteFile(ctx, obj.Key); err != nil {
			return err
		}
	}
}
