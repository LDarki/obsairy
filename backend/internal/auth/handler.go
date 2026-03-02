package auth

import "context"

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Register(ctx context.Context, email, password string) (*User, error) {
	return h.service.Register(ctx, email, password)
}

func (h *Handler) Login(ctx context.Context, email, password, userAgent, ip string) (*Session, error) {
	return h.service.Login(ctx, email, password, userAgent, ip)
}

func (h *Handler) GetSessions(ctx context.Context, userID string) ([]Session, error) {
	return h.service.GetUserSessions(ctx, userID)
}

func (h *Handler) RevokeSession(ctx context.Context, userID, sessionID string) error {
	return h.service.RevokeSession(ctx, userID, sessionID)
}

func (h *Handler) RevokeAllSessions(ctx context.Context, userID string) error {
	return h.service.RevokeAllUserSessions(ctx, userID)
}
