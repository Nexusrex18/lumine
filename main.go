package main

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"lumine/configs"

	"log"
	"lumine/integrations/docker"
	"lumine/integrations/monitoring"
	"lumine/integrations/providers"
	"lumine/templates"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {

	const (
		Logo = `
██╗     ██╗   ██╗███╗   ███╗██╗███╗   ██╗███████╗
██║     ██║   ██║████╗ ████║██║████╗  ██║██╔════╝
██║     ██║   ██║██╔████╔██║██║██╔██╗ ██║█████╗  
██║     ██║   ██║██║╚██╔╝██║██║██║╚██╗██║██╔══╝  
███████╗╚██████╔╝██║ ╚═╝ ██║██║██║ ╚████║███████╗
╚══════╝ ╚═════╝ ╚═╝     ╚═╝╚═╝╚═╝  ╚═══╝╚══════╝
                                                
	`
	)
	fmt.Println(Logo)

	for {
		choice, err := configs.MainMenu()
		if err != nil {
			fmt.Println(configs.FormatError(err))
			break
		}

		switch choice {
		case "Setup CI/CD":
			fmt.Println("Setting up CI/CD...")

		case "Cloud providers":
			cloudChoice, err := configs.AWSMenu()
			if err != nil {
				fmt.Println(configs.FormatError(err))
				break
			}

			switch cloudChoice {
			case "ECR":
				prompt := promptui.Prompt{
					Label: "Enter the name of your ECR repository",
				}

				name, _ := prompt.Run()
				prompt2 := promptui.Prompt{
					Label: "Enter address :",
				}
				dirname, _ := prompt2.Run()

				err := providers.GenerateECRConfig(name, dirname)
				if err != nil {
					fmt.Println("Error generating ECR config:", err)
				} else {
					fmt.Println("Successfully generated ECR Terraform config at ./outputs/aws/ecr")
				}

			case "S3":
				prompt := promptui.Prompt{
					Label: "Enter the name of your S3 bucket",
				}
				bucketName, _ := prompt.Run()

				prompt2 := promptui.Prompt{
					Label: "Enter address :",
				}

				dirname, _ := prompt2.Run()

				err := providers.GenerateS3Config(bucketName, dirname)
				if err != nil {
					fmt.Println("Error generating S3 config:", err)
				} else {
					fmt.Println("Successfully generated S3 Terraform config at ./outputs/aws/s3")
				}

			case "EKS":
				prompt := promptui.Prompt{
					Label: "Enter the name of your EKS cluster",
				}
				clusterName, _ := prompt.Run()

				serverName, _ := configs.AWSServerMenu()

				prompt2 := promptui.Prompt{
					Label: "Enter the directory you want the terraform configs",
				}

				dirname, _ := prompt2.Run()

				err := providers.GenerateEKSConfig(clusterName, serverName, dirname)
				if err != nil {
					fmt.Println("Error generating EKS config:", err)
				} else {
					fmt.Println("Successfully generated EKS Terraform config at ./outputs/aws/eks")
				}
			}

		case "Setup Monitoring":
			monitoringChoice, err := configs.InputPrompt("Would you like to set up Prometheus monitoring? (y/n)")
			if err != nil {
				fmt.Println(configs.FormatError(err))
				break
			}

			if monitoringChoice == "y" || monitoringChoice == "Y" {
				err := monitoring.SetupPrometheusMonitoring()
				if err != nil {
					fmt.Println("Error setting up Prometheus:", err)
				} else {
					fmt.Println("Prometheus setup complete!")
				}
			} else {
				fmt.Println("Skipping monitoring setup.")
			}
		case "Generate Dockerfile":
			// Prompt the user for the root directory
			rootDir, err := configs.InputPrompt("Enter the root directory for the Dockerfile")
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				continue
			}

			// Generate the Dockerfile in the specified directory
			err = docker.GenerateDockerfile(rootDir)
			if err != nil {
				fmt.Printf("Failed to generate Dockerfile: %v\n", err)
			} else {
				fmt.Printf("Dockerfile successfully generated in: %s\n", rootDir)
			}
		case "Exit":
			return
		case "Fetch from github repo":
			if err := CheckAndCreateFiles(); err != nil {
				log.Fatal(err)
			}
			if err := templates.SetupPrometheusMonitoring(); err != nil {
				log.Fatal(err)
			} else {
				fmt.Println("Prometheus Monitoring has been set up successfully!")
			}

		default:
			fmt.Println("Invalid choice, returning to main menu...")
		}
	}
}

// FindProjectFile searches for a go.mod or package.json file in the repository
func FindProjectFile() (string, error) {
	var foundFile string
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(path, "package.json") || strings.HasSuffix(path, "go.mod") {
			foundFile = path
			return fmt.Errorf("file found") // Stop the walk once we find the file
		}
		return nil
	})
	if err != nil && err.Error() != "file found" {
		return "", err
	}
	return foundFile, nil
}

