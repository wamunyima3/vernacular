package releases

import (
	"bytes"
	"context"
	"fmt"

	"github.com/wamunyima3/vernacular/apps/api/internal/gen"
)

// Handler implements the release operations of the generated strict server.
type Handler struct {
	svc *Service
}

// NewHandler returns a releases handler.
func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// ListReleases implements GET /v1/releases.
func (h *Handler) ListReleases(ctx context.Context, _ gen.ListReleasesRequestObject) (gen.ListReleasesResponseObject, error) {
	rs, err := h.svc.List(ctx)
	if err != nil {
		return nil, err
	}
	data := make([]gen.Release, len(rs))
	for i, r := range rs {
		data[i] = toAPI(r)
	}
	return gen.ListReleases200JSONResponse{Data: data}, nil
}

// CreateRelease implements POST /v1/releases.
func (h *Handler) CreateRelease(ctx context.Context, req gen.CreateReleaseRequestObject) (gen.CreateReleaseResponseObject, error) {
	r, err := h.svc.Build(ctx, req.Body.Notes)
	if err != nil {
		return nil, err
	}
	return gen.CreateRelease201JSONResponse(toAPI(r)), nil
}

func toAPI(r Release) gen.Release {
	langs := make(map[string]gen.KindCounts, len(r.Languages))
	for code, k := range r.Languages {
		langs[code] = gen.KindCounts{Entries: k["entry"], Senses: k["sense"], Examples: k["example"], Phrases: k["phrase"]}
	}
	return gen.Release{
		Name: r.Name, PublishedAt: r.PublishedAt, Notes: r.Notes, Languages: langs,
		Changes: gen.ReleaseChanges{Added: r.Changes.Added, Changed: r.Changes.Changed, Removed: r.Changes.Removed},
	}
}

// DownloadRelease implements GET /v1/releases/{name}/download.
func (h *Handler) DownloadRelease(ctx context.Context, req gen.DownloadReleaseRequestObject) (gen.DownloadReleaseResponseObject, error) {
	name, body, err := h.svc.Export(ctx, req.Name)
	if err != nil {
		return nil, err
	}
	disposition := fmt.Sprintf(`attachment; filename="vernacular-%s.zip"`, name)
	return gen.DownloadRelease200ApplicationzipResponse{
		Body: bytes.NewReader(body), ContentLength: int64(len(body)),
		Headers: gen.DownloadRelease200ResponseHeaders{ContentDisposition: &disposition},
	}, nil
}
