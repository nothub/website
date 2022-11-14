.PHONY: default
default: clean build tidy

.PHONY: clean
clean:
	-rm -rf \
		"work" \
		"result" \
		".hugo_build.lock"

.PHONY: build
build:
	hugo \
		--cleanDestinationDir \
		--ignoreCache \
		--verbose \
		--destination "work"

.PHONY: tidy
tidy:
	find "work" \
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
