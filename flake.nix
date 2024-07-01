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
  let
    pkgs = import nixpkgs{system = "x86_64-linux";};
  in {
    
    devShell = pkgs.mkShell {
      buildInputs = with pkgs; [
        go
      ];
    };
  };
}
