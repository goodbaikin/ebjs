{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-23.11";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem
      (system:
        let
          pkgs = nixpkgs.legacyPackages.${system};

        in
        {
          devShells.default = pkgs.stdenv.mkDerivation {
            name = "ebjs";
            src = self;
            buildInputs = with pkgs; [
              go
              gopls
              gotests
              gomodifytags
              impl
              delve
              go-tools
              gnumake
              protobuf
              docker
            ];
          };
        }
      );
}