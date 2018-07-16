package container

import (
	"log"
	"os"
	"os/exec"
	"syscall"
)

var InitProcessChannel = make(chan func(string, []string))

// Schedule Init Process
func init() {

	if len(os.Args) > 1 && os.Args[1] == "_container_init_" {

		// Do schedule
		// WIP
		command := os.Args[2]

		var args []string

		for i := 2; i < len(os.Args); i ++ {

			args = append(args, os.Args[i])

		}

		// Mount proc filesystem
		// defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
		// syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")

		go func() {

			// Get the process function
			callback := <- InitProcessChannel

			// Call the function
			callback(command, args)

			// Do execve
			if err := syscall.Exec(command, args, os.Environ()); err != nil {
				log.Fatalf("Initialize error: %v", err)
			}

		}()

	}

}

func NewInitProcess(ns *Namespace, tty bool, command string, arguments ...string) *exec.Cmd {

	var args = []string {
		"_container_init_",
		command,
	}

	for _, v := range arguments {
		args = append(args, v)
	}

	cmd := ns.Command("/proc/self/exe", args...)

	if tty {
		cmd = ns.RedirectStd(cmd)
	}

	return cmd

}
