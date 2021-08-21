package pgrep

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// find pid using by ppid
// this function using pgrep command
func GetPidToPPid(ppid int) (pidList []os.Process, err error) {
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

// get process list
func GetPidList() (pList []os.Process, err error) {
	files, err := ioutil.ReadDir("/proc")
	if err != nil {
		return
	}

	for _, file := range files {
		if file.IsDir() {
			if pid, typeErr := strconv.Atoi(file.Name()); typeErr == nil {
				pList = append(pList, os.Process{Pid: pid})
			}
		}
	}
	return
}

// find pid by ppid and kill all pid
func KillPidToPPID(ppid int) (err error) {
	pidList, err := GetPidToPPid(ppid)
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
