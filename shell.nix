#!/usr/bin/env -S nix-shell --pure
let url = "https://github.com/NixOS/nixpkgs/archive/0c4800d579af4ed98ecc47d464a5e7b0870c4b1f.tar.gz";
in { pkgs ? import (fetchTarball url) { } }:
with pkgs; mkShell {
  nativeBuildInputs = [
    hugo
    html-tidy
  ];
}
