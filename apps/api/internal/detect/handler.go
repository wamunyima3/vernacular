package detect

import (
	"context"

	"github.com/wamunyima3/vernacular/apps/api/internal/gen"
)

// maxResults is how many languages a detection returns.
const maxResults = 5

// Handler implements POST /v1/detect.
type Handler struct {
	model *Model
}

// NewHandler returns a detect handler serving the model.
func NewHandler(m *Model) *Handler {
	return &Handler{model: m}
}

// DetectLanguage implements POST /v1/detect.
func (h *Handler) DetectLanguage(_ context.Context, req gen.DetectLanguageRequestObject) (gen.DetectLanguageResponseObject, error) {
	scores := h.model.Detect(req.Body.Text)
	data := make([]gen.LanguageScore, min(len(scores), maxResults))
	for i := range data {
		data[i] = gen.LanguageScore{Language: scores[i].Language, Confidence: scores[i].Confidence}
	}
	return gen.DetectLanguage200JSONResponse{Data: data}, nil
}
