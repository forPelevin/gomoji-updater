package common

import (
	"context"

	"github.com/forPelevin/gomoji"
)

type (
	// Provider is an emoji provider.
	Provider interface {
		AllEmojis(ctx context.Context) ([]gomoji.Emoji, error)
	}
)

// Service is a service to work with emojis using a provider.
// Deprecated: new consumers can call provider directly.
type Service struct {
	p Provider
}

// NewService creates a new instance of service.
// Deprecated: prefer using providers without Service to avoid extra allocations.
func NewService(p Provider) *Service {
	return &Service{p: p}
}

// AllEmojis gets all emojis from provider.
func (s *Service) AllEmojis(ctx context.Context) ([]gomoji.Emoji, error) {
	return s.p.AllEmojis(ctx)
}
