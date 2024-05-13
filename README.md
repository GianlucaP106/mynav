# MyNav

A user friendly TUI workspace manager

https://github.com/GianlucaP106/mynav/assets/93693693/e14917f3-6b0d-4b0e-a23d-52b610e3fcf0

## Description
Mynav is a TUI written in go. It aims to allow for an easy view of all your workspaces, notes or programming projects.

## Installation

### Try with docker first

```bash
docker run -it --name mynav --rm ubuntu bash -c '
        apt update &&
        apt install -y git golang-go neovim tmux curl unzip &&
        cd &&
        (curl -fsSL https://raw.githubusercontent.com/GianlucaP106/mynav/main/install.sh | bash) &&
        export PATH="$PATH:$HOME/.mynav" &&
        mkdir nav && cd nav &&
        mynav
    '
```

> Note: The installation uses go and git, and the application uses git, nvim, and tmux.

### Build from source

```bash
curl -fsSL https://raw.githubusercontent.com/GianlucaP106/mynav/main/install.sh | bash
```

### Add to PATH
```bash
export PATH="$PATH:$HOME/.mynav"
```

## Usage
```bash
# The first time this is ran, it will initialize the directory
mynav
```

> Use '?' in the UI to see the key maps
