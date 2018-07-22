package main

import (
	"os"
	"fmt"
	"github.com/xtlsoft/container"
)

func usage() {

	fmt.Println(`USAGE: aufs [COMMAND] [ARGUMENT]

COMMANDS:
	init		Init an empty aufs environment in this directory.
	new-layer	Create an empty layer.
	new-image	Create an empty image layer.
	delete		Delete a layer.
	mount		Mount current aufs to some place.`)

}

func main() {

	if len(os.Args) < 2 {
		usage()
		return
	}

	aufs := container.NewAUFS("./")

	switch (os.Args[1]) {

		case "init":
			container.NewAUFS("./");
		case "new-layer":
			aufs.NewLayer(os.Args[2])
		case "new-image":
			aufs.NewImageLayer(os.Args[2])
		case "delete":
			aufs.RemoveLayer(os.Args[2])
		case "mount":
			aufs.Mount(os.Args[2], os.Args[3])

	}

}