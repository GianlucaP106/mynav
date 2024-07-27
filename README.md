# MyNav

A user friendly TUI workspace manager

## Description
Mynav is a TUI workspace manager. It aims to allow for an easy view of all your workspaces, notes or programming projects. It integrates with tmux and neovim for a great workspace management experience.

<img width="1937" alt="Screenshot 2024-07-26 at 9 13 48 PM" src="https://github.com/user-attachments/assets/9fee6ed3-eddb-4e30-bad6-c4a7a06207bf">
<img width="1937" alt="Screenshot 2024-07-26 at 9 14 01 PM" src="https://github.com/user-attachments/assets/884fd525-66db-4d13-9e2a-1b94daad1cea">


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
