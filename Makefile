
APP=go2dot

GIT_VERSION=$(shell git describe --tags --abbrev=0)
BUILD_DATE=$(shell date +"%Y-%m-%dT%H:%M")
COMMIT=$(shell git show -s --format="%h" HEAD)

DIRTY_COUNT=$(shell git status --porcelain | wc -l)
ifeq ($(DIRTY_COUNT),0)
  VERSION="$(GIT_VERSION)"
else
  VERSION="$(GIT_VERSION)-dirty"
endif

all: build

build:				## Build binary
	go build \
	--ldflags "-X 'pehrs.com/go2dot/cmd.version=$(VERSION)' -X 'pehrs.com/go2dot/cmd.date=$(BUILD_DATE)' -X 'pehrs.com/go2dot/cmd.commit=$(COMMIT)'" \
	-o bin/$(APP) main.go

generate-png: build		## Generate a PNG from ./pkg/golang
	rm -f samples/test*
	./bin/go2dot dot ./pkg/golang > samples/test-golang-public.dot \
		&& dot -Gfontname="Sans" -Nfontname="Serif" -Gsize=4,3 -Gdpi=1000 -Tpng samples/test-golang-public.dot -o samples/test-golang-public.png \
		&& dot -Gfontname="Courier" -Nfontname="Courier" -Tsvg samples/test-golang-public.dot -o samples/test-golang-public.svg
	./bin/go2dot dot -p ./pkg/golang > samples/test-golang.dot && dot -T png samples/test-golang.dot -o samples/test-golang.png
	./bin/go2dot dot -p ./cmd > samples/test-cmd.dot && dot -T png samples/test-cmd.dot -o samples/test-cmd.png
	./bin/go2dot dot -p ./ > samples/test-main.dot && dot -T png samples/test-main.dot -o samples/test-main.png


help:				## Show this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m\033[0m\n"} /^[$$()% a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
