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
      packages.default = pkgs.buildGoModule rec {
        pname = "createp";
        version = "0.1.10";

        src = ./src/.;

        vendorHash = "sha256-i2FG/Dlw0r5PVHak+37VBeRwG7Vf7qWNlYzNyJUIURg=";

        subPackages = [ "." ];

        buildPhase  = ''
          runHook preBuild  
          go build -o createp
          runHook postBuild
        '';

        installPhase = ''
          runHook preInstall
          mkdir -p $out/bin
          cp createp $out/bin
          runHook postInstall
          '';
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