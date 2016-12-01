EXECUTABLE=bin/rpiservo-linux-arm
REMOTE_EXECUTABLE=rpiservo
RPI_USER=pi
RPI_PASS=goncalonsp12
RPI_IP=192.168.1.137

all: build

build: src/
	env GOOS=linux GOARCH=arm gb build all

deploy: build
	./deploy-to-pi.sh $(EXECUTABLE) $(REMOTE_EXECUTABLE) $(RPI_USER) $(RPI_IP) $(RPI_PASS)

clean:
	-rm $(EXECUTABLE) || true

.PHONY: clean