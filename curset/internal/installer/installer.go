package installer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bilgehannal/cursor-config/curset/internal/collection"
	"github.com/bilgehannal/cursor-config/curset/internal/github"
	"github.com/bilgehannal/cursor-config/curset/internal/manifest"
)

// Installer handles installing collections into the current directory.
type Installer struct {
	client   *github.Client
	manifest *manifest.Manifest
}

// NewInstaller creates a new Installer.
func NewInstaller(client *github.Client) (*Installer, error) {
	m, err := manifest.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load manifest: %w", err)
	}

	return &Installer{
		client:   client,
		manifest: m,
	}, nil
}

// Install installs a named collection into the current directory's .cursor/ folder.
func (inst *Installer) Install(col collection.Collection, name string) error {
	fmt.Printf("Installing collection: %s\n\n", name)

	inst.manifest.Collection = name

	for objType, entries := range col {
		for _, entry := range entries {
			if err := inst.installEntry(objType, entry); err != nil {
				return fmt.Errorf("failed to install %s/%s: %w", objType, entry, err)
			}
		}
	}

	// Save the manifest after successful install.
	if err := inst.manifest.Save(); err != nil {
		return fmt.Errorf("failed to save manifest: %w", err)
	}

	fmt.Println("\nDone.")
	return nil
}

// installEntry installs a single entry (which may be a directory or a file).
func (inst *Installer) installEntry(objType, entry string) error {
	// Use GitHub Contents API to determine if this is a file or directory.
	remotePath := fmt.Sprintf("%s/%s", objType, entry)
	result, err := inst.client.ListContents(remotePath)

	if err != nil {
		// If not found as a direct path, it might be a file without extension.
		// List the parent directory and find matching files.
		return inst.installFileByName(objType, entry)
	}

	if result.IsDir {
		// It's a directory - install all files in it.
		return inst.installDirectory(objType, entry, result.Entries)
	}

	// It's a single file.
	return inst.installSingleFile(objType, result.Entries[0], entry)
}

// installDirectory installs all files from a remote directory.
func (inst *Installer) installDirectory(objType, entry string, contents []github.ContentEntry) error {
	localDir := filepath.Join(".cursor", objType, entry)
	managed := inst.manifest.IsManaged(objType, entry)

	// Check if directory already exists locally and is NOT managed by curset.
	if _, err := os.Stat(localDir); err == nil && !managed {
		fmt.Printf("  skipped: %s (already exists, not managed by curset)\n", localDir)
		return nil
	}

	if managed {
		fmt.Printf("  updating: %s\n", localDir)
	}

	// Create the directory.
	if err := os.MkdirAll(localDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", localDir, err)
	}

	var installedFiles []string
	for _, c := range contents {
		if c.Type != "file" {
			continue
		}

		remotePath := fmt.Sprintf("%s/%s/%s", objType, entry, c.Name)
		data, err := inst.client.DownloadFile(remotePath)
		if err != nil {
			return err
		}

		localPath := filepath.Join(localDir, c.Name)
		if err := os.WriteFile(localPath, data, 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", localPath, err)
		}
		installedFiles = append(installedFiles, filepath.Join(objType, entry, c.Name))
	}

	if !managed {
		fmt.Printf("  installed: %s (%d files)\n", localDir, len(installedFiles))
	} else {
		fmt.Printf("  updated: %s (%d files)\n", localDir, len(installedFiles))
	}

	// Track in manifest.
	inst.manifest.AddOrUpdate(manifest.Entry{
		Type:  objType,
		Name:  entry,
		IsDir: true,
		Files: installedFiles,
	})

	return nil
}

// installSingleFile installs a single file that was found directly by path.
func (inst *Installer) installSingleFile(objType string, entry github.ContentEntry, entryName string) error {
	localDir := filepath.Join(".cursor", objType)
	if err := os.MkdirAll(localDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", localDir, err)
	}

	localPath := filepath.Join(localDir, entry.Name)
	managed := inst.manifest.IsManaged(objType, entryName)

	// Check if file already exists locally and is NOT managed by curset.
	if _, err := os.Stat(localPath); err == nil && !managed {
		fmt.Printf("  skipped: %s (already exists, not managed by curset)\n", localPath)
		return nil
	}

	// The entry.Path is relative to the repo root (e.g. "data/.cursor/commands/file.md").
	// We need the path relative to data/.cursor/.
	relativePath := strings.TrimPrefix(entry.Path, "data/.cursor/")
	data, err := inst.client.DownloadFile(relativePath)
	if err != nil {
		return err
	}

	if err := os.WriteFile(localPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", localPath, err)
	}

	action := "installed"
	if managed {
		action = "updated"
	}
	fmt.Printf("  %s: %s\n", action, localPath)

	// Track in manifest.
	inst.manifest.AddOrUpdate(manifest.Entry{
		Type:  objType,
		Name:  entryName,
		IsDir: false,
		Files: []string{filepath.Join(objType, entry.Name)},
	})

	return nil
}

// installFileByName searches the parent directory for files matching the entry name
// (without extension) and installs them.
func (inst *Installer) installFileByName(objType, entry string) error {
	// List the parent directory (e.g. "commands").
	result, err := inst.client.ListContents(objType)
	if err != nil {
		return fmt.Errorf("failed to list %s: %w", objType, err)
	}

	localDir := filepath.Join(".cursor", objType)
	if err := os.MkdirAll(localDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", localDir, err)
	}

	managed := inst.manifest.IsManaged(objType, entry)

	found := false
	var installedFiles []string
	for _, c := range result.Entries {
		if c.Type != "file" {
			continue
		}

		// Match by name without extension.
		nameWithoutExt := strings.TrimSuffix(c.Name, filepath.Ext(c.Name))
		if nameWithoutExt == entry {
			localPath := filepath.Join(localDir, c.Name)

			// Check if file already exists locally and is NOT managed by curset.
			if _, err := os.Stat(localPath); err == nil && !managed {
				fmt.Printf("  skipped: %s (already exists, not managed by curset)\n", localPath)
				found = true
				continue
			}

			relativePath := strings.TrimPrefix(c.Path, "data/.cursor/")
			data, err := inst.client.DownloadFile(relativePath)
			if err != nil {
				return err
			}

			if err := os.WriteFile(localPath, data, 0644); err != nil {
				return fmt.Errorf("failed to write %s: %w", localPath, err)
			}

			action := "installed"
			if managed {
				action = "updated"
			}
			fmt.Printf("  %s: %s\n", action, localPath)

			installedFiles = append(installedFiles, filepath.Join(objType, c.Name))
			found = true
		}
	}

	if !found {
		return fmt.Errorf("entry %s/%s not found in remote repository", objType, entry)
	}

	// Track in manifest.
	inst.manifest.AddOrUpdate(manifest.Entry{
		Type:  objType,
		Name:  entry,
		IsDir: false,
		Files: installedFiles,
	})

	return nil
}
