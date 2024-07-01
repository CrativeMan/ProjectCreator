let
  pkgs = import <nixpkgs> {};
in 

with pkgs;

stdenv.mkDerivation {
  buildInputs = [ go ];
  pname = "createp";
  version = "0.1.0";

  src = ./.;

  buildPhase = ''
    go build -o main
  '';

  installPhase = ''
    mkdir -p $out/bin/createp
    cp main $out/bin/createp
  '';
}
