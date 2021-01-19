
build:
	go build -ldflags "-X 'main.BRANCH=`git rev-parse --abbrev-ref HEAD`' -X 'main.VERSION=`git log --pretty=format:"%h" -1`' -X 'main.BUILD_TIME=`date`' -X 'main.GO_VERSION=`go version`'"
	@echo build success

reset: build dropdb
	@./run.sh

run: build
	@./run.sh

clean:
	@rm -f *.gz
	@rm -f *.rdf

linux:
	@GOOS=linux GOARCH=amd64 go build -o ledger_committer_linux -ldflags "-X 'main.BRANCH=`git rev-parse --abbrev-ref HEAD`' -X 'main.VERSION=`git log --pretty=format:"%h" -1`' -X 'main.BUILD_TIME=`date`' -X 'main.GO_VERSION=`go version`'" main.go
	@echo build success


dropdb:
	@./dropdb.sh
