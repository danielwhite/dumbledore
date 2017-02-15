 # Options

SOURCES = $(wildcard *.go)
PLUGINS = prune.so clone.so

# Targets

.PHONY: all
all: dumbledore plugins

.PHONY: clean
clean:
	-rm -f dumbledore *.so

dumbledore:
	go build

.PHONY: plugins
plugins: ${PLUGINS}

%.so: plugin/%
	go build -buildmode=plugin ./$^

.PHONY: run
run:
	docker build -t dumbledore .
	docker run -it --rm -v $(CURDIR)/testdata:/etc/dumbledore -p 8080:8080 dumbledore
