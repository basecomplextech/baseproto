# baseProto Makefile
.PHONY: proto

install:
	@ go install ./cmd/baseprotoc

generate:
	@ go generate ./...

test:
	@ go test ./...

clean:
	@ find . -name '*pb.go' -delete
	@ find . -name '*_generated.go' -delete

proto:
	@ go generate ./proto/...
