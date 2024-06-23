package ui

var globalKeyBindings []*KeyBindingMapping = []*KeyBindingMapping{
	{
		key:    "t",
		action: "Focus tmux session view",
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

var workspaceKeyBindings []*KeyBindingMapping = []*KeyBindingMapping{
	{
		key:    "j",
		action: "Move down",
	},
	{
		key:    "k",
		action: "Move up",
	},
	{
		key:    "down arrow | ctrl-j",
		action: "Focus Tmux view",
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
}

var topicKeyBindings []*KeyBindingMapping = []*KeyBindingMapping{
	{
		key:    "j",
		action: "Move down",
	},
	{
		key:    "k",
		action: "Move up",
	},
	{
		key:    "down arrow | ctrl-j",
		action: "Focus Port View",
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
}

var portKeyBindings []*KeyBindingMapping = []*KeyBindingMapping{
	{
		key:    "esc | up arrow | ctrl-k",
		action: "Focus Topic View",
	},
	{
		key:    "right arrow | ctrl-l",
		action: "Focus Tmux View",
	},
	{
		key:    "enter",
		action: "Open associated tmux session (if it exists)",
	},
	{
		key:    "D",
		action: "Kill port",
	},
}

var githubPrViewKeyBindings []*KeyBindingMapping = []*KeyBindingMapping{
	{
		key:    "j",
		action: "Move down",
	},
	{
		key:    "k",
		action: "Move up",
	},
	{
		key:    "esc | up arrow | ctrl-k",
		action: "Focus Workspace View",
	},
	{
		key:    "left arrow | ctrl-h",
		action: "Focus Tmux View",
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
}

func getTmuxKeyBindings(standalone bool) []*KeyBindingMapping {
	tmuxKeyBindings := []*KeyBindingMapping{
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
	}

	if !standalone {
		tmuxKeyBindings = append([]*KeyBindingMapping{
			{
				key:    "esc | up arrow | ctrl-k",
				action: "Focus Workspace View",
			},
			{
				key:    "left arrow | ctrl-h",
				action: "Focus Port View",
			},
			{
				key:    "right arrow | ctrl-l",
				action: "Focus Github View",
			},
		}, tmuxKeyBindings...)
	}

	return tmuxKeyBindings
}
