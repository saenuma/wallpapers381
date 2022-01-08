
build:
	rm -rf bin/
	mkdir -p bin/
	go build -o bin/ubuntu ./ubuntu
	go build -o bin/ubuntu_switch ./ubuntu_switch
	go build -o bin/windows ./windows

