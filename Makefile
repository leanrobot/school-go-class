PACKAGES=...
TEST_PACKAGE=...

GODOC_PORT=:8080

all: fmt install

build: 
	GOPATH="$(CURDIR)"
	go build $(PACKAGES)

install:
	GOPATH="$(CURDIR)"
	go install $(PACKAGES)

test:
	GOPATH="$(CURDIR)"
	go test $(PACKAGES)

fmt:
	GOPATH="$(CURDIR)"
	go fmt $(PACKAGES)

doc:
	GOPATH="$(CURDIR)"
	godoc -v --http=$(GODOC_PORT) --index=true

clean:
	rm -f bin/*
	rm -rf pkg/*
	rm -rf log/*
