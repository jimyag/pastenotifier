version=$(shell git describe --tags --always)
ldflags="-s -w -X github.com/jimyag/version-go.version=$(version) -X github.com/jimyag/version-go.enableCmd=true"
binary="pastenotifier"
build:
	CGO_ENABLED=0 go build -o ${binary} -v --trimpath -ldflags ${ldflags} ./