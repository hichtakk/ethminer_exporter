.PHONY: build clean

bin/ethminer_exporter:
	@go build -o bin/ethminer_exporter ./main.go

build: bin/ethminer_exporter

clean:
	@rm -rf ./bin
