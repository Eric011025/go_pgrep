package pgrep

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// find pid using by ppid
// this function using pgrep command
func FindPID2PPID(ppid int) (pidList []os.Process, err error) {
	pidCmdOut, err := exec.Command("sh", "-c", fmt.Sprintf("pgrep -P %d | awk '{ print $1 }' ", ppid)).Output()
	if err != nil {
		return
	}
	tmp := strings.Split(string(pidCmdOut), "\n")
	for _, pidStr := range tmp[0 : len(tmp)-1] {
		var pid *os.Process
		var pidInt int
		pidInt, err = strconv.Atoi(pidStr)
		if err != nil {
			return
		}
		pid, err = os.FindProcess(pidInt)
		if err != nil {
			return
		}
		pidList = append(pidList, *pid)
	}
	return
}

// find pid using ppid and kill all of pid
func KillChildProcess(ppid int) (err error) {
	pidList, err := FindPID2PPID(ppid)
	if err != nil {
		return
	}
	for _, pid := range pidList {
		err = pid.Kill()
		if err != nil {
			return
		}
	}
	return err
}
