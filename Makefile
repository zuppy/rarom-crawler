#########
#
# To start, just run "make" in this directory.
# You can override the params defined below, like this:
# make DOCKER_IMG_TAG=my_image_name
#
#########

# the following params can be overriden at runtime (ie: make APP_NAME=xx), these are the default values
APP_NAME?=$(lastword $(subst /, ,$(shell pwd)))
DOCKER_IMG_TAG?=$(APP_NAME)_build:1.11
BIN_PATH?=$(shell pwd)/../../bin
PKG_PATH?=$(shell pwd)/../../pkg

# these can not be overriden
CUSTOM_RUN_FILE=run-app.sh
GO_ARG_LINUX=GOOS=linux GOARCH=amd64
GO_ARG_MAC=GOOS=darwin GOARCH=amd64
GO_ARG_WIN=GOOS=windows GOARCH=amd64


all: build

build-linux: build-linux-docker

build-mac: build-mac-local

test: test-docker

info:
	@echo "\n████▓▓▓▒▒░ current arguments ░▒▒▓▓▓████"
	@echo "APP_NAME       = $(APP_NAME)"
	@echo "DOCKER_IMG_TAG = $(DOCKER_IMG_TAG)"
	@echo "BIN_PATH       = $(BIN_PATH)"
	@echo "PKG_PATH       = $(PKG_PATH)"

build:
	@echo "\n████▓▓▓▒▒░ building unspecified image ░▒▒▓▓▓████"
	@echo "You have to select one of the following build variants:\n"
	@echo "make build-linux-docker # will build through docker a linux binary"
	@echo "make build-linux-local  # will build through local go a linux binary"
	@echo "make build-mac-local    # will build through local go a mac binary\n"
	@echo "make build-windows-docker # will build through docker a windows binary"
	@echo "The recommended method is through Jenkins, via Jenkinsfile\n"

build-docker-image:
	@echo "\n████▓▓▓▒▒░ building docker image ░▒▒▓▓▓████"
	docker build -t $(DOCKER_IMG_TAG) .

clean-docker-image:
	@echo "\n████▓▓▓▒▒░ removing docker image ░▒▒▓▓▓████"
	docker rmi $(DOCKER_IMG_TAG)

test-local:
	@echo "\n████▓▓▓▒▒░ running local tests ░▒▒▓▓▓████"
	dep ensure && env go test -race -v

test-docker: build-docker-image
	@echo "\n████▓▓▓▒▒░ running docker tests ░▒▒▓▓▓████"
	docker run --rm -v $(shell pwd):/go/src/app \
		-v $(PKG_PATH):/go/pkg \
		-v $(BIN_PATH):/go/bin \
	    -w /go/src/app $(DOCKER_IMG_TAG) /bin/bash \
	    -c "dep ensure && env go test -race -v"

build-linux-docker: info clean build-docker-image
	@echo "\n████▓▓▓▒▒░ building docker linux binary ░▒▒▓▓▓████"
	echo $(CWD)
	docker run --rm -v $(shell pwd):/go/src/app \
		-v $(PKG_PATH):/go/pkg \
		-v $(BIN_PATH):/go/bin \
	    -w /go/src/app $(DOCKER_IMG_TAG) /bin/bash \
	    -c "dep ensure && env $(GO_ARG_LINUX) go build -race -v \
	    -o /go/bin/$(APP_NAME) -ldflags \
	    	\"-X main.version=tag=$(git tag -l --points-at HEAD);commit=$(shell git rev-parse --short HEAD)\""
	@echo "\nOutput path for binary: $(BIN_PATH)/$(APP_NAME)\n"

build-windows-docker: info clean build-docker-image
	@echo "\n████▓▓▓▒▒░ building docker windows binary ░▒▒▓▓▓████"
	echo $(CWD)
	docker run --rm -v $(shell pwd):/go/src/app \
		-v $(PKG_PATH):/go/pkg \
		-v $(BIN_PATH):/go/bin \
	    -w /go/src/app $(DOCKER_IMG_TAG) /bin/bash \
	    -c "dep ensure && env $(GO_ARG_WIN) go build -race -v \
	    -o /go/bin/$(APP_NAME) -ldflags \
	    	\"-X main.version=tag=$(git tag -l --points-at HEAD);commit=$(shell git rev-parse --short HEAD)\""
	@echo "\nOutput path for binary: $(BIN_PATH)/$(APP_NAME)\n"

build-linux-local: info clean
	@echo "\n████▓▓▓▒▒░ building local linux binary ░▒▒▓▓▓████"
	dep ensure && env $(GO_ARG_LINUX) go build -race -v -o $(BIN_PATH)/$(APP_NAME) -ldflags \
	    "-X main.version=tag=$(git tag -l --points-at HEAD);commit=$(shell git rev-parse --short HEAD)"
	@echo "\nOutput path for binary: $(BIN_PATH)/$(APP_NAME)\n"

build-mac-docker:
	@echo "\n████▓▓▓▒▒░ building docker mac binary ░▒▒▓▓▓████"
	@echo "This is not a supported build option. See https://github.com/golang/go/issues/29170

build-mac-local: info clean
	@echo "\n████▓▓▓▒▒░ building local mac binary ░▒▒▓▓▓████"
	dep ensure && env $(GO_ARG_MAC) go build -race -v -o $(BIN_PATH)/$(APP_NAME) -ldflags \
	    "-X main.version=tag=$(git tag -l --points-at HEAD);commit=$(shell git rev-parse --short HEAD)"
	@echo "\nOutput path for binary: $(BIN_PATH)/$(APP_NAME)\n"

run:
	@echo "\n████▓▓▓▒▒░ run ░▒▒▓▓▓████"
ifneq (,$(wildcard $(shell pwd)/$(CUSTOM_RUN_FILE)))
	@echo Custom run file found. Running $(CUSTOM_RUN_FILE) $(BIN_PATH)/$(APP_NAME):
	$(shell pwd)/$(CUSTOM_RUN_FILE) $(BIN_PATH)/$(APP_NAME)
else
	@echo Custom run file not found, running -h on the built binary:
	$(BIN_PATH)/$(APP_NAME) -h
endif

clean:
	@echo "\n████▓▓▓▒▒░ deleting output binary ░▒▒▓▓▓████"
	@echo "Output file to delete: $(BIN_PATH)/$(APP_NAME)"
# this also covers /, as it is not a readable file
ifneq (,$(wildcard $(BIN_PATH)/$(APP_NAME)))
	@rm $(BIN_PATH)/$(APP_NAME)
	@echo Deleted
else
	@echo "Nothing to delete, file does not exist"
endif
