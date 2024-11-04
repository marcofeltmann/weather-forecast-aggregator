let
  # one should sha-fix this tarball as there will be a time with Go 1.23 as default
  nixpkgs = fetchTarball "https://github.com/NixOS/nixpkgs/tarball/nixos-unstable";
  pkgs = import nixpkgs { config = {}; overlays = []; };
in

pkgs.mkShellNoCC {
  packages = with pkgs; [
    git           # for version control, kinda redundant as you might have cloned
    go_1_23       # version requirement for new http route registering
    golangci-lint # run a collection of linters on the code, v1.61 is for Go 1.23
    gnumake       # GNU make build system to automate tool usage
    k6            # End-to-End and performance testing suite from Grafana
  ];
}
