SOURCES := $(shell find ./ -type f -not -path "./ui*" -not -path '*/.*' -not -path './build/*' -not -name '*_test.go' -name '*.go') \
	   static/asset-manifest.json static/*.png static/index.html static/static/*/* static/resources/* static/manifest.json static/robots.txt

all:
	$(MAKE) -C ui
	rsync -Pizza ui/build/ static/ --exclude=/resources --delete
	$(MAKE) wows-whaling-simulator wows-whaling-simulator-static wows-whaling-simulator-cli wows-whaling-simulator-cli-static

wows-whaling-simulator: $(SOURCES)
	go build

wows-whaling-simulator-cli: $(SOURCES)
	go build -o wows-whaling-simulator-cli misc/cli/main.go

wows-whaling-simulator-static: $(SOURCES)
	CGO_ENABLED=0 go build -ldflags "-s -w" -o wows-whaling-simulator-static

wows-whaling-simulator-cli-static: $(SOURCES)
	CGO_ENABLED=0 go build -ldflags "-s -w" -o wows-whaling-simulator-cli-static misc/cli/main.go

test:
	go test

clean:
	$(MAKE) -C ui clean
	rm -f wows-whaling-simulator
	rm -f wows-whaling-simulator-static
	rm -f static/asset-manifest.json
	rm -f static/christmas.png
	rm -f static/index.html
	rm -f static/logo192.png
	rm -f static/logo512.png
	rm -f static/manifest.json
	rm -f static/robots.txt
	rm -rf static/static/
	rm -rf ui/build/

clean-all:
	$(MAKE) clean
	$(MAKE) -C ui clean-all

.PHONY: clean test clean-all all