// CheckAndCreateFiles clones the repo, detects project type, and generates the GitHub Actions workflow and Jenkinsfile
func CheckAndCreateFiles() error {
	// Prompt the user for the repository URL
	prompt := promptui.Prompt{
		Label: "Enter GitHub Repository URL (e.g. https://github.com/owner/repo)",
	}
	repoURL, err := prompt.Run()
	if err != nil {
		return fmt.Errorf("error reading repository URL: %v", err)
	}

	// Extract the repository name from the URL
	parts := strings.Split(repoURL, "/")
	repoName := strings.Split(parts[len(parts)-1], ".")[0]

	// Clone the repository into the current directory
	fmt.Println("Cloning repository:", repoURL)
	cmd := exec.Command("git", "clone", repoURL)
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("error cloning repository: %v", err)
	}

	// Change directory into the cloned repo
	err = os.Chdir(repoName)
	if err != nil {
		return fmt.Errorf("error changing directory into cloned repo: %v", err)
	}

	// Detect the project type by finding package.json or go.mod
	projectFile, err := FindProjectFile()
	if err != nil {
		return fmt.Errorf("error finding project file: %v", err)
	}

	if projectFile == "" {
		return fmt.Errorf("neither package.json nor go.mod found in repository")
	}

	// Detect project type
	var buildCmd, testCmd, cacheKey, goVersion, nodeVersion string
	var isGoProject bool
	if strings.HasSuffix(projectFile, "package.json") {
		// Node.js project
		buildCmd = "npm run build"
		testCmd = "npm test"
		cacheKey = "node-modules"
		nodeVersion = "16" // Default Node.js version
	} else if strings.HasSuffix(projectFile, "go.mod") {
		// Go project
		buildCmd = "go build ./..."
		testCmd = "go test ./..."
		cacheKey = "go-modules"
		goVersion = "1.18" // Default Go version
		isGoProject = true
	}

	// Go back to the original directory
	err = os.Chdir("..")
	if err != nil {
		return fmt.Errorf("error changing back to the original directory: %v", err)
	}

	// Generate the GitHub Actions workflow in the current directory
	fmt.Println("Generating GitHub Actions workflow")
	if err := templates.GenerateGitHubActions("CI Workflow", "push", goVersion, nodeVersion, buildCmd, testCmd, cacheKey, repoURL, isGoProject); err != nil {
		return fmt.Errorf("error generating GitHub Actions workflow: %v", err)
	}

	// Generate Jenkinsfile
	fmt.Println("Generating Jenkinsfile")
	if err := templates.GenerateJenkinsfile(repoURL, isGoProject); err != nil {
		return fmt.Errorf("error generating Jenkinsfile: %v", err)
	}

	// Ask the user if they want to generate AWS configurations
	awsPrompt := promptui.Select{
		Label: "Do you want to generate AWS resource configurations (ECR, EKS, S3)?",
		Items: []string{"No", "Yes"},
	}
	_, awsChoice, err := awsPrompt.Run()
	if err != nil {
		return fmt.Errorf("error reading AWS choice: %v", err)
	}

	// If the user chooses 'Yes', prompt for AWS resource type
	if awsChoice == "Yes" {
		// Prompt for AWS resource type (ECR, EKS, S3)
		resourcePrompt := promptui.Prompt{
			Label: "Enter AWS resource type (ecr, eks, s3)",
		}
		resourceType, err := resourcePrompt.Run()
		if err != nil {
			return fmt.Errorf("error reading resource type: %v", err)
		}

		var resourceName string
		switch resourceType {
		case "ecr":
			// Prompt for ECR repository name
			resourcePrompt := promptui.Prompt{
				Label: "Enter the ECR repository name",
			}
			resourceName, err = resourcePrompt.Run()
			if err != nil {
				return fmt.Errorf("error reading repository name: %v", err)
			}

			// Generate ECR configuration
			if err := templates.GenerateECRConfig(resourceName); err != nil {
				return fmt.Errorf("error generating ECR config: %v", err)
			}
			fmt.Println("ECR config generated successfully.")

		case "eks":
			// Prompt for EKS cluster name and region
			clusterPrompt := promptui.Prompt{
				Label: "Enter the EKS cluster name",
			}
			clusterName, err := clusterPrompt.Run()
			if err != nil {
				return fmt.Errorf("error reading cluster name: %v", err)
			}

			regionPrompt := promptui.Prompt{
				Label: "Enter the AWS region",
			}
			region, err := regionPrompt.Run()
			if err != nil {
				return fmt.Errorf("error reading region: %v", err)
			}

			// Generate EKS configuration
			if err := templates.GenerateEKSConfig(clusterName, region); err != nil {
				return fmt.Errorf("error generating EKS config: %v", err)
			}
			fmt.Println("EKS config generated successfully.")

		case "s3":
			// Prompt for S3 bucket name
			bucketPrompt := promptui.Prompt{
				Label: "Enter the S3 bucket name",
			}
			bucketName, err := bucketPrompt.Run()
			if err != nil {
				return fmt.Errorf("error reading bucket name: %v", err)
			}

			// Generate S3 configuration
			if err := templates.GenerateS3Config(bucketName); err != nil {
				return fmt.Errorf("error generating S3 config: %v", err)
			}
			fmt.Println("S3 config generated successfully.")

		default:
			return fmt.Errorf("Invalid AWS resource type. Please enter 'ecr', 'eks', or 's3'.")
		}
	}

	// Clean up: Delete the cloned repository
	err = os.RemoveAll(repoName)
	if err != nil {
		return fmt.Errorf("error deleting cloned repository: %v", err)
	}

	fmt.Println("GitHub Actions workflow, Jenkinsfile, and AWS resource configurations (if selected) created successfully and local repository cleaned up.")
	return nil
}
