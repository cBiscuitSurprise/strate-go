pipeline {
    agent any

    environment {
        TMP_DIR = "${env.WORKSPACE}@2@tmp"
        BIN_DIR = "bazel-bin/strate-go_"
        OUTPUT_DIR = ".jenkins-artifacts"
        OUTPUT_COV_DIR = ".jenkins-cov"
        OUTPUT_BIN = "strate-go-dev-${env.BUILD_ID}"
        BAZEL_OUTPUT_PATH = ""
    }

    stages {
        stage('Pull') {
            steps {
                git branch: 'main', url: 'https://github.com/cBiscuitSurprise/strate-go.git'
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

                sh "rm -rf ${OUTPUT_DIR}/*"
                sh "mkdir -p \$(dirname ${OUTPUT_DIR}/${OUTPUT_BIN})"
                sh "cp ${BIN_DIR}/strate-go ${OUTPUT_DIR}/${OUTPUT_BIN}"
                sh "chmod +x ${OUTPUT_DIR}/${OUTPUT_BIN}"
                stash includes: "${OUTPUT_DIR}/**/*", name: "artifacts"
            }

            post {
                success {
                    archiveArtifacts "${OUTPUT_DIR}/${OUTPUT_BIN}"
                }
            }
        }

        stage('Test') {
            agent {
                docker {
                    image 'bazel-public/bazel:latest'
                    registryUrl 'https://gcr.io/'
                    args '--entrypoint='
                }
            }

            steps {
                sh "bazel --output_user_root=${TMP_DIR}/build_output coverage --combined_report=lcov //..."

                sh "rm -rf ${OUTPUT_COV_DIR}/*"
                sh "mkdir -p ${OUTPUT_COV_DIR}"
                sh "cp -f bazel-out/_coverage/_coverage_report.dat \"${OUTPUT_COV_DIR}/\""
                stash includes: "${OUTPUT_COV_DIR}/*", name: "cov"
            }
        }

        stage('Coverage') {
            agent {
                docker {
                    image 'alpine:latest'
                    args '-u root'
                }
            }

            steps {
                // TODO: create image that has lcov in it
                sh "apk add \
                    --no-cache \
                    --repository http://dl-cdn.alpinelinux.org/alpine/edge/testing \
                    --repository http://dl-cdn.alpinelinux.org/alpine/edge/main \
                    lcov"

                sh "rm -rf ${OUTPUT_COV_DIR}/*"
                unstash "cov"
                sh "mkdir -p ${OUTPUT_COV_DIR}/report/"
                sh "genhtml --output ${OUTPUT_COV_DIR}/report/ \"${OUTPUT_COV_DIR}/_coverage_report.dat\""
                sh "chmod a+rw -R ${OUTPUT_COV_DIR}/report"
            }

            post {
                success {
                    publishHTML (
                        target : [
                            reportName: 'Unit Test Coverage',
                            allowMissing: false,
                            alwaysLinkToLastBuild: true,
                            keepAll: false,
                            reportDir: "${OUTPUT_COV_DIR}/report/",
                            reportFiles: 'index.html',
                        ]
                    )
                }
            }
        }
        
        stage('Package') {
            agent {
                docker {
                    image 'amazon/aws-cli:latest'
                    args '--entrypoint='
                }
            }

            steps {
                sh "rm -rf ${OUTPUT_DIR}/*"

                unstash "artifacts"

                echo "TODO: push to S3"
                sh "ls -la ${OUTPUT_DIR}/${OUTPUT_BIN}"
            }
        }
        
        stage('Image') {
            steps {
                unstash "artifacts"

                script {
                    docker.withRegistry('', 'dockerhub_cbiscuit87') {
                        img = docker.build(
                            "cbiscuit87/strate-go:dev-${env.BUILD_ID}",
                            "--build-arg=\"BINARY=${OUTPUT_DIR}/${OUTPUT_BIN}\" -f build/prod.Dockerfile ."
                        )
                        img.push()
                    }
                }
            }
        }
    }
}
