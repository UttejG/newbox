package pkgmgr

import (
	"context"
	"fmt"

	"github.com/uttejg/newbox/internal/core/domain"
	"github.com/uttejg/newbox/internal/core/port"
)

// CompositeManager tries each sub-manager in order, routing to the one that
// can handle the given PackageRef based on which fields are set.
type CompositeManager struct {
	managers []port.PackageManager
}

// NewComposite creates a CompositeManager from the provided sub-managers.
func NewComposite(managers ...port.PackageManager) *CompositeManager {
	return &CompositeManager{managers: managers}
}

func (c *CompositeManager) Name() string { return "composite" }

func (c *CompositeManager) IsAvailable(ctx context.Context) bool {
	for _, m := range c.managers {
		if m.IsAvailable(ctx) {
			return true
		}
	}
	return false
}

func (c *CompositeManager) IsInstalled(ctx context.Context, ref domain.PackageRef) (bool, error) {
	mgr := c.managerFor(ref)
	if mgr == nil {
		return false, nil
	}
	return mgr.IsInstalled(ctx, ref)
}

func (c *CompositeManager) Install(ctx context.Context, ref domain.PackageRef) (*port.RunResult, error) {
	mgr := c.managerFor(ref)
	if mgr == nil {
		return nil, fmt.Errorf("no package manager available for this package")
	}
	return mgr.Install(ctx, ref)
}

// BuildCommand delegates to the appropriate sub-manager for plan display.
func (c *CompositeManager) BuildCommand(ref domain.PackageRef) string {
	mgr := c.managerFor(ref)
	if mgr == nil {
		return ""
	}
	return mgr.BuildCommand(ref)
}

// managerFor selects the right manager based on which fields are set in ref.
func (c *CompositeManager) managerFor(ref domain.PackageRef) port.PackageManager {
	if ref.MAS != "" {
		for _, m := range c.managers {
			if m.Name() == "mas" {
				return m
			}
		}
	}
	if ref.Formula != "" || ref.Cask != "" {
		for _, m := range c.managers {
			if m.Name() == "brew" {
				return m
			}
		}
	}
	return nil
}
