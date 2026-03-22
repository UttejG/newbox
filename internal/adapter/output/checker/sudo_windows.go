//go:build windows

package checker

import "context"

// CheckSudo is a no-op on Windows; winget handles elevation internally.
func (c *SystemChecker) CheckSudo(_ context.Context) error {
	return nil
}
