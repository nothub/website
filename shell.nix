{ pkgs ? import (fetchTarball "https://nixos.org/channels/nixos-22.05/nixexprs.tar.xz") { } }:
with pkgs; mkShell { nativeBuildInputs = [ hugo html-tidy ]; }
