all: 
	@make build-node 
	@make build-tracker

build-node:
	@go build -o ../bin/node ./node
	@mkdir client_files/downloaded -p


build-tracker:
	@go build -o ../bin/tracker ./tracker

clean:
	@rm ../bin -r
	@rm client_files/downloaded -r
