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
        sh '''export XDG_CACHE_HOME="/tmp/.cache"
go get -u golang.org/x/lint/golint
go build'''
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