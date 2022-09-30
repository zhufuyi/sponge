pipeline {
    agent any

    stages  {
        stage("检查构建分支") {
            steps {
                echo "检查构建分支中......"
                script {
                    if (env.GIT_BRANCH ==~ /^v([0-9])+\.([0-9])+\.([0-9])+.*/)  {
                        echo "构建生产环境，tag=${env.GIT_BRANCH}"
                    } else if (env.GIT_BRANCH ==~ /^test-([0-9])+\.([0-9])+\.([0-9])+.*/) {
                        echo "构建测试环境，tag=${env.GIT_BRANCH}"
                    } else if (env.GIT_BRANCH ==~ /(origin\/develop)/) {
                        echo "构建开发环境，/origin/develop"
                    } else {
                        echo "构建分支${env.GIT_BRANCH}不合法，允许构建生产环境分支(例如：v1.0.0)，开发产环境分支(例如：test-1.0.0)，开发环境分支(/origin/develop)"
                        sh 'exit 1'
                    }
                }
                echo "检查构建分支完成."
            }
        }

        stage("代码检查") {
            steps {
                echo "代码检查中......"
                sh 'make ci-lint'
                echo "代码检查完成."
            }
        }

        stage("单元测试") {
            steps {
                echo "单元测试中......"
                sh 'make test'
                echo "单元测试完成."
            }
        }

        stage("编译代码") {
            steps {
                echo "编译代码中......"
                sh 'make build'
                echo "编译代码完成."
            }
        }

        stage("构建镜像") {
            steps {
                echo "构建镜像中......"
                // 兼容自动构建和参数构建
                script {
                    registryHost=""
                    tagName=""
                    if (env.GIT_BRANCH ==~ /^v([0-9])+\.([0-9])+\.([0-9])+.*/) {
                        if (env.PROD_REPO_HOST == null) {
                            echo "环境变量PROD_REPO_HOST值为空，请在【Jenkins管理】--> 【系统设置】-->【环境变量】设置PROD_REPO_HOST值"
                            sh 'exit 1'
                        }
                        echo "使用生产环境镜像仓库 ${env.PROD_REPO_HOST}"
                        registryHost=env.PROD_REPO_HOST
                        tagName=env.GIT_BRANCH
                    }
                    else if (env.GIT_BRANCH ==~ /^test-([0-9])+\.([0-9])+\.([0-9])+.*/) {
                          if (env.TEST_REPO_HOST == null) {
                              echo "环境变量TEST_REPO_HOST值为空，请在【Jenkins管理】--> 【系统设置】-->【环境变量】设置TEST_REPO_HOST值"
                              sh 'exit 1'
                          }
                          echo "使用测试环境镜像仓库 ${env.TEST_REPO_HOST}"
                          registryHost=env.TEST_REPO_HOST
                          tagName=env.GIT_BRANCH
                    }
                    else {
                        if (env.DEV_REPO_HOST == null) {
                            echo "环境变量DEV_REPO_HOST值为空，请在【Jenkins管理】--> 【系统设置】-->【环境变量】设置DEV_REPO_HOST值"
                            sh 'exit 1'
                        }
                        echo "使用开发环境 ${env.DEV_REPO_HOST}"
                        registryHost=env.DEV_REPO_HOST
                    }
                    sh "make image-build REPO_HOST=$registryHost TAG=$tagName"
                }
                echo "构建镜像完成"
            }
        }

        stage("上传镜像") {
            steps {
                echo "上传镜像中......"
                // 兼容自动构建和参数构建
                script {
                    registryHost=""
                    tagName=""
                    if (env.GIT_BRANCH ==~ /^v([0-9])+\.([0-9])+\.([0-9])+.*/) {
                        if (env.PROD_REPO_HOST == null) {
                            echo "环境变量PROD_REPO_HOST值为空，请在【Jenkins管理】--> 【系统设置】-->【环境变量】设置PROD_REPO_HOST值"
                            sh 'exit 1'
                        }
                        echo "使用生产环境镜像仓库 ${env.PROD_REPO_HOST}"
                        registryHost=env.PROD_REPO_HOST
                        tagName=env.GIT_BRANCH
                    }
                    else if (env.GIT_BRANCH ==~ /^test-([0-9])+\.([0-9])+\.([0-9])+.*/) {
                          if (env.TEST_REPO_HOST == null) {
                              echo "环境变量TEST_REPO_HOST值为空，请在【Jenkins管理】--> 【系统设置】-->【环境变量】设置TEST_REPO_HOST值"
                              sh 'exit 1'
                          }
                          echo "使用测试环境镜像仓库 ${env.TEST_REPO_HOST}"
                          registryHost=env.TEST_REPO_HOST
                          tagName=env.GIT_BRANCH
                    }
                    else {
                        if (env.DEV_REPO_HOST == null) {
                            echo "环境变量DEV_REPO_HOST值为空，请在【Jenkins管理】--> 【系统设置】-->【环境变量】设置DEV_REPO_HOST值"
                            sh 'exit 1'
                        }
                        echo "使用开发环境 ${env.DEV_REPO_HOST}"
                        registryHost=env.DEV_REPO_HOST
                    }
                    sh "make image-push REPO_HOST=$registryHost TAG=$tagName"
                }
                echo "上传镜像完成，清除镜像完成。"
            }
        }

        stage("部署到k8s") {
            // 生产环境和测试环境跳过部署，手动部署
            when { expression { return env.GIT_BRANCH ==~ /(origin\/staging|origin\/develop)/ } }
            steps {
                echo "部署到k8s..."
                sh 'make deploy-k8s'
                echo "部署到k8s完成"
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


void SendDingding(res)
{
	// 输入相应的手机号码，在钉钉群指定通知某个人
	tel_num="xxxxxxxxxxx"

	// 钉钉机器人的地址
	dingding_url="https://oapi.dingtalk.com/robot/send\\?access_token\\=你的钉钉机器人token"

    branchName=""
    if (env.GIT_BRANCH ==~ /^v([0-9])+\.([0-9])+\.([0-9])+.*/) {
        branchName="${env.SERVER_PLATFORM}生产环境 tag=${env.GIT_BRANCH},  ${env.JOB_NAME}"
    }
    else if (env.GIT_BRANCH ==~ /^test-([0-9])+\.([0-9])+\.([0-9])+.*/){
        branchName="${env.SERVER_PLATFORM}测试环境 tag=${env.GIT_BRANCH},  ${env.JOB_NAME}"
    }
    else {
        branchName="${env.SERVER_PLATFORM}开发环境 branch=${env.GIT_BRANCH},  ${env.JOB_NAME}"
    }

    // 发送内容
	json_msg=""
	if( res == "success" ) {
		json_msg='{\\"msgtype\\":\\"text\\",\\"text\\":{\\"content\\":\\"@' + tel_num +' [OK] ' + "${branchName} 第${env.BUILD_NUMBER}次构建，"  + '构建成功。 \\"},\\"at\\":{\\"atMobiles\\":[\\"' + tel_num + '\\"],\\"isAtAll\\":false}}'
	}
	else {
		json_msg='{\\"msgtype\\":\\"text\\",\\"text\\":{\\"content\\":\\"@' + tel_num +' [大哭] ' + "${branchName} 第${env.BUILD_NUMBER}次构建，"  + '构建失败，请及时处理！ \\"},\\"at\\":{\\"atMobiles\\":[\\"' + tel_num + '\\"],\\"isAtAll\\":false}}'
	}

    post_header="Content-Type:application/json;charset=utf-8"
    sh_cmd="curl -X POST " + dingding_url + " -H " + "\'" + post_header + "\'" + " -d " + "\""  + json_msg + "\""
	sh sh_cmd
}

void SendEmail(res)
{
	//在这里定义邮箱地址
	addr="xxx@xxx.com"
	if( res == "success" )
	{
		mail to: addr,
		subject: "构建成功 ：${currentBuild.fullDisplayName}",
		body: "\n发布成功。 \n\n任务名称： ${env.JOB_NAME} 第 ${env.BUILD_NUMBER} 次构建 \n\n 更多信息请查看 : ${env.BUILD_URL}"
	}
	else
	{
		mail to: addr,
		subject: "构建失败 ：${currentBuild.fullDisplayName}",
		body: "\n发布失败！ \n\n任务名称： ${env.JOB_NAME} 第 ${env.BUILD_NUMBER} 次构建 \n\n 更多信息请查看 : ${env.BUILD_URL}"
	}
}
