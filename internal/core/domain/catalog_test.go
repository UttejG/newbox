package domain_test

import (
	"testing"

	"github.com/uttejg/newbox/internal/core/domain"
)

func TestPackageRef_IsEmpty(t *testing.T) {
	tests := []struct {
		name string
		ref  domain.PackageRef
		want bool
	}{
		{"empty ref", domain.PackageRef{}, true},
		{"formula set", domain.PackageRef{domain.PackageManagerFormula: "git"}, false},
		{"cask set", domain.PackageRef{domain.PackageManagerCask: "signal"}, false},
		{"winget set", domain.PackageRef{domain.PackageManagerWinget: "OpenWhisperSystems.Signal"}, false},
		{"apt set", domain.PackageRef{domain.PackageManagerApt: "signal-desktop"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ref.IsEmpty(); got != tt.want {
				t.Errorf("PackageRef.IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTool_IsAvailableOn(t *testing.T) {
	tool := domain.Tool{
		Name:    "Signal",
		MacOS:   domain.PackageRef{domain.PackageManagerCask: "signal"},
		Linux:   domain.PackageRef{domain.PackageManagerApt: "signal-desktop"},
		Windows: nil,
	}

	tests := []struct {
		name string
		os   domain.OS
		want bool
	}{
		{"available on macOS", domain.OSMacOS, true},
		{"available on Linux", domain.OSLinux, true},
		{"not available on Windows", domain.OSWindows, false},
		{"not available on unknown", domain.OSUnknown, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tool.IsAvailableOn(tt.os); got != tt.want {
				t.Errorf("Tool.IsAvailableOn(%v) = %v, want %v", tt.os, got, tt.want)
			}
		})
	}
}

func TestTool_PackageRefFor(t *testing.T) {
	macRef := domain.PackageRef{domain.PackageManagerCask: "signal"}
	linuxRef := domain.PackageRef{domain.PackageManagerApt: "signal-desktop"}
	tool := domain.Tool{
		Name:  "Signal",
		MacOS: macRef,
		Linux: linuxRef,
	}

	if got := tool.PackageRefFor(domain.OSMacOS); got[domain.PackageManagerCask] != macRef[domain.PackageManagerCask] {
		t.Errorf("PackageRefFor(macOS) returned wrong ref")
	}
	if got := tool.PackageRefFor(domain.OSLinux); got[domain.PackageManagerApt] != linuxRef[domain.PackageManagerApt] {
		t.Errorf("PackageRefFor(linux) returned wrong ref")
	}
	if got := tool.PackageRefFor(domain.OSWindows); !got.IsEmpty() {
		t.Errorf("PackageRefFor(windows) = %v, want empty", got)
	}
}

func TestCategory_FilteredTools(t *testing.T) {
	cat := domain.Category{
		ID:   "messaging",
		Name: "💬 Messaging",
		Tools: []domain.Tool{
			{Name: "Signal", MacOS: domain.PackageRef{domain.PackageManagerCask: "signal"}, Linux: domain.PackageRef{domain.PackageManagerApt: "signal-desktop"}},
			{Name: "iMessage", MacOS: domain.PackageRef{domain.PackageManagerMAS: "12345"}}, // macOS only
			{Name: "WinTool", Windows: domain.PackageRef{domain.PackageManagerWinget: "Win.Tool"}}, // windows only
		},
	}

	macTools := cat.FilteredTools(domain.OSMacOS)
	if len(macTools) != 2 {
		t.Errorf("FilteredTools(macOS) = %d tools, want 2", len(macTools))
	}

	linuxTools := cat.FilteredTools(domain.OSLinux)
	if len(linuxTools) != 1 {
		t.Errorf("FilteredTools(linux) = %d tools, want 1", len(linuxTools))
	}

	winTools := cat.FilteredTools(domain.OSWindows)
	if len(winTools) != 1 {
		t.Errorf("FilteredTools(windows) = %d tools, want 1", len(winTools))
	}
}

func TestCategory_IsAvailableOn(t *testing.T) {
	cat := domain.Category{
		ID: "macos-only",
		Tools: []domain.Tool{
			{Name: "Alfred", MacOS: domain.PackageRef{domain.PackageManagerCask: "alfred"}},
		},
	}

	if !cat.IsAvailableOn(domain.OSMacOS) {
		t.Error("category should be available on macOS")
	}
	if cat.IsAvailableOn(domain.OSLinux) {
		t.Error("category should NOT be available on Linux")
	}
}

func TestUserSelection_TotalCount(t *testing.T) {
	sel := &domain.UserSelection{
		ToolsByCategory: map[string][]domain.Tool{
			"messaging": {
				{Name: "Signal"},
				{Name: "Telegram"},
			},
			"browsers": {
				{Name: "Firefox"},
			},
		},
	}

	if got := sel.TotalCount(); got != 3 {
		t.Errorf("TotalCount() = %d, want 3", got)
	}
}
