HOST := rastejante.local
USER := pi
DEPLOY_DIR := ~/robot-control
VENDOR_DIR := $(PWD)/vendor
LIBJPEG_DIR := $(VENDOR_DIR)/libjpeg/
CFLAGS := "-I$(LIBJPEG_DIR)/include  -ljpeg"
LDFLAGS := "-L$(LIBJPEG_DIR)/lib -Wl,-rpath=\$$ORIGIN/vendor/libjpeg/lib/"


.PHONY=build
build: libjpeg_x86 install_libjpeg
	CGO_ENABLED=1 GOOS=linux CGO_CFLAGS=$(CFLAGS) CGO_LDFLAGS=$(LDFLAGS) go build -v

.PHONY=cross
cross: libjpeg_arm install_libjpeg
	CGO_ENABLED=1 CC=arm-linux-gnueabi-gcc GOOS=linux GOARCH=arm CGO_CFLAGS=$(CFLAGS) CGO_LDFLAGS=$(LDFLAGS) go build -v

.PHONY=deploy
deploy: cross
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
	cd $(VENDOR_DIR)/jpeg-9 && ./configure --host=arm-linux CC=arm-linux-gnueabi-gcc
	cd $(VENDOR_DIR)/jpeg-9 && make install DESTDIR=`pwd`

.PHONY=libjpeg_x86
libjpeg_x86: download_libjpeg
	cd $(VENDOR_DIR)/jpeg-9 && ./configure
	cd $(VENDOR_DIR)/jpeg-9 && make install DESTDIR=`pwd`

.PHONY=install_libjpeg
install_libjpeg: 
	cd $(VENDOR_DIR)/jpeg-9 && cp -r usr/local/include/ $(LIBJPEG_DIR)
	cd $(VENDOR_DIR)/jpeg-9 && cp -r usr/local/lib/ $(LIBJPEG_DIR)

.PHONY=download_libjpeg
download_libjpeg:
	mkdir -p $(LIBJPEG_DIR)
	cd $(VENDOR_DIR)  && wget -c http://www.ijg.org/files/jpegsrc.v9.tar.gz
	cd $(VENDOR_DIR)  && tar -xzvf jpegsrc.v9.tar.gz

.PHONY=clean
clean:
	git clean -dfx
