{
  description = "LazyNet - Terminal network monitor";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
        
        lazynet = pkgs.buildGoModule  {
          pname = "lazynet";
          version = "0.0.1";
          
          src = ./.;
          
          vendorHash = "sha256-Z18bTupa1NTOFANEUSOykNJ6jxxneUOLHB3K54yiYs4=";
          
          subPackages = [ "cmd/lazynet" ];
          
          meta = with pkgs.lib; {
            description = "Terminal network monitor";
            homepage = "https://github.com/sureshkrishnan-v/lazynet";
            license = licenses.mit;
            maintainers = [ ];
            platforms = platforms.unix;
          };
        };
      in
      {
        packages = {
          default = lazynet;
          lazynet = lazynet;
        };
        
        apps = {
          default = flake-utils.lib.mkApp {
            drv = lazynet;
            name = "lazynet";
          };
        };
        
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            gopls
            gotools
            go-tools
          ];
        };
      });
}
