NIXPKGS := https://github.com/NixOS/nixpkgs/archive/f294325aed382b66c7a188482101b0f336d1d7db.tar.gz

.PHONY: default
default: clean build tidy

.PHONY: clean
clean:
	-rm -rf \
		"public" \
		"result" \
		".hugo_build.lock"

define build
hugo                      \
    --cleanDestinationDir \
    --ignoreCache         \
    --verbose
endef
.PHONY: build
build:
	nix-shell -I $(NIXPKGS) -p hugo --run "$(build)"

define tidy
find "public"       \
    -type "f"       \
    -iname "*.html" \
    -exec tidy -config "tidy.conf" {} \;
endef
.PHONY: tidy
tidy:
	nix-shell -I $(NIXPKGS) -p html-tidy --run "$(tidy)"

define serve
hugo serve                \
    --cleanDestinationDir \
    --ignoreCache         \
    --verbose             \
    --buildDrafts         \
    --buildFuture         \
    --watch               \
    --bind "127.0.0.1"    \
    --port "8080"
endef
.PHONY: serve
serve:
	nix-shell -I $(NIXPKGS) -p hugo --run "$(serve)"
