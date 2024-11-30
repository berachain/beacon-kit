{ self, ... }: {
  perSystem = { pkgs, self', system, ... }:
    let
      NAME = "beacond";
      APP_NAME = "beacond";
      DB_BACKEND = "pebbledb";

      gitRev =
        if (builtins.hasAttr "rev" self) then self.rev else "dirty";

      buildTags = [ "netgo" "muslc" "blst" "bls12381" "pebbledb" ];

      fs = pkgs.lib.fileset;
    in
    {
      packages = {
        beacond =
          pkgs.pkgsStatic.buildGo122Module {
            name = NAME;
            src = fs.toSource {
              root = ../.;
              fileset = fs.union
                # unconditionally include...
                (fs.unions (pkgs.lib.flatten [ ../go.work.sum ../go.work ]))
                # ...and include the go sources of all go files in the repo
                (fs.fileFilter
                  (file: (file.name == "go.mod") || (builtins.any file.hasExt [ "go" ]))
                  ../.
                );
            };

            # if any dependencies change, set this to an empty string and rebuild; nix will complain about a hash mismatch (it uses the empty hash if an empty string is provided) and will print the expected hash. copy paste that in here and you're good to go!
            vendorHash = "sha256-fWH0rKwpoqeVNMVRUuCmNJZDx2nqCEF1WVP8qgWGSy8=";
            CGO_ENABLED = 1;
            subPackages = [ "./beacond/cmd" ];
            inherit buildTags;

            # https://github.com/NixOS/nixpkgs/issues/299096
            overrideModAttrs = (_: {
              buildPhase = ''
                go work vendor -e
              '';
            });

            # go bearishness with vendoring means it doesn't vendor files it doesn't think are required for build
            # https://github.com/golang/go/issues/26366
            proxyVendor = true;

            ldflags = [
              "-X github.com/cosmos/cosmos-sdk/version.Name=${NAME}"
              "-X github.com/cosmos/cosmos-sdk/version.AppName=${APP_NAME}"
              # this one is a bit trickier, `git describe --tags --always --dirty` won't work due to ifd; i would recommend tracking this value internally
              # "-X github.com/cosmos/cosmos-sdk/version.Version=${}"
              "-X github.com/cosmos/cosmos-sdk/version.Commit=${gitRev}"
              "-X github.com/cosmos/cosmos-sdk/version.BuildTags=${pkgs.lib.concatStringsSep "," buildTags}"
              "-X github.com/cosmos/cosmos-sdk/types.DBBackend=${DB_BACKEND}"
              "-w"
              "-s"
              "-linkmode=external"
              "-extldflags"
              "'-Wl,-z,muldefs -static'"
            ];

            postInstall = ''
              mv $out/bin/cmd $out/bin/beacond
            '';

            meta.mainProgram = "beacond";
          };
      };

      # build.mk is a bit cursed right now (impure bc it uses git, and for some reason tput?)
      # checks = {
      #   go-test = pkgs.go.stdenv.mkDerivation {
      #     name = "go-test";
      #     buildInputs = [ pkgs.go ];
      #     src = ../.;
      #     doCheck = true;
      #     checkPhase = ''
      #       # Go will try to create a .cache/ dir in $HOME.
      #       # We avoid this by setting $HOME to the builder directory
      #       export HOME=$(pwd)

      #       go version
      #       go test ./...
      #       touch $out
      #     '';
      #   };

      #   go-vet = pkgs.go.stdenv.mkDerivation {
      #     name = "go-vet";
      #     buildInputs = [ pkgs.go ];
      #     src = ../.;
      #     doCheck = true;
      #     checkPhase = ''
      #       # Go will try to create a .cache/ dir in $HOME.
      #       # We avoid this by setting $HOME to the builder directory
      #       export HOME=$(pwd)

      #       go version
      #       go vet ./...
      #       touch $out
      #     '';
      #   };

      #   go-staticcheck = pkgs.go.stdenv.mkDerivation {
      #     name = "go-staticcheck";
      #     buildInputs = [ pkgs.go pkgs.go-tools ];
      #     src = ../.;
      #     doCheck = true;
      #     checkPhase = ''
      #       # Go will try to create a .cache/ dir in $HOME.
      #       # We avoid this by setting $HOME to the builder directory
      #       export HOME=$(pwd)

      #       staticcheck ./...
      #       touch $out
      #     '';
      #   };
      # };
    };
}
