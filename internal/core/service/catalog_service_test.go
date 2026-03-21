package service_test

import (
	"errors"
	"testing"

	"github.com/uttejg/newbox/internal/core/domain"
	"github.com/uttejg/newbox/internal/core/service"
	"github.com/uttejg/newbox/internal/testutil"
)

func newService() *service.CatalogService {
	cats, profiles := testutil.DefaultTestCatalog()
	return service.NewCatalogService(&testutil.FakeCatalogProvider{
		Categories: cats,
		Profiles:   profiles,
	})
}

func TestCatalogService_GetCategories_FiltersForOS(t *testing.T) {
	svc := newService()

	tests := []struct {
		os        domain.OS
		wantCount int
		wantID    string
	}{
		{domain.OSMacOS, 3, "messaging"},   // all 3 categories have macOS tools
		{domain.OSLinux, 3, "messaging"},   // all 3 have Linux tools
		{domain.OSWindows, 3, "messaging"}, // all 3 have Windows tools (WinOnly in messaging, Firefox in browsers, git+jq in cli)
	}

	for _, tt := range tests {
		t.Run(tt.os.String(), func(t *testing.T) {
			cats, err := svc.GetCategories(tt.os)
			if err != nil {
				t.Fatalf("GetCategories(%v) error: %v", tt.os, err)
			}
			if len(cats) != tt.wantCount {
				t.Errorf("GetCategories(%v) = %d categories, want %d", tt.os, len(cats), tt.wantCount)
			}
			found := false
			for _, c := range cats {
				if c.ID == tt.wantID {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("GetCategories(%v) missing category with ID %q", tt.os, tt.wantID)
			}
		})
	}
}

func TestCatalogService_GetCategories_ToolsFiltered(t *testing.T) {
	svc := newService()

	// On macOS, messaging should have Signal + Telegram (not WinOnly)
	cats, err := svc.GetCategories(domain.OSMacOS)
	if err != nil {
		t.Fatal(err)
	}

	var msgCat *domain.Category
	for i := range cats {
		if cats[i].ID == "messaging" {
			msgCat = &cats[i]
			break
		}
	}
	if msgCat == nil {
		t.Fatal("messaging category not found")
	}
	if len(msgCat.Tools) != 2 {
		t.Errorf("messaging on macOS has %d tools, want 2 (Signal, Telegram)", len(msgCat.Tools))
	}
}

func TestCatalogService_GetCategories_PropagatesError(t *testing.T) {
	svc := service.NewCatalogService(&testutil.FakeCatalogProvider{
		Err: errors.New("disk read error"),
	})
	_, err := svc.GetCategories(domain.OSMacOS)
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestCatalogService_GetProfile_Found(t *testing.T) {
	svc := newService()

	p, err := svc.GetProfile("developer")
	if err != nil {
		t.Fatalf("GetProfile(developer) error: %v", err)
	}
	if p.ID != "developer" {
		t.Errorf("got profile ID %q, want developer", p.ID)
	}
	if len(p.Categories) != 3 {
		t.Errorf("developer profile has %d categories, want 3", len(p.Categories))
	}
}

func TestCatalogService_GetProfile_NotFound(t *testing.T) {
	svc := newService()

	_, err := svc.GetProfile("nonexistent")
	if err == nil {
		t.Error("expected error for unknown profile, got nil")
	}
}

func TestCatalogService_GetAllProfiles(t *testing.T) {
	svc := newService()

	profiles, err := svc.GetAllProfiles()
	if err != nil {
		t.Fatalf("GetAllProfiles() error: %v", err)
	}
	if len(profiles) == 0 {
		t.Error("GetAllProfiles() returned empty list")
	}
}
