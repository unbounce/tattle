.DEFAULT_GOAL: build
.PHONY: build test dist

DIST_DIR := $(GOPATH)/dist/tattle

build:
	go build -o tattle github.com/unbounce/tattle

# only works on Linux/EC2 machines
test:
	PAM_TYPE=open_session PAM_USER=johndoe PAM_RHOST=office ./tattle

dist:
	GOOS=linux GOARCH=amd64 go build -o $(DIST_DIR)/linux/amd64/tattle github.com/unbounce/tattle
	tar cfz $(DIST_DIR)/tattle.linux-amd64.tgz -C $(DIST_DIR)/linux/amd64 tattle
