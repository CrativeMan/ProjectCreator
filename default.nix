let
  pkgs = import <nixpkgs> { };
  go = pkgs.go_1_21;
in

with pkgs;

buildGoModule rec {
  pname = "createp";
  version = "0.1.1";

  src = ./.;

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