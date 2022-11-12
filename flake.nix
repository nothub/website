{
  inputs = {
    nixpkgs.url = "nixpkgs/nixos-22.05";
    flake-compat = {
      url = "github:edolstra/flake-compat";
      flake = false;
    };
    nix-filter.url = "github:numtide/nix-filter";
  };
  outputs = { self, nixpkgs, flake-compat, nix-filter }: {
    packages.x86_64-linux.default =
      let
        pkgs = import nixpkgs {
          system = "x86_64-linux";
        };
      in
        pkgs.stdenv.mkDerivation {
          name = "website";
          meta = with pkgs.stdenv.lib; {
            homepage = "https://hub.lol";
            licenses = licenses.cc-by-40;
          };
          src = nix-filter.lib {
            root = ./.;
            include = [
              "archetypes"
              "content"
              "layouts"
              "static"
              "LICENSE"
              "Makefile"
              "config.yaml"
              "tidy.conf"
            ];
          };
          nativeBuildInputs = with pkgs; [ hugo html-tidy ];
          shellHook = "source <(hugo completion bash)";
        };
  };
}
