package ui

func (ui *UI) getKeyBindings(viewName string) []*KeyBindingMappings {
	keys := map[string][]*KeyBindingMappings{}
	keys["global"] = []*KeyBindingMappings{
		{
			key:    "q | Ctrl-c",
			action: "Quit",
		},
		{
			key:    "?",
			action: "Toggle help view",
		},
	}

	keys[ui.workspaces.viewName] = []*KeyBindingMappings{
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
			key:    "r",
			action: "Add/change description",
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
			key:    "t",
			action: "Open in terminal",
		},
		{
			key:    "s",
			action: "See workspace information",
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

	keys[ui.topics.viewName] = []*KeyBindingMappings{
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

	return keys[viewName]
}
