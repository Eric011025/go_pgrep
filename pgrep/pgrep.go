package pgrep

import (
	"fmt"
	"io/fs"
	"os"
	"strconv"
	"syscall"

	"github.com/Eric011025/go_pgrep/process"
)

// get process list
func GetPidList() ([]process.Process, error) {
	var (
		p     process.Process
		pList []process.Process
		files []fs.DirEntry
		pid   int
		err   error
	)

	if files, err = os.ReadDir("/proc"); err != nil {
		return nil, fmt.Errorf("GetPidList::ReadDir::files: %w", err)
	}

	for _, file := range files {
		if file.IsDir() {
			if pid, err = strconv.Atoi(file.Name()); err != nil {
				// When the file name is not a number, it is not a process, so skip it
				continue
			}

			if p, err = process.NewProcess(pid); err == nil {
				return nil, fmt.Errorf("GetPidList::NewProcess::p: %w", err)
			}

			pList = append(pList, p)
		}
	}

	return pList, nil
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
func GetPidToPPid(ppid int) (pid []process.Process, err error) {
	pList, err := GetPidList()
	if err != nil {
		return
	}
	for _, p := range pList {
		if p.PPid == ppid {
			pid = append(pid, p)
		}
	}
	return
}

// find pid using by cmd
func GetPidToCmd(cmd string) (pid []process.Process, err error) {
	pList, err := GetPidList()
	if err != nil {
		return
	}
	for _, p := range pList {
		if p.Cmd == cmd {
			pid = append(pid, p)
		}
	}
	if len(pid) == 0 {
		return nil, process.ProcessNotFound
	}
	return
}

func SelfPid() (p process.Process, err error) {
	p = process.Process{Pid: syscall.Getpid()}
	return
}
