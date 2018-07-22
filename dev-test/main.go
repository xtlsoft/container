package main

import (
	container ".."
	"github.com/containerd/cgroups"
	specs "github.com/opencontainers/runtime-spec/specs-go"
	"os"
	"fmt"
)

func main() {

	if container.IsInitProcess() {
		
		cmd, args := container.ParseInit()

		fmt.Println("Starting From This")
		container.MountBasicFileSystems()

		container.Execve(cmd, args)
		
		return
	}

	if len(os.Args) > 1 && os.Args[1] == "new" {
		r1, _ := os.Open("/")
		r, _ := r1.Readdir(1)
		for _, v := range r {
			fmt.Println(v.Name())
		}
		return
	}

	ns := container.NewNS()

	ns.ApplyUTS().
		ApplyUser().
		ApplyPID().
		ApplyMount().
		ApplyNet().
		ApplyIPC().
		SetPS1("-[Hello Namespace]-")

	cmd := container.NewInitProcess(ns, true, "/bin/sh")

	var shares uint64 = 100
	// var quota int64 = 100
	// var limit int64 = 10 * 1024 * 1024

	cgroup, err := container.NewOrLoadCgroups(cgroups.V1, cgroups.StaticPath("/test01"), &specs.LinuxResources{
		CPU: &specs.LinuxCPU{
			Shares: &shares,
			// Quota: &quota,
		},
		Memory: &specs.LinuxMemory{
			// Limit: &limit,
		},
	})

	if err != nil {
		panic(err)
	}

	defer cgroup.Delete()

	cmd = ns.RedirectStd(cmd)

	cmd.Start()

	cgroup.Add(cgroups.Process{Pid: cmd.Process.Pid})

	cmd.Wait()

}