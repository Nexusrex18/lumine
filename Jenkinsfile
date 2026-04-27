pipeline {
    agent any

    // ── Trigger: run on every commit to any branch ───────────────────────────
    triggers {
        githubPush()   // requires the "GitHub" Jenkins plugin + webhook configured
    }

    tools {
        // Must match the name in: Manage Jenkins → Tools → Go installations
        go 'Go'
    }

    environment {
        CGO_ENABLED = '0'
        GOPROXY     = 'https://proxy.golang.org,direct'
    }

    stages {

        // ── Stage 1: Checkout ────────────────────────────────────────────────
        stage('Checkout') {
            steps {
                checkout scm
            }
        }

        // ── Stage 2: Build ───────────────────────────────────────────────────
        stage('Build') {
            steps {
                sh '''
                    go version
                    go mod tidy
                    go build ./...
                '''
            }
        }

        // ── Stage 3: Test (integrations/ only) ──────────────────────────────
        stage('Test – Integrations') {
            steps {
                sh 'go test -v -count=1 ./integrations/...'
            }
        }

        // ── Stage 4: Build Docker Image ──────────────────────────────────────
        stage('Docker – Build') {
            steps {
                sh '''
                    # Point Docker to Minikube's internal daemon — builds the image
                    # directly inside the cluster, no image load/cache step needed
                    eval $(minikube docker-env)
                    docker build -t lumine-backend:latest -f ./backend/Dockerfile .
                    echo "Image built directly into Minikube daemon: lumine-backend:latest"
                '''
            }
        }

    }

    post {
        success {
            echo '✅ Build + integration tests passed.'
        }
        failure {
            echo '❌ Pipeline failed — check the stage logs above.'
        }
        always {
            cleanWs()
        }
    }
}
