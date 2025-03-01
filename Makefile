# This Makefile enables one-line test, build, and binary signing (Mac) Go
# apps on many platforms.
#
# For signing and notarizing on Mac, the keychain profile name should be
# the aop (module) name. The codesigning identity should already exist.
#
# This is always run as part of CI, so .PHONY doesn't really make sense
# here.

APP=$(shell grep module go.mod | rev | cut -f1 -d'/' | rev)

# Apple doesn't write tools for automateed consumption, so this works but may not be futureproof
CODESIGNINGID=$(shell security find-identity -v -p codesigning | grep -Eo '[0-9A-Z]{40}')

# Some values to make binary verification easier at runtime
VERSION=$(shell git describe --tags --always --abbrev=0 --match='v[0-9]*.[0-9]*.[0-9]*' 2> /dev/null | sed 's/^.//')
COMMIT=$(shell git rev-parse --short HEAD)
DATE=$(shell date '+%Y-%m-%d %H:%M:%S')
	
# Inject these values into the compliled binary
LDFLAGS=-X 'main.Version=$(VERSION)' -X 'main.commit=$(COMMIT)' -X 'main.date=$(DATE)'

test:
	go test -count=1 -tags dev ./...

# shortcut for local builds
build: test
	go build -o build/$(APP) -ldflags="$(LDFLAGS)" *.go

# Make releases in zip format intended for website downloads.
releases: releaselinux releasefreebsd releasemac

# MAC releases
releasemac: releasedarwinarm64 releasedarwinamd64

checksigningid:
	@[ "${CODESIGNINGID}" ] || ( echo ">> CODESIGNINGID is not set"; exit 1 )

releasedarwinarm64: checksigningid test
	mkdir -p build/darwin/arm64/${APP}_${VERSION}
	rm -f build/darwin/arm64/${APP}_${VERSION}/${APP}
	GOOS=darwin GOARCH=arm64 go build -ldflags="${LDFLAGS}" -o build/darwin/arm64/${APP}_${VERSION}/${APP} *.go
	codesign -s ${CODESIGNINGID} -o runtime -v build/darwin/arm64/${APP}_${VERSION}/${APP}
	ditto -c -k --keepParent build/darwin/arm64/${APP}_${VERSION}/${APP} build/releases/${APP}_${VERSION}_darwin_arm64.zip
	xcrun notarytool submit build/releases/${APP}_${VERSION}_darwin_arm64.zip --keychain-profile ${APP}

releasedarwinamd64: checksigningid test 
	mkdir -p build/darwin/amd64
	rm -f build/darwin/amd64/${APP}_${VERSION}/${APP}
	GOOS=darwin GOARCH=amd64 go build -o build/darwin/amd64/${APP}_${VERSION}/${APP} *.go
	codesign -s ${CODESIGNINGID} -o runtime -v build/darwin/amd64/${APP}_${VERSION}/${APP}
	ditto -c -k --keepParent build/darwin/amd64/${APP}_${VERSION}/${APP} build/releases/${APP}_${VERSION}_darwin_amd64.zip
	xcrun notarytool submit build/releases/${APP}_${VERSION}_darwin_amd64.zip --keychain-profile ${APP} 

# LINUX releases
releaselinux: test
	for arch in 386 arm arm64 amd64; do \
		mkdir -p build/linux/$$arch; \
		rm -f build/linux/$$arch/${APP}_${VERSION}/${APP}; \
		GOOS=freebsd GOARCH=$$arch go build -ldflags="${LDFLAGS}" -o build/linux/$$arch/${APP}_${VERSION}/${APP} *.go; \
		ditto -c -k --keepParent build/linux/$$arch/${APP}_${VERSION}/${APP} build/releases/${APP}_${VERSION}_linux_$$arch.zip; \
	done

# FREEBSD releases
releasefreebsd: test
	for arch in 386 arm amd64; do \
		mkdir -p build/freebsd/$$arch; \
		rm -f build/freebsd/$$arch/${APP}_${VERSION}/${APP}; \
		GOOS=freebsd GOARCH=$$arch go build -ldflags="${LDFLAGS}" -o build/freebsd/$$arch/${APP}_${VERSION}/${APP} *.go; \
		ditto -c -k --keepParent build/freebsd/$$arch/${APP}_${VERSION}/${APP} build/releases/${APP}_${VERSION}_freebsd_$$arch.zip; \
	done

