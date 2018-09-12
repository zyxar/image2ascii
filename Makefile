GO          = go
PRODUCT     = image2ascii
GOARCH     := amd64
# VERSION    := $(shell git describe --all --always --dirty --long)
# BUILD_TIME := $(shell date +%FT%T%z)
# BUILDER    := $(shell whoami)
# GOVERSION  := $(shell $(GO) version | cut -d ' ' -f 3)
# LDFLAGS     = -ldflags "-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME} -X main.Builder=${BUILDER} -X main.GoVersion=${GOVERSION}"
GO111MODULE = on

all: $(shell $(GO) env GOOS)

build:
	env GO111MODULE=${GO111MODULE} GOOS=${GOOS} GOARCH=$(GOARCH) $(GO) build ${LDFLAGS} -mod vendor -v -o $(PRODUCT)$(EXT) .

linux: export GOOS=linux
linux: EXT=.elf
linux: build

darwin: export GOOS=darwin
darwin: build

.PHONY: clean
clean:
	@rm -f $(PRODUCT) $(PRODUCT).elf $(PRODUCT).mach
