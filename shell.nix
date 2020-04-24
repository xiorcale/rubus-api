with import <nixpkgs> {};

pkgs.mkShell rec {

  name = "rubus-shell";

  buildInputs = with pkgs; [
    # Editor dependencies (auto-completion, syntax checking, ...)
    go

    # Project dependencies
  ];

  shellHook = ''
    # gopls, ...
    go get golang.org/x/tools/...
    export GOCACHE=$TMPDIR/go-cache
    export GOPATH=$HOME/go
    export PATH=$PATH:$HOME/go/bin
  '';

}
