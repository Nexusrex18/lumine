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
                    # Build context = repo root (.) so go.mod + go.sum are visible
                    # -f points Jenkins at the Dockerfile inside backend/
                    docker build -t lumine-backend:latest -f ./backend/Dockerfile .

                    # Load the image into Minikube's internal registry
                    minikube image load lumine-backend:latest

                    echo "Image built and loaded into Minikube: lumine-backend:latest"
                '''
            }
        }

        // ── Stage 5: Deploy to local Minikube ────────────────────────────────
        stage('Minikube – Deploy') {
            steps {
                sh '''
                    # Apply (or update) the Deployment + Service
                    kubectl apply -f backend/k8s/deployment.yaml

                    # Block until the pod is Running (max 2 minutes)
                    kubectl rollout status deployment/lumine-backend --timeout=120s

                    # Print pod status and the URL to hit the service
                    kubectl get pods -l app=lumine-backend
                    echo "Service available at: http://$(minikube ip):30080"
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
