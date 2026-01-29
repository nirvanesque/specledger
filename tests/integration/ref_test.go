package ref_test

import (
	"os"
	"strings"
	"testing"
	"specledger/internal/ref"
	"specledger/internal/spec"
)

func TestParseSpecMarkdownLinks(t *testing.T) {
	content := `
# Test Document

See [GitHub Repository](https://github.com/example/repo) for more info.
Also check out [Another Link](https://github.com/other/repo/blob/main/spec.md).

![Image](https://example.com/image.png)
`

	resolver := ref.NewResolver("test.lock")
	references, err := resolver.ParseSpec(content)

	if err != nil {
		t.Fatalf("ParseSpec() error = %v", err)
	}

	// Should find 3 markdown links
	markdownLinks := 0
	for _, ref := range references {
		if ref.Type == "markdown" {
			markdownLinks++
			if ref.URL != "https://github.com/example/repo" && ref.URL != "https://github.com/other/repo/blob/main/spec.md" {
				t.Errorf("ParseSpec() got URL %s, want one of the expected URLs", ref.URL)
			}
		}
	}

	if markdownLinks != 2 {
		t.Errorf("ParseSpec() found %d markdown links, want 2", markdownLinks)
	}
}

func TestParseSpecInlineReferences(t *testing.T) {
	content := `
# Test Document

This references spec[main#section].
Also uses spec.release#api-endpoints.
`

	resolver := ref.NewResolver("test.lock")
	references, err := resolver.ParseSpec(content)

	if err != nil {
		t.Fatalf("ParseSpec() error = %v", err)
	}

	inlineRefs := 0
	for _, ref := range references {
		if ref.Type == "inline" {
			inlineRefs++
			if !strings.Contains(ref.URL, "spec.") {
				t.Errorf("ParseSpec() got URL %s, should start with 'spec.'", ref.URL)
			}
		}
	}

	if inlineRefs != 2 {
		t.Errorf("ParseSpec() found %d inline references, want 2", inlineRefs)
	}
}

func TestValidateReferences(t *testing.T) {
	// Create lockfile with test entries
	lockfile := spec.NewLockfile("1.0.0")
	lockfile.AddEntry(spec.LockfileEntry{
		RepositoryURL: "https://github.com/example/repo",
		Branch:        "#main",
		SpecPath:      "specs/test.md",
	})
	lockfile.AddEntry(spec.LockfileEntry{
		RepositoryURL: "https://github.com/other/repo",
		Branch:        "#release",
		SpecPath:      "specs/other.md",
	})

	// Build dependencies map
	dependencies := make(map[string]string)
	for _, entry := range lockfile.Entries {
		alias := entry.Branch
		if alias != "" {
			dependencies[alias] = entry.RepositoryURL
		}
	}

	resolver := ref.NewResolver("test.lock")
	resolver.Dependencies = dependencies

	content := `
# Test Document

References: spec.main#section
And: spec.release#api
`

	references, _ := resolver.ParseSpec(content)

	errors := resolver.ValidateReferences(references)

	if len(errors) > 0 {
		t.Errorf("ValidateReferences() got %d errors, want 0", len(errors))
		for _, err := range errors {
			t.Errorf("Error: %s", err.Error())
		}
	}
}

func TestValidateReferencesUnknown(t *testing.T) {
	dependencies := map[string]string{
		"main": "https://github.com/example/repo",
	}

	resolver := ref.NewResolver("test.lock")
	resolver.Dependencies = dependencies

	content := `
# Test Document

References: spec.main#section
And: spec.unknown#section
`

	references, _ := resolver.ParseSpec(content)

	errors := resolver.ValidateReferences(references)

	if len(errors) != 1 {
		t.Errorf("ValidateReferences() got %d errors, want 1", len(errors))
	}

	if errors[0].Field != "url" {
		t.Errorf("ValidateReferences() error field = %s, want 'url'", errors[0].Field)
	}

	if !strings.Contains(errors[0].Message, "unknown dependency") {
		t.Errorf("ValidateReferences() error message = %s, want 'unknown dependency'", errors[0].Message)
	}
}

func TestValidateReferencesImageSkipped(t *testing.T) {
	dependencies := map[string]string{
		"main": "https://github.com/example/repo",
	}

	resolver := ref.NewResolver("test.lock")
	resolver.Dependencies = dependencies

	content := `
# Test Document

Check this image: ![Logo](https://example.com/logo.png)
`

	references, _ := resolver.ParseSpec(content)

	errors := resolver.ValidateReferences(references)

	// Images should be skipped
	if len(errors) != 0 {
		t.Errorf("ValidateReferences() got %d errors for images (should be skipped)", len(errors))
	}
}

func TestReferenceResolver(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantType string
	}{
		{
			name:    "markdown link",
			input:   "[Text](https://github.com/example/repo)",
			wantType: "markdown",
		},
		{
			name:    "image",
			input:   "![](https://example.com/image.png)",
			wantType: "image",
		},
		{
			name:    "inline reference",
			input:   "spec.main#section",
			wantType: "inline",
		},
		{
			name:    "inline reference with section",
			input:   "spec.release#api-endpoints",
			wantType: "inline",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resolver := ref.NewResolver("test.lock")
			references, err := resolver.ParseSpec(tt.input)

			if err != nil {
				t.Fatalf("ParseSpec() error = %v", err)
			}

			if len(references) == 0 {
				t.Errorf("ParseSpec() got no references, want at least 1")
				return
			}

			if references[0].Type != tt.wantType {
				t.Errorf("ParseSpec() got type %s, want %s", references[0].Type, tt.wantType)
			}
		})
	}
}
