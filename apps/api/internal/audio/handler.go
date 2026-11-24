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
	return gen.Speaker{UserId: s.UserID, Languages: s.Languages, ConsentedAt: s.ConsentedAt, WithdrawnAt: s.WithdrawnAt}
}

// RecordConsent implements PUT /v1/speakers/{user_id}/consent.
func (h *Handler) RecordConsent(ctx context.Context, req gen.RecordConsentRequestObject) (gen.RecordConsentResponseObject, error) {
	s, err := h.svc.RecordConsent(ctx, req.UserId, req.Params.Languages, req.Body)
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

func clip(c Clip) gen.Clip {
	return gen.Clip{Id: c.ID, EntryId: c.EntryID, Url: "/v1/audio/" + c.ID.String(), ContentType: c.ContentType, RecordedAt: c.RecordedAt}
}

// ListClips implements GET /v1/audio.
func (h *Handler) ListClips(ctx context.Context, req gen.ListClipsRequestObject) (gen.ListClipsResponseObject, error) {
	cs, err := h.svc.Clips(ctx, req.Params.EntryId)
	if err != nil {
		return nil, err
	}
	out := make([]gen.Clip, len(cs))
	for i, c := range cs {
		out[i] = clip(c)
	}
	return gen.ListClips200JSONResponse{Data: out}, nil
}

// RecordClip implements POST /v1/audio.
func (h *Handler) RecordClip(ctx context.Context, req gen.RecordClipRequestObject) (gen.RecordClipResponseObject, error) {
	c, err := h.svc.Record(ctx, req.Params.EntryId, req.Body)
	if err != nil {
		return nil, err
	}
	return gen.RecordClip201JSONResponse(clip(c)), nil
}

// GetClip implements GET /v1/audio/{id}.
func (h *Handler) GetClip(ctx context.Context, req gen.GetClipRequestObject) (gen.GetClipResponseObject, error) {
	rd, c, err := h.svc.Open(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	// The generated writer copies Body to the response and closes it.
	return gen.GetClip200AudioResponse{Body: rd, ContentType: c.ContentType, ContentLength: int64(c.Size)}, nil
}

// DeleteClip implements DELETE /v1/audio/{id}.
func (h *Handler) DeleteClip(ctx context.Context, req gen.DeleteClipRequestObject) (gen.DeleteClipResponseObject, error) {
	if err := h.svc.Delete(ctx, req.Id); err != nil {
		return nil, err
	}
	return gen.DeleteClip204Response{}, nil
}
