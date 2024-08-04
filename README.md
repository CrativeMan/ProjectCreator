# Createp

A project creator written in the [golang](https://go.dev/) using the [lipgloss](https://github.com/charmbracelet/lipgloss)
cli framework. It is optimized for the nix/nixos flake workflow.

# Features

- Create C,Go,C++ programs
- Choose between different project types
- Simple fast and efficient

## Command line flags

- `-h` for showing a help
- `-v` for showing the version
- `-nf` for not generating a flake.nix file and using nix/flake dev envs.

# How to install

## Manual install

```
git clone https://github.com/CrativeMan/ProjectCreator createp
cd createp
cd src
go build
./createp
```

## Nix/Flake install

```
git clone https://github.com/CrativeMan/ProjectCreator createp
cd createp
nix build
```

Then either use the `nix-env -i ./result`
or add this to your `flake.nix` file

```
createp = {
    url = "github:CrativeMan/ProjectCreator";
};
```

and then install it like your other programs/inputs in your `configuration.nix` like this

```
environment.systemPackages = [
    pkgs.inputs.createp.packages.x86_64-linux.default
];
```

# Todo

- [ ] Make the user change direcotrys after project creation
- [ ] Java implementation
