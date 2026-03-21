package pkgmgr

import (
	"context"
	"fmt"
	"strings"

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

func (c *CompositeManager) CanHandle(_ domain.PackageRef) bool { return true }

// SubManagers returns the underlying sub-managers for delegation.
func (c *CompositeManager) SubManagers() []port.PackageManager {
	return c.managers
}

func (c *CompositeManager) IsAvailable(ctx context.Context) error {
	if len(c.managers) == 0 {
		return fmt.Errorf("no supported package manager found for this platform")
	}
	var errs []string
	for _, m := range c.managers {
		if err := m.IsAvailable(ctx); err != nil {
			errs = append(errs, err.Error())
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("package managers unavailable: %s", strings.Join(errs, "; "))
	}
	return nil
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

// managerFor selects the right manager based on CanHandle.
func (c *CompositeManager) managerFor(ref domain.PackageRef) port.PackageManager {
	for _, m := range c.managers {
		if m.CanHandle(ref) {
			return m
		}
	}
	return nil
}
