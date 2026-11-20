package sources

import (
	"context"

	"github.com/wamunyima3/vernacular/apps/api/internal/gen"
	"github.com/wamunyima3/vernacular/apps/api/internal/platform/httpx"
)

// Handler implements the source register operations of the generated strict server.
type Handler struct {
	reg *Register
}

// NewHandler returns a sources handler.
func NewHandler(reg *Register) *Handler { return &Handler{reg: reg} }

// ListSources implements GET /v1/sources.
func (h *Handler) ListSources(ctx context.Context, _ gen.ListSourcesRequestObject) (gen.ListSourcesResponseObject, error) {
	list, err := h.reg.List(ctx)
	if err != nil {
		return nil, err
	}
	data := make([]gen.Source, len(list))
	for i, s := range list {
		data[i] = toAPI(s)
	}
	return gen.ListSources200JSONResponse{Data: data}, nil
}

// GetSource implements GET /v1/sources/{id}.
func (h *Handler) GetSource(ctx context.Context, req gen.GetSourceRequestObject) (gen.GetSourceResponseObject, error) {
	s, err := h.reg.Get(ctx, req.Id)
	return gen.GetSource200JSONResponse(toAPI(s)), err
}

// CreateSource implements POST /v1/sources.
func (h *Handler) CreateSource(ctx context.Context, req gen.CreateSourceRequestObject) (gen.CreateSourceResponseObject, error) {
	s, err := h.reg.Create(ctx, fromAPI(*req.Body))
	return gen.CreateSource201JSONResponse(toAPI(s)), err
}

// UpdateSource implements PUT /v1/sources/{id}.
func (h *Handler) UpdateSource(ctx context.Context, req gen.UpdateSourceRequestObject) (gen.UpdateSourceResponseObject, error) {
	s, err := h.reg.Update(ctx, req.Id, fromAPI(*req.Body))
	return gen.UpdateSource200JSONResponse(toAPI(s)), err
}

// DeleteSource implements DELETE /v1/sources/{id}.
func (h *Handler) DeleteSource(ctx context.Context, req gen.DeleteSourceRequestObject) (gen.DeleteSourceResponseObject, error) {
	return gen.DeleteSource204Response{}, h.reg.Delete(ctx, req.Id)
}

// UploadPermission implements PUT /v1/sources/{id}/permission.
func (h *Handler) UploadPermission(ctx context.Context, req gen.UploadPermissionRequestObject) (gen.UploadPermissionResponseObject, error) {
	s, err := h.reg.UploadPermission(ctx, req.Id, req.Body)
	return gen.UploadPermission200JSONResponse(toAPI(s)), err
}

// DownloadPermission implements GET /v1/sources/{id}/permission.
func (h *Handler) DownloadPermission(ctx context.Context, req gen.DownloadPermissionRequestObject) (gen.DownloadPermissionResponseObject, error) {
	rd, size, err := h.reg.Permission(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	// The generated writer copies Body to the response and closes it.
	return gen.DownloadPermission200ApplicationoctetStreamResponse{Body: rd, ContentLength: size}, nil
}

func fromAPI(in gen.SourceInput) Input {
	uses := make([]string, len(in.AllowedUses))
	for i, u := range in.AllowedUses {
		uses[i] = string(u)
	}
	return Input{
		Kind: string(in.Kind), Name: in.Name, RightsHolder: in.RightsHolder, Licence: in.Licence,
		URL: deref(in.Url), AllowedUses: uses, Notes: deref(in.Notes),
	}
}

func toAPI(s Source) gen.Source {
	uses := make([]gen.AllowedUse, len(s.AllowedUses))
	for i, u := range s.AllowedUses {
		uses[i] = gen.AllowedUse(u)
	}
	return gen.Source{
		Id: s.ID, Kind: gen.SourceKind(s.Kind), Name: s.Name, RightsHolder: s.RightsHolder, Licence: s.Licence,
		Url: httpx.Optional(s.URL), AllowedUses: uses, Notes: httpx.Optional(s.Notes),
		HasPermissionDocument: s.HasPermissionDocument, CreatedAt: s.CreatedAt, UpdatedAt: s.UpdatedAt,
	}
}

func deref(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}
