package main

const (
	CFLAKECONTENT string = `{
description = "A very basic flake";

inputs = {
  nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
  flake-utils.url = "github:numtide/flake-utils";
};

outputs = {
  self,
  nixpkgs,
  flake-utils,
}:
  flake-utils.lib.eachDefaultSystem (
    system: let
      pkgs = import nixpkgs {
        inherit system;
      };
    in
      with pkgs; {
        formatter = pkgs.alejandra;
        devShell = mkShell.override {stdenv = clangStdenv;} {
          packages = [ `

	CMAINCONTENTS string = `#include <stdio.h>
	
int main() {
	printf("Hello, World!\n");
	return 0;
}`

	GOFLAKECONTENT string = `{
description = "Basic Go Flake";

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
    devShell = pkgs.mkShell {
      buildInputs = with pkgs; [`

	GOMAINCONTENTS string = `package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello, World!")
}`

	ENVRCCONTENT string = "use flake"

	FLAKECONTENT_END string = `          ];
        };
      }
  );
}`
)
