package system

import (
	"os"

	"github.com/shirou/gopsutil/process"
)

func IsProcessChildOf(child int, parent int) bool {
	var pid int32
	pid = int32(child)
	parentPid := int32(parent)

	for {
		if pid == parentPid {
			return true
		}

		proc, err := process.NewProcess(pid)
		if err != nil {
			return false
		}

		ppid, err := proc.Ppid()
		if err != nil {
			return false
		}

		pid = ppid
	}
}

func KillProcess(pid int) error {
	proc, err := os.FindProcess(pid)
	if err != nil {
		return err
	}

	if err := proc.Kill(); err != nil {
		return err
	}

	return nil
}
