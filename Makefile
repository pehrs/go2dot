
APP=go2dot

VERSION=$(shell git describe --tags --abbrev=0)
ifndef $$VERSION
VERSION := "0.0.1"
endif
BUILD_DATE=$(shell date +"%Y-%m-%dT%H:%M")


all: build

build:
	go build \
	--ldflags "-X 'pehrs.com/go2dot/cmd.version=$(VERSION)' -X 'pehrs.com/go2dot/cmd.date=$(BUILD_DATE)' -X 'pehrs.com/go2dot/cmd.commit=e1638af69'" \
	-o bin/$(APP) main.go

