# Makefile

QUIET := @
MKDIR := mkdir -p
RM    := rm -rf
GO    := go
HASH  := $(shell git rev-parse --short HEAD)

main := main.go
source := $(filter-out version_gen.go,$(wildcard cmd/*.go)) \
  $(wildcard pkg/*/*.go) $(main)
bin := mkpw


.PHONY: all clean update-version

all: update-version $(bin)

update-version:
	$(QUIET) sed -e 's/{hash}/${HASH}/g' ./cmd/version_gen.go.tpl > ./cmd/version_gen.go

$(bin): $(source) $(update-version)
	$(QUIET) $(GO) build -o $(notdir $@) $(main)

clean:
	$(QUIET) $(RM) $(bin)

