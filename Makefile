NAME=madmock
VERSION=0.0.1
BUILD=`git rev-parse HEAD`
PLATFORMS=darwin linux windows
ARCHITECTURES=386 amd64

FLAGS=-ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD}"

build:
	go build ${FLAGS} -o ${NAME}

all:
	$(foreach GOOS, $(PLATFORMS),\
	$(foreach GOARCH, $(ARCHITECTURES), $(shell export GOOS=$(GOOS); export GOARCH=$(GOARCH); go build -v -o ./build/$(NAME)-$(GOOS)-$(GOARCH))))

clean:
	rm -f madmock
	rm -rf build

.PHONY: build all clean
