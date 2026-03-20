package detector_test

import (
	"runtime"
	"strings"
	"testing"

	"github.com/uttejg/newbox/internal/adapter/output/detector"
	"github.com/uttejg/newbox/internal/core/domain"
)

func TestSystemDetector_Detect(t *testing.T) {
	d := &detector.SystemDetector{}
	p, err := d.Detect()
	if err != nil {
		t.Fatalf("Detect() returned error: %v", err)
	}

	if p == nil {
		t.Fatal("Detect() returned nil platform")
	}

	// Verify OS matches runtime
	switch runtime.GOOS {
	case "darwin":
		if p.OS != domain.OSMacOS {
			t.Errorf("expected OSMacOS on darwin, got %v", p.OS)
		}
	case "linux":
		if p.OS != domain.OSLinux {
			t.Errorf("expected OSLinux on linux, got %v", p.OS)
		}
	case "windows":
		if p.OS != domain.OSWindows {
			t.Errorf("expected OSWindows on windows, got %v", p.OS)
		}
	}

	// Verify arch matches runtime
	switch runtime.GOARCH {
	case "amd64":
		if p.Arch != domain.ArchAMD64 {
			t.Errorf("expected ArchAMD64 on amd64, got %v", p.Arch)
		}
	case "arm64":
		if p.Arch != domain.ArchARM64 {
			t.Errorf("expected ArchARM64 on arm64, got %v", p.Arch)
		}
	}
}

func TestSystemDetector_Detect_SummaryNotEmpty(t *testing.T) {
	d := &detector.SystemDetector{}
	p, err := d.Detect()
	if err != nil {
		t.Fatalf("Detect() returned error: %v", err)
	}

	summary := p.Summary()
	if summary == "" {
		t.Error("Platform.Summary() returned empty string")
	}
}

func TestFormatDetectionInfo(t *testing.T) {
	tests := []struct {
		name     string
		platform domain.Platform
		wantSub  string
	}{
		{
			name: "macOS with brew",
			platform: domain.Platform{
				OS:             domain.OSMacOS,
				Arch:           domain.ArchARM64,
				PackageManager: domain.PkgMgrBrew,
			},
			wantSub: "Package manager: brew",
		},
		{
			name: "no package manager",
			platform: domain.Platform{
				OS:             domain.OSMacOS,
				Arch:           domain.ArchARM64,
				PackageManager: domain.PkgMgrNone,
			},
			wantSub: "No supported package manager found",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info := detector.FormatDetectionInfo(&tt.platform)
			if info == "" {
				t.Error("FormatDetectionInfo returned empty string")
			}
			if !strings.Contains(info, tt.wantSub) {
				t.Errorf("FormatDetectionInfo() = %q, want substring %q", info, tt.wantSub)
			}
		})
	}
}

