package container

import (
	"syscall"
	"fmt"
	"os/exec"
	"os"
)

func NewNS() *Namespace {

	ns := new(Namespace)

	ns.P_UID = -1;

	return ns

}

type Namespace struct {

	Cloneflags uintptr
	P_PS1 string
	P_UID int
	P_GID int

}

func (ns *Namespace) ApplyUTS() *Namespace {
	
	ns.Cloneflags = ns.Cloneflags | syscall.CLONE_NEWUTS

	return ns

}

func (ns *Namespace) ApplyUser() *Namespace {

	if ns.P_UID == -1 {
		P_UID = os.Getuid()
		P_GID = os.Getgid()
	}

	ns.Cloneflags = ns.Cloneflags | syscall.CLONE_NEWUSER

	return ns

}

func (ns *Namespace) SetUID(uid int, gid int) *Namespace {

	ns.P_UID = uid
	ns.P_GID = gid

	return ns

}

func (ns *Namespace) ApplyMount() *Namespace {

	ns.Cloneflags = ns.Cloneflags | syscall.CLONE_NEWNS

	return ns

}

func (ns *Namespace) ApplyNet() *Namespace {

	ns.Cloneflags = ns.Cloneflags | syscall.CLONE_NEWNET

	return ns

}

func (ns *Namespace) ApplyPID() *Namespace {

	ns.Cloneflags = ns.Cloneflags | syscall.CLONE_NEWPID

	return ns

}

func (ns *Namespace) ApplyIPC() *Namespace {

	ns.Cloneflags = ns.Cloneflags | syscall.CLONE_NEWIPC

	return ns

}

func (ns *Namespace) SetPS1(name string) *Namespace {

	ns.P_PS1 = name

	return ns

}

func (ns *Namespace) Command(cmdl string, arg ...string) *exec.Cmd {

	cmd := exec.Command(cmdl, arg...)

	var env []string

	if ns.P_PS1 != "" {
		env = append(env, fmt.Sprintf("PS1=%s # ", ns.P_PS1))
	}

	cmd.Env = env

	cmd.SysProcAttr = &syscall.SysProcAttr {
		Cloneflags: ns.Cloneflags,
	}

	if ns.P_UID != -1 {
		cmd.SysProcAttr.UidMappings = []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID: ns.P_UID,
				Size: 1,
			},
		}
		cmd.SysProcAttr.GidMappings = []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID: ns.P_GID,
				Size: 1,
			},
		}
	}

	return cmd

}

func (ns *Namespace) RedirectStd(cmd *exec.Cmd) *exec.Cmd {

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd

}