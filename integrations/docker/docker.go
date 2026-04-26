package docker

import (
	"fmt"
	"os"
	"path/filepath"
)

// DetectProjectType checks for go.mod or package.json to determine the project language
func DetectProjectType(projectPath string) (string, error) {
	if _, err := os.Stat(filepath.Join(projectPath, "go.mod")); err == nil {
		return "go", nil
	}
	if _, err := os.Stat(filepath.Join(projectPath, "package.json")); err == nil {
		return "node", nil
	}
	return "", fmt.Errorf("could not detect project type: no go.mod or package.json found")
}

func GenerateDockerfile(projectPath string) error {
	projectType, err := DetectProjectType(projectPath)
	if err != nil {
		return err
	}

	switch projectType {
	case "go":
		return generateGoDockerfile(projectPath)
	case "node":
		return generateNodeDockerfile(projectPath)
	default:
		return fmt.Errorf("unsupported project type: %s", projectType)
	}
}

func generateGoDockerfile(projectPath string) error {
	// Dockerfile content
	dockerfileContent := `# Stage 1: Build
FROM golang:1.22 AS builder
WORKDIR /app
COPY . .
RUN go mod tidy && go build -o app

# Stage 2: Run
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/app .
CMD ["./app"]
`

	// Ensure the directory exists
	err := os.MkdirAll(projectPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create directory %s: %w", projectPath, err)
	}

	// Write the Dockerfile
	dockerfilePath := filepath.Join(projectPath, "Dockerfile")
	file, err := os.Create(dockerfilePath)
	if err != nil {
		return fmt.Errorf("failed to create Dockerfile: %w", err)
	}
	defer file.Close()

	_, err = file.WriteString(dockerfileContent)
	if err != nil {
		return fmt.Errorf("failed to write Dockerfile content: %w", err)
	}

	return nil
}

func generateNodeDockerfile(projectPath string) error {
	dockerfileContent := `# Stage 1: Build
FROM node:20-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

# Stage 2: Run
FROM node:20-alpine
WORKDIR /app
COPY --from=builder /app/dist ./dist
COPY --from=builder /app/node_modules ./node_modules
COPY --from=builder /app/package.json ./
EXPOSE 3000
CMD ["npm", "start"]
`

	err := os.MkdirAll(projectPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create directory %s: %w", projectPath, err)
	}

	dockerfilePath := filepath.Join(projectPath, "Dockerfile")
	file, err := os.Create(dockerfilePath)
	if err != nil {
		return fmt.Errorf("failed to create Dockerfile: %w", err)
	}
	defer file.Close()

	_, err = file.WriteString(dockerfileContent)
	if err != nil {
		return fmt.Errorf("failed to write Dockerfile content: %w", err)
	}

	return nil
}
