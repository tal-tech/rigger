
build:
	export GOPROXY=http://goproxy.xesv5.com && \
	export GO111MODULE=on && \
	go build -o bin/rigger
clean:
	rm bin/*