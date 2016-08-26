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

CFLAGS := "-I$(LIBJPEG_DIR)/include  -ljpeg -O3"
LDFLAGS := "-L$(LIBJPEG_DIR)/lib -Wl,-rpath=\$$ORIGIN/vendor/libjpeg/lib/ -O3"


.PHONY=build
build: libjpeg_x86 install_libjpeg
	CGO_ENABLED=1 GOOS=linux CGO_CFLAGS=$(CFLAGS) CGO_LDFLAGS=$(LDFLAGS) go build -v

.PHONY=cross_arm
cross_arm: libjpeg_arm install_libjpeg
	CGO_ENABLED=1 CC=arm-linux-gnueabi-gcc GOOS=linux GOARCH=arm CGO_CFLAGS=$(CFLAGS) CGO_LDFLAGS=$(LDFLAGS) go build -v

.PHONY=deploy
deploy: clean cross_arm
	ssh $(USER)@$(HOST) "mkdir -p $(DEPLOY_DIR)"
	ssh $(USER)@$(HOST) "mkdir -p $(DEPLOY_DIR)/vendor"
	scp -r robot-control $(USER)@$(HOST):$(DEPLOY_DIR)
	scp -r webapp $(USER)@$(HOST):$(DEPLOY_DIR)
	scp -r $(LIBJPEG_DIR) $(USER)@$(HOST):$(DEPLOY_DIR)/vendor

.PHONY=stop
stop:
	ssh $(USER)@$(HOST) "eval killall -9 robot-control; echo 1"

.PHONY=start
start: stop
	ssh $(USER)@$(HOST) "cd $(DEPLOY_DIR) && nohup ./robot-control && disown"

.PHONY=deps
deps:
	sudo aptitude install -y gccgo-arm-linux-gnueabi

.PHONY=libjpeg_arm
libjpeg_arm: download_libjpeg
	cd $(VENDOR_DIR)/$(LIBJPEG_SRC_NAME) && ./configure --host=arm-linux CC=arm-linux-gnueabi-gcc
	cd $(VENDOR_DIR)/$(LIBJPEG_SRC_NAME) && make install DESTDIR=`pwd`

.PHONY=libjpeg_x86
libjpeg_x86: download_libjpeg
	cd $(VENDOR_DIR)/$(LIBJPEG_SRC_NAME) && ./configure
	cd $(VENDOR_DIR)/$(LIBJPEG_SRC_NAME) && make install DESTDIR=`pwd`

.PHONY=install_libjpeg
install_libjpeg:
	cd $(VENDOR_DIR)/$(LIBJPEG_SRC_NAME) && cp -r $(LIBJPEG_INSTALL_DIR)/include/ $(LIBJPEG_DIR)
	cd $(VENDOR_DIR)/$(LIBJPEG_SRC_NAME) && cp -r $(LIBJPEG_INSTALL_DIR)/lib*/ $(LIBJPEG_DIR)/lib/

.PHONY=download_libjpeg
download_libjpeg:
	mkdir -p $(LIBJPEG_DIR)
	cd $(VENDOR_DIR)  && wget -c $(LIBJPEG_SRC_URL)
	cd $(VENDOR_DIR)  && tar -xzf $(LIBJPEG_SRC_TGZ)

.PHONY=clean
clean:
	git clean -dfx
