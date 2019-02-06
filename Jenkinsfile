#!/usr/bin/env groovy

def getDockerImageTag() {
  if (env.GIT_COMMIT == "") {
    error "GIT_COMMIT value was empty at usage. "
  }
  return "${env.BUILD_ID}-${env.GIT_COMMIT}"
}

pipeline {
  agent none

  environment {
    ENVIRONMENT = 'ops'
    DOCKER_IMAGE_TAG = "${getDockerImageTag()}"
  }

  stages {

    stage('GPG Verify') {
      // agent is require to check out source and get GIT_URL
      agent {
        docker { image 'docker.io/controlplane/gcloud-sdk:latest' }
      }

      options {
        timeout(time: 15, unit: 'MINUTES')
        retry(1)
        timestamps()
      }

      steps {
        ansiColor('xterm') {
          sh "echo Verifying GPG signatures for: ${env.GIT_URL}"
        }

        build job: 'git-gpg-verify', parameters: [
          string(
            name: 'GIT_REPO',
            value: "${env.GIT_URL}"
          )
        ]
      }
    }

    stage('Test') {
      agent {
        docker { image 'docker.io/controlplane/gcloud-sdk:latest' }
      }

      options {
        timeout(time: 15, unit: 'MINUTES')
        retry(1)
        timestamps()
      }

      steps {
        ansiColor('xterm') {
          sh 'make test'
        }
      }
    }

    stage('Push') {
      agent {
        docker {
          image 'docker.io/controlplane/gcloud-sdk:latest'
          args '-v /var/run/docker.sock:/var/run/docker.sock ' +
            '--user=root ' +
            '--cap-drop=ALL ' +
            '--cap-add=DAC_OVERRIDE'
        }
      }

      environment {
        DOCKER_REGISTRY_CREDENTIALS = credentials("${ENVIRONMENT}_docker_credentials")
      }

      options {
        timeout(time: 15, unit: 'MINUTES')
        retry(1)
        timestamps()
      }

      steps {
        ansiColor('xterm') {
          sh """
            echo '${DOCKER_REGISTRY_CREDENTIALS_PSW}' \
            | docker login \
              --username '${DOCKER_REGISTRY_CREDENTIALS_USR}' \
              --password-stdin

            make push DOCKER_IMAGE_TAG="${DOCKER_IMAGE_TAG}"
          """
        }
      }
    }

  }
}
