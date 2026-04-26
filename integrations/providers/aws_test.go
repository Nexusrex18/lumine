package providers

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateECRConfig(t *testing.T) {
	tmpDir := t.TempDir()
	outputDir := filepath.Join(tmpDir, "ecr-config")

	err := GenerateECRConfig("my-repo", outputDir)
	if err != nil {
		t.Fatalf("failed to generate ECR config: %v", err)
	}

	// Verify main.tf was created
	mainTfPath := filepath.Join(outputDir, "main.tf")
	content, err := os.ReadFile(mainTfPath)
	if err != nil {
		t.Fatalf("failed to read main.tf: %v", err)
	}

	contentStr := string(content)

	// Verify content contains expected elements
	if !strings.Contains(contentStr, `name = "my-repo"`) {
		t.Error("ECR config missing repository name")
	}
	if !strings.Contains(contentStr, "aws_ecr_repository") {
		t.Error("ECR config missing aws_ecr_repository resource")
	}
	if !strings.Contains(contentStr, "image_tag_mutability = \"MUTABLE\"") {
		t.Error("ECR config missing image_tag_mutability")
	}
	if !strings.Contains(contentStr, "scan_on_push = true") {
		t.Error("ECR config missing scan_on_push")
	}
}

func TestGenerateEKSConfig(t *testing.T) {
	tmpDir := t.TempDir()
	outputDir := filepath.Join(tmpDir, "eks-config")

	err := GenerateEKSConfig("my-cluster", "us-west-2", outputDir)
	if err != nil {
		t.Fatalf("failed to generate EKS config: %v", err)
	}

	// Verify main.tf was created
	mainTfPath := filepath.Join(outputDir, "main.tf")
	content, err := os.ReadFile(mainTfPath)
	if err != nil {
		t.Fatalf("failed to read main.tf: %v", err)
	}

	contentStr := string(content)

	// Verify content contains expected elements
	// Note: template has `name     = "{{.ClusterName}}"` with spaces
	if !strings.Contains(contentStr, `name     = "my-cluster"`) {
		t.Error("EKS config missing cluster name")
	}
	if !strings.Contains(contentStr, "aws_eks_cluster") {
		t.Error("EKS config missing aws_eks_cluster resource")
	}
	if !strings.Contains(contentStr, "vpc_config") {
		t.Error("EKS config missing vpc_config")
	}
}

func TestGenerateS3Config(t *testing.T) {
	tmpDir := t.TempDir()
	outputDir := filepath.Join(tmpDir, "s3-config")

	err := GenerateS3Config("my-bucket", outputDir)
	if err != nil {
		t.Fatalf("failed to generate S3 config: %v", err)
	}

	// Verify main.tf was created
	mainTfPath := filepath.Join(outputDir, "main.tf")
	content, err := os.ReadFile(mainTfPath)
	if err != nil {
		t.Fatalf("failed to read main.tf: %v", err)
	}

	contentStr := string(content)

	// Verify content contains expected elements
	if !strings.Contains(contentStr, `bucket = "my-bucket"`) {
		t.Error("S3 config missing bucket name")
	}
	if !strings.Contains(contentStr, "aws_s3_bucket") {
		t.Error("S3 config missing aws_s3_bucket resource")
	}
	if !strings.Contains(contentStr, `acl    = "private"`) {
		t.Error("S3 config missing private ACL")
	}
	if !strings.Contains(contentStr, "versioning") {
		t.Error("S3 config missing versioning block")
	}
	if !strings.Contains(contentStr, "server_side_encryption_configuration") {
		t.Error("S3 config missing server_side_encryption_configuration")
	}
	if !strings.Contains(contentStr, "AES256") {
		t.Error("S3 config missing AES256 encryption")
	}
}

func TestGenerateECRConfig_CreatesDirectoryIfNotExists(t *testing.T) {
	tmpDir := t.TempDir()
	outputDir := filepath.Join(tmpDir, "nested", "deeply", "ecr-config")

	err := GenerateECRConfig("test-repo", outputDir)
	if err != nil {
		t.Fatalf("failed to generate ECR config: %v", err)
	}

	mainTfPath := filepath.Join(outputDir, "main.tf")
	if _, err := os.Stat(mainTfPath); os.IsNotExist(err) {
		t.Error("main.tf was not created in nested directory")
	}
}

func TestGenerateECRConfig_OverwritesExisting(t *testing.T) {
	tmpDir := t.TempDir()
	outputDir := filepath.Join(tmpDir, "ecr-config")

	// Create directory with existing file
	os.MkdirAll(outputDir, 0755)
	existingPath := filepath.Join(outputDir, "main.tf")
	os.WriteFile(existingPath, []byte("old content"), 0644)

	err := GenerateECRConfig("new-repo", outputDir)
	if err != nil {
		t.Fatalf("failed to generate ECR config: %v", err)
	}

	content, _ := os.ReadFile(existingPath)
	if strings.Contains(string(content), "old content") {
		t.Error("existing file was not overwritten")
	}
}
