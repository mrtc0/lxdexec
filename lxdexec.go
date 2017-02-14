package lxdexec

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/mrtc0/lxdexec/unixsocket"
)

func ContainerExist(containerName string) bool {
	cli := unixsocket.NewClient("")
	var client *http.Client
	client = cli
	re, err := client.Get("http://unix.socket/1.0/containers/" + containerName)
	defer re.Body.Close()
	if err != nil {
		return false
	}
	var response map[string]interface{}
	bytes, err := ioutil.ReadAll(re.Body)
	if err != nil {
		return false
	}
	err = json.Unmarshal(bytes, &response)
	if err != nil {
		return false
	}
	code, ok := response["error_code"]
	if ok && code.(float64) == 404 {
		return false
	}
	return true
}

func ContainerExec(containerName string, command []string) (error, string) {
	cli := unixsocket.NewClient("")
	var client *http.Client
	client = cli
	str := `
    {
        "command": %s,       
        "environment": {},              
        "wait-for-websocket": false,    
		"record-output": true,
        "interactive": true
    }
    `
	commands, err := json.Marshal(&command)
	if err != nil {
		return err, ""
	}
	post := fmt.Sprintf(str, string(commands))
	re, err := client.Post("http://unix.socket/1.0/containers/"+containerName+"/exec", "application/json", strings.NewReader(post))
	if err != nil {
		return err, ""
	}
	defer re.Body.Close()
	bytes, err := ioutil.ReadAll(re.Body)
	if err != nil {
		return err, ""
	}
	var response map[string]interface{}
	err = json.Unmarshal(bytes, &response)
	if err != nil {
		return err, ""
	}
	code, ok := response["status_code"]
	if !ok || code.(float64) != 100 {
		return errors.New(response["error"].(string)), ""
	}
	uuid := response["metadata"].(map[string]interface{})["id"].(string)
	return nil, uuid
}

func ContainerGetStd(container_name string, uuid string) (error, string, string) {
	cli := unixsocket.NewClient("")
	var client *http.Client
	client = cli

	// Get Stdout
	re, err := client.Get("http://unix.socket/1.0/containers/" + container_name + "/logs/exec_" + uuid + ".stdout")
	if err != nil {
		return err, "", ""
	}
	defer re.Body.Close()
	bytes, err := ioutil.ReadAll(re.Body)
	if err != nil {
		fmt.Println(err)
		return err, "", ""
	}
	stdout := string(bytes)

	// Get Stderr
	re, err = client.Get("http://unix.socket/1.0/containers/" + container_name + "/logs/exec_" + uuid + ".stderr")
	if err != nil {
		return err, "", ""
	}
	defer re.Body.Close()
	bytes, err = ioutil.ReadAll(re.Body)
	if err != nil {
		fmt.Println(err)
		return err, "", ""
	}
	stderr := string(bytes)
	return nil, stdout, stderr
}

func Wait(uuid string) {
	cli := unixsocket.NewClient("")
	var client *http.Client
	client = cli
	re, _ := client.Get("http://unix.socket/1.0/operations/" + uuid + "/wait")
	defer re.Body.Close()
}
