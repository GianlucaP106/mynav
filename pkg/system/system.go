package system

import (
	"os/exec"
	"runtime"
	"strconv"
	"time"

	"github.com/atotto/clipboard"
)

type OS = uint

const (
	Darwin OS = iota
	Linux
	Unsuported
)

func DetectOS() OS {
	switch runtime.GOOS {
	case "darwin":
		return Darwin
	case "linux":
		return Linux
	default:
		return Unsuported
	}
}

func IsWarpInstalledMac() bool {
	if DetectOS() != Darwin {
		return false
	}

	return Exists("/Applications/Warp.app")
}

func IsItermInstalledMac() bool {
	if DetectOS() != Darwin {
		return false
	}

	return Exists("/Applications/iTerm.app")
}

func OpenBrowser(url string) error {
	var cmd string

	switch DetectOS() {
	case Linux:
		cmd = "xdg-open"
	case Darwin:
		cmd = "open"
	}

	return exec.Command(cmd, url).Start()
}

func CopyToClip(s string) error {
	return clipboard.WriteAll(s)
}

func TimeFormat() string {
	return "Mon, 02 Jan 15:04:05"
}

func DoesProgramExist(program string) bool {
	_, err := exec.LookPath(program)
	return err == nil
}

func OpenLazygit(path string) error {
	return exec.Command("lazygit", "-p", path).Run()
}

func UnixTime(t string) time.Time {
	activityInt, err := strconv.Atoi(t)
	if err != nil {
	}

	time := time.Unix(int64(activityInt), 0)
	return time
}
