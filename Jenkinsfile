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

        // ════════════════════════════════════════════════════════════════════
        //  TODO – Stage 4: Build Docker Image (backend service)
        //
        //  Uncomment when ready. Requires:
        //    • A Dockerfile at backend/Dockerfile
        //    • Docker daemon accessible by Jenkins agent
        //    • (Optional) Docker registry credentials stored in Jenkins
        //
        // stage('Docker – Build') {
        //     steps {
        //         sh 'docker build -t lumine-backend:${GIT_COMMIT[0..6]} ./backend'
        //     }
        // }
        // ════════════════════════════════════════════════════════════════════

        // ════════════════════════════════════════════════════════════════════
        //  TODO – Stage 5: Deploy to local Minikube
        //
        //  Uncomment when ready. Requires:
        //    • Minikube running locally: `minikube start`
        //    • kubectl configured to point at minikube context
        //    • A k8s manifest at backend/k8s/deployment.yaml
        //
        // stage('Minikube – Deploy') {
        //     steps {
        //         sh '''
        //             # Point Docker to Minikube's registry so the image is available inside the cluster
        //             eval $(minikube docker-env)
        //             docker build -t lumine-backend:latest ./backend
        //
        //             # Apply the deployment manifest
        //             kubectl apply -f backend/k8s/deployment.yaml
        //
        //             # Wait until the pod is Running
        //             kubectl rollout status deployment/lumine-backend --timeout=120s
        //         '''
        //     }
        // }
        // ════════════════════════════════════════════════════════════════════

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
