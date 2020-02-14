ifndef ${GOROOT}
	export GOROOT=/usr/local/go
endif

PWD=$(shell pwd)
CFLAGS=-I${PWD}/include

export GOOS=darwin
export GOARCH=arm64
export CC=${GOROOT}/misc/ios/clangwrap.sh

all:
	CGO_ENABLED=1 CGO_CFLAGS="-I${PWD}/include" CGO_LDFLAGS="-framework Foundation" go build
	ldid -Sglobal.entitlements ibreaker

clean:
	rm -rf ibreaker

install: all
	ssh ipad "rm -rf /usr/libexec/ibreaker"
	scp ibreaker ipad:/usr/libexec/

