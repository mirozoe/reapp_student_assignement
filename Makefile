GOBIN=$(shell go env GOBIN)

secret:=$(shell cat secret)

local: localbuild hidekey

localbuild: test fillinkey
	go build -o bin/abc main.go

test:
	go test -v -trace ./...

fillinkey:
	sed -i 's/KEY_PLACEHOLDER/$(secret)/' pkg/distance.go

hidekey:
	sed -i 's/key=[a-zA-Z0-9]\+"/key=KEY_PLACEHOLDER"/' pkg/distance.go

