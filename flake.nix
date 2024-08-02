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
        version = "1.1.0";

        src = ./src;

        vendorHash = "sha256-UnOebGoWD33c+1dNiXpbS0kSZWRCAsZcQnXFTXIZFpk=";
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
