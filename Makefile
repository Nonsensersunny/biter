# GOOS：darwin、freebsd、linux、windows
# GOARCH：386、amd64、arm、s390x

.PHONY: prepare darwin linux windows publish clean

all: publish

publish:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o bin/biter-mac ./cmd/main.go
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/biter ./cmd/main.go
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/biter.exe ./cmd/main.go

darwin:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o bin/biter ./cmd/main.go

linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/biter ./cmd/main.go

windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/biter.exe ./cmd/main.go

prepare:
	go mod vendor

clean:
	rm -rf ./bin