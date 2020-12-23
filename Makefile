.PHONY: build
build:
	go build -o ./build/server ./cmd/server/main.go
	go build -o ./build/fileserver ./cmd/fileserver/main.go
	go build -o ./build/session ./cmd/session/main.go

.PHONY: test
test:
	go test ./...

.PHONY: cover
cover:
	go test -coverprofile=coverage.out -coverpkg=./... -cover ./... && cat coverage.out | grep -v _mock | grep -v _easyjson.go | grep -v _easyjson.go | grep -v pb.go > cover.out && go tool cover -func=cover.out
	rm -f *.out

# .PHONY: build
# build:
# 	go build -o app ./cmd/server/main.go
