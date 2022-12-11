SOURCES := $(shell find ./ -type f -not -path "./node_modules/*" -not -path '*/.*' -not -path './build/*')

all: build node_modules

node_modules: yarn.lock package.json
	yarn install
	touch node_modules

build: ${SOURCES}
	yarn build
clean:
	rm -rf build/

clean-all:
	rm -rf node_modules/
	rm -rf build/

.PHONY: clean all clean-all
