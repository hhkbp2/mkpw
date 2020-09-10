# Makefile

QUIET   := @
MKDIR   := mkdir -p
RM      := rm -rf
GO      := GO111MODULE=on go
VERSION := 0.1
HASH    := $(shell git rev-parse --short HEAD)

VERSIONFLAGS += -X 'mkpw/cmd.Version=$(VERSION)'
VERSIONFLAGS += -X 'mkpw/cmd.GitHash=$(HASH)'
LDFLAGS += -ldflags "$(VERSIONFLAGS)"

main := main.go
source := $(wildcard cmd/*.go) $(wildcard pkg/*/*.go) $(main)
bin := mkpw


.PHONY: all clean

all: $(bin)

$(bin): $(source)
	$(QUIET) $(GO) build $(LDFLAGS) -o $(notdir $@) $(main)

clean:
	$(QUIET) $(RM) $(bin)

