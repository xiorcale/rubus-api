with import <nixpkgs> {};

pkgs.mkShell rec {

  name = "rubus-shell";

  buildInputs = with pkgs; [
    go
  ];

  shellHook = ''
    go get golang.org/x/tools/...
    export GOCACHE=$TMPDIR/go-cache
    export GOPATH=$HOME/go
    export PATH=$PATH:$HOME/go/bin
  '';

}
