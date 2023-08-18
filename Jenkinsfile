pipeline {
    agent any
    stages {
        stage("deploy") {
            steps {
                sh "$WORKSPACE/deploy.sh"
            }
        }
    }
}
