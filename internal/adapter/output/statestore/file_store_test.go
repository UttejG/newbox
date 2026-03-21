package statestore_test

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/uttejg/newbox/internal/adapter/output/statestore"
	"github.com/uttejg/newbox/internal/core/domain"
)

func newTestStore(t *testing.T) *statestore.FileStore {
	t.Helper()
	dir := t.TempDir()
	return &statestore.FileStore{Path: filepath.Join(dir, "state.json")}
}

func TestFileStore_ExistsReturnsFalseInitially(t *testing.T) {
	fs := newTestStore(t)
	if fs.Exists() {
		t.Error("expected Exists() = false before any Save")
	}
}

func TestFileStore_SaveLoad_RoundTrip(t *testing.T) {
	fs := newTestStore(t)

	state := &domain.InstallState{
		CompletedIDs: []string{"git", "signal"},
		FailedIDs:    []string{"mas-app"},
		StartedAt:    time.Now().Round(time.Second),
		UpdatedAt:    time.Now().Round(time.Second),
	}

	if err := fs.Save(state); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	if !fs.Exists() {
		t.Error("expected Exists() = true after Save")
	}

	loaded, err := fs.Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if loaded == nil {
		t.Fatal("Load() returned nil")
	}
	if len(loaded.CompletedIDs) != 2 {
		t.Errorf("CompletedIDs len = %d, want 2", len(loaded.CompletedIDs))
	}
	if loaded.CompletedIDs[0] != "git" || loaded.CompletedIDs[1] != "signal" {
		t.Errorf("CompletedIDs = %v", loaded.CompletedIDs)
	}
	if len(loaded.FailedIDs) != 1 || loaded.FailedIDs[0] != "mas-app" {
		t.Errorf("FailedIDs = %v", loaded.FailedIDs)
	}
}

func TestFileStore_Load_ReturnsNilWhenMissing(t *testing.T) {
	fs := newTestStore(t)
	state, err := fs.Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if state != nil {
		t.Error("expected nil when file does not exist")
	}
}

func TestFileStore_Clear(t *testing.T) {
	fs := newTestStore(t)

	if err := fs.Save(&domain.InstallState{CompletedIDs: []string{"git"}}); err != nil {
		t.Fatalf("Save() error = %v", err)
	}
	if !fs.Exists() {
		t.Fatal("expected file to exist after Save")
	}

	if err := fs.Clear(); err != nil {
		t.Fatalf("Clear() error = %v", err)
	}
	if fs.Exists() {
		t.Error("expected Exists() = false after Clear")
	}
}

func TestFileStore_Clear_IdempotentWhenMissing(t *testing.T) {
	fs := newTestStore(t)
	// Clearing a non-existent file should not error.
	if err := fs.Clear(); err != nil {
		t.Errorf("Clear() on missing file error = %v", err)
	}
}

func TestFileStore_IsCompleted_MarkCompleted(t *testing.T) {
	state := &domain.InstallState{}
	if state.IsCompleted("git") {
		t.Error("expected IsCompleted to return false initially")
	}

	state.MarkCompleted("git")
	if !state.IsCompleted("git") {
		t.Error("expected IsCompleted to return true after MarkCompleted")
	}

	// Idempotent: second mark should not duplicate.
	state.MarkCompleted("git")
	if len(state.CompletedIDs) != 1 {
		t.Errorf("CompletedIDs len = %d after duplicate mark, want 1", len(state.CompletedIDs))
	}
}
