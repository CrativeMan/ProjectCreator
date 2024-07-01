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

  outputs = {
    self,
    flake-utils,
    nixpkgs,
    ...
  }:
    flake-utils.lib.eachDefaultSystem (system: let
      pkgs = nixpkgs.legacyPackages.${system};
    in {

      createp = pkgs.buildGoModule rec {
        pname = "createp";
        version = "0.1.0";
        src = ./.;
      };

      devShell = pkgs.mkShell {
        buildInputs = with pkgs; [
          go
        ];
      };
    });
}
