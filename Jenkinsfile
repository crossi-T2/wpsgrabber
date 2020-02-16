def labels = ['c6-jenkins-go', 'c7-jenkins-go']
def builders = [:]

for (x in labels) {
    def label = x // Need to bind the label variable before the closure - can't do 'for (label in labels)'

    // Create a map to pass in to the 'parallel' step so we can fire all the builds at once
    builders[label] = {
        node(label) {
        
            env.GOPATH = env.WORKSPACE
            env.PATH="/usr/local/go/bin/:${GOPATH}:${PATH}"

            stage('Checkout'){
                echo 'Checking out SCM'
                checkout scm
            }

            stage('Init'){
                echo 'Pulling Dependencies'
                
                sh 'go version'
                sh 'go get -d ./...'

                echo 'Preparing rpmbuild workspace'
                echo '(requires rpm-build redhat-rpm-config rpmdevtools libtool)'
        
                sh 'mkdir -p $WORKSPACE/build/{BUILD,RPMS,SOURCES,SPECS,SRPMS}'
                sh 'cp build/package/wpsgrabber.spec $WORKSPACE/build/SPECS/'
                sh 'cp -r init $WORKSPACE/build/SOURCES/'
                sh 'cp -r configs $WORKSPACE/build/SOURCES/'
                sh 'spectool -g -R --directory $WORKSPACE/build/SOURCES $WORKSPACE/build/SPECS/wpsgrabber.spec'

                script {
                    def sdf = sh(returnStdout: true, script: 'date -u +%Y%m%dT%H%M%S').trim()
                    if (env.BRANCH_NAME == 'master')
                        env.release = env.BUILD_NUMBER
                    else
                        env.release = "SNAPSHOT" + sdf
                }
            }

            stage('Test'){
                echo 'Vetting'
                sh """cd $GOPATH && go vet cmd/wpsgrabber/*.go"""

                echo 'Testing'
                sh """cd $GOPATH && go test -race -cover cmd/wpsgrabber/*.go"""
            }

            stage('Build'){
                echo 'Building Executable'
                sh """go build -o $WORKSPACE/build/SOURCES/wpsgrabber -ldflags '-s' cmd/wpsgrabber/*.go"""
            }

            stage('Package') {
                echo 'Packaging'
                sh """rpmbuild --define \"_topdir $WORKSPACE/build\" -ba --define '_branch ${env.BRANCH_NAME}' --define '_release ${env.release}' $WORKSPACE/build/SPECS/wpsgrabber.spec"""
                sh 'rpm -qpl $WORKSPACE/build/RPMS/*/*.rpm'
            }

            stage('Publish') {
                echo 'Publishing'

                script {
                    // Obtain an Artifactory server instance, defined in Jenkins --> Manage:
                    def server = Artifactory.server "repository.terradue.com"

                    // Read the upload specs
                    def release = sh(returnStdout: true, script: 'rpm -E %{rhel}').trim()

                    def uploadSpec
                    
                    if (release == '6')	
                        uploadSpec = readFile 'build/deploy/artifactdeploy-el6.json'
                    else
                        if (release == '7')
                            uploadSpec = readFile 'build/deploy/artifactdeploy-el7.json'

                    // Upload files to Artifactory:
                    def buildInfo = server.upload spec: uploadSpec

                    // Publish the merged build-info to Artifactory
                    server.publishBuildInfo buildInfo
                }
            }
        }
    }
}

parallel builders
