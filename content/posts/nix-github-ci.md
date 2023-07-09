---
title: "Nix packages for reproducible environments (lazy-mode)"
description: "How to pin nix packages in shell scripts, ci pipelines and other build processes"
date: "2023-07-09"
tags: ["nix", "make", "github", "ci"]
draft: false
---

Reproducible software environments are very nice to have, they can immensely reduce the amount of migraine induced by software updates.
After some trial-and-error I came up with a few ways to ~~ab~~use [*Nix*](https://nixos.org/) to create reproducible linux environments without writing any [*Nix*](https://nixos.wiki/wiki/Overview_of_the_Nix_Language) code.
Plugging into the Nix ecosystem for this has a few advantages, for one the package collection contains everything I would ever need:

> The Nix Packages collection (Nixpkgs) is a set of over 80 000 packages for the Nix package manager.
> ― [nixos.org](https://nixos.org/)

Also, the [*nix-shell*](https://nixos.org/manual/nix/stable/command-ref/nix-shell.html) application can be used in multiple handy ways to execute commands inside these environments.

To pin packages to a specific version, a commit hash of the [NixOS/nixpkgs](https://github.com/NixOS/nixpkgs) git repository is used.
The commit hash is set either via the [`NIX_PATH`](https://nixos.org/manual/nix/stable/command-ref/env-common.html#env-NIX_PATH) environment variable or the nix-shell [`-I`](https://nixos.org/manual/nix/stable/command-ref/nix-shell.html#opt-I) flag in the form of an [archive URL](https://docs.github.com/en/repositories/working-with-files/using-files/downloading-source-code-archives#source-code-archive-urls):
`https://github.com/NixOS/nixpkgs/archive/3c7487575d9445185249a159046cc02ff364bff8.tar.gz`

With the [Nix package manager](https://nixos.org/manual/nix/stable/installation/multi-user.html) installed, the latest commit hash can be grabbed like this:

```sh
cmd="curl -sSL https://api.github.com/repos/NixOS/nixpkgs/commits/nixos-unstable | jq -r '.sha'"
url="https://github.com/NixOS/nixpkgs/archive/3c7487575d9445185249a159046cc02ff364bff8.tar.gz"
nix-shell --pure -I "nixpkgs=${url}" -p "cacert" "curl" "jq" --run "${cmd}"
```

In this example, the [`-I`](https://nixos.org/manual/nix/stable/command-ref/nix-shell.html#opt-I) flag is used to pin a specific version of the package collection and the [`-p`](https://nixos.org/manual/nix/stable/command-ref/nix-shell.html#opt-p) flag is used to declared packages required in the environment.
Because the [`--pure`](https://nixos.org/manual/nix/stable/command-ref/nix-shell.html#opt--pure) flag is set, the environment will not allow invocations of software, that was not declared for the environment.

### Make

The command can be wrapped in different ways, e.g. called from a [makefile](https://pubs.opengroup.org/onlinepubs/9699919799/utilities/make.html):

```Makefile
NIXPKGS := https://github.com/NixOS/nixpkgs/archive/3c7487575d9445185249a159046cc02ff364bff8.tar.gz

# compact command
.PHONY: build
build:
	nix-shell --pure -I nixpkgs=$(NIXPKGS) -p hugo --run "hugo --verbose"

# long command
define serve
hugo serve                \
    --buildDrafts         \
    --buildFuture         \
    --watch               \
    --bind "127.0.0.1"    \
    --port "8080"
endef
.PHONY: serve
serve:
	nix-shell --pure -I nixpkgs=$(NIXPKGS) -p hugo --run "$(serve)"
```

### Shebang

The *nix-shell* command and the package declarations can also be set directly in a scripts shebang header:

```sh
#!/usr/bin/env nix-shell
#! nix-shell -I nixpkgs=https://github.com/NixOS/nixpkgs/archive/3c7487575d9445185249a159046cc02ff364bff8.tar.gz
#! nix-shell -p cacert hugo
#! nix-shell -i sh --pure
# shellcheck shell=sh
set -e
hugo --verbose
```

The 2 packages [*cacert*](https://search.nixos.org/packages?channel=23.05&show=cacert) and [*hugo*](https://search.nixos.org/packages?channel=23.05&show=hugo) are both available in the whole scripts scope.

With the [shellcheck directive](https://www.shellcheck.net/wiki/Directive) `shell=sh`, we can make sure that shellcheck does not misbehave.

### GitHub Actions

With the following [GitHub Actions](https://docs.github.com/en/actions) steps, it is possible to pin packages and cache them too, to drastically reduce the traffic required for downloads:

```yaml
- name: 'Checkout repo'
  uses: actions/checkout@v3.5.3

  # do pre-build stuff

- name: 'Install Nix'
  uses: cachix/install-nix-action@v22
  with:
    nix_path: "nixpkgs=channel:nixos-23.05"

- name: 'Activate Nix store cache'
  uses: actions/cache@v3.3.1
  id: nix-cache
  with:
    path: "/tmp/nixcache"
    key: "nix-store-${{ hashFiles(format('{0}/build.sh', github.workspace)) }}"

- name: "Import Nix store cache"
  if: "steps.nix-cache.outputs.cache-hit == 'true'"
  run: nix-store --import < /tmp/nixcache

- name: 'Build application'
  run: ./build.sh

- name: "Export Nix store cache"
  if: >
    ( github.ref == 'refs/heads/main' || github.ref == 'refs/heads/master' )
    && steps.nix-cache.outputs.cache-hit != 'true'
  run: nix-store --export $(find /nix/store -maxdepth 1 -name '*-*') > /tmp/nixcache

  # do post-build stuff
```

The workflow makes use of [cachix/install-nix-action](https://github.com/cachix/install-nix-action) to bootstrap the Nix package manager in the CI environment.
Actually, [this websites deployment](https://github.com/nothub/website/blob/7da49d566c756efd9553e86bc10cc70bef076bb9/.github/workflows/ci.yaml) is also done by a GitHub Actions pipeline using Nix.
