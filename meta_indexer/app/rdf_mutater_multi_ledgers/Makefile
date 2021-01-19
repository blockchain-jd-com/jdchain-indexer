
build:
	@go build
	@echo build success

run: build
	@./run.sh

rerun: build drop schema-update
	@./run.sh

clean:
	@rm -f *.gz
	@rm -f *.rdf

schema-update:
	@./schema_update.sh

drop:
	@./dropdb.sh
