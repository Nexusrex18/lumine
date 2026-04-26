package monitoring

import (
	"os"
	"strings"
	"testing"
)

func TestGenerateDockerCompose(t *testing.T) {
	// Use temp directory to avoid polluting current directory
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(origDir)

	err := GenerateDockerCompose()
	if err != nil {
		t.Fatalf("failed to generate docker-compose.yml: %v", err)
	}

	// Verify file was created
	content, err := os.ReadFile("docker-compose.yml")
	if err != nil {
		t.Fatalf("failed to read docker-compose.yml: %v", err)
	}

	contentStr := string(content)

	// Verify content
	if !strings.Contains(contentStr, "prom/prometheus") {
		t.Error("docker-compose.yml missing prometheus image")
	}
	if !strings.Contains(contentStr, "9090:9090") {
		t.Error("docker-compose.yml missing prometheus port mapping")
	}
	if !strings.Contains(contentStr, "prometheus.yml") {
		t.Error("docker-compose.yml missing prometheus config volume mount")
	}
}

func TestGeneratePrometheusConfig(t *testing.T) {
	// Use temp directory
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(origDir)

	err := GeneratePrometheusConfig()
	if err != nil {
		t.Fatalf("failed to generate prometheus.yml: %v", err)
	}

	// Verify file was created
	content, err := os.ReadFile("prometheus.yml")
	if err != nil {
		t.Fatalf("failed to read prometheus.yml: %v", err)
	}

	contentStr := string(content)

	// Verify content
	if !strings.Contains(contentStr, "scrape_interval: 15s") {
		t.Error("prometheus.yml missing scrape_interval")
	}
	if !strings.Contains(contentStr, "node_exporter") {
		t.Error("prometheus.yml missing node_exporter job")
	}
	if !strings.Contains(contentStr, "localhost:9100") {
		t.Error("prometheus.yml missing node exporter target")
	}
}

func TestStatusMessage(t *testing.T) {
	// StatusMessage just prints - we can at least verify it doesn't panic
	StatusMessage()
}

func TestGenerateDockerCompose_AppendsToExisting(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(origDir)

	// Create existing docker-compose.yml
	os.WriteFile("docker-compose.yml", []byte("version: '3'\nservices:\n  existing: true\n"), 0644)

	err := GenerateDockerCompose()
	if err != nil {
		t.Fatalf("failed to generate docker-compose.yml: %v", err)
	}

	content, _ := os.ReadFile("docker-compose.yml")
	if !strings.Contains(string(content), "existing: true") {
		t.Error("existing content was not preserved")
	}
}
