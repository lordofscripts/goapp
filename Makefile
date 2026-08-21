# - ⚙️ Builder settings
#   BFILTER :: gocolor|gofilter|""
#	MODE :: release|debug		(Example: make MODE=debug target)
BFILTER=gocolor
MODE ?= release

# - ⚙️ For each target recipe, the source directory and the output binary
BIN=./bin
SOURCES=cmd/examples

SRC_DEMO_APP=$(SOURCES)/demo_app
BIN_DEMO_APP=$(BIN)/demo_goapp
SRC_DEMO_MLOG=$(SOURCES)/demo_mlog
BIN_DEMO_MLOG=$(BIN)/demo_mlog
SRC_DEMO_ZLOG=$(SOURCES)/demo_zlog
BIN_DEMO_ZLOG=$(BIN)/demo_zlog

# ------------------------ 🚫 ---------------------------

# - 🚫 Select colored compiler errors or plain unreadable GO default
GO=go
GOPRETTY=""
ifeq ($(BFILTER),gocolor)
    GOPRETTY="$(HOME)/go/bin/gocolor"
else ifeq ($(BFILTER),gofilter)
    GOPRETTY="$(HOME)/go/bin/gofilter"
endif

SUFFIX=""
ifeq ($(MODE),debug)
    SUFFIX="_dev"
endif

# - 🚫 Packagers only (DO NOT MODIFY)
PKG_SEMANTIC_VERSION=$(shell sed -n '/>>>BEGIN/,/>>>END/p' version.go > /tmp/mainversion_gap.go && go run /tmp/mainversion_gap.go short)
PKG_FULL_VERSION:= $(patsubst v%,%,$(PKG_SEMANTIC_VERSION))
PKG_PUBLIC_NAME=$(shell go list -m')
PKG_NAME=goapp
PKG_PREFERRED_ARCH="amd64"
PKG_REVISION=1
PKG_FULLNAME=${PKG_NAME}_${PKG_FULL_VERSION}-${PKG_REVISION}-${PKG_PREFERRED_ARCH}


# ------------------------ 🚫 ---------------------------

dapp:
ifeq ($(BFILTER),)
	$(GO) build -v -o $(BIN_DEMO_APP)$(SUFFIX) $(SRC_DEMO_APP)/*.go
else
	$(GO) build  -v -o $(BIN_DEMO_APP)$(SUFFIX) $(SRC_DEMO_APP)/*.go 2>&1 | $(GOPRETTY) -color -width 75 -version
endif	

dmlog:
ifeq ($(BFILTER),)
	$(GO) build -tags mlog,$(MODE) -v -o $(BIN_DEMO_MLOG)$(SUFFIX) $(SRC_DEMO_MLOG)/*.go
else
	$(GO) build -tags mlog,$(MODE) -v -o $(BIN_DEMO_MLOG)$(SUFFIX) $(SRC_DEMO_MLOG)/*.go 2>&1 | $(GOPRETTY) -color -width 75 -version
endif

dzlog:
ifeq ($(BFILTER),)
	$(GO) build -tags zlog,$(MODE) -v -o $(BIN_DEMO_ZLOG)$(SUFFIX) $(SRC_DEMO_ZLOG)/*.go
else
	$(GO) build -tags zlog,$(MODE) -v -o $(BIN_DEMO_ZLOG)$(SUFFIX) $(SRC_DEMO_ZLOG)/*.go 2>&1 | $(GOPRETTY) -color -width 75 -version
endif

dummy:
	@echo "BFILTER $(BFILTER)"
	@echo "PRETTY $(GOPRETTY)"
	@echo "MODE $(MODE)"
	@echo "BIN bin/demo$(SUFFIX)"
ifeq ($(BFILTER),)
	@echo "empty"
else
	@echo "use post poroces"
endif	

# ------------------------ 🚫 ---------------------------

# Publish package info to GO Package Repository
proxy:
	GOPROXY=proxy.golang.org go list -m $(PKG_PUBLIC_NAME)@v$(PKG_FULL_VERSION)

update:
	go get -u all

clean:
	go clean

testall:
	go test ./...

testfull:
	go test -v tests/*_test.go

version:
	@echo $(PKG_FULL_VERSION)

name:
	@echo $(PKG_FULLNAME)	
