package system

import (
	"os"

	"github.com/shirou/gopsutil/process"
)

func GetChildProcesses(pid int) ([]*process.Process, error) {
	proc, err := process.NewProcess(int32(pid))
	if err != nil {
		return nil, err
	}

	children, err := proc.Children()
	if err != nil {
		return nil, err
	}

	return children, nil
}

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

type ProcessInfo struct {
	Exe string
}

func GetProcessInfo(pid int) *ProcessInfo {
	proc, err := process.NewProcess(int32(pid))
	if err != nil {
		return nil
	}

	exe, err := proc.Exe()
	if err != nil {
		return nil
	}

	return &ProcessInfo{
		Exe: exe,
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

func IsCurrentProcessHomeDir() bool {
	homeDir, _ := os.UserHomeDir()
	cwd, _ := os.Getwd()
	return homeDir == cwd
}
