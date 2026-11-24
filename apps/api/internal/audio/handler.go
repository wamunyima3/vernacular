package audio

import (
	"context"

	"github.com/wamunyima3/vernacular/apps/api/internal/gen"
)

// Handler implements the audio operations of the generated strict server.
type Handler struct {
	svc *Service
}

// NewHandler returns an audio handler.
func NewHandler(svc *Service) *Handler { return &Handler{svc: svc} }

func toAPI(s Speaker) gen.Speaker {
	return gen.Speaker{UserId: s.UserID, ConsentedAt: s.ConsentedAt, WithdrawnAt: s.WithdrawnAt}
}

// RecordConsent implements PUT /v1/speakers/{user_id}/consent.
func (h *Handler) RecordConsent(ctx context.Context, req gen.RecordConsentRequestObject) (gen.RecordConsentResponseObject, error) {
	s, err := h.svc.RecordConsent(ctx, req.UserId, req.Body)
	if err != nil {
		return nil, err
	}
	return gen.RecordConsent200JSONResponse(toAPI(s)), nil
}

// DownloadConsent implements GET /v1/speakers/{user_id}/consent.
func (h *Handler) DownloadConsent(ctx context.Context, req gen.DownloadConsentRequestObject) (gen.DownloadConsentResponseObject, error) {
	rd, size, err := h.svc.ConsentForm(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	// The generated writer copies Body to the response and closes it.
	return gen.DownloadConsent200ApplicationoctetStreamResponse{Body: rd, ContentLength: size}, nil
}

// GetMySpeaker implements GET /v1/me/speaker.
func (h *Handler) GetMySpeaker(ctx context.Context, _ gen.GetMySpeakerRequestObject) (gen.GetMySpeakerResponseObject, error) {
	s, err := h.svc.Me(ctx)
	if err != nil {
		return nil, err
	}
	return gen.GetMySpeaker200JSONResponse(toAPI(s)), nil
}

// WithdrawConsent implements DELETE /v1/me/speaker.
func (h *Handler) WithdrawConsent(ctx context.Context, _ gen.WithdrawConsentRequestObject) (gen.WithdrawConsentResponseObject, error) {
	s, err := h.svc.Withdraw(ctx)
	if err != nil {
		return nil, err
	}
	return gen.WithdrawConsent200JSONResponse(toAPI(s)), nil
}
