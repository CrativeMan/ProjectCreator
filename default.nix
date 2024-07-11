let
  pkgs = import <nixpkgs> { };
in

with pkgs;

buildGoModule rec {
  pname = "createp";
  version = "0.1.10";

  src = ./src/.;

  vendorSha256 = "sha256-i2FG/Dlw0r5PVHak+37VBeRwG7Vf7qWNlYzNyJUIURg=";

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
}
