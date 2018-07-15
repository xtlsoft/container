package container

import (
	"syscall"
	"fmt"
	"os/exec"
)

func NewNS() *Namespace {

	ns := new(Namespace)

	return ns

}

type Namespace struct {

	Cloneflags uintptr
	UTSHostname string

}

func (ns *Namespace) ApplyUTS() *Namespace {
	
	ns.Cloneflags = ns.Cloneflags | syscall.CLONE_NEWUTS

	return ns

}

func (ns *Namespace) SetUTSHostName(name string) *Namespace {

	ns.UTSHostname = name

	return ns

}

func (ns *Namespace) Command(cmdl string, arg ...string) *exec.Cmd {

	cmd := exec.Command(cmdl, arg...)

	cmd.Env = []string{fmt.Sprintf("PS1=%s # ", ns.UTSHostname)}

	cmd.SysProcAttr = &syscall.SysProcAttr {
		Cloneflags: ns.Cloneflags,
	}

}