package statestore

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/uttejg/newbox/internal/core/domain"
)

// FileStore persists InstallState as JSON at ~/.newbox/state.json.
type FileStore struct {
	Path string
}

// NewFileStore creates a FileStore pointing to ~/.newbox/state.json,
// creating the directory if needed.
func NewFileStore() (*FileStore, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	dir := filepath.Join(home, ".newbox")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	return &FileStore{Path: filepath.Join(dir, "state.json")}, nil
}

func (f *FileStore) Save(state *domain.InstallState) error {
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(f.Path, data, 0644)
}

func (f *FileStore) Load() (*domain.InstallState, error) {
	data, err := os.ReadFile(f.Path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var state domain.InstallState
	return &state, json.Unmarshal(data, &state)
}

func (f *FileStore) Clear() error {
	err := os.Remove(f.Path)
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

func (f *FileStore) Exists() bool {
	_, err := os.Stat(f.Path)
	return err == nil
}
