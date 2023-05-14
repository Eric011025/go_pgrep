# go_pgrep: A PID Management Library in Go

## Introduction

The `go_pgrep` library is a Go package that enables you to conveniently handle Linux process IDs (PIDs) and parent process IDs (PPIDs). Inspired by the Linux command `pgrep`, this package was developed out of the necessity to find and manage PIDs associated with a given PPID. 

Often, when working with processes in Go, we encounter situations where killing a parent process does not terminate its child processes. For instance:

```go
cmd, err := exec.Command("sh", "-c", "watch cat test.txt").Output()
// The command above creates a child process.
cmd.Process.Kill()  
// Killing the parent process does not terminate the child process.
```

To handle such cases, `go_pgrep` provides a set of functions that lets you manage processes more effectively.

## Features

`go_pgrep` provides the following functions for managing processes:

- `GetPidList`: Retrieves all PIDs.
- `GetPidToPPid`: Retrieves PIDs associated with a given PPID.
- `GetPidToCmd`: Retrieves PIDs associated with a given command.
- `SelfPid`: Retrieves the PID of the current process.
- `KillPidToPPid`: Terminates processes associated with a given PPID.
- `(process).Kill`: Terminates a process.

## Usage

Here's an example demonstrating how to use `go_pgrep`:

```go
package main

import (
	"fmt"

	"github.com/Eric011025/go_pgrep/pgrep"
)

func main() {
  	// Get PID by PPID
	ppidList, err := pgrep.GetPidToPPid(1)
	if err != nil {
		panic(err)
	}
	for _, ppidListItem := range ppidList {
    fmt.Println("PID: ", ppidListItem.Pid)
	}
  
  	// Kill PID by PPID
  	err = pgrep.KillPidToPPid(6684)
	if err != nil {
		panic(err)
	}
  
  	// Get PID by command
  	result, err := pgrep.GetPidToCmd("watch")
	if err != nil {
		panic(err)
	}
	for _, value := range result {
		fmt.Println("PID: ", value.Pid)
    	// Kill the process
		value.Kill()
	}
  
  	// Get the current process's PID
  	selfPid, err := pgrep.SelfPid()
	if err != nil {
		panic(err)
	}
  	fmt.Println("Self PID: ", selfPid.Pid)
}
```

## Contact
For any queries or suggestions, feel free to reach out to Eric Park at ericpark011025@gmail.com.
