#!/bin/bash

PROJECT_PDMQD=pdmqd
PROJECT_PDMQD_DIC=server/app/pdmqd
PROJECT_PDMQLOOPD=pdmqloopd
PROJECT_PDMQLOOPD_DIC=server/app/pdmqloopd
CURRENT_DIR=$(shell pwd)
UNAME=$(shell uname)

.PHONY:run

clean:
	rm -rf $(PROJECT_PDMQD_DIC)/$(PROJECT_PDMQD)
	rm -rf $(PROJECT_PDMQLOOPD_DIC)/$(PROJECT_PDMQLOOPD)

run:
	rm -rf $(PROJECT_PDMQD_DIC)/$(PROJECT_PDMQD)
	rm -rf $(PROJECT_PDMQLOOPD_DIC)/$(PROJECT_PDMQLOOPD)
	go build -o $(PROJECT_PDMQD_DIC)/$(PROJECT_PDMQD) $(PROJECT_PDMQD_DIC)/main.go
	go build -o $(PROJECT_PDMQLOOPD_DIC)/$(PROJECT_PDMQLOOPD) $(PROJECT_PDMQLOOPD_DIC)/main.go
	open -a Terminal.app $(PROJECT_PDMQLOOPD_DIC)/$(PROJECT_PDMQLOOPD)
	$(PROJECT_PDMQD_DIC)/$(PROJECT_PDMQD)