let url = "https://github.com/NixOS/nixpkgs/archive/814f8f3363cecee7cce314286846e4da5e92d689.tar.gz";
in { pkgs ? import (fetchTarball url) { } }:
with pkgs; mkShell {
  nativeBuildInputs = [
    hugo
    html-tidy
  ];
}
