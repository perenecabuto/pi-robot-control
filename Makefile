HOST := rastejante.local
USER := pi
DEPLOY_DIR := ~/robot-control

VENDOR_DIR := $(PWD)/vendor
LIBJPEG_DIR := $(VENDOR_DIR)/libjpeg
LIBJPEG_SRC_NAME := jpeg-9
#LIBJPEG_SRC_NAME := libjpeg-turbo-1.5.0
LIBJPEG_INSTALL_DIR := $(VENDOR_DIR)/$(LIBJPEG_SRC_NAME)/usr/local
#LIBJPEG_INSTALL_DIR := $(VENDOR_DIR)/$(LIBJPEG_SRC_NAME)/opt/libjpeg-turbo
LIBJPEG_SRC_TGZ := jpegsrc.v9.tar.gz
#LIBJPEG_SRC_TGZ := libjpeg-turbo-1.5.0.tar.gz
LIBJPEG_SRC_URL := http://www.ijg.org/files/$(LIBJPEG_SRC_TGZ)
#LIBJPEG_SRC_URL := http://ufpr.dl.sourceforge.net/project/libjpeg-turbo/1.5.0/$(LIBJPEG_SRC_TGZ)
ARM_COMPILER := arm-linux-gnueabihf-gcc

CFLAGS := "-I$(LIBJPEG_DIR)/include  -ljpeg -O3"
LDFLAGS := "-L$(LIBJPEG_DIR)/lib -Wl,-rpath=\$$ORIGIN/vendor/libjpeg/lib/ -O3"

remote = @ssh $(USER)@$(HOST)

.PHONY=build
build: libjpeg-x86 install-libjpeg
	CGO_ENABLED=1 GOOS=linux CGO_CFLAGS=$(CFLAGS) CGO_LDFLAGS=$(LDFLAGS) go build -v

.PHONY=test
test:
	CompileDaemon -color -command "go test ./..."

.PHONY=copy-ssh-key
copy-ssh-key:
	ssh-copy-id -i ~/.ssh/id_rsa.pub $(USER)@$(HOST)

.PHONY=cross-arm
cross-arm: libjpeg-arm install-libjpeg
	CGO_ENABLED=1 CC=$(ARM_COMPILER) CGO_CFLAGS=$(CFLAGS) CGO_LDFLAGS=$(LDFLAGS) GOOS=linux GOARCH=arm go build

.PHONY=deploy upload
deploy: cross-arm upload
	@echo Deployed to $(HOST)

upload:
	$(remote) "mkdir -p $(DEPLOY_DIR)"
	$(remote) "mkdir -p $(DEPLOY_DIR)/vendor"
	scp -r webapp $(USER)@$(HOST):$(DEPLOY_DIR)
	scp -r robot-control $(USER)@$(HOST):$(DEPLOY_DIR)
	#scp -r $(LIBJPEG_DIR) $(USER)@$(HOST):$(DEPLOY_DIR)/vendor

.PHONY=stop
stop:
	$(remote) "eval killall -9 robot-control; echo killed"

.PHONY=start
start: stop
	$(remote) "cd $(DEPLOY_DIR) && eval nohup ./robot-control && disown"

.PHONY=remote-deps
remote-deps:
	$(remote) sudo aptitude install -y i2c-tools
	$(remote) sudo aptitude install -y libjpeg9-dev

.PHONY=deps
deps:
	sudo aptitude install -y gccgo-arm-linux-gnueabi

.PHONY=libjpeg-arm
libjpeg-arm: download-libjpeg
	cd $(VENDOR_DIR)/$(LIBJPEG_SRC_NAME) && ./configure --host=arm-linux CC=$(ARM_COMPILER)
	cd $(VENDOR_DIR)/$(LIBJPEG_SRC_NAME) && make clean install DESTDIR=`pwd`

.PHONY=libjpeg-x86
libjpeg-x86: download-libjpeg
	cd $(VENDOR_DIR)/$(LIBJPEG_SRC_NAME) && ./configure
	cd $(VENDOR_DIR)/$(LIBJPEG_SRC_NAME) && make clean install DESTDIR=`pwd`

.PHONY=install-libjpeg
install-libjpeg:
	cd $(VENDOR_DIR)/$(LIBJPEG_SRC_NAME) && cp -r $(LIBJPEG_INSTALL_DIR)/include/ $(LIBJPEG_DIR)
	cd $(VENDOR_DIR)/$(LIBJPEG_SRC_NAME) && cp -r $(LIBJPEG_INSTALL_DIR)/lib*/ $(LIBJPEG_DIR)/lib/

.PHONY=download-libjpeg
download-libjpeg:
	mkdir -p $(LIBJPEG_DIR)
	cd $(VENDOR_DIR)  && wget -c $(LIBJPEG_SRC_URL)
	cd $(VENDOR_DIR)  && tar -xzf $(LIBJPEG_SRC_TGZ)

.PHONY=clean
clean:
	git clean -dfx
