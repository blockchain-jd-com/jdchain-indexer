
all: build
	@echo "all build OK"

build:
	@go build -ldflags "-X 'main.BRANCH=`git rev-parse --abbrev-ref HEAD`' -X 'main.VERSION=`git log --pretty=format:"%h" -1`' -X 'main.BUILD_TIME=`date`' -X 'main.GO_VERSION=`go version`'"
	@echo build success

release: build
	@mkdir -p ../../../bin
	@mv api_server ../../../bin/api_server
	@echo release *api_server* success

run: build
	@./api_server

