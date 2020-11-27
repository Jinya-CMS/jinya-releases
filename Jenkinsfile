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
  imagePullSecrets:
    - name: dev-imanuel-jenkins-regcred
  volumes:
    - name: docker-sock
      hostPath:
        path: /var/run/docker.sock
  containers:
  - name: golang
    image: registry.imanuel.dev/library/golang:1.15-buster
    command:
    - sleep
    args:
    - infinity
  - name: docker
    image: registry.imanuel.dev/library/docker:stable
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
                    script {
                        def image = docker.build "registry-hosted.imanuel.dev/jinya/jinya-releases:$BUILD_NUMBER"
                        docker.withRegistry('https://registry-hosted.imanuel.dev', 'nexus.imanuel.dev') {
                            image.push()
                        }
                    }
                }
            }
        }
    }
}
