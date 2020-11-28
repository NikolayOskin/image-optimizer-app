pipeline {

  environment {
    imagename = "leakso86/tinyimg"
    credentialsId = 'dockerhub'
  }

  agent any

  stages {
    stage('Build & Test') {
      agent {
        docker {
          image 'golang:latest'
        }
      }

      environment {
        XDG_CACHE_HOME = '/tmp/.cache'
      }

      steps {
        sh 'go get -u golang.org/x/lint/golint'
        sh 'go build'

        echo 'Running vetting'
        sh 'go vet .'

        echo 'Running linting'
        sh 'golint .'

        echo 'Running tests'
        sh 'go test -v'
      }
    }

    stage('Build & Push Image') {
      steps {
        script {
          echo "Building production image"
          def prodImage = docker.build(imagename, "-f prod.Dockerfile --rm .")

          echo "Pushing image to registry"
          docker.withRegistry('', credentialsId) {
            prodImage.push()
          }

          echo 'Cleaning intermediate images'
          sh 'docker image prune -f --filter label=stage=builder'
        }
      }
    }

    stage('Deploy') {
      steps {
        script {
          sshPublisher(
          continueOnError: false, failOnError: true, publishers: [
          sshPublisherDesc(
          configName: "tinyimg", verbose: true, transfers: [
          sshTransfer(execCommand: "cd /var/www/tinyimg && make pull"), sshTransfer(execCommand: "cd /var/www/tinyimg && make update")])])
        }
      }
    }

  }
}