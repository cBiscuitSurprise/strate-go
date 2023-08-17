pipeline {
    agent any

    environment { 
        TMP_DIR = "${env.WORKSPACE}@tmp"
        BIN_DIR = "bazel-bin/strate-go_"
        OUTPUT_BIN = "strate-go-dev-${env.BUILD_ID}"
        BAZEL_OUTPUT_PATH = ""
    }

    stages {
        stage('Pull') {
            steps {
                git branch: 'main', url: 'https://github.com/cBiscuitSurprise/strate-go.git'
            }
        }

        stage('Unit Test') {
            agent {
                docker {
                    image 'bazel-public/bazel:latest'
                    registryUrl 'https://gcr.io/'
                    args '--entrypoint='
                }
            }

            steps {
                sh "bazel test --test_output=errors //..."
                sh "bazel coverage --combined_report=lcov //..."
                // sh "genhtml --output ${TMP_DIR}/coverage/unit  \"\$(bazel info output_path)/_coverage/_coverage_report.dat\""

                script {
                    BAZEL_OUTPUT_PATH = sh(returnStdout: true, script: 'bazel info output_path')
                }
            }
        }

        stage('Build') {
            agent {
                docker {
                    image 'bazel-public/bazel:latest'
                    registryUrl 'https://gcr.io/'
                    args '--entrypoint='
                }
            }
            
            steps {
                sh "bazel --output_user_root=${TMP_DIR}/build_output build //:strate-go"

                sh "cp ${BIN_DIR}/strate-go ${BIN_DIR}/${OUTPUT_BIN}"
                sh "chmod +x ${BIN_DIR}/${OUTPUT_BIN}"
            }

            post {
                success {
                    archiveArtifacts "${BIN_DIR}/${OUTPUT_BIN}"
                    // publishHTML (target : [allowMissing: false,
                    //     alwaysLinkToLastBuild: true,
                    //     keepAll: true,
                    //     reportDir: "${TMP_DIR}/coverage/unit",
                    //     reportName: 'Unit Test Coverage',
                    //     reportTitles: 'Unit Test Coverage'])
                }
            }
        }

        stage('Package') {
            steps {
                echo "TODO: push to S3"
                sh "ls -la ${env.BIN_DIR}/${env.OUTPUT_BIN}"
            }
        }
    }
}
