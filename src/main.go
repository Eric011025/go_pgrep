package main

import (
	"fmt"

	"github.com/Eric011025/go_pgrep/pgrep"
)

func main() {
	pidList, err := pgrep.FindPID2PPID(1)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("len : ", len(pidList))
	for _, pid := range pidList {
		fmt.Println("pid : ", pid.Pid)
	}
	pgrep.KillChildProcess(6684)
}
