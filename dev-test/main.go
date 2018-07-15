package main

import (
	container ".."
)

func main() {

	ns := container.NewNS()

	ns.ApplyUTS().
		ApplyUser().
		ApplyPID().
		ApplyMount().
		ApplyNet().
		ApplyIPC().
		SetPS1("-[Hello Namespace]-")

	cmd := ns.Command("bash")

	ns.RedirectStd(cmd).Run()

}