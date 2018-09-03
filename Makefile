GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=speedtest
BINARY_LINUX=$(BINARY_NAME)_linux
BINARY_ARM5=$(BINARY_NAME)_arm5
BINARY_ARM7=$(BINARY_NAME)_arm7

all: test build
build: 
		$(GOBUILD) -o $(BINARY_NAME) -v
test: 
		$(GOTEST) -v ./...
clean: 
		$(GOCLEAN)
		rm -f $(BINARY_NAME)
		rm -f $(BINARY_LINUX)
		rm -f $(BINARY_ARM5)
		rm -f $(BINARY_ARM7)
run:
		$(GOBUILD) -o $(BINARY_NAME) -v ./...
		./$(BINARY_NAME)
# deps:
# 		$(GOGET) github.com/markbates/goth
# 		$(GOGET) github.com/markbates/pop

# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_LINUX) -v

# docker:
# 	docker run --rm -it -v "$(GOPATH)":/go -w /go/src/bitbucket.org/rsohlich/makepost golang:latest go build -o "$(BINARY_LINUX)" -v

build-arm5:
	GOOS=linux GOARCH=arm GOARM=5 $(GOBUILD) -o $(BINARY_ARM5)

build-arm7:
	GOOS=linux GOARCH=arm GOARM=7 $(GOBUILD) -o $(BINARY_ARM7)
