node {
    checkout scm

    stage('Build') {
        try {
            // build
            sh 'make docker-compile'
        } finally {
            // remove docker image
            sh 'make docker-compile-rmi'
        }
    }
}
