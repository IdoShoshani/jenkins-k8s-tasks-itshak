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
                git branch: "main",
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
                    
                    def apps = ['goApp', 'nodeApp', 'pythonApp']
                    
                    apps.each { app ->
                        def currentTag = sh(
                            script: "yq eval '.${app}.container.tag' ${CHART_DIRECTORY}/values.yaml",
                            returnStdout: true
                        ).trim()

                        def tagParts = currentTag.tokenize('.')
                        def newTag = ""

                        if (tagParts.size() == 3) {
                            newTag = "${tagParts[0]}.${tagParts[1]}.${tagParts[2]}.${env.BUILD_NUMBER}"
                        } else if (tagParts.size() == 4) {
                            newTag = "${tagParts[0]}.${tagParts[1]}.${tagParts[2]}.${env.BUILD_NUMBER}"
                        } else {
                            error "Invalid tag format for ${app}: ${currentTag}"
                        }

                        sh """
                            yq -i ".${app}.container.tag = \\"${newTag}\\"" ${CHART_DIRECTORY}/values.yaml
                        """
                        
                        echo "Updated ${app} tag from ${currentTag} to ${newTag}"
                    }
                    
                    withCredentials([usernamePassword(credentialsId: 'github-creds-pat', usernameVariable: 'GIT_USERNAME', passwordVariable: 'GIT_PASSWORD')]) {
                        sh """
                            git config --global user.email "ci@jenkins"
                            git config --global user.name "Jenkins CI"
                            git add ${CHART_DIRECTORY}/values.yaml
                            git diff --cached --name-only
                            git commit -m "Update application tags to build number \${BUILD_NUMBER}" || echo "No changes to commit"
                            git push https://\${GIT_USERNAME}:\${GIT_PASSWORD}@github.com/IdoShoshani/jenkins-k8s-tasks-itshak.git main
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

        stage('Build and Push Images') {
            steps {
                script {
                    // Define the apps and their respective directories
                    def apps = [
                        [name: 'pythonApp', dir: 'python-app'],
                        [name: 'nodeApp', dir: 'node-app'],
                        [name: 'goApp', dir: 'go-app']
                    ]
                    def parallelStages = [:]
                    
                    apps.each { app ->
                        parallelStages[app.name] = {
                            dir(app.dir) {
                                script {
                                    // Fetch version from values.yaml
                                    def appVersion = sh(
                                        script: "yq eval '.${app.name}.container.tag' ../${CHART_DIRECTORY}/values.yaml",
                                        returnStdout: true
                                    ).trim()
                                    echo "Building and pushing ${app.name} with version: ${appVersion}"
                                    sh """
                                        docker build -t ${DOCKER_USERNAME}/${app.dir}:${appVersion} .
                                        docker tag ${DOCKER_USERNAME}/${app.dir}:${appVersion} ${DOCKER_USERNAME}/${app.dir}:latest
                                        docker push ${DOCKER_USERNAME}/${app.dir}:${appVersion}
                                        docker push ${DOCKER_USERNAME}/${app.dir}:latest
                                    """
                                }
                            }
                        }
                    }
                    parallel parallelStages
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
            echo "Successfully built, pushed, and deployed applications."
        }
    }
}
