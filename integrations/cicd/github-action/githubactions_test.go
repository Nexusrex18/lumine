package githubaction

import (
	"os"
	"strings"
	"testing"
)

func TestHashFiles(t *testing.T) {
	// Create a temp directory with test files
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(origDir)

	// Create go.sum files with known content
	os.WriteFile("go.sum", []byte("module1 v1.0.0 h1:abc123\n"), 0644)
	os.WriteFile("subdir/go.sum", []byte("module2 v2.0.0 h1:def456\n"), 0644)

	hash, err := hashFiles("**/go.sum")
	if err != nil {
		t.Fatalf("hashFiles failed: %v", err)
	}

	if hash == "" {
		t.Error("hash should not be empty")
	}

	// Hash should be consistent
	hash2, _ := hashFiles("**/go.sum")
	if hash != hash2 {
		t.Error("hash should be deterministic")
	}

	// Hash should be hexadecimal
	for _, c := range hash {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')) {
			t.Error("hash should be hexadecimal")
			break
		}
	}
}

func TestHashFiles_NoMatch(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(origDir)

	hash, err := hashFiles("**/nonexistent.txt")
	if err != nil {
		t.Fatalf("hashFiles should not error on no match: %v", err)
	}

	if hash == "" {
		t.Error("hash should be empty string for no matches")
	}
}

func TestGenerateGitHubActionYAML(t *testing.T) {
	params := GitHubActionParams{
		WorkflowName:  "Test Workflow",
		TriggerEvents: "push",
		GoVersion:     "1.21",
		BuildCommand:  "go build",
		TestCommand:   "go test ./...",
		CacheKey:      "test-cache-key",
	}

	yaml := generateGitHubActionYAML(params)

	// Verify content contains expected elements
	if !strings.Contains(yaml, "Test Workflow") {
		t.Error("YAML missing workflow name")
	}
	if !strings.Contains(yaml, "push:") {
		t.Error("YAML missing trigger event")
	}
	if !strings.Contains(yaml, "go-version: 1.21") {
		t.Error("YAML missing go version")
	}
	if !strings.Contains(yaml, "go build") {
		t.Error("YAML missing build command")
	}
	if !strings.Contains(yaml, "go test ./...") {
		t.Error("YAML missing test command")
	}
	if !strings.Contains(yaml, "go-") {
		t.Error("YAML missing cache key prefix")
	}
	if !strings.Contains(yaml, "<runner_os>") {
		t.Error("YAML should still contain placeholder before manipulation")
	}
}

func TestManipulateYAMLToAddRunnerOS(t *testing.T) {
	input := `
key: go-<runner_os>-{{ .CacheKey }}
another: some-<runner_os>-value
`
	result := manipulateYAMLToAddRunnerOS(input)

	if strings.Contains(result, "<runner_os>") {
		t.Error("placeholder should be replaced")
	}
	if !strings.Contains(result, "${{ runner.os }}") {
		t.Error("replacement value not found")
	}

	// Verify count - input has 2 occurrences of <runner_os>
	count := strings.Count(result, "${{ runner.os }}")
	if count != 2 {
		t.Errorf("expected 2 replacements, got %d", count)
	}
}

func TestManipulateYAMLToAddRunnerOS_NoPlaceholders(t *testing.T) {
	input := "no placeholders here"
	result := manipulateYAMLToAddRunnerOS(input)

	if result != input {
		t.Error("string should remain unchanged when no placeholders")
	}
}

func TestGitHubActionParams_Struct(t *testing.T) {
	params := GitHubActionParams{
		WorkflowName:  "My Workflow",
		TriggerEvents: "pull_request",
		GoVersion:     "1.22",
		BuildCommand:  "make build",
		TestCommand:   "make test",
		CacheKey:      "abc123",
	}

	if params.WorkflowName != "My Workflow" {
		t.Error("WorkflowName not set correctly")
	}
	if params.TriggerEvents != "pull_request" {
		t.Error("TriggerEvents not set correctly")
	}
	if params.GoVersion != "1.22" {
		t.Error("GoVersion not set correctly")
	}
	if params.BuildCommand != "make build" {
		t.Error("BuildCommand not set correctly")
	}
	if params.TestCommand != "make test" {
		t.Error("TestCommand not set correctly")
	}
	if params.CacheKey != "abc123" {
		t.Error("CacheKey not set correctly")
	}
}

func TestGenerateGitHubActionYAML_BranchMain(t *testing.T) {
	params := GitHubActionParams{
		WorkflowName:  "Test",
		TriggerEvents: "push",
		GoVersion:     "1.20",
		BuildCommand:  "go build",
		TestCommand:   "go test",
		CacheKey:      "key",
	}

	yaml := generateGitHubActionYAML(params)

	// Should default to main branch
	if !strings.Contains(yaml, "- main") {
		t.Error("YAML should specify main branch")
	}
}

