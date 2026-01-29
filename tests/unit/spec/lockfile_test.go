package spec_test

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"
	"specledger/internal/spec"
)

func TestLockfileWriteRead(t *testing.T) {
	tests := []struct {
		name    string
		version string
		entries []spec.LockfileEntry
		wantErr bool
	}{
		{
			name: "basic lockfile",
			version: "1.0.0",
			entries: []spec.LockfileEntry{
				{
					RepositoryURL: "https://github.com/example/repo",
					CommitHash:    "abc123",
					ContentHash:   "def456",
					SpecPath:      "specs/test.md",
					Size:          100,
				},
			},
			wantErr: false,
		},
		{
			name: "lockfile with multiple entries",
			version: "1.0.0",
			entries: []spec.LockfileEntry{
				{
					RepositoryURL: "https://github.com/example/repo",
					CommitHash:    "abc123",
					ContentHash:   "def456",
					SpecPath:      "specs/test.md",
					Size:          100,
				},
				{
					RepositoryURL: "https://github.com/example/other",
					CommitHash:    "ghi789",
					ContentHash:   "jkl012",
					SpecPath:      "specs/other.md",
					Size:          200,
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpFile, err := os.CreateTemp("", "lockfile-*.sum")
			if err != nil {
				t.Fatal(err)
			}
			defer os.Remove(tmpFile.Name())

			lockfile := spec.NewLockfile(tt.version)
			for _, entry := range tt.entries {
				lockfile.AddEntry(entry)
			}

			err = lockfile.Write(tmpFile.Name())
			if (err != nil) != tt.wantErr {
				t.Errorf("Lockfile.Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Read back
			got, err := spec.ReadLockfile(tmpFile.Name())
			if err != nil {
				t.Errorf("Lockfile.Read() error = %v", err)
				return
			}

			if got.Version != tt.version {
				t.Errorf("Lockfile.Read() version = %v, want %v", got.Version, tt.version)
			}

			if len(got.Entries) != len(tt.entries) {
				t.Errorf("Lockfile.Read() got %d entries, want %d", len(got.Entries), len(tt.entries))
			}

			if got.TotalSize != 300 {
				t.Errorf("Lockfile.Read() total size = %d, want 300", got.TotalSize)
			}
		})
	}
}

func TestLockfileAddEntry(t *testing.T) {
	lockfile := spec.NewLockfile("1.0.0")

	entry := spec.LockfileEntry{
		RepositoryURL: "https://github.com/example/repo",
		CommitHash:    "abc123",
		ContentHash:   "def456",
		SpecPath:      "specs/test.md",
		Size:          100,
	}

	lockfile.AddEntry(entry)

	if len(lockfile.Entries) != 1 {
		t.Errorf("AddEntry() got %d entries, want 1", len(lockfile.Entries))
	}

	if lockfile.TotalSize != 100 {
		t.Errorf("AddEntry() total size = %d, want 100", lockfile.TotalSize)
	}
}

func TestLockfileRemoveEntry(t *testing.T) {
	lockfile := spec.NewLockfile("1.0.0")

	lockfile.AddEntry(spec.LockfileEntry{
		RepositoryURL: "https://github.com/example/repo",
		SpecPath:      "specs/test.md",
		Size:          100,
	})

	lockfile.AddEntry(spec.LockfileEntry{
		RepositoryURL: "https://github.com/example/other",
		SpecPath:      "specs/other.md",
		Size:          200,
	})

	removed := lockfile.RemoveEntry("https://github.com/example/repo", "specs/test.md")
	if !removed {
		t.Error("RemoveEntry() should return true")
	}

	if len(lockfile.Entries) != 1 {
		t.Errorf("RemoveEntry() got %d entries, want 1", len(lockfile.Entries))
	}

	removed = lockfile.RemoveEntry("https://github.com/example/repo", "specs/test.md")
	if removed {
		t.Error("RemoveEntry() should return false for non-existent entry")
	}
}

func TestLockfileGetEntry(t *testing.T) {
	lockfile := spec.NewLockfile("1.0.0")

	lockfile.AddEntry(spec.LockfileEntry{
		RepositoryURL: "https://github.com/example/repo",
		SpecPath:      "specs/test.md",
	})

	entry, found := lockfile.GetEntry("https://github.com/example/repo", "specs/test.md")
	if !found {
		t.Error("GetEntry() should find the entry")
	}

	if entry.RepositoryURL != "https://github.com/example/repo" {
		t.Errorf("GetEntry() got URL %s, want https://github.com/example/repo", entry.RepositoryURL)
	}

	_, found = lockfile.GetEntry("https://github.com/example/repo", "specs/other.md")
	if found {
		t.Error("GetEntry() should return false for non-existent entry")
	}
}

func TestLockfileGetRepositoryEntries(t *testing.T) {
	lockfile := spec.NewLockfile("1.0.0")

	lockfile.AddEntry(spec.LockfileEntry{
		RepositoryURL: "https://github.com/example/repo",
		SpecPath:      "specs/test.md",
	})

	lockfile.AddEntry(spec.LockfileEntry{
		RepositoryURL: "https://github.com/example/repo",
		SpecPath:      "specs/other.md",
	})

	lockfile.AddEntry(spec.LockfileEntry{
		RepositoryURL: "https://github.com/other/repo",
		SpecPath:      "specs/test.md",
	})

	entries := lockfile.GetRepositoryEntries("https://github.com/example/repo")
	if len(entries) != 2 {
		t.Errorf("GetRepositoryEntries() got %d entries, want 2", len(entries))
	}

	entries = lockfile.GetRepositoryEntries("https://github.com/other/repo")
	if len(entries) != 1 {
		t.Errorf("GetRepositoryEntries() got %d entries, want 1", len(entries))
	}

	entries = lockfile.GetRepositoryEntries("https://github.com/nonexistent/repo")
	if len(entries) != 0 {
		t.Errorf("GetRepositoryEntries() got %d entries, want 0", len(entries))
	}
}

func TestLockfileVerify(t *testing.T) {
	tests := []struct {
		name    string
		manifest *spec.Manifest
		wantErr bool
		wantIssues int
	}{
		{
			name: "valid manifest",
			manifest: &spec.Manifest{
				Version: "1.0.0",
				Dependecies: []models.Dependency{
					{
						RepositoryURL: "https://github.com/example/repo",
						SpecPath:      "specs/test.md",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "missing entry",
			manifest: &spec.Manifest{
				Version: "1.0.0",
				Dependecies: []models.Dependency{
					{
						RepositoryURL: "https://github.com/example/repo",
						SpecPath:      "specs/test.md",
					},
				},
			},
			wantErr: true,
			wantIssues: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lockfile := spec.NewLockfile("1.0.0")
			issues, err := lockfile.Verify(tt.manifest)

			if tt.wantErr {
				if err == nil {
					t.Error("Verify() should return error")
				}
				if len(issues) != tt.wantIssues {
					t.Errorf("Verify() got %d issues, want %d", len(issues), tt.wantIssues)
				}
			} else {
				if err != nil {
					t.Errorf("Verify() error = %v", err)
				}
			}
		})
	}
}

func TestCalculateSHA256(t *testing.T) {
	data := []byte("test content")
	lockfile := spec.NewLockfile("1.0.0")
	entry := spec.LockfileEntry{
		ContentHash: spec.CalculateSHA256FromBytes(data),
	}

	if entry.ContentHash == "" {
		t.Error("CalculateSHA256() should not return empty string")
	}

	if len(entry.ContentHash) != 64 {
		t.Errorf("CalculateSHA256() got length %d, want 64", len(entry.ContentHash))
	}
}

func TestCalculateSHA256FromBytes(t *testing.T) {
	data := []byte("test content")
	hash := spec.CalculateSHA256FromBytes(data)
	expected := hex.EncodeToString(sha256.Sum256(data))

	if hash != expected {
		t.Errorf("CalculateSHA256FromBytes() = %s, want %s", hash, expected)
	}
}

// Add missing import
import (
	"os"
	"specledger/pkg/models"
)
