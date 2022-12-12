pipeline {
    agent any

    stages  {
        stage("Check Build Branch") {
            steps {
                echo "Checking build branch in progress ......"
                script {
                    if (env.GIT_BRANCH ==~ /^v([0-9])+\.([0-9])+\.([0-9])+.*/)  {
                        echo "building production environment, tag=${env.GIT_BRANCH}"
                    } else if (env.GIT_BRANCH ==~ /^test-([0-9])+\.([0-9])+\.([0-9])+.*/) {
                        echo "building test environment, tag=${env.GIT_BRANCH}"
                    } else if (env.GIT_BRANCH ==~ /(origin\/develop)/) {
                        echo "building development environment, /origin/develop"
                    } else {
                        echo "The build branch ${env.GIT_BRANCH} is not legal, allowing to build the development environment branch (/origin/develop), the test environment branch (e.g. test-1.0.0), and the production environment branch (e.g. v1.0.0)"
                        sh 'exit 1'
                    }
                }
                echo "Check build branch complete."
            }
        }

        stage("Check Code") {
            steps {
                echo "Checking code in progress ......"
                sh 'make ci-lint'
                echo "Check code complete."
            }
        }

        stage("Unit Testing") {
            steps {
                echo "Unit testing in progress ......"
                sh 'make test'
                echo "Unit testing complete."
            }
        }

        stage("Compile Code") {
            steps {
                echo "Compiling code  in progress ......"
                sh 'make build'
                echo "compile code complete."
            }
        }

        stage("Build Image") {
            steps {
                echo "building image in progress ......"
                script {
                    registryHost=""
                    tagName=""
                    if (env.GIT_BRANCH ==~ /^v([0-9])+\.([0-9])+\.([0-9])+.*/) {
                        if (env.PROD_REPO_HOST == null) {
                            echo "The value of environment variable PROD_REPO_HOST is empty, please set the value of PROD_REPO_HOST in [Jenkins Management] --> [System Settings] --> [Environment Variables]."
                            sh 'exit 1'
                        }
                        echo "Use the production environment image repository ${env.PROD_REPO_HOST}"
                        registryHost=env.PROD_REPO_HOST
                        tagName=env.GIT_BRANCH
                    }
                    else if (env.GIT_BRANCH ==~ /^test-([0-9])+\.([0-9])+\.([0-9])+.*/) {
                          if (env.TEST_REPO_HOST == null) {
                              echo "The value of environment variable TEST_REPO_HOST is empty, please set the value of TEST_REPO_HOST in [Jenkins Management] --> [System Settings] --> [Environment Variables]."
                              sh 'exit 1'
                          }
                          echo "Use the test environment image repository ${env.TEST_REPO_HOST}"
                          registryHost=env.TEST_REPO_HOST
                          tagName=env.GIT_BRANCH
                    }
                    else {
                        if (env.DEV_REPO_HOST == null) {
                            echo "The value of environment variable DEV_REPO_HOST is empty, please set the value of DEV_REPO_HOST in [Jenkins Management] --> [System Settings] --> [Environment Variables]."
                            sh 'exit 1'
                        }
                        echo "Using the development environment ${env.DEV_REPO_HOST}"
                        registryHost=env.DEV_REPO_HOST
                    }
                    sh "make image-build REPO_HOST=$registryHost TAG=$tagName"
                }
                echo "Build image complete"
            }
        }

        stage("Push Image") {
            steps {
                echo "pushing image in progress ......"
                script {
                    registryHost=""
                    tagName=""
                    if (env.GIT_BRANCH ==~ /^v([0-9])+\.([0-9])+\.([0-9])+.*/) {
                        if (env.PROD_REPO_HOST == null) {
                            echo "The value of environment variable PROD_REPO_HOST is empty, please set the value of PROD_REPO_HOST in [Jenkins Management] --> [System Settings] --> [Environment Variables]."
                            sh 'exit 1'
                        }
                        echo "Use the production environment image repository ${env.PROD_REPO_HOST}"
                        registryHost=env.PROD_REPO_HOST
                        tagName=env.GIT_BRANCH
                    }
                    else if (env.GIT_BRANCH ==~ /^test-([0-9])+\.([0-9])+\.([0-9])+.*/) {
                          if (env.TEST_REPO_HOST == null) {
                              echo "The value of environment variable TEST_REPO_HOST is empty, please set the value of TEST_REPO_HOST in [Jenkins Management] --> [System Settings] --> [Environment Variables]."
                              sh 'exit 1'
                          }
                          echo "Use the test environment image repository ${env.TEST_REPO_HOST}"
                          registryHost=env.TEST_REPO_HOST
                          tagName=env.GIT_BRANCH
                    }
                    else {
                        if (env.DEV_REPO_HOST == null) {
                            echo "The value of environment variable DEV_REPO_HOST is empty, please set the value of DEV_REPO_HOST in [Jenkins Management] --> [System Settings] --> [Environment Variables]."
                            sh 'exit 1'
                        }
                        echo "Using the development environment ${env.DEV_REPO_HOST}"
                        registryHost=env.DEV_REPO_HOST
                    }
                    sh "make image-push REPO_HOST=$registryHost TAG=$tagName"
                }
                echo "push image complete, clear image complete."
            }
        }

        stage("Deploy to k8s") {
            when { expression { return env.GIT_BRANCH ==~ /(origin\/staging|origin\/develop)/ } }
            steps {
                echo "Deploying to k8s in progress ......"
                sh 'make deploy-k8s'
                echo "Deploy to k8s complete."
            }
        }
    }

    post {
		always {
			echo 'One way or another, I have finished'
			echo sh(returnStdout: true, script: 'env')
			deleteDir() /* clean up our workspace */
		}
		success {
			SendDingding("success")
			//SendEmail("success")
			echo 'structure success'
		}
		failure {
			SendDingding("failure")
			//SendEmail("failure")
			echo 'structure failure'
		}
   }
}

