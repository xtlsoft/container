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

	var shares, quota, limit uint64 = 100, 100, 10 * 1024 * 1024

	cgroup, _ := container.NewOrLoadCgroups(cgroups.V1, cgroups.StaticPath("/test01"), &specs.LinuxResources{
		CPU: &specs.LinuxCPU{
			Shares: &shares,
			Quota: &quota,
		},
		Memory: &specs.LinuxMemory{
			Limit: &limit,
		},
	})

	defer cgroup.Delete()

	cmd.Start()

	cgroup.Add(cgroups.Process{Pid: cmd.Process.Pid})

	cmd.Wait()

}