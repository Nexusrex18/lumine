package docker

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetectProjectType_GoProject(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a go.mod file
	err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte("module test"), 0644)
	if err != nil {
		t.Fatalf("failed to create go.mod: %v", err)
	}

	projectType, err := DetectProjectType(tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if projectType != "go" {
		t.Errorf("expected 'go', got '%s'", projectType)
	}
}

func TestDetectProjectType_NodeProject(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a package.json file
	err := os.WriteFile(filepath.Join(tmpDir, "package.json"), []byte("{}"), 0644)
	if err != nil {
		t.Fatalf("failed to create package.json: %v", err)
	}

	projectType, err := DetectProjectType(tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if projectType != "node" {
		t.Errorf("expected 'node', got '%s'", projectType)
	}
}

func TestDetectProjectType_GoTakesPrecedence(t *testing.T) {
	tmpDir := t.TempDir()

	// Create both files - go.mod should take precedence
	err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte("module test"), 0644)
	if err != nil {
		t.Fatalf("failed to create go.mod: %v", err)
	}
	err = os.WriteFile(filepath.Join(tmpDir, "package.json"), []byte("{}"), 0644)
	if err != nil {
		t.Fatalf("failed to create package.json: %v", err)
	}

	projectType, err := DetectProjectType(tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if projectType != "go" {
		t.Errorf("expected 'go' to take precedence, got '%s'", projectType)
	}
}

func TestDetectProjectType_NotFound(t *testing.T) {
	tmpDir := t.TempDir()

	_, err := DetectProjectType(tmpDir)
	if err == nil {
		t.Error("expected error when no project type found")
	}
}

func TestDetectProjectType_NonExistentPath(t *testing.T) {
	_, err := DetectProjectType("/nonexistent/path/12345")
	if err == nil {
		t.Error("expected error for non-existent path")
	}
}

func TestGenerateGoDockerfile(t *testing.T) {
	tmpDir := t.TempDir()

	err := generateGoDockerfile(tmpDir)
	if err != nil {
		t.Fatalf("failed to generate Dockerfile: %v", err)
	}

	// Verify Dockerfile was created
	dockerfilePath := filepath.Join(tmpDir, "Dockerfile")
	content, err := os.ReadFile(dockerfilePath)
	if err != nil {
		t.Fatalf("failed to read Dockerfile: %v", err)
	}

	// Verify content contains expected elements
	if !contains(string(content), "FROM golang:1.22") {
		t.Error("Dockerfile missing golang base image")
	}
	if !contains(string(content), "FROM alpine:latest") {
		t.Error("Dockerfile missing alpine runtime stage")
	}
	if !contains(string(content), "go mod tidy") {
		t.Error("Dockerfile missing go mod tidy")
	}
	if !contains(string(content), "go build") {
		t.Error("Dockerfile missing go build")
	}
}

func TestGenerateNodeDockerfile(t *testing.T) {
	tmpDir := t.TempDir()

	err := generateNodeDockerfile(tmpDir)
	if err != nil {
		t.Fatalf("failed to generate Dockerfile: %v", err)
	}

	// Verify Dockerfile was created
	dockerfilePath := filepath.Join(tmpDir, "Dockerfile")
	content, err := os.ReadFile(dockerfilePath)
	if err != nil {
		t.Fatalf("failed to read Dockerfile: %v", err)
	}

	// Verify content contains expected elements
	if !contains(string(content), "FROM node:20-alpine") {
		t.Error("Dockerfile missing node base image")
	}
	if !contains(string(content), "npm ci") {
		t.Error("Dockerfile missing npm ci")
	}
	if !contains(string(content), "npm run build") {
		t.Error("Dockerfile missing npm run build")
	}
	if !contains(string(content), "EXPOSE 3000") {
		t.Error("Dockerfile missing EXPOSE 3000")
	}
}

func TestGenerateDockerfile_GoProject(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a go.mod file
	err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte("module test"), 0644)
	if err != nil {
		t.Fatalf("failed to create go.mod: %v", err)
	}

	err = GenerateDockerfile(tmpDir)
	if err != nil {
		t.Fatalf("failed to generate Dockerfile: %v", err)
	}

	dockerfilePath := filepath.Join(tmpDir, "Dockerfile")
	if _, err := os.Stat(dockerfilePath); os.IsNotExist(err) {
		t.Error("Dockerfile was not created")
	}
}

func TestGenerateDockerfile_NodeProject(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a package.json file
	err := os.WriteFile(filepath.Join(tmpDir, "package.json"), []byte("{}"), 0644)
	if err != nil {
		t.Fatalf("failed to create package.json: %v", err)
	}

	err = GenerateDockerfile(tmpDir)
	if err != nil {
		t.Fatalf("failed to generate Dockerfile: %v", err)
	}

	dockerfilePath := filepath.Join(tmpDir, "Dockerfile")
	if _, err := os.Stat(dockerfilePath); os.IsNotExist(err) {
		t.Error("Dockerfile was not created")
	}
}

func TestGenerateDockerfile_UnsupportedProject(t *testing.T) {
	tmpDir := t.TempDir()

	// Create an empty directory with no project files
	err := GenerateDockerfile(tmpDir)
	if err == nil {
		t.Error("expected error for unsupported project type")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
