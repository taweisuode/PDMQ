#!/bin/bash

PROJECT=pdmqd
PROJECT_DIC=server/app/pdmqd
CURRENT_DIR=$(shell pwd)
UNAME=$(shell uname)

.PHONY:common

common:
	rm -rf $(PROJECT_DIC)/$(PROJECT)

	cd _publish_dir #test dir is exsit

	go build -o $(PROJECT_DIC)/$(PROJECT) $(PROJECT_DIC)/main.go
clean:
	rm -rf $(PROJECT_DIC)/$(PROJECT)

run:
	rm -rf $(PROJECT_DIC)/$(PROJECT)
	go build -o $(PROJECT_DIC)/$(PROJECT) $(PROJECT_DIC)/main.go
	$(PROJECT_DIC)/$(PROJECT)
