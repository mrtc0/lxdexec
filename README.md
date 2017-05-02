# lxdexec

LXD Exec API wrapper

# Example

```
_, uuid := lxdexec.ContainerExec('container1', ["echo", "hello"])
lxdexec.Wait(uuid)
_, stdout, stderr := lxdexec.COntainerGetStd('container1', uuid)
fmt.Println(stdout) // => hello
```
