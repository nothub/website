# invoked by nix stdenv derivation default buildPhase
.PHONY: default
default: clean build tidy

# invoked by nix stdenv derivation default installPhase
.PHONY: install
install:
	mkdir -p "${out}/srv"
	cp -vR "build/"* "${out}/srv"

.PHONY: clean
clean:
	-rm -rf \
		"build" \
		"result" \
		".hugo_build.lock"

.PHONY: build
build:
	hugo \
		--cleanDestinationDir \
		--ignoreCache \
		--verbose \
		--destination "build"

.PHONY: tidy
tidy:
	find "build" \
	    -type "f" \
	    -iname "*.html" \
	    -exec tidy -config "tidy.conf" {} \;

.PHONY: serve
serve:
	hugo serve \
		--cleanDestinationDir \
		--ignoreCache \
		--verbose \
		--buildDrafts \
		--buildFuture \
		--watch \
		--bind "127.0.0.1" \
		--port "8080"