// Notifications using dingding
void SendDingding(res)
{
	// Fill in the corresponding cell phone number and specify a person to be notified in the pinned group
	tel_num="xxxxxxxxxxx"
	dingding_url="https://oapi.dingtalk.com/robot/send\\?access_token\\=your dingding robot token"

    branchName=""
    if (env.GIT_BRANCH ==~ /^v([0-9])+\.([0-9])+\.([0-9])+.*/) {
        branchName="${env.SERVER_PLATFORM} production environment, tag=${env.GIT_BRANCH},  ${env.JOB_NAME}"
    }
    else if (env.GIT_BRANCH ==~ /^test-([0-9])+\.([0-9])+\.([0-9])+.*/){
        branchName="${env.SERVER_PLATFORM} test environment, tag=${env.GIT_BRANCH},  ${env.JOB_NAME}"
    }
    else {
        branchName="${env.SERVER_PLATFORM} develop environment, branch=${env.GIT_BRANCH},  ${env.JOB_NAME}"
    }

	json_msg=""
	if( res == "success" ) {
		json_msg='{\\"msgtype\\":\\"text\\",\\"text\\":{\\"content\\":\\"@' + tel_num +' [OK] ' + "${branchName} ${env.BUILD_NUMBER}th "  + 'build success. \\"},\\"at\\":{\\"atMobiles\\":[\\"' + tel_num + '\\"],\\"isAtAll\\":false}}'
	}
	else {
		json_msg='{\\"msgtype\\":\\"text\\",\\"text\\":{\\"content\\":\\"@' + tel_num +' [cry] ' + "${branchName} ${env.BUILD_NUMBER}th "  + 'build failed, please deal with it promptly! \\"},\\"at\\":{\\"atMobiles\\":[\\"' + tel_num + '\\"],\\"isAtAll\\":false}}'
	}

    post_header="Content-Type:application/json;charset=utf-8"
    sh_cmd="curl -X POST " + dingding_url + " -H " + "\'" + post_header + "\'" + " -d " + "\""  + json_msg + "\""
	sh sh_cmd
}

// Notifications using email
void SendEmail(res)
{
	emailAddr="xxx@xxx.com"
	if( res == "success" )
	{
		mail to: emailAddr,
		subject: "Build Success: ${currentBuild.fullDisplayName}",
		body: "\nJob name: ${env.JOB_NAME} ${env.BUILD_NUMBER}th build. \n\n For more information, please see: ${env.BUILD_URL}"
	}
	else
	{
		mail to: emailAddr,
		subject: "Build Failed: ${currentBuild.fullDisplayName}",
		body: "\nJob name: ${env.JOB_NAME} ${env.BUILD_NUMBER}th build. \n\n For more information, please see: ${env.BUILD_URL}"
	}
}
