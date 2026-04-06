# Lumine - The DaaS (Devops-as-a-service) CLI

Lumine is a CLI tool that simplifies the process of setting up DevOps automation for startups and developers. With Lumine, you can easily configure CI/CD pipelines, manage cloud infrastructure (AWS, ECR, S3, EKS), and estimate costs—all without the need for a dedicated DevOps team. The tool integrates with cloud providers, CI/CD services, and monitoring tools, allowing developers to focus on building applications rather than worrying about infrastructure setup.

## Features
- **AWS Integration**: Automate infrastructure setup with Terraform for AWS services like ECR, S3, and EKS.
- **CI/CD Setup**: Set up Continuous Integration and Continuous Deployment pipelines using GitHub Actions.
- **Monitoring Setup**: Set up monitoring with Prometheus for easy tracking and alerting.
- **Cost Estimation**: Estimate the cost of your infrastructure deployment.
- **Interactive CLI**: A user-friendly command-line interface that asks for your preferences and generates configuration files.

## Getting Started

### Prerequisites

- Go 1.18+ installed
- AWS account with necessary permissions
- Terraform installed for AWS integration

### Installing

1. Clone the repository:

    ```bash
    git clone https://github.com/Nexusrex18/lumine.git
    cd lumine
    ```

2. Install the required Go dependencies:

    ```bash
    go mod tidy
    ```

### Running the Application

To start the CLI tool, simply run:

```bash
go run main.go
```

## OR, run the one-click install command ! 

```bash
curl -fsSL https://raw.githubusercontent.com/nexusrex18/lumine/main/install.sh | bash
```

The program will prompt you to select the services you want to configure. Follow the interactive menu to:

- Set up CI/CD
- Configure cloud infrastructure (ECR, S3, EKS)
- Estimate costs
- Set up monitoring

### Example Workflow

1. **Setup AWS Service**: Choose `Cloud providers` > `ECR`, and provide a repository name.
2. **Setup CI/CD**: Choose `Setup CI/CD` and follow the prompts to generate a GitHub Actions pipeline.
3. **Generate Terraform Config**: The CLI will automatically generate the necessary configuration files for AWS infrastructure.
4. **Monitoring**: Set up Prometheus for monitoring with a few simple prompts.

## How It Works

- The tool interacts with AWS services using the AWS SDK and Terraform to generate configuration files.
- CI/CD pipelines are automatically set up using GitHub Actions by creating `.github` workflow files.
- Monitoring is configured using Prometheus for infrastructure health tracking.
- Cost estimation can be done using AWS pricing APIs..

## AWS Regions Supported

The tool supports the following AWS regions:

- **US East (N. Virginia)**: `us-east-1`
- **US East (Ohio)**: `us-east-2`
- **US West (N. California)**: `us-west-1`
- **US West (Oregon)**: `us-west-2`
- **Canada (Central)**: `ca-central-1`
- **South America (São Paulo)**: `sa-east-1`

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

