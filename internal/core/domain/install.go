package domain

import "time"

// InstallStatus is the status of an individual tool in the install plan.
type InstallStatus string

const (
StatusPending    InstallStatus = "pending"
StatusInstalling InstallStatus = "installing"
StatusDone       InstallStatus = "done"
StatusSkipped    InstallStatus = "skipped" // already installed
StatusFailed     InstallStatus = "failed"
StatusDryRun     InstallStatus = "dry_run" // would install (dry-run mode)
)

// ExecutionStep represents one installation action.
type ExecutionStep struct {
Tool     Tool
Ref      PackageRef
Command  string // e.g., "brew install --cask signal"
Status   InstallStatus
Error    error
Output   string
Duration time.Duration
}

// InstallPlan is the full computed plan before execution.
type InstallPlan struct {
Steps  []ExecutionStep
DryRun bool
}

// PendingSteps returns steps that need to be executed (pending or dry-run).
func (p *InstallPlan) PendingSteps() []ExecutionStep {
var out []ExecutionStep
for _, s := range p.Steps {
if s.Status == StatusPending || s.Status == StatusDryRun {
out = append(out, s)
}
}
return out
}

// SkippedSteps returns steps that are already installed and will be skipped.
func (p *InstallPlan) SkippedSteps() []ExecutionStep {
var out []ExecutionStep
for _, s := range p.Steps {
if s.Status == StatusSkipped {
out = append(out, s)
}
}
return out
}

// ProgressEvent is sent on a channel during execution.
type ProgressEvent struct {
Step  ExecutionStep
Index int // 0-based index into PendingSteps
Total int
}

// InstallState persists between runs for resume support.
type InstallState struct {
CompletedIDs []string  `json:"completed_ids"` // tool names already installed
FailedIDs    []string  `json:"failed_ids"`    // tool names that failed
StartedAt    time.Time `json:"started_at"`
UpdatedAt    time.Time `json:"updated_at"`
}

func (s *InstallState) IsCompleted(toolID string) bool {
for _, id := range s.CompletedIDs {
if id == toolID {
return true
}
}
return false
}

func (s *InstallState) MarkCompleted(toolID string) {
if !s.IsCompleted(toolID) {
s.CompletedIDs = append(s.CompletedIDs, toolID)
s.UpdatedAt = time.Now()
}
}

// PreflightResult captures system readiness.
type PreflightResult struct {
InternetOK       bool
DiskSpaceOK      bool
SudoOK           bool
PackageManagerOK bool
Errors           []string
}

// OK returns true if all required checks passed.
func (r *PreflightResult) OK() bool {
return r.InternetOK && r.DiskSpaceOK && r.PackageManagerOK
}
