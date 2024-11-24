# mynav

A TUI workspace and session manager, aiming to allow for an easy view of all your workspaces and sessions in a terminal environment.

![demo2](https://github.com/user-attachments/assets/4fa82356-26a8-4260-9b3f-479b0bc18914)

## Elevator pitch

Before mynav, I would often get annoyed when working on multiple projects using tmux, I would manually navigate from project directory to project directory. Of course, if you have tmux sessions active, you can use choose tree to bounce from session to session, but this persists only if the tmux server is alive and is not what I had hoped for as a workspace manager. With mynav, I can combine the features of tmux, with a workspace management system, allowing for a effecient development workflow in a terminal environment.

## Installation

### Try with docker first

```bash
docker run -it --name mynav --rm ubuntu bash -c '
        apt update &&
        apt install -y git golang-go neovim tmux curl unzip &&
        cd &&
        (curl -fsSL https://raw.githubusercontent.com/GianlucaP106/mynav/main/install.bash | bash) &&
        export PATH="$PATH:$HOME/.mynav" &&
        mkdir nav && cd nav &&
        mynav
    '
```

### Binary installation

```bash
curl -fsSL https://raw.githubusercontent.com/GianlucaP106/mynav/main/install.bash | bash
```

> **See all binaries on the [latest release](https://github.com/GianlucaP106/mynav/releases/latest) page.**

### Build from source

```bash
# go to mynav's configuration directory
mkdir ~/.mynav/ 2>/dev/null | true
cd ~/.mynav

# clone the repo
git clone https://github.com/GianlucaP106/mynav src
cd src

# build the project
go build -o ~/.mynav/mynav

# add to path
export PATH="$PATH:$HOME/.mynav"

# optionally delete the source code
cd ~/.mynav
rm -rf src
```

#### Supported platforms

mynav is only supported for Linux and MacOS.

## Usage

```bash
mynav
```

#### Configuration

- mynav will look for a configuration in the current or any parent directory, otherwise will ask to initialize the current directory.
- mynav can be initialized in multiple independant directories, but not nested.
- mynav cannot be initialized in the user home directory.

## Features

#### Workspace and session management

- Organize workspaces by topic.
- Create, view, update and delete workspaces and topics.
- View information about workspaces, information about its session, preview...
- Enter a session for each workspace, allowing to swap between workspaces easilty (uses tmux).
- Create, view, update and delete workspace sessions.


## Keymaps

#### Use '?' in the TUI to see all the key maps
