with import <nixpkgs> {};

pkgs.mkShell rec {

  name = "rubus-shell";

  buildInputs = with pkgs; [
    go
  ];

  shellHook = ''
    export GOCACHE=$TMPDIR/go-cache
    export GOPATH=$HOME/go
    export PATH=$PATH:$HOME/go/bin
    go get golang.org/x/tools/...
  '';

}
