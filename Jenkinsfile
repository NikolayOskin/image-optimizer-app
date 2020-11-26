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
        sh 'go get -u golang.org/x/lint/golint'
      }
    }

  }
}