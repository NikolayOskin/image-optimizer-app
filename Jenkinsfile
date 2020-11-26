pipeline {
  agent {
    docker {
      image 'golang:latest'
      args '-p 3000:3000'
    }

  }
  stages {
    stage('Build') {
      steps {
        sh '''export XDG_CACHE_HOME="/tmp/.cache"'''

        echo 'Getting golint'
        sh 'go get -u golang.org/x/lint/golint'

        echo 'Building'
        sh 'go build'
      }
    }

    stage('Lint') {
      steps {
        withEnv(["PATH+GO=${GOPATH}/bin"]) {
          echo 'Running vetting'
          sh 'go vet .'
          echo 'Running linting'
          sh 'golint .'
        }
      }
    }

  }
}