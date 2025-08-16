# MyNav 🧭

A powerful terminal-based workspace navigator and session manager built in Go. MyNav helps developers organize and manage multiple projects through an intuitive interface, seamlessly integrating with tmux sessions.

[![Version](https://img.shields.io/badge/version-v2.1.1-blue)](https://github.com/GianlucaP106/mynav/releases)
[![Go Version](https://img.shields.io/badge/go-1.22.3+-00ADD8?logo=go)](https://golang.org/)
[![Platform](https://img.shields.io/badge/platform-macOS%20%7C%20Linux-lightgrey)](https://github.com/GianlucaP106/mynav#prerequisites)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](CONTRIBUTING.md)

![MyNav Demo](https://github.com/user-attachments/assets/c2482080-6c1d-4fda-a3d5-e0ae6d8a916b)

## Overview

MyNav addresses the common challenge of managing multiple development projects in a terminal environment. While tmux provides excellent session management, it lacks robust workspace organization capabilities. MyNav bridges this gap by combining tmux's powerful features with an intuitive workspace management system, enabling developers to efficiently navigate between projects and maintain organized development workflows.

## Features

### 🏢 Workspace Management
- **Topic-based organization**: Group related workspaces into logical topics
- **Rapid workspace creation**: Quick setup and navigation between projects
- **Filesystem-based storage**: Direct integration with your existing directory structure

### 💻 Session Management
- **Comprehensive session control**: Create, modify, delete, and enter sessions seamlessly
- **Live session preview**: Real-time display of window and pane information
- **Instant session switching**: Fast navigation between active development sessions

### 🛠️ Developer Experience
- **Fuzzy search**: Intelligent search across workspaces and sessions
- **tmux integration**: Built on top of tmux for maximum compatibility
- **Extensive shortcuts**: Comprehensive keyboard navigation and shortcuts
- **Git awareness**: Integration with Git repositories and status
- **Modern UI**: Clean, intuitive terminal interface inspired by Lazygit
- **Vim-style navigation**: Familiar navigation patterns for power users

## Installation

### Quick Start with Docker

Experience MyNav immediately with our Docker setup:

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

### Automated Installation

Install MyNav with a single command:

```bash
curl -fsSL https://raw.githubusercontent.com/GianlucaP106/mynav/main/install.bash | bash
```

### Manual Installation

For users who prefer manual installation:

```bash
# Clone the repository
git clone https://github.com/GianlucaP106/mynav.git

# Navigate to project directory
cd mynav

# Build the project
go build
```

### Prerequisites

- **tmux 3.0+**: Required for session management
- **Git**: Optional, enables repository-specific features
- **Terminal**: UTF-8 support required for proper display

## Usage

### Getting Started

MyNav operates within a designated root directory. You can initialize multiple independent directories, but nested configurations are not supported.

Launch MyNav from any location:

```bash
mynav
```

> **Note**: MyNav automatically detects existing configurations in the current directory or any parent directory.

Specify a custom root directory:

```bash
mynav -path /your/root/path
```

### Navigation

Press `?` within the interface to view all available keyboard shortcuts for your current context.

## tmux Integration

MyNav seamlessly integrates with **tmux** to provide robust session management. When creating a session from a workspace, MyNav uses the workspace's directory path as the tmux session name, maintaining transparency and familiarity.

### Key Integration Features

- **Full tmux compatibility**: All standard tmux features remain available
- **Session detachment**: Press `Leader + D` to detach and return to MyNav
- **State synchronization**: MyNav stays in sync with your development workflow

## Keyboard Shortcuts

### Navigation Controls

| Key | Action | Context |
|-----|--------|---------|
| `h` / `←` | Focus left panel | Global |
| `l` / `→` | Focus right panel | Global |
| `j` / `↓` | Move down | List views |
| `k` / `↑` | Move up | List views |
| `Tab` | Toggle focus | Search dialog |
| `Esc` | Close/cancel | Dialogs |

### Action Commands

| Key | Action | Context |
|-----|--------|---------|
| `Enter` | Open/select item | Global |
| `a` | Create new topic/workspace | Topics/Workspaces view |
| `D` | Delete item | Topics/Workspaces/Sessions view |
| `r` | Rename item | Topics/Workspaces view |
| `X` | Kill session | Workspaces/Sessions view |
| `s` | Search workspaces | Global |
| `?` | Toggle help menu | Global |
| `q` | Quit application | Global |
| `<` | Cycle preview left | Global |
| `>` | Cycle preview right | Global |
| `Ctrl+C` | Quit application | Global |

## Configuration

MyNav employs a flexible configuration system designed for multi-project development:

- **Multiple workspace support**: Initialize independent workspaces in different directories
- **Automatic configuration detection**: Searches current and parent directories for existing configurations
- **Non-nested architecture**: Nested configurations are not allowed; MyNav will use the parent configuration
- **Home directory protection**: The home directory cannot be initialized as a MyNav workspace

## Development

### Environment Setup

MyNav is designed for minimal configuration requirements. To begin development:

1. Ensure Go runtime is installed
2. Clone the repository
3. Run `go build` to compile

The project structure is straightforward and requires no additional dependencies beyond the Go standard library.

## Contributing

We welcome contributions from the community! Please ensure your commits follow the [Conventional Commits](https://www.conventionalcommits.org/) specification.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for complete details.

---

<div align="center">

[⭐ Star on GitHub](https://github.com/GianlucaP106/mynav/stargazers) • [📫 Report Bug](https://github.com/GianlucaP106/mynav/issues) • [💬 Discussions](https://github.com/GianlucaP106/mynav/discussions)

</div>
