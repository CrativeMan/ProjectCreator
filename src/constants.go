package main

import lip "github.com/charmbracelet/lipgloss"

var version = "createp: 0.3.2"

const (
	C     = 0
	CPP   = 1
	GO    = 2
	JAVA  = 3
	ENVRC = 4
	CLOSE = 9
)

var (
	language     int
	path         string
	sty          styles
	Hostname     string
	GoModuleName string
)

type Args struct {
	Help    bool
	Version bool
	Flake   bool
}

type styles struct {
	success lip.Style
	fail    lip.Style
	warning lip.Style
}

// =============================
// FILES
// =============================
const (
	FLAKECONTENT string = `{
  description = "A very basic flake";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    { self
    , nixpkgs
    , flake-utils
    ,
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs {
          inherit system;
        };
      in
      with pkgs; {
        formatter = pkgs.alejandra;
        devShell = mkShell.override { stdenv = clangStdenv; } {
          packages = [
            %s
          ];
        };
      }
    );
}
`

	CMAINCONTENTS string = `#include <stdio.h>
	
int main() {
	printf("Hello, World!\n");
	return 0;
}`

	CMAKECONTENTS string = `CC = gcc
CFLAGS = -Wall -Wextra -g
SRCS = main.c
OBJS = $(SRCS:.c=.o)
TARGET = main

all: $(TARGET)

%.o: %.c
	$(CC) $(CFLAGS) -c $< -o $@

$(TARGET): $(OBJS)
	$(CC) $(CFLAGS) $^ -o $@

clean:
	rm -f $(OBJS) $(TARGET)
`

	CPPMAINCONTENTS string = `#include <iostream>

int main() {
    std::cout << "Hello, World!\n";
    return 0;
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
