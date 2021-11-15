VERSION=v0.2.2

release:
	@git tag -a ${VERSION} -m "Release ${VERSION}" && git push origin ${VERSION}
	@goreleaser release --rm-dist

build:
	@goreleaser build --rm-dist --skip-validate
