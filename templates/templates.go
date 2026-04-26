package templates

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

// Jenkinsfile template
const JenkinsfileTemplate = `
pipeline {
    agent any

    stages {
        stage('Checkout') {
            steps {
                git '{{ .RepoURL }}'
            }
        }
        stage('Build') {
            steps {
                {{- if .IsGoProject }}
                sh 'go mod tidy && go build ./...'
                {{- else }}
                sh 'npm ci && npm run build'
                {{- end }}
            }
        }
        stage('Test') {
            steps {
                {{- if .IsGoProject }}
                sh 'go test -v ./...'
                {{- else }}
                sh 'npm test'
                {{- end }}
            }
        }
    }

    post {
        success {
            echo 'Pipeline succeeded!'
        }
        failure {
            echo 'Pipeline failed. Check the console output.'
        }
    }
}
`

// GitHub Actions Template for Go Project
const GitHubActionsGoTemplate = `
name: {{ .WorkflowName }}

on:
  {{ .TriggerEvents }}:
    branches:
      - main

jobs:
  build:
    runs-on: ubuntu-latest  # Static value for now

    steps:
      - name: Checkout code
        uses: actions/checkout@v2

      - name: Set up Go {{ .GoVersion }}
        uses: actions/setup-go@v3
        with:
          go-version: {{ .GoVersion }}

      - name: Cache Go Modules
        uses: actions/cache@v2
        with:
          path: ~/go/pkg/mod
          key: go-<runner_os>-go-{{ .CacheKey }}
          restore-keys: |
            go-<runner_os>-

      - name: Install dependencies
        run: go mod tidy

      - name: Build Project
        run: {{ .BuildCommand }}

      - name: Run Tests
        run: {{ .TestCommand }}
`

// GitHub Actions Template for Node.js Project
const GitHubActionsNodeTemplate = `
name: {{ .WorkflowName }}

on:
  {{ .TriggerEvents }}:
    branches:
      - main

jobs:
  build:
    runs-on: ubuntu-latest  # Static value for now

    steps:
      - name: Checkout code
        uses: actions/checkout@v2

      - name: Set up Node.js {{ .NodeVersion }}
        uses: actions/setup-node@v3
        with:
          node-version: {{ .NodeVersion }}

      - name: Cache Node Modules
        uses: actions/cache@v2
        with:
          path: ~/.npm
          key: node-<runner_os>-npm-{{ .CacheKey }}
          restore-keys: |
            node-<runner_os>-

      - name: Install dependencies
        run: npm install

      - name: Build Project
        run: npm run build

      - name: Run Tests
        run: npm test
`

// ReplaceRunnerOSInFile replaces the <runner_os> placeholder with the actual value of GitHub Actions variable
func ReplaceRunnerOSInFile(filePath string) error {
	// Read the content of the file
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading the file: %w", err)
	}

	// Convert the content to a string
	fileContent := string(content)

	// Replace <runner_os> with GitHub Actions runner OS variable
	updatedContent := strings.ReplaceAll(fileContent, "<runner_os>", "${{ runner.os }}")

	// Write the updated content back to the same file
	err = os.WriteFile(filePath, []byte(updatedContent), 0644)
	if err != nil {
		return fmt.Errorf("error writing to the file: %w", err)
	}

	fmt.Printf("Replaced <runner_os> with ${{ runner.os }} in %s\n", filePath)
	return nil
}

