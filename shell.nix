#!/usr/bin/env -S nix-shell --pure
let url = "https://github.com/NixOS/nixpkgs/archive/b69883faca9542d135fa6bab7928ff1b233c167f.tar.gz";
in { pkgs ? import (fetchTarball url) { } }:
with pkgs; mkShell {
  nativeBuildInputs = [
    hugo
    html-tidy
  ];
}
