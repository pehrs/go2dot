
APP=go2dot

GIT_VERSION=$(shell git describe --tags --abbrev=0)
BUILD_DATE=$(shell date +"%Y-%m-%dT%H:%M")
COMMIT=$(shell git show -s --format="%h" HEAD)

DIRTY_COUNT=$(shell git status --porcelain | wc -l)
ifeq ($(DIRTY_COUNT),"0")
  VERSION="$(GIT_VERSION)-dirty"
else
  VERSION="$(GIT_VERSION)-dirty"
endif

all: build

build:
	go build \
	--ldflags "-X 'pehrs.com/go2dot/cmd.version=$(VERSION)' -X 'pehrs.com/go2dot/cmd.date=$(BUILD_DATE)' -X 'pehrs.com/go2dot/cmd.commit=$(COMMIT)'" \
	-o bin/$(APP) main.go

