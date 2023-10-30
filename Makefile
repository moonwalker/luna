VERSION=v0.3.2-beta

release:
	@git tag -a ${VERSION} -m "Release ${VERSION}" && git push origin ${VERSION}
	@goreleaser release --clean

build:
	@goreleaser build --skip=validate --clean
