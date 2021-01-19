
release:clean bin_dir macos
	@echo "build OK"
	@tree -al bin

clean:
	@rm -rf bin

bin_dir:
	@mkdir -p bin/macos
	@mkdir -p bin/linux

macos:bin_dir
	@go build -o bin/macos/argus -ldflags "-X 'main.BRANCH=`git rev-parse --abbrev-ref HEAD`' -X 'main.VERSION=`git log --pretty=format:"%h" -1`' -X 'main.BUILD_TIME=`date`' -X 'main.GO_VERSION=`go version`'" app/argus/main.go
	@echo arugs build success

linux:bin_dir
	@GOOS=linux GOARCH=amd64 go build -o bin/linux/argus -ldflags "-X 'main.BRANCH=`git rev-parse --abbrev-ref HEAD`' -X 'main.VERSION=`git log --pretty=format:"%h" -1`' -X 'main.BUILD_TIME=`date`' -X 'main.GO_VERSION=`go version`'" app/argus/main.go
	@echo arugs for **linux** build success
	@tree -al bin


