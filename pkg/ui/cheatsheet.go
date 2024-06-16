package ui

func getKeyBindings(viewName string) []*KeyBindingMapping {
	keys := map[string][]*KeyBindingMapping{}
	keys["global"] = []*KeyBindingMapping{
		{
			key:    "q | Ctrl-c",
			action: "Quit",
		},
		{
			key:    "?",
			action: "Toggle help view",
		},
	}

	keys[WorkspacesViewName] = []*KeyBindingMapping{
		{
			key:    "j",
			action: "Move down",
		},
		{
			key:    "k",
			action: "Move up",
		},
		{
			key:    "down arrow",
			action: "Focus Tmux view",
		},
		{
			key:    "left arrow",
			action: "Go back",
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
			key:    "x",
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

	keys[TopicViewName] = []*KeyBindingMapping{
		{
			key:    "j",
			action: "Move down",
		},
		{
			key:    "k",
			action: "Move up",
		},
		{
			key:    "arrow down",
			action: "Focus Port View",
		},
		{
			key:    "enter | arrow right",
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

	keys[TmuxSessionViewName] = []*KeyBindingMapping{
		{
			key:    "esc | arrow up",
			action: "Focus Workspace View",
		},
		{
			key:    "arrow left",
			action: "Focus Port View",
		},
		{
			key:    "d",
			action: "Delete session",
		},
		{
			key:    "x",
			action: "Kill ALL tmux sessions",
		},
		{
			key:    "w",
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

	keys[PortViewName] = []*KeyBindingMapping{
		{
			key:    "esc | arrow up",
			action: "Focus Topic View",
		},
		{
			key:    "arrow right",
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

	return keys[viewName]
}
