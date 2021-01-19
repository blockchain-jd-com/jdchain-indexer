
build:
	@go build
	@echo build success

reset: build dropdb
	@./run.sh

run: build
	@./run.sh

clean:
	@rm -f *.gz
	@rm -f *.rdf


dropdb:
	@./dropdb.sh
