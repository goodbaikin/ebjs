REPO=goodbaikin
MODULE=github.com/${REPO}/ebjs
API=./api
BIN=bin/ebjs bin/ebjclient

DOCKER=sudo ${shell which docker}

GIT_COMMIT:=${shell git rev-parse HEAD || echo '?'}
LD_FLAGS+=-X ${MODULE}/version.gitCommit=${GIT_COMMIT}

all: build

build: ${BIN}

bin/ebjs: cmd/server/main.go
	CGO_ENABLED=0 go build -trimpath -ldflags "${LD_FLAGS}" -o $@ $^

bin/ebjclient: cmd/client/main.go
	CGO_ENABLED=0 go build -trimpath -ldflags "${LD_FLAGS}" -o $@ $^

proto:
	@protoc --go_out=. --go-grpc_out . --go_opt=module=${MODULE} --go-grpc_opt=module=${MODULE} `find ${API} -name *.proto`

docker-build: docker-build-ebjs docker-build-epgstation

docker-build-ebjs:
	@${DOCKER} build -t ${REPO}/ebjs .

docker-build-epgstation:
	@${DOCKER} build -t ${REPO}/epgstation -f epgstation/docker/Dockerfile .

clean:
	@rm -rf bin

.PHONY: dummy