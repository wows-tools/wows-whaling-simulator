build:
	go build
	$(MAKE) -C ui
	rsync -Pizza ui/build/ static/

build-static:
	CGO_ENABLED=0 go build -ldflags "-s -w"
	$(MAKE) -C ui
	rsync -Pizza ui/build/ static/

test:
	go test

clean:
	rm -f wows-whaling-simulator
	$(MAKE) -C ui
