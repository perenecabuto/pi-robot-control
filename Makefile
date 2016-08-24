HOST := rastejante.local
USER := pi
DEPLOY_DIR := ~/robot-control
CFLAGS := "-ljpeg"


.PHONY=build
build:
	go build

.PHONY=cross
cross:
	CGO_ENABLED=1 CC=arm-linux-gnueabi-gcc GOOS=linux GOARCH=arm go build -v

.PHONY=deploy
deploy: cross
	ssh $(USER)@$(HOST) "mkdir -p $(DEPLOY_DIR)"
	scp -r robot-control $(USER)@$(HOST):$(DEPLOY_DIR)
	scp -r webapp $(USER)@$(HOST):$(DEPLOY_DIR)

.PHONY=stop
stop:
	ssh $(USER)@$(HOST) "eval killall -9 robot-control; echo 1"

.PHONY=start
start: stop
	ssh $(USER)@$(HOST) "cd $(DEPLOY_DIR) && nohup ./robot-control && disown"

.PHONY=deps
deps:
	sudo aptitude install -y gccgo-arm-linux-gnueabi
