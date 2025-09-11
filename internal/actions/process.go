package actions

import "github.com/shirou/gopsutil/v3/process"

// KillProcess attempts to gracefully terminate a process by its PID using SIGTERM.
func KillProcess(pid int32) error {
	if pid == 0 {
		return nil // Nothing to kill
	}
	p, err := process.NewProcess(pid)
	if err != nil {
		return err
	}
	return p.Terminate()
}

// ForceKillProcess forcefully terminates a process by its PID using SIGKILL.
func ForceKillProcess(pid int32) error {
	if pid == 0 {
		return nil
	}
	p, err := process.NewProcess(pid)
	if err != nil {
		return err
	}
	return p.Kill()
}

