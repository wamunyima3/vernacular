package languages

import (
	"context"

	"github.com/wamunyima3/vernacular/apps/api/internal/gen"
)

// Lister is what the handler needs from the service.
type Lister interface {
	List(ctx context.Context) ([]Language, error)
}

// Handler implements the language operations of the generated strict server.
type Handler struct {
	svc Lister
}

// NewHandler returns a languages handler.
func NewHandler(svc Lister) *Handler {
	return &Handler{svc: svc}
}

// ListLanguages implements GET /v1/languages.
func (h *Handler) ListLanguages(ctx context.Context, _ gen.ListLanguagesRequestObject) (gen.ListLanguagesResponseObject, error) {
	langs, err := h.svc.List(ctx)
	if err != nil {
		return nil, err
	}
	data := make([]gen.Language, len(langs))
	for i, l := range langs {
		data[i] = gen.Language{
			Code:    l.Code,
			Name:    l.Name,
			OwnName: l.OwnName,
			Wave:    l.Wave,
			Status:  gen.LanguageStatus(l.Status),
			Entries: gen.EntryCounts{Verified: l.Entries.Verified, Total: l.Entries.Total},
		}
	}
	return gen.ListLanguages200JSONResponse{Data: data}, nil
}
