pipeline {
    agent {
        kubernetes {
            yaml '''
                apiVersion: v1
                kind: Pod
                metadata:
                  labels:
                    jenkins: agent
                spec:
                  containers:
                  - name: docker-container
                    image: idoshoshani123/docker_and_helm:latest
                    command:
                    - sleep
                    args:
                    - 9999999
                    securityContext:
                      privileged: true
                    volumeMounts:
                    - name: docker-sock
                      mountPath: /var/run/docker.sock
                  volumes:
                  - name: docker-sock
                    hostPath:
                      path: /var/run/docker.sock
            '''
            defaultContainer 'docker-container'
        }
    }
    
    environment {
        DOCKERHUB_CREDENTIALS = credentials('docker-creds')
        DOCKER_USERNAME = 'idoshoshani123'
        CHART_DIRECTORY = "./helm-charts"
    }
    
    stages {
        stage('Checkout') {
            steps {
                git branch: "3-create-jenkins-pipeline",
                    url: 'https://github.com/IdoShoshani/jenkins-k8s-tasks-itshak.git',
                    credentialsId: 'github-creds-pat'
            }
        }

        stage('Version Management') {
            steps {
                script {
                    // Define the path to version.txt
                    def versionFile = "${CHART_DIRECTORY}/version.txt"

                    // Check if version.txt exists; if not, initialize it
                    if (!fileExists(versionFile)) {
                        echo "version.txt not found. Initializing with version 1.0.0."
                        writeFile file: versionFile, text: "1.0.0"
                    }

                    // Read the current version
                    def version = readFile(versionFile).trim()
                    echo "Current Version: ${version}"

                    // Split the version into major, minor, and patch
                    def (major, minor, patch) = version.tokenize('.').collect { it as int }

                    // Increment the patch version by default
                    patch += 1

                    // Set the new version
                    env.VERSION = "${major}.${minor}.${patch}"
                    echo "New Version: ${env.VERSION}"

                    // Save the new version back to version.txt
                    writeFile file: versionFile, text: env.VERSION

                    // Commit the updated version.txt to Git
                    sh """
                        git config --global user.email "ci@jenkins"
                        git config --global user.name "Jenkins CI"
                        git add ${versionFile}
                        git commit -m "Update version to ${env.VERSION}" || echo "No changes to commit"
                        git push origin ${env.BRANCH_NAME}
                    """
                }
            }
        }

        stage('Login to DockerHub') {
            steps {
                sh 'echo $DOCKERHUB_CREDENTIALS_PSW | docker login -u $DOCKERHUB_CREDENTIALS_USR --password-stdin'
            }
        }
        
        stage('Build Images') {
            parallel {
                stage('Build Python App') {
                    steps {
                        dir('python-app') {
                            sh """
                                docker build -t ${DOCKER_USERNAME}/python-app:${VERSION} .
                                docker tag ${DOCKER_USERNAME}/python-app:${VERSION} ${DOCKER_USERNAME}/python-app:latest
                            """
                        }
                    }
                }
                
                stage('Build Node App') {
                    steps {
                        dir('node-app') {
                            sh """
                                docker build -t ${DOCKER_USERNAME}/node-app:${VERSION} .
                                docker tag ${DOCKER_USERNAME}/node-app:${VERSION} ${DOCKER_USERNAME}/node-app:latest
                            """
                        }
                    }
                }
                
                stage('Build Go App') {
                    steps {
                        dir('go-app') {
                            sh """
                                docker build -t ${DOCKER_USERNAME}/go-app:${VERSION} .
                                docker tag ${DOCKER_USERNAME}/go-app:${VERSION} ${DOCKER_USERNAME}/go-app:latest
                            """
                        }
                    }
                }
            }
        }
        
        stage('Push Images') {
            steps {
                script {
                    // Push Python App
                    sh """
                        docker push ${DOCKER_USERNAME}/python-app:${VERSION}
                        docker push ${DOCKER_USERNAME}/python-app:latest
                    """
                    
                    // Push Node App
                    sh """
                        docker push ${DOCKER_USERNAME}/node-app:${VERSION}
                        docker push ${DOCKER_USERNAME}/node-app:latest
                    """
                    
                    // Push Go App
                    sh """
                        docker push ${DOCKER_USERNAME}/go-app:${VERSION}
                        docker push ${DOCKER_USERNAME}/go-app:latest
                    """
                }
            }
        }

        stage('Deploy with Helm') {
            steps {
                script {
                    sh """
                        helm upgrade --install helm-charts ${CHART_DIRECTORY} \
                            --set pythonApp.container.tag=${VERSION} \
                            --set nodeApp.container.tag=${VERSION} \
                            --set goApp.container.tag=${VERSION}
                    """
                }
            }
        }
    }
    
    post {
        always {
            sh 'docker logout'
        }
        failure {
            echo 'Pipeline failed! Check the logs for details.'
        }
        success {
            echo "Successfully built, pushed, and deployed version: ${VERSION}"
        }
    }
}
