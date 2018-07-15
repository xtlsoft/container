package container

import (
	"github.com/containerd/cgroups"
	specs "github.com/opencontainers/runtime-spec/specs-go"
)

// For futher use,
// Please import github.com/containerd/cgroups .
// This is only a wrapper of some useful functions.

func NewOrLoadCgroups(hierarchy cgroups.Hierarchy, path cgroups.Path, resources specs.LinuxResources) (cgroups.Cgroup, error) {

	// Try load...
	rslt, err := cgroups.Load(hierarchy, path)
	if err == nil {
		return rslt, nil
	}

	if err != cgroups.ErrCgroupDeleted {
		return rslt, err
	}

	return cgroups.New(hierarchy, path, resources)

}