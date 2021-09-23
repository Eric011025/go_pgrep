package pgrep

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

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
func KillPidToPPid(ppid int) (err error) {
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

// find pid using by ppid
func GetPidToPPid(ppid int) (pid []os.Process, err error) {
	pList, err := GetPidList()
	for pIndex, p := range pList {
		statByte, err := ioutil.ReadFile(fmt.Sprintf("/proc/%d/stat", p.Pid))
		if err != nil {
			return nil, err
		}
		if strings.Split(string(statByte), " ")[3] == strconv.Itoa(ppid) {
			pid = append(pid, pList[pIndex])
		}
	}
	return
}
