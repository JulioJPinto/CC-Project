all: build-node build-tracker

build-node:
	go build -o ../bin/node ./node


build-tracker:
	go build -o ../bin/tracker ./tracker
