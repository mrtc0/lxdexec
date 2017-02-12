package lxdexec

import (
	"testing"
)

func TestContainerExec(t *testing.T) {
	var containerName string = "admin"

	if ContainerExist(containerName) {
		var expectStdout string = "hello\n"
		_, uuid := ContainerExec(containerName, []string{"echo", "hello"})
		Wait(uuid)
		_, stdout, stderr := ContainerGetStd(containerName, uuid)
		if stdout != expectStdout {
			t.Errorf("stdout should be %s, but %s", expectStdout, stdout)
		}
		_, uuid = ContainerExec(containerName, []string{"hoge"})
		Wait(uuid)
		_, stdout, stderr = ContainerGetStd(containerName, uuid)
		if stderr != "" {
			t.Error("stderr should be not empty, but empty")
		}
	}
}
