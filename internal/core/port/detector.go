package port

import "github.com/uttejg/newbox/internal/core/domain"

// PlatformDetector detects the current system platform.
type PlatformDetector interface {
	Detect() (*domain.Platform, error)
}
