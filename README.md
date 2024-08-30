# mynav

A user friendly TUI workspace manager, aiming to allow for an easy view of all your workspaces in a terminal environment.

### Main tab

![main-tab](https://github.com/user-attachments/assets/3e340077-1cd5-41e3-a5c0-4ee4bff6cf4a)

### Tmux tab

![tmux-view2](https://github.com/user-attachments/assets/f139408a-8855-40fb-8411-8e9de8bdd947)

## Elevator pitch

Before mynav, I would often get annoyed when working on multiple projects using tmux, I would manually `cd` from project directory to project directory. Of course, if you have tmux sessions active, you can use choose tree to bounce from session to session, but this persists only if the tmux server is alive. With mynav, I can combine the features of tmux, with a workspace management system, allowing for a effecient development workflow in a terminal environment.

## Installation

### Try with docker first

```bash
docker run -it --name mynav --rm ubuntu bash -c '
        apt update &&
        apt install -y git golang-go neovim tmux curl unzip libx11-dev &&
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
- View information about workspaces, such as git information, activity, descriptions...
- Enter a session for each workspace, allowing to swap between workspaces easilty (uses tmux).
- Create, view, update and delete workspace sessions.

#### Tmux session, windows and panes

- View tmux session, windows and panes.
- Create, view. update and delete tmux sessions.
- View a preview of tmux sessions.
- A number of tmux commands as keymaps.

#### Simple Github client

- Authenticate using device authentication or personal access token authentication.
- View github profile info, repos and PRs.
- Open browser/Copy url of PRs and repos.
- Clone repo directly to a workspace, avoiding the need to use your browser.

## Keymaps

#### Use '?' in the TUI to see all the key maps

#### Global

| Key       | Action          |
| ---       | --------------- |
| q         | quit            |
| S         | Open settings   |
| [         | Cycle left tab  |
| ]         | Cycle right tab |
| j        | Move down in lists |
| k        | Move up in lists   |
| arrows   | Move arround panes |
| ctrl-vim | Move around panes  |
| ?        | Open cheatsheet  |

#### Topics

| Key   | Action                        |
| ----- | ----------------------------- |
| enter | Open topic                    |
| /     | Search by name                |
| a     | Create topic                  |
| r     | Rename topic                  |
| s     | Search for workspace globally |
| D     | Delete topic                  |

#### Workspaces

| Key   | Action                                     |
| ----- | ------------------------------------------ |
| esc   | Go back                                    |
| s     | See workspace information                  |
| L     | Open lazygit at workspace                  |
| g     | git clone                                  |
| G     | Open browser to git repo                   |
| u     | Copy git remote url to clipboard           |
| /     | Search by name                             |
| enter | Open workspace                             |
| v     | Open workspace using neovim                |
| t     | Open workspace using native terminal       |
| m     | Move workspace to a different topic        |
| D     | Delete workspace                           |
| r     | Rename workspace                           |
| e     | Add/change description                     |
| a     | Create workspace                           |
| X     | Kill tmux session (if any)                 |

#### Tmux sessions

| Key   | Action                                 |
| ----- | -------------------------------------- |
| o     | Attach to session                      |
| enter | Focus windows view                     |
| D     | Delete session                         |
| X     | Kill the tmux server                   |
| W     | Kill all workspace-associated sessions |
| c     | Open choose tree on session            |
| a     | New session                            |

#### Tmux windows

| Key   | Action              |
| ----- | ------------------- |
| o     | Attach to session   |
| X     | Kill this window    |
| esc   | Focus sessions view |
| enter | Focus panes view    |

#### Github profile

| Key | Action                             |
| --- | ---------------------------------- |
| L   | Login with device code and browser |
| P   | Login with personal access token   |
| o   | Open profile in browser            |
| u   | Copy profile url to clipboard      |
| O   | Logout                             |
| R   | Refetch all Github data            |

#### Github repos

| Key | Action                   |
| --- | ------------------------ |
| c   | Clone repo to workspace  |
| o   | Open repo in browser     |
| u   | Copy repo url to browser |
| R   | Refetch all Github data |

#### Github pull requests

| Key | Action                  |
| --- | ----------------------- |
| o   | Open PR in browser      |
| u   | Copy PR url to browser  |
| R   | Refetch all Github data |

#### Settings dialog

| Key   | Action                   |
| ----- | ------------------------ |
| enter | Change setting           |
| D     | Reset setting to default |
| esc   | Close settings dialog    |
