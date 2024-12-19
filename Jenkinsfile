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
        VERSION = sh(script: 'date +%Y.%m.%d.%H%M', returnStdout: true).trim()
        DOCKER_USERNAME = 'idoshoshani123'
    }
    
    stages {
        stage('Checkout') {
            steps {
                git branch: 'main',
                    url: 'https://github.com/IdoShoshani/jenkins-k8s-tasks-itshak.git',
                    credentialsId: 'github-creds-pat'
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
    }
    
    post {
        always {
            sh 'docker logout'
        }
        failure {
            echo 'Pipeline failed! Check the logs for details.'
        }
        success {
            echo "Successfully built and pushed version: ${VERSION}"
        }
    }
}