# go_pgrep document





## Overview 

> find PID by PPID

This go_pgrep package motivated by linux command pgrep
I felt PID find with PPID was necessary, so I developed it.

```go
cmd, err := exec.Command("sh", "-c", "watch cat test.txt").output()
// "sh -c " command create child process
cmd.Process.Kill()  
// This kill function just kill parnet process (child process is still alive)
// We need other function to kill process
```





## Support functions 

- GetPidList : Get all Pid
- GetPidToPPid : Get Pid by Pid
- GetPidToCmd : Get Pid by cmd
- SelfPid : Get self Pid
- KillPidToPPid : Kill processes by PPid
- process.NewProcess : Create Process object
- (process).Kill : Kill process





## Usage

```go
package main

import (
	"fmt"

	"github.com/Eric011025/go_pgrep/pgrep"
)

func main() {
  // get pid by ppid
	ppidList, err := pgrep.GetPidToPPid(1)
	if err != nil {
		panic(err)
	}
	for _, ppidListItem := range ppidList {
    fmt.Println("pid : ", ppidListItem.Pid)
	}
  
  
  // kill pid by ppid
  err := pgrep.KillPidToPPid(6684)
	if err != nil {
		panic(err)
	}
  
  
  // get pid by cmd 
  result, err := pgrep.GetPidToCmd("watch")
	if err != nil {
		panic(err)
	}
	for _, value := range result {
		fmt.Println(value.Pid)
    // process kill
		value.Kill()
	}
  
  
  // get my pid
  selfPid, err := pgrep.SelfPid()
	if err != nil {
		panic(err)
	}
  fmt.Println("self pid : ", selfPid.Pid)
}

```


Contact email : ericpark011025@gmail.com
