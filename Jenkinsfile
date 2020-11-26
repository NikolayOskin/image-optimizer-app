pipeline {
  agent any
  stages {
    stage('Build') {
      steps {
        git(url: 'https://github.com/NikolayOskin/image-optimizer-app', branch: 'master')
      }
    }

  }
  environment {
    Prod = ''
  }
}