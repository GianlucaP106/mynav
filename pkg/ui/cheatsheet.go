package ui

func getKeyBindings(viewName string) []*KeyBindingMapping {
	keys := map[string][]*KeyBindingMapping{}
	keys["global"] = []*KeyBindingMapping{
		{
			key:    "t",
			action: "Tmux session view",
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
			key:    "a",
			action: "Create a workspace",
		},
		{
			key:    "d",
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
			key:    "r",
			action: "Rename topic",
		},
		{
			key:    "a",
			action: "Create a topic",
		},
		{
			key:    "d",
			action: "Delete a topic",
		},
		{
			key:    "enter",
			action: "Open topic",
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
			key:    "esc",
			action: "Exit view",
		},
		{
			key:    "d",
			action: "Delete session",
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

	return keys[viewName]
}
