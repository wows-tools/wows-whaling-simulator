build:
	go build
	$(MAKE) -C ui
	rsync -Pizza ui/build/ static/

static:
	CGO_ENABLED=0 go build -ldflags "-s -w"
	$(MAKE) -C ui
	rsync -Pizza ui/build/ static/

test:
	go test

clean:
	rm -f wows-whaling-simulator
	rm -f wows-whaling-simulator
	rm -f static/asset-manifest.json
	rm -f static/christmas.png
	rm -f static/index.html
	rm -f static/logo192.png
	rm -f static/logo512.png
	rm -f static/manifest.json
	rm -f static/robots.txt
	rm -rf static/static/
	rm -rf ui/build/

.PHONY: build clean static test
