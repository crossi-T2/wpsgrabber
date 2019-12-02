node('c7-jenkins-go') {

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

		//List all our project files with 'go list ./... | grep -v /vendor/ | grep -v github.com | grep -v golang.org'
		//Push our project files relative to ./src
		sh 'cd $GOPATH && go list ./... | grep -v /vendor/ | grep -v github.com | grep -v golang.org > projectPaths'

		//Print them with 'awk '$0="./src/"$0' projectPaths' in order to get full relative path to $GOPATH
		def paths = sh returnStdout: true, script: """awk '\$0="./src/"\$0' projectPaths"""
		
		echo 'Vetting'

		//sh """cd $GOPATH && go tool vet ${paths}"""

		echo 'Linting'
		//sh """cd $GOPATH && golint ${paths}"""

		echo 'Testing'
		//sh """cd $GOPATH && go test -race -cover ${paths}"""
	}

	stage('Build'){
		echo 'Building Executable'

		sh """go build -o wpsgrabber -ldflags '-s' cmd/wpsgrabber/*.go"""
	}

	stage('Package') {
        	echo 'Packaging'
        	sh """sudo rpmbuild --define \"_topdir $WORKSPACE/build\" -ba --define '_branch ${env.BRANCH_NAME}' --define '_release ${env.release}' $WORKSPACE/build/SPECS/wpsgrabber.spec"""
        	sh 'rpm -qpl $WORKSPACE/build/RPMS/*/*.rpm'
    	}

	stage('Publish') {
        	echo 'Publishing'
        	script {
            		// Obtain an Artifactory server instance, defined in Jenkins --> Manage:
            		def server = Artifactory.server "repository.terradue.com"

            		// Read the upload specs:
            		def uploadSpec = readFile 'artifactdeploy.json'

            		// Upload files to Artifactory:
            		def buildInfo = server.upload spec: uploadSpec

            		// Publish the merged build-info to Artifactory
            		server.publishBuildInfo buildInfo
        	}
    	}
}
