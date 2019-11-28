package common

import (
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

func RunInBackground(cmd string) (*os.Process, error) {
	devnull, _ := os.Open(os.DevNull)
	defer devnull.Close()

	var err error
	args := []string{"sh", "-c", "exec " + cmd}
	args[0], err = exec.LookPath(args[0])
	if err != nil {
		return nil, err
	}

	// ignoreing SIGHUP
	signal.Ignore(syscall.SIGHUP)

	attr := &os.ProcAttr{
		Files: []*os.File{
			devnull, // stdin
			devnull,
			devnull,
		},
		Sys: &syscall.SysProcAttr{
			Setsid:     true,
			Foreground: false,
		},
	}
	p, err := os.StartProcess(args[0], args, attr)
	if err != nil {
		return nil, err
	}

	return p, err
}
