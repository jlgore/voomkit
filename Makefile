.DEFAULT_GOAL := all
BIN_FILE=voomkit

export DOCKER_BUILDKIT=1

DIVOOM_ADDR=192.168.88.53

.PHONY: help build run all

help: 
	@echo "Makefile arguments:"
	@echo "DIVOOM_ADDR=192.168.88.53"
	@echo "Makefile commands:"
	@echo "build"
	@echo "run"
	@echo "all"

build:
	@docker build --network=host -t ${BIN_FILE} .
run:
	@docker run --net=host -e DIVOOM_ADDR=${DIVOOM_ADDR} ${BIN_FILE}
all: build run
