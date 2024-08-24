# MyNav

A user friendly TUI workspace manager

## Description
Mynav is a TUI workspace and session manager. It aims to allow for an easy view of all your workspaces, notes or programming projects.

### Main tab
![basic-workspace](https://github.com/user-attachments/assets/50d52d3e-ac73-43cc-8ad8-d6d8d22daf82)

### Tmux tab
![tmux-view](https://github.com/user-attachments/assets/8e9ae3cd-0338-4c1f-ba48-db674124c1b5)


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

> ### Use '?' in the TUI to see all the key maps!
