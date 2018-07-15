package container

import (
	"syscall"
	"fmt"
	"os/exec"
	"os"
)

func NewNS() *Namespace {

	ns := new(Namespace)

	return ns

}

type Namespace struct {

	Cloneflags uintptr
	P_PS1 string

}

func (ns *Namespace) ApplyUTS() *Namespace {
	
	ns.Cloneflags = ns.Cloneflags | syscall.CLONE_NEWUTS

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

	return cmd

}

func (ns *Namespace) RedirectStd(cmd *exec.Cmd) *exec.Cmd {

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd

}