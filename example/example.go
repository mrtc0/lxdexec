package main

import (
	"fmt"

	"github.com/mrtc0/lxdexec"
)

func main() {
	var containerName string = "admin"
	if lxdexec.ContainerExist(containerName) {
		_, uuid := lxdexec.ContainerExec(containerName, []string{"echo", "hello"})
		lxdexec.Wait(uuid)
		_, stdout, stderr := lxdexec.ContainerGetStd(containerName, uuid)
		fmt.Println(stdout)
		fmt.Println(stderr)
	}
}
