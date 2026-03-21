package statestore

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/uttejg/newbox/internal/core/domain"
)

// installStateDTO is the persistence representation of domain.InstallState.
// Keeping JSON tags here prevents coupling the domain model to serialisation.
type installStateDTO struct {
	CompletedIDs []string  `json:"completed_ids"`
	FailedIDs    []string  `json:"failed_ids"`
	StartedAt    time.Time `json:"started_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func toDTO(s *domain.InstallState) installStateDTO {
	return installStateDTO{
		CompletedIDs: s.CompletedIDs,
		FailedIDs:    s.FailedIDs,
		StartedAt:    s.StartedAt,
		UpdatedAt:    s.UpdatedAt,
	}
}

func fromDTO(d installStateDTO) *domain.InstallState {
	return &domain.InstallState{
		CompletedIDs: d.CompletedIDs,
		FailedIDs:    d.FailedIDs,
		StartedAt:    d.StartedAt,
		UpdatedAt:    d.UpdatedAt,
	}
}

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
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, err
	}
	return &FileStore{Path: filepath.Join(dir, "state.json")}, nil
}

func (f *FileStore) Save(state *domain.InstallState) error {
	data, err := json.MarshalIndent(toDTO(state), "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(f.Path, data, 0600)
}

func (f *FileStore) Load() (*domain.InstallState, error) {
	data, err := os.ReadFile(f.Path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var dto installStateDTO
	if err := json.Unmarshal(data, &dto); err != nil {
		return nil, err
	}
	return fromDTO(dto), nil
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
