# mynav

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

## Features
### Workspace and session management
- Organize workspaces by topic.
- Create, view, update and delete workspaces and topics.
- View information about workspaces, such as git information, activity, descriptions...
- Enter a session for each workspace, allowing to swap between workspaces easilty (uses tmux).
- Create, view, update and delete workspace sessions.

### Tmux session, windows and panes
- View tmux session, windows and panes.
- Create, view. update and delete tmux sessions.
- View a preview of tmux sessions.
- A number of tmux commands as keymaps.

### Simple dev oriented Github client
- Authenticate using device authentication or personal access token authentication.
- View github profile info, repos and PRs.
- Open browser/Copy url of PRs and repos.
- Clone repo directly to a workspace, avoiding the need to use your browser.















