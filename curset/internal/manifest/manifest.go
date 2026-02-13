package manifest

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const manifestPath = ".cursor/.curset.json"

// Entry represents a single installed item tracked by curset.
type Entry struct {
	Type  string   `json:"type"`  // object type, e.g. "rules", "commands"
	Name  string   `json:"name"`  // entry name, e.g. "common", "get-conflict-responsible"
	IsDir bool     `json:"is_dir"`
	Files []string `json:"files"` // list of file paths relative to .cursor/
}

// Manifest tracks all items installed by curset.
type Manifest struct {
	Collections []string `json:"collections"` // list of installed collection names
	Entries     []Entry  `json:"entries"`
}

// Load reads the manifest from .cursor/.curset.json.
// Returns an empty manifest if the file does not exist.
func Load() (*Manifest, error) {
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		if os.IsNotExist(err) {
			return &Manifest{}, nil
		}
		return nil, fmt.Errorf("failed to read manifest: %w", err)
	}

	var m Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("failed to parse manifest: %w", err)
	}

	return &m, nil
}

// Save writes the manifest to .cursor/.curset.json.
func (m *Manifest) Save() error {
	dir := filepath.Dir(manifestPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal manifest: %w", err)
	}

	if err := os.WriteFile(manifestPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write manifest: %w", err)
	}

	return nil
}

// IsManaged checks if a given type/name entry is tracked by curset.
func (m *Manifest) IsManaged(objType, name string) bool {
	for _, e := range m.Entries {
		if e.Type == objType && e.Name == name {
			return true
		}
	}
	return false
}

// HasCollection checks if a collection is already installed.
func (m *Manifest) HasCollection(name string) bool {
	for _, c := range m.Collections {
		if c == name {
			return true
		}
	}
	return false
}

// AddCollection adds a collection name to the installed list if not already present.
func (m *Manifest) AddCollection(name string) {
	if !m.HasCollection(name) {
		m.Collections = append(m.Collections, name)
	}
}

// RemoveCollection removes a collection name from the installed list.
func (m *Manifest) RemoveCollection(name string) {
	for i, c := range m.Collections {
		if c == name {
			m.Collections = append(m.Collections[:i], m.Collections[i+1:]...)
			return
		}
	}
}

// AddOrUpdate adds a new entry or updates an existing one in the manifest.
func (m *Manifest) AddOrUpdate(entry Entry) {
	for i, e := range m.Entries {
		if e.Type == entry.Type && e.Name == entry.Name {
			m.Entries[i] = entry
			return
		}
	}
	m.Entries = append(m.Entries, entry)
}

// RemoveEntry removes an entry from the manifest by type and name.
func (m *Manifest) RemoveEntry(objType, name string) {
	for i, e := range m.Entries {
		if e.Type == objType && e.Name == name {
			m.Entries = append(m.Entries[:i], m.Entries[i+1:]...)
			return
		}
	}
}

// GetEntry returns an entry by type and name, or nil if not found.
func (m *Manifest) GetEntry(objType, name string) *Entry {
	for i, e := range m.Entries {
		if e.Type == objType && e.Name == name {
			return &m.Entries[i]
		}
	}
	return nil
}
