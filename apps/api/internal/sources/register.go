package sources

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/wamunyima3/vernacular/apps/api/internal/accounts"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/blob"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/database"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/httpx"
	"github.com/wamunyima3/vernacular/apps/api/internal/sources/dbgen"
)

// MaxPermissionSize caps a permission document upload.
const MaxPermissionSize = 10 << 20

// Input is what an admin sends to create or update a source.
type Input struct {
	Kind         string
	Name         string
	RightsHolder string
	Licence      string
	URL          string
	AllowedUses  []string
	Notes        string
}

// Source is a register row.
type Source struct {
	Input
	ID                    uuid.UUID
	HasPermissionDocument bool
	CreatedAt, UpdatedAt  time.Time
}

// Register keeps the source register (DQ-1). Reads are public; writes and
// permission documents need an admin.
type Register struct {
	db     database.Beginner
	q      *dbgen.Queries
	bucket *blob.Bucket
}

// DB is what the register needs from the pool.
type DB interface {
	database.Beginner
	dbgen.DBTX
}

// NewRegister returns a register over the pool and the object store.
func NewRegister(db DB, bucket *blob.Bucket) *Register {
	return &Register{db: db, q: dbgen.New(db), bucket: bucket}
}

// List returns every source by name.
func (r *Register) List(ctx context.Context) ([]Source, error) {
	rows, err := r.q.ListSources(ctx)
	if err != nil {
		return nil, fmt.Errorf("list sources: %w", err)
	}
	out := make([]Source, len(rows))
	for i, row := range rows {
		out[i] = fromRow(row)
	}
	return out, nil
}

// Get returns one source.
func (r *Register) Get(ctx context.Context, id uuid.UUID) (Source, error) {
	row, err := r.q.GetSource(ctx, id)
	return r.result(row, err)
}

// Create adds a source.
func (r *Register) Create(ctx context.Context, in Input) (Source, error) {
	return r.write(ctx, func(q *dbgen.Queries) (dbgen.Source, error) {
		return q.CreateSource(ctx, dbgen.CreateSourceParams{
			Kind: dbgen.SourceKind(in.Kind), Name: in.Name, RightsHolder: in.RightsHolder, Licence: in.Licence,
			Url: database.Text(in.URL), AllowedUses: in.AllowedUses, Notes: in.Notes,
		})
	})
}

// Update replaces a source's details.
func (r *Register) Update(ctx context.Context, id uuid.UUID, in Input) (Source, error) {
	return r.write(ctx, func(q *dbgen.Queries) (dbgen.Source, error) {
		return q.UpdateSource(ctx, dbgen.UpdateSourceParams{
			ID: id, Kind: dbgen.SourceKind(in.Kind), Name: in.Name, RightsHolder: in.RightsHolder, Licence: in.Licence,
			Url: database.Text(in.URL), AllowedUses: in.AllowedUses, Notes: in.Notes,
		})
	})
}

// Delete removes a source that no item cites.
func (r *Register) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.write(ctx, func(q *dbgen.Queries) (dbgen.Source, error) {
		n, err := q.DeleteSource(ctx, id)
		if err == nil && n == 0 {
			err = pgx.ErrNoRows
		}
		return dbgen.Source{}, err
	})
	return err
}

// UploadPermission stores the signed permission document for a source.
func (r *Register) UploadPermission(ctx context.Context, id uuid.UUID, body io.Reader) (Source, error) {
	actor, err := accounts.RequireAdmin(ctx)
	if err != nil {
		return Source{}, err
	}
	data, contentType, err := httpx.ReadUpload(body, MaxPermissionSize, httpx.DocumentTypes, "Permission documents", httpx.DocumentAccept)
	if err != nil {
		return Source{}, err
	}
	if _, err := r.Get(ctx, id); err != nil {
		return Source{}, err
	}
	key := "sources/" + id.String() + "/permission"
	if err := r.bucket.WriteAll(ctx, key, data, &blob.WriterOptions{ContentType: contentType}); err != nil {
		return Source{}, fmt.Errorf("store permission document: %w", err)
	}
	var row dbgen.Source
	err = database.InTx(ctx, r.db, actor.ID, func(tx pgx.Tx) error {
		row, err = dbgen.New(tx).SetPermissionKey(ctx, dbgen.SetPermissionKeyParams{ID: id, PermissionKey: database.Text(key)})
		return err
	})
	return r.result(row, err)
}

// Permission opens the permission document. The caller closes the reader.
func (r *Register) Permission(ctx context.Context, id uuid.UUID) (io.ReadCloser, int64, error) {
	if _, err := accounts.RequireAdmin(ctx); err != nil {
		return nil, 0, err
	}
	row, err := r.q.GetSource(ctx, id)
	if _, err := r.result(row, err); err != nil {
		return nil, 0, err
	}
	if !row.PermissionKey.Valid {
		return nil, 0, httpx.ErrNotFound
	}
	rd, err := r.bucket.NewReader(ctx, row.PermissionKey.String, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("open permission document: %w", err)
	}
	return rd, rd.Size(), nil
}

// write runs an admin-only change in an audited transaction and maps
// database errors to API problems.
func (r *Register) write(ctx context.Context, fn func(*dbgen.Queries) (dbgen.Source, error)) (Source, error) {
	actor, err := accounts.RequireAdmin(ctx)
	if err != nil {
		return Source{}, err
	}
	var row dbgen.Source
	err = database.InTx(ctx, r.db, actor.ID, func(tx pgx.Tx) error {
		row, err = fn(dbgen.New(tx))
		return err
	})
	return r.result(row, err)
}

func (r *Register) result(row dbgen.Source, err error) (Source, error) {
	switch {
	case err == nil:
		return fromRow(row), nil
	case errors.Is(err, pgx.ErrNoRows):
		return Source{}, httpx.ErrNotFound
	case database.Violates(err, database.UniqueViolation):
		return Source{}, httpx.Conflict("A source with this name already exists")
	case database.Violates(err, database.ForeignKeyViolation):
		return Source{}, httpx.Conflict("Items cite this source, so it cannot be deleted")
	default:
		return Source{}, fmt.Errorf("source: %w", err)
	}
}

func fromRow(row dbgen.Source) Source {
	return Source{
		Input: Input{
			Kind: string(row.Kind), Name: row.Name, RightsHolder: row.RightsHolder, Licence: row.Licence,
			URL: row.Url.String, AllowedUses: row.AllowedUses, Notes: row.Notes,
		},
		ID: row.ID, HasPermissionDocument: row.PermissionKey.Valid, CreatedAt: row.CreatedAt, UpdatedAt: row.UpdatedAt,
	}
}
