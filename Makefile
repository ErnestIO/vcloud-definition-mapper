install:
	go install -v

build:
	go build -v ./...

lint:
	golint ./...
	go vet ./...

test:
	go test -v ./... --cover

deps: dev-deps
	go get -u github.com/nats-io/nats
	go get -u github.com/r3labs/binary-prefix
	go get -u github.com/r3labs/workflow
	go get -u github.com/r3labs/graph
	go get -u github.com/ErnestIO/ernest-config-client

dev-deps:
	go get -u github.com/golang/lint/golint
	go get -u github.com/smartystreets/goconvey/convey

clean:
	go clean

dist-clean:
	rm -rf pkg src bin
