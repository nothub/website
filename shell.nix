let url = "https://github.com/NixOS/nixpkgs/archive/596a8e828c5dfa504f91918d0fa4152db3ab5502.tar.gz";
in { pkgs ? import (fetchTarball url) { } }:
with pkgs; mkShell {
  nativeBuildInputs = [
    hugo
    html-tidy
  ];
}
