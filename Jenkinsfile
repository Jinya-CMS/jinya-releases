// Uses Declarative syntax to run commands inside a container.
pipeline {
    triggers {
        pollSCM("*/5 * * * *")
    }
    agent {
        kubernetes {
            yaml '''
apiVersion: v1
kind: Pod
spec:
  volumes:
    - name: docker-sock
      hostPath:
        path: /var/run/docker.sock
  containers:
  - name: golang
    image: quay.imanuel.dev/dockerhub/library---golang:latest
    command:
    - sleep
    args:
    - infinity
  - name: docker
    image: quay.imanuel.dev/dockerhub/library---docker:latest
    command:
    - cat
    tty: true
    volumeMounts:
    - mountPath: /var/run/docker.sock
      name: docker-sock
'''
            defaultContainer 'golang'
        }
    }
    stages {
        stage('Build container') {
            when {
              branch 'main'
            }
            steps {
                container('docker') {
                    sh "docker build -t quay.imanuel.dev/jinya/jinya-releases:$BUILD_NUMBER -f ./Dockerfile ."
                    sh "docker tag quay.imanuel.dev/jinya/jinya-releases:$BUILD_NUMBER quay.imanuel.dev/jinya/jinya-releases:latest"
                    sh "docker tag quay.imanuel.dev/jinya/jinya-releases:$BUILD_NUMBER jinyacms/jinya-releases:$BUILD_NUMBER"
                    sh "docker tag quay.imanuel.dev/jinya/jinya-releases:$BUILD_NUMBER jinyacms/jinya-releases:latest"

                    withDockerRegistry(credentialsId: 'quay.imanuel.dev', url: 'https://quay.imanuel.dev') {
                        sh "docker push quay.imanuel.dev/jinya/jinya-releases:$BUILD_NUMBER"
                        sh "docker push quay.imanuel.dev/jinya/jinya-releases:latest"
                    }
                    withDockerRegistry(credentialsId: 'hub.docker.com', url: '') {
                        sh "docker push jinyacms/jinya-releases:$BUILD_NUMBER"
                        sh "docker push jinyacms/jinya-releases:latest"
                    }
                }
            }
        }
    }
}
