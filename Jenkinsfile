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
                  serviceAccountName: helm-service-account
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
        DOCKER_USERNAME = 'idoshoshani123'
        CHART_DIRECTORY = "./helm-charts"
        RELEASE_NAME = "my-app-release"
    }
    
    stages {
        stage('Checkout') {
            steps {
                git branch: "5-improved-jenkins-version-update",
                    url: 'https://github.com/IdoShoshani/jenkins-k8s-tasks-itshak.git',
                    credentialsId: 'github-creds-pat'
            }
        }

    stage('Version Management') {
        steps {
            script {
                sh """
                    echo "Installing yq..."
                    if wget https://github.com/mikefarah/yq/releases/latest/download/yq_linux_\$(dpkg --print-architecture) -O /usr/bin/yq && chmod +x /usr/bin/yq; then
                        echo "yq was installed successfully."
                    else
                        echo "Failed to install yq."
                        exit 1
                    fi
                """

                sh 'git config --global --add safe.directory ${WORKSPACE}'

                def versionFile = "${CHART_DIRECTORY}/version.txt"
                if (!fileExists(versionFile)) {
                    echo "version.txt not found. Initializing with version 1.0.0."
                    writeFile file: versionFile, text: "1.0.0"
                }

                def version = readFile(versionFile).trim()
                echo "Current Version: ${version}"

                def versionParts = version.tokenize('-')
                def (major, minor, patch) = versionParts[0].tokenize('.').collect { it as int }
                patch += 1

                env.VERSION = "${major}.${minor}.${patch}-${env.BUILD_NUMBER}"
                echo "New Version: ${env.VERSION}"

                writeFile file: versionFile, text: env.VERSION

                sh """
                    set -x
                    yq -i '.goApp.container.tag = "\${VERSION}"' ${CHART_DIRECTORY}/values.yaml
                    yq -i '.pythonApp.container.tag = "\${VERSION}"' ${CHART_DIRECTORY}/values.yaml
                    yq -i '.nodeApp.container.tag = "\${VERSION}"' ${CHART_DIRECTORY}/values.yaml
                    yq -i '.appVersion = "\${VERSION}"' ${CHART_DIRECTORY}/Chart.yaml
                    yq -i '.version = "\${VERSION}"' ${CHART_DIRECTORY}/Chart.yaml
                """

                withCredentials([usernamePassword(credentialsId: 'github-creds-pat', usernameVariable: 'GIT_USERNAME', passwordVariable: 'GIT_PASSWORD')]) {
                    sh """
                        git config --global user.email "ci@jenkins"
                        git config --global user.name "Jenkins CI"
                        git add ${versionFile} ${CHART_DIRECTORY}/values.yaml ${CHART_DIRECTORY}/Chart.yaml
                        git diff --cached --name-only
                        git commit -m "Update version to \${VERSION}" || echo "No changes to commit"
                        git push https://\${GIT_USERNAME}:\${GIT_PASSWORD}@github.com/IdoShoshani/jenkins-k8s-tasks-itshak.git 5-improved-jenkins-version-update
                    """
                }
            }
        }
    }


        stage('Login to DockerHub') {
            steps {
                withCredentials([usernamePassword(credentialsId: 'docker-creds', usernameVariable: 'DOCKERHUB_USER', passwordVariable: 'DOCKERHUB_PASS')]) {
                    sh 'echo $DOCKERHUB_PASS | docker login -u $DOCKERHUB_USER --password-stdin'
                }
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
                    sh """
                        docker push ${DOCKER_USERNAME}/python-app:${VERSION}
                        docker push ${DOCKER_USERNAME}/python-app:latest
                    """
                    sh """
                        docker push ${DOCKER_USERNAME}/node-app:${VERSION}
                        docker push ${DOCKER_USERNAME}/node-app:latest
                    """
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
                        helm upgrade --install ${RELEASE_NAME} ${CHART_DIRECTORY} \
                            --values ${CHART_DIRECTORY}/values.yaml \
                            --namespace default
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