// GenerateGitHubActions generates the GitHub Actions workflow file dynamically for Go or Node.js project
func GenerateGitHubActions(workflowName, triggerEvents, goVersion, nodeVersion, buildCmd, testCmd, cacheKey, repoURL string, isGoProject bool) error {
	var tmpl *template.Template
	var err error

	if isGoProject {
		// Parse Go-specific GitHub Actions template
		tmpl, err = template.New("github-actions-go").Parse(GitHubActionsGoTemplate)
	} else {
		// Parse Node.js-specific GitHub Actions template
		tmpl, err = template.New("github-actions-node").Parse(GitHubActionsNodeTemplate)
	}

	if err != nil {
		return fmt.Errorf("error parsing GitHub Actions template: %w", err)
	}

	// Data to pass to the template
	data := map[string]interface{}{
		"WorkflowName":  workflowName,
		"TriggerEvents": triggerEvents,
		"CacheKey":      cacheKey,
		"BuildCommand":  buildCmd,
		"TestCommand":   testCmd,
		"RepoURL":       repoURL,
	}

	if isGoProject {
		data["GoVersion"] = goVersion
	} else {
		data["NodeVersion"] = nodeVersion
	}

	// Create the ci.yml file in the current directory (not .github directory)
	fileName := "ci.yml"
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("error creating GitHub Actions file: %w", err)
	}
	defer file.Close()

	// Write the processed template to the file
	err = tmpl.Execute(file, data)
	if err != nil {
		return fmt.Errorf("error executing GitHub Actions template: %w", err)
	}
	fmt.Printf("GitHub Actions workflow file (%s) generated successfully\n", fileName)

	// Call function to replace <runner_os> with ${{ runner.os }} in the YAML file
	err = ReplaceRunnerOSInFile(fileName)
	if err != nil {
		return fmt.Errorf("error replacing <runner_os> with ${{ runner.os }}: %w", err)
	}

	return nil
}

// GenerateJenkinsfile generates the Jenkinsfile dynamically for Go or Node.js project
func GenerateJenkinsfile(repoURL string, isGoProject bool) error {
	// Parse the Jenkinsfile template
	tmpl, err := template.New("jenkinsfile").Parse(JenkinsfileTemplate)
	if err != nil {
		return fmt.Errorf("error parsing Jenkinsfile template: %w", err)
	}

	// Data to pass to the template
	data := map[string]interface{}{
		"RepoURL":     repoURL,
		"IsGoProject": isGoProject,
	}

	// Create the Jenkinsfile in the current directory
	fileName := "Jenkinsfile"
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("error creating Jenkinsfile: %w", err)
	}
	defer file.Close()

	// Write the processed template to the file
	err = tmpl.Execute(file, data)
	if err != nil {
		return fmt.Errorf("error executing Jenkinsfile template: %w", err)
	}

	fmt.Printf("Jenkinsfile generated successfully: %s\n", fileName)
	return nil
}

//aws

const eksTemplate = `
resource "aws_eks_cluster" "{{.ClusterName}}" {
  name     = "{{.ClusterName}}"
  role_arn = "arn:aws:iam::your_account_id:role/eks-service-role"

  vpc_config {
    subnet_ids = ["subnet-0123456789abcdef0", "subnet-abcdef0123456789"]
  }
}
`

const ecrTemplate = `
resource "aws_ecr_repository" "{{.Name}}" {
  name = "{{.Name}}"
  image_tag_mutability = "MUTABLE"
  image_scanning_configuration {
    scan_on_push = true
  }
}
`

const s3Template = `
resource "aws_s3_bucket" "{{.BucketName}}" {
  bucket = "{{.BucketName}}"
  acl    = "private"

  versioning {
    enabled = true
  }

  server_side_encryption_configuration {
    rule {
      apply_server_side_encryption_by_default {
        sse_algorithm = "AES256"
      }
    }
  }
}
`

// GenerateECRConfig generates a Terraform configuration file for AWS ECR in the current working directory
func GenerateECRConfig(name string) error {
	tmpl, err := template.New("ecr").Parse(ecrTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse ECR template: %w", err)
	}

	data := map[string]string{"Name": name}

	// Get current working directory
	outputDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Path for the new file
	filePath := filepath.Join(outputDir, "main.tf")

	// Create or overwrite the file
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	return tmpl.Execute(file, data)
}

