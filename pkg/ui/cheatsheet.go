package ui

var globalKeyBindings []*KeyBindingMapping = []*KeyBindingMapping{
	{
		key:    "[",
		action: "Prev tab",
	},
	{
		key:    "]",
		action: "Next tab",
	},
	{
		key:    "q | Ctrl-c",
		action: "Quit",
	},
	{
		key:    "?",
		action: "Toggle help view",
	},
}

var workspaceKeyBindings []*KeyBindingMapping = append([]*KeyBindingMapping{
	{
		key:    "j",
		action: "Move down",
	},
	{
		key:    "k",
		action: "Move up",
	},
	{
		key:    "left arrow | ctrl-h",
		action: "Focus topic view",
	},
	{
		key:    "a",
		action: "Create a workspace",
	},
	{
		key:    "D",
		action: "Delete a workspace",
	},
	{
		key:    "e",
		action: "Add/change description",
	},
	{
		key:    "r",
		action: "Rename workspace",
	},
	{
		key:    "g",
		action: "Clone git repo",
	},
	{
		key:    "G",
		action: "Open browser to git repo",
	},
	{
		key:    "enter",
		action: "Open in tmux/open in neovim",
	},
	{
		key:    "v",
		action: "Open in neovim",
	},
	{
		key:    "m",
		action: "Open in terminal",
	},
	{
		key:    "s",
		action: "See workspace information",
	},
	{
		key:    "X",
		action: "Kill tmux session",
	},
	{
		key:    "/",
		action: "Search by name",
	},
	{
		key:    "esc",
		action: "Escape search / Go back",
	},
}, globalKeyBindings...)

var topicKeyBindings []*KeyBindingMapping = append([]*KeyBindingMapping{
	{
		key:    "j",
		action: "Move down",
	},
	{
		key:    "k",
		action: "Move up",
	},
	{
		key:    "enter | right arrow | ctrl-l",
		action: "Open topic",
	},
	{
		key:    "r",
		action: "Rename topic",
	},
	{
		key:    "a",
		action: "Create a topic",
	},
	{
		key:    "D",
		action: "Delete a topic",
	},
	{
		key:    "/",
		action: "Search by name",
	},
	{
		key:    "esc",
		action: "Escape search",
	},
}, globalKeyBindings...)

var portKeyBindings []*KeyBindingMapping = append([]*KeyBindingMapping{
	{
		key:    "enter",
		action: "Open associated tmux session (if it exists)",
	},
	{
		key:    "D",
		action: "Kill port",
	},
}, globalKeyBindings...)

var githubPrViewKeyBindings []*KeyBindingMapping = append([]*KeyBindingMapping{
	{
		key:    "j",
		action: "Move down",
	},
	{
		key:    "k",
		action: "Move up",
	},
	{
		key:    "left arrow | ctrl-h",
		action: "Focus Repo View",
	},
	{
		key:    "o",
		action: "Open Browser to PR",
	},
	{
		key:    "L",
		action: "Login with device code and browser",
	},
	{
		key:    "P",
		action: "Login with personal access token",
	},
	{
		key:    "O",
		action: "Logout",
	},
}, globalKeyBindings...)

var tmuxKeyBindings []*KeyBindingMapping = append([]*KeyBindingMapping{
	{
		key:    "D",
		action: "Delete session",
	},
	{
		key:    "X",
		action: "Kill ALL tmux sessions",
	},
	{
		key:    "W",
		action: "Kill ALL non-external (has a workspace) tmux sessions",
	},
	{
		key:    "a",
		action: "New external session (not associated to a workspace)",
	},
	{
		key:    "enter",
		action: "Attach to session",
	},
}, globalKeyBindings...)
