package system

import (
	"os/exec"
	"strings"
)

func IsFzfInstalled() bool {
	_, err := exec.LookPath("fzf")
	return err == nil
}

func FuzzyFind(lst []string, search string) []string {
	input := strings.Join(lst, "\n")
	cmd := exec.Command("fzf", "--filter", search)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.Output()
	if err != nil {
		return []string{}
	}

	items := strings.Split(string(out), "\n")
	return items
}