// GenerateEKSConfig generates a Terraform configuration file for AWS EKS in the current working directory
func GenerateEKSConfig(clusterName, region string) error {
	tmpl, err := template.New("eks").Parse(eksTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse EKS template: %w", err)
	}

	data := map[string]string{
		"ClusterName": clusterName,
		"Region":      region,
	}

	// Get current working directory
	outputDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Path for the new file
	filePath := filepath.Join(outputDir, "main.tf")

	// Create or overwrite the file
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	return tmpl.Execute(file, data)
}

// GenerateS3Config generates a Terraform configuration file for AWS S3 in the current working directory
func GenerateS3Config(bucketName string) error {
	tmpl, err := template.New("s3").Parse(s3Template)
	if err != nil {
		return fmt.Errorf("failed to parse S3 template: %w", err)
	}

	data := map[string]string{"BucketName": bucketName}

	// Get current working directory
	outputDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Path for the new file
	filePath := filepath.Join(outputDir, "main.tf")

	// Create or overwrite the file
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	return tmpl.Execute(file, data)
}

// Docker Compose template
const dockerComposeTemplate = `
version: '3'
services:
  prometheus:
    image: prom/prometheus
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
    ports:
      - "9090:9090"
`

// Prometheus config template
const prometheusConfigTemplate = `
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'node_exporter'
    static_configs:
      - targets: ['localhost:9100']
`

// GenerateDockerCompose will append to docker-compose.yml if it exists, or create it if it doesn't.
func GenerateDockerCompose() error {
	// Open the docker-compose.yml file (append mode if it exists)
	file, err := os.OpenFile("docker-compose.yml", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening docker-compose.yml: %w", err)
	}
	defer file.Close()

	// Parse the template
	tmpl, err := template.New("dockerCompose").Parse(dockerComposeTemplate)
	if err != nil {
		return fmt.Errorf("error parsing docker-compose template: %w", err)
	}

	// Write the content to the file
	return tmpl.Execute(file, nil)
}

// GeneratePrometheusConfig creates prometheus.yml if it doesn't exist.
func GeneratePrometheusConfig() error {
	// Create the file if it doesn't exist
	file, err := os.Create("prometheus.yml")
	if err != nil {
		return fmt.Errorf("error creating prometheus.yml: %w", err)
	}
	defer file.Close()

	// Parse the template for Prometheus config
	tmpl, err := template.New("prometheusConfig").Parse(prometheusConfigTemplate)
	if err != nil {
		return fmt.Errorf("error parsing prometheus config template: %w", err)
	}

	// Write the Prometheus configuration to the file
	return tmpl.Execute(file, nil)
}

// StartPrometheus starts Prometheus using Docker Compose.
func StartPrometheus() error {
	cmd := exec.Command("docker-compose", "up", "-d")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// StartNodeExporter starts the Node Exporter.
func StartNodeExporter() error {
	cmd := exec.Command("docker", "run", "-d", "-p", "9100:9100", "prom/node-exporter")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// StatusMessage displays the URLs where Prometheus and Node Exporter can be accessed.
func StatusMessage() {
	fmt.Println("Prometheus running at http://localhost:9090")
	fmt.Println("Node Exporter running at http://localhost:9100")
}

// SetupPrometheusMonitoring sets up Docker Compose and Prometheus monitoring.
func SetupPrometheusMonitoring() error {
	// Generate Docker Compose file
	err := GenerateDockerCompose()
	if err != nil {
		return fmt.Errorf("failed to generate docker-compose.yml: %v", err)
	}

	// Generate Prometheus config file
	err = GeneratePrometheusConfig()
	if err != nil {
		return fmt.Errorf("failed to generate prometheus.yml: %v", err)
	}

	// Start Prometheus and Node Exporter services
	err = StartPrometheus()
	if err != nil {
		return fmt.Errorf("failed to start Prometheus: %v", err)
	}

	err = StartNodeExporter()
	if err != nil {
		return fmt.Errorf("failed to start Node Exporter: %v", err)
	}

	// Print the status message
	StatusMessage()

	return nil
}
