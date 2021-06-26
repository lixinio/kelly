REPO = github.com/lixinio/kelly
BINARIES:=helloworld write route plugin static binder jwtauth openidauth session auth multi trace trace2 swagger

all: $(BINARIES)
build: $(BINARIES)

$(BINARIES):
	mkdir -p build
	echo "$@"
	find ./examples/ -mindepth 1 -maxdepth 1 -type d -name "*_$@" | \
		grep "$@" |\
		xargs go build -mod=vendor -o build

.PHONY: mod
mod:
	go mod tidy && go mod vendor
