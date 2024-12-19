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
                git branch: 'main',
                    url: 'https://github.com/IdoShoshani/jenkins-k8s-tasks-itshak.git',
                    credentialsId: 'github-creds-pat'
            }
        }

        stage('Version Management') {
            steps {
                script {
                    // Read and increment version
                    def versionFile = "${CHART_DIRECTORY}/version.txt"
                    def version = fileExists(versionFile) ? readFile(versionFile).trim() : "1.0.0"
                    
                    def (major, minor, patch) = version.tokenize('.').collect { it as int }
                    patch += 1 // Increment patch version by default
                    env.NEW_VERSION = "${major}.${minor}.${patch}"

                    // Save new version back to file
                    writeFile file: versionFile, text: env.NEW_VERSION
                }
            }
        }

        stage('Login to DockerHub') {
            steps {
                sh 'echo $DOCKERHUB_CREDENTIALS_PSW | docker login -u $DOCKERHUB_CREDENTIALS_USR --password-stdin'
            }
        }
        
        stage('Build and Push Images') {
            parallel {
                stage('Python App') {
                    steps {
                        dir('python-app') {
                            script {
                                sh """
                                    docker build -t ${DOCKER_USERNAME}/python-app:${NEW_VERSION} .
                                    docker push ${DOCKER_USERNAME}/python-app:${NEW_VERSION}
                                    docker tag ${DOCKER_USERNAME}/python-app:${NEW_VERSION} ${DOCKER_USERNAME}/python-app:latest
                                    docker push ${DOCKER_USERNAME}/python-app:latest
                                """
                            }
                        }
                    }
                }
                stage('Node App') {
                    steps {
                        dir('node-app') {
                            script {
                                sh """
                                    docker build -t ${DOCKER_USERNAME}/node-app:${NEW_VERSION} .
                                    docker push ${DOCKER_USERNAME}/node-app:${NEW_VERSION}
                                    docker tag ${DOCKER_USERNAME}/node-app:${NEW_VERSION} ${DOCKER_USERNAME}/node-app:latest
                                    docker push ${DOCKER_USERNAME}/node-app:latest
                                """
                            }
                        }
                    }
                }
                stage('Go App') {
                    steps {
                        dir('go-app') {
                            script {
                                sh """
                                    docker build -t ${DOCKER_USERNAME}/go-app:${NEW_VERSION} .
                                    docker push ${DOCKER_USERNAME}/go-app:${NEW_VERSION}
                                    docker tag ${DOCKER_USERNAME}/go-app:${NEW_VERSION} ${DOCKER_USERNAME}/go-app:latest
                                    docker push ${DOCKER_USERNAME}/go-app:latest
                                """
                            }
                        }
                    }
                }
            }
        }
        
        stage('Deploy with Helm') {
            steps {
                script {
                    sh """
                        helm upgrade --install helm-charts ${CHART_DIRECTORY} \
                            --set pythonApp.container.tag=${NEW_VERSION} \
                            --set nodeApp.container.tag=${NEW_VERSION} \
                            --set goApp.container.tag=${NEW_VERSION}
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
            echo "Successfully built, pushed, and deployed version: ${NEW_VERSION}"
        }
    }
}
