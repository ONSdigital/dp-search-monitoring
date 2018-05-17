BINPATH ?= build

build: generate
	go build -tags 'production' -o $(BINPATH)/dp-search-monitoring

debug: generate
	go build -tags 'debug' -o $(BINPATH)/dp-search-monitoring
	HUMAN_LOG=1 DEBUG=1 $(BINPATH)/dp-search-monitoring

generate: ${GOPATH}/bin/go-bindata

package:
	gox -osarch="linux/amd64" -output="dp_search_monitoring_linux_amd64"
	zip dp-search-monitoring.zip dp_search_monitoring_linux_amd64 Dockerfile

test:
	go test -tags 'production' ./...

${GOPATH}/bin/go-bindata:
	go get -u github.com/jteeuwen/go-bindata/go-bindata

.PHONY: build debug
