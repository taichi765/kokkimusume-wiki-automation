{
  description = "A very basic flake";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixpkgs-unstable";
  };

  outputs = {self,nixpkgs}: 
    let
      system = "x86_64-linux";
      pkgs = import nixpkgs { inherit system;};
    in {
      devShells.${system}.default =  pkgs.mkShell {
        buildInputs = [
          pkgs.go
          pkgs.azure-cli
          pkgs.terraform
        ];
      };
  };
}
