let
  pkgs = import <nixpkgs> { };
in

with pkgs;

stdenv.mkDerivation {
  pname = "createp";
  version = "0.1.0";

  src = ./.;

  buildInputs = [ go ];

  buildPhase  = ''
    export GOPATH=$(mktemp -d)
    export GOCACHE=$GOPATH/cache
    go build -o createp
  '';

  installPhase = ''
    mkdir -p $out/bin
    cp createp $out/bin
  '';
}