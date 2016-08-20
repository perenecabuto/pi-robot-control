HOST := rastejante.local
USER := pi
DEPLOY_DIR := ~/robot-control
CFLAGS := "-ljpeg"


.PHONY=build
build:
	go build

.PHONY=cross
cross: deps
	CGO_ENABLED=1 CC=arm-linux-gnueabi-gcc GOOS=linux GOARCH=arm go build -v

.PHONY=deploy
deploy: cross
	ssh $(USER)@$(HOST) "mkdir -p $(DEPLOY_DIR)"
	scp -r robot-control $(USER)@$(HOST):$(DEPLOY_DIR)
	scp -r webapp $(USER)@$(HOST):$(DEPLOY_DIR)

.PHONY=stop
stop:
	ssh $(USER)@$(HOST) "killall -9 robot_control"

.PHONY=start
start: deploy stop
	ssh $(USER)@$(HOST) "cd $(DEPLOY_DIR); nohup ./robot_control &"

.PHONY=deps
deps:
	sudo aptitude install -y gccgo-arm-linux-gnueabi
