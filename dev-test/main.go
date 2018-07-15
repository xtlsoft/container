package main

import (
	container "../../Container"
)

func main() {

	ns := container.NewNS()

	ns.ApplyUTS().SetUTSHostname("-[Hello Namespace]-")

	ns.Exec("bash").Run()

}