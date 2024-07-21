{
  description = "Basic Go development environment";

  inputs = {
    flake-compat = {
      url = "github:edolstra/flake-compat";
      flake = false;
    };
    flake-utils.url = "github:numtide/flake-utils";
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
  };

  outputs =
    { self
    , flake-utils
    , nixpkgs
    , ...
    }:
    flake-utils.lib.eachDefaultSystem (system:
    let
      pkgs = nixpkgs.legacyPackages.${system};
    in
    {
      packages.default = pkgs.buildGoModule rec {
        pname = "createp";
        version = "0.1.12";

        src = ./src;

        vendorHash = "sha256-FcztgJAxge1uQLUZwOHOW+vPsfI67wtt9lDFwZaP4wU=";
      };

      devShell = pkgs.mkShell {
        buildInputs = with pkgs; [
          go
          gofumpt
          goimports-reviser
          golines
          delve
        ];
      };
    });
}
