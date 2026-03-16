package port

import "github.com/uttejg/newbox/internal/core/domain"

// StateStore persists installation progress for resume support.
type StateStore interface {
	Save(state *domain.InstallState) error
	Load() (*domain.InstallState, error)
	Clear() error
	Exists() bool
}
