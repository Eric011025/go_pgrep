package main

import (
	"fmt"

	"github.com/Eric011025/eric_go_pgrep/go_pgrep"
)

func main() {
	pidList, err := go_pgrep.FindPID2PPID(6684)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("len : ", len(pidList))
	for _, pid := range pidList {
		fmt.Println("pid : ", pid.Pid)
	}
	go_pgrep.KillChildProcess(6684)
}
