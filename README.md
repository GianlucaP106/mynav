# MyNav 🧭

A powerful terminal-based workspace navigator and session manager built in Go. MyNav helps developers organize and manage multiple projects through an intuitive interface, seamlessly integrating with tmux sessions.

![Version](https://img.shields.io/badge/version-v2.1.1-blue)
![Go Version](https://img.shields.io/badge/go-1.22.3+-00ADD8?logo=go)
![Platform](https://img.shields.io/badge/platform-macOS%20%7C%20Linux-lightgrey)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](CONTRIBUTING.md)

![demo](https://github.com/user-attachments/assets/c2482080-6c1d-4fda-a3d5-e0ae6d8a916b)

## 🎤 Elevator Pitch

Before creating mynav, I often found myself frustrated when working on multiple projects using tmux, as I had to manually navigate between project directories. While tmux’s choose-tree feature allows jumping between active sessions, it relies on the tmux server staying alive and doesn't fully meet the needs of a robust workspace manager. mynav bridges this gap by combining tmux's powerful features with a workspace management system, enabling a more efficient and streamlined development workflow in a terminal environment.

## ✨ Features

- 📁 **Workspace Management**
  - Group workspaces into topics
  - Quick workspace creation and navigation
  - Lives directly on your filesystem

- 💻 **Advanced Session Management**
  - Session creation and switching
  - Create, modify, delete and enter sessions seamlessly
  - Live session preview with window/pane information

- 🔧 **Developer Experience**
  - Fuzzy search workspaces and sessions
  - Built on tmux
  - Extensive keyboard shortcuts
  - Git integration
  - Clean, intuitive Lazygit-like terminal UI
  - Vim-style navigation

## 🚀 Quick Start

### Try with docker

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

### One-Line Installation

```bash
curl -fsSL https://raw.githubusercontent.com/GianlucaP106/mynav/main/install.bash | bash
```

### Manual Installation

```bash
# Clone the repository
git clone https://github.com/GianlucaP106/mynav.git

# Navigate to project directory
cd mynav

# Build project
go build
```

### Prerequisites

- Tmux 3.0+
- Git (optional, for repository features)
- Terminal with UTF-8 support

---

## 📖 Usage

Mynav requires a root directory to initialize in. You may initialize multiple directories but not nested. You can start mynav anywhere with:

```bash
mynav
```

> This will look for an existing configuration if it exists (in the current or any parent directory).

You may specify a directory to launch in using:

```bash
mynav -path /your/root/path
```

You can use the `?` key in the TUI to view all the key bindings that are available in your context.

## 📺 Tmux Integration

Mynav integrates seamlessly with **tmux**, using it to manage sessions efficiently. When a session is created from a workspace, the workspace’s directory path is used as the tmux session name. This design keeps the state transparent and familiar, rather than hidden behind abstraction.

Once inside a tmux session, you can use all your usual tmux features. One key feature that enhances the mynav experience is the ability to **detach from the session** and return to the mynav interface by pressing **`Leader + D`**.

This tight integration gives you the full power of tmux while keeping mynav in sync with your development workflow.

## ⌨️ Key Bindings

### Navigation

| Key | Action | Context |
|-----|--------|---------|
| `h/←` | Focus left panel | Global |
| `l/→` | Focus right panel | Global |
| `j/↓` | Move down | List views |
| `k/↑` | Move up | List views |
| `Tab` | Toggle focus | Search dialog |
| `Esc` | Close/cancel | Dialogs |

### Actions

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

## ⚙️ Configuration

- MyNav uses a configuration system that supports multiple independent workspaces
- MyNav looks for configuration in the current or any parent directory
- Multiple independent directories can be initialized with MyNav
- Nested configurations are not allowed (invoking mynav nestedly will simply open the parent configuration)
- Home directory cannot be initialized as a MyNav workspace

## 🛠️ Development

### Setting Up Development Environment

Mynav is a straightforward, low-configuration project that only requires the Go runtime to get started in development.

## 🤝 Contributing

Ensure commits use conventional commits.

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

<p align="center">
  <a href="https://github.com/GianlucaP106/mynav/stargazers">⭐ Star on GitHub</a> •
  <a href="https://github.com/GianlucaP106/mynav/issues">📫 Report Bug</a> •
  <a href="https://github.com/GianlucaP106/mynav/discussions">💬 Discussions</a>
</p>
