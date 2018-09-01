################
# Configuration
##############

EXECUTABLE_BASENAME = hopgoblin
DUMMYPROXY_BASENAME = dummyproxy

EXECUTABLE_DIR = build

DOC_PORT = 64079


############
# Internals
##########

EXECUTABLE_EXT :=
ifeq ($(OS),Windows_NT)
	EXECUTABLE_EXT := .exe
endif

EXECUTABLE_NAME = $(EXECUTABLE_DIR)/$(EXECUTABLE_BASENAME)$(EXECUTABLE_EXT)
DUMMYPROXY_NAME = $(EXECUTABLE_DIR)/$(DUMMYPROXY_BASENAME)$(EXECUTABLE_EXT)

ifdef ARGS
	RUN_ARGS = $(ARGS)
else
	RUN_ARGS =
endif


##########
# Targets
########

all: clean build test

clean:
	rm -f $(EXECUTABLE_DIR)/*

$(EXECUTABLE_NAME):
	go build -o $(EXECUTABLE_NAME) github.com/rinusser/hopgoblin/main

$(DUMMYPROXY_NAME):
	go build -o $(DUMMYPROXY_NAME) github.com/rinusser/hopgoblin/http/dummyproxy/main

build: $(EXECUTABLE_NAME)

run: $(EXECUTABLE_NAME)
	./$(EXECUTABLE_NAME) $(RUN_ARGS)

test: $(EXECUTABLE_NAME) $(DUMMYPROXY_NAME)
	go test ./... -failfast -bench . $(ARGS)

check:
	./check.sh

doc:
	godoc -http :$(DOC_PORT)

todo:
	/usr/bin/find . -name "*.go" | xargs grep -Pin --color=auto "(?<!log\.)(xxx|debug|todo)"
