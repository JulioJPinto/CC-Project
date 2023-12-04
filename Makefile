all: build-node build-tracker

build-node:
	mkdir bin
	go build -o bin/node ./node
	mkdir client_files/downloaded


build-tracker:
	go build -o bin/tracker ./tracker

clean:
	@rm bin -r
	@rm client_files/downloaded -r
