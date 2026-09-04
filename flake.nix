{
  description = "A very basic flake";

  inputs = {
    nixpkgs = {
      url = "github:nixos/nixpkgs/nixpkgs-unstable";
    };
  };

  outputs = {self,nixpkgs}: 
    let
      system = "x86_64-linux";
      pkgs = import nixpkgs {
        inherit system;
        config.allowUnfreePredicate = pkg: builtins.elem (pkgs.lib.getName pkg) [
          "terraform"
       ];
      };
    in {
      devShells.${system}.default =  pkgs.mkShell {
        buildInputs = [
          pkgs.go
          # Use system-installed Azure CLI because nixpkgs' one is not latest and extension can't be installed successfully.
          #(with pkgs.azure-cli; withExtensions [ extensions.log-analytics ])
          pkgs.terraform
          pkgs.just
        ];
      };
  };
}
