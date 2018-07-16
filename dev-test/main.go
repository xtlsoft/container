package main

import (
	container ".."
	"github.com/containerd/cgroups"
	specs "github.com/opencontainers/runtime-spec/specs-go"
)

func main() {

	if container.IsInitProcess() {
		container.InitProcessChannel <- func(cmd string, args []string) {
			container.MountBasicFileSystems()
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

	cmd := container.NewInitProcess(ns, true, "sh")

	cgroup, _ := container.NewOrLoadCgroups(cgroups.V1, cgroups.StaticPath("/test01"), &specs.LinuxResources{
		CPU: &specs.LinuxCPU{
			Shares: &uint64(100),
			Quota: &uint64(100),
		},
		Memory: &specs.LinuxMemory{
			Limit: &int64(10 * 1024 * 1024),	
		},
	})

	defer cgroup.Delete()

	cmd.Start()

	cgroup.Add(cgroups.Process{Pid: cmd.Process.Pid})

	cmd.Wait()

}