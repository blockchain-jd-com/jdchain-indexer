
all:raw_searcher rdf_generator rdf_mutater_with_status rdf_mutater_multi_ledgers
	@echo "build OK"

clean:
	@rm -f app/raw_searcher/raw_searcher
	@rm -f app/rdf_generator/rdf_generator
	@rm -f app/rdf_mutater/rdf_mutater
	@rm -f app/rdf_mutater_with_status/rdf_mutater_with_status
	@rm -f app/rdf_mutater_multi_ledgers/rdf_mutater_multi_ledgers

raw_searcher:
	@go build -o raw_searcher -ldflags "-X 'main.BRANCH=`git rev-parse --abbrev-ref HEAD`' -X 'main.VERSION=`git log --pretty=format:"%h" -1`' -X 'main.BUILD_TIME=`date`' -X 'main.GO_VERSION=`go version`'" app/raw_searcher/main.go
	@mv raw_searcher app/raw_searcher/raw_searcher
	@echo raw_searcher build success

rdf_generator:
	@go build -o rdf_generator -ldflags "-X 'main.BRANCH=`git rev-parse --abbrev-ref HEAD`' -X 'main.VERSION=`git log --pretty=format:"%h" -1`' -X 'main.BUILD_TIME=`date`' -X 'main.GO_VERSION=`go version`'" app/rdf_generator/main.go
	@mv rdf_generator app/rdf_generator/rdf_generator
	@echo rdf_generator build success

rdf_mutater_with_status:
	@go build -o rdf_mutater_with_status -ldflags "-X 'main.BRANCH=`git rev-parse --abbrev-ref HEAD`' -X 'main.VERSION=`git log --pretty=format:"%h" -1`' -X 'main.BUILD_TIME=`date`' -X 'main.GO_VERSION=`go version`'" app/rdf_mutater_with_status/main.go
	@mv rdf_mutater_with_status app/rdf_mutater_with_status/rdf_mutater_with_status
	@echo rdf_mutater_with_status build success


rdf_mutater_multi_ledgers:
	@go build -o rdf_mutater_multi_ledgers -ldflags "-X 'main.BRANCH=`git rev-parse --abbrev-ref HEAD`' -X 'main.VERSION=`git log --pretty=format:"%h" -1`' -X 'main.BUILD_TIME=`date`' -X 'main.GO_VERSION=`go version`'" app/rdf_mutater_multi_ledgers/main.go
	@mv rdf_mutater_multi_ledgers app/rdf_mutater_multi_ledgers/rdf_mutater_multi_ledgers
	@echo rdf_mutater_multi_ledgers build success

