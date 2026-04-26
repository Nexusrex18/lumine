package jenkins

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateJenkinsfile(t *testing.T) {
	params := JenkinsParams{
		PipelineName: "Test Pipeline",
		BranchName:   "main",
		BuildCommand: "go build",
		TestCommand:  "go test ./...",
		AgentLabel:   "linux",
	}

	jenkinsfile := generateJenkinsfile(params)

	// Verify content contains expected elements
	if !strings.Contains(jenkinsfile, "Test Pipeline") {
		t.Error("Jenkinsfile missing pipeline name")
	}
	if !strings.Contains(jenkinsfile, "main") {
		t.Error("Jenkinsfile missing branch name")
	}
	if !strings.Contains(jenkinsfile, "go build") {
		t.Error("Jenkinsfile missing build command")
	}
	if !strings.Contains(jenkinsfile, "go test ./...") {
		t.Error("Jenkinsfile missing test command")
	}
	if !strings.Contains(jenkinsfile, "linux") {
		t.Error("Jenkinsfile missing agent label")
	}
	if !strings.Contains(jenkinsfile, "pipeline {") {
		t.Error("Jenkinsfile missing pipeline block")
	}
	if !strings.Contains(jenkinsfile, "stage(") {
		t.Error("Jenkinsfile missing stages")
	}
}

func TestManipulateJenkinsfile(t *testing.T) {
	input := "some placeholder content"
	result := manipulateJenkinsfile(input)

	if !strings.Contains(result, "dynamic_value") {
		t.Error("placeholder not replaced with dynamic_value")
	}
	if strings.Contains(result, "placeholder") {
		t.Error("placeholder should be replaced")
	}
}

func TestManipulateJenkinsfile_NoPlaceholders(t *testing.T) {
	// Note: "placeholder" gets replaced with "dynamic_value"
	// So we test with a string that has NO occurrence of "placeholder"
	input := "no-match-here"
	result := manipulateJenkinsfile(input)

	if result != input {
		t.Error("content without placeholder should remain unchanged")
	}
}

func TestJenkinsParams_Struct(t *testing.T) {
	params := JenkinsParams{
		PipelineName: "My Pipeline",
		BranchName:   "develop",
		BuildCommand: "make",
		TestCommand:  "make test",
		AgentLabel:   "docker",
	}

	if params.PipelineName != "My Pipeline" {
		t.Error("PipelineName not set correctly")
	}
	if params.BranchName != "develop" {
		t.Error("BranchName not set correctly")
	}
	if params.BuildCommand != "make" {
		t.Error("BuildCommand not set correctly")
	}
	if params.TestCommand != "make test" {
		t.Error("TestCommand not set correctly")
	}
	if params.AgentLabel != "docker" {
		t.Error("AgentLabel not set correctly")
	}
}

func TestGenerateJenkinsfile_ContainsStages(t *testing.T) {
	params := JenkinsParams{
		PipelineName: "Test",
		BranchName:   "main",
		BuildCommand: "build",
		TestCommand:  "test",
		AgentLabel:   "linux",
	}

	jenkinsfile := generateJenkinsfile(params)

	// Should have Checkout, Build, Test stages
	if !strings.Contains(jenkinsfile, "stage('Checkout')") {
		t.Error("Jenkinsfile missing Checkout stage")
	}
	if !strings.Contains(jenkinsfile, "stage('Build')") {
		t.Error("Jenkinsfile missing Build stage")
	}
	if !strings.Contains(jenkinsfile, "stage('Test')") {
		t.Error("Jenkinsfile missing Test stage")
	}
}

func TestGenerateJenkinsfile_PostAlways(t *testing.T) {
	params := JenkinsParams{
		PipelineName: "Test",
		BranchName:   "main",
		BuildCommand: "build",
		TestCommand:  "test",
		AgentLabel:   "linux",
	}

	jenkinsfile := generateJenkinsfile(params)

	// Should have post always block
	if !strings.Contains(jenkinsfile, "post {") {
		t.Error("Jenkinsfile missing post block")
	}
	if !strings.Contains(jenkinsfile, "always {") {
		t.Error("Jenkinsfile missing always post condition")
	}
}

func TestGenerateJenkinsfile_CheckoutScms(t *testing.T) {
	params := JenkinsParams{
		PipelineName: "Test",
		BranchName:   "main",
		BuildCommand: "build",
		TestCommand:  "test",
		AgentLabel:   "linux",
	}

	jenkinsfile := generateJenkinsfile(params)

	// Should use checkout scm
	if !strings.Contains(jenkinsfile, "checkout scm") {
		t.Error("Jenkinsfile missing checkout scm")
	}
}

func TestWriteJenkinsfileToDisk(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "Jenkinsfile")

	// Create Jenkinsfile content
	params := JenkinsParams{
		PipelineName: "Test",
		BranchName:   "main",
		BuildCommand: "go build",
		TestCommand:  "go test",
		AgentLabel:   "docker",
	}

	content := generateJenkinsfile(params)
	err := os.WriteFile(outputFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("failed to write Jenkinsfile: %v", err)
	}

	// Verify file was written correctly
	readContent, _ := os.ReadFile(outputFile)
	if string(readContent) != content {
		t.Error("written content does not match original")
	}
}
