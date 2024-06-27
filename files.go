package main

const (
	CFLAKECONTENT string = `{
  description = "Basic C development environment";

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
        buildInputs = with pkgs; [
					clang
					gcc
					llvmPackages.clangUseLLVM
        ];
      };
    });
}` 

	CMAINCONTENTS string = `#include <stdio.h>
	
int main() {
	printf("Hello, World!\n");
	return 0;
}`

	GOFLAKECONTENT string = `{
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
      devShell = pkgs.mkShell {
        buildInputs = with pkgs; [
          go
        ];
      };
    });
}`
	
	GOMAINCONTENTS string = `package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello, World!")
}`

	ENVRCCONTENT string = "use flake"
)