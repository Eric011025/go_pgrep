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
				// When the filename is not a number, it is not a process, so skip it
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

// find pid by using ppid and kill all pid
func KillPidToPPid(ppid int) error {
	var (
		pidList []process.Process
		err     error
	)

	if pidList, err = GetPidToPPid(ppid); err != nil {
		return fmt.Errorf("KillPidToPPid::GetPidToPPid::pidList: %w", err)
	}

	for _, pid := range pidList {
		if err = pid.Kill(); err != nil {
			return fmt.Errorf("KillPidToPPid::pid.Kill: %w", err)
		}
	}

	return nil
}

// find pid by using ppid
func GetPidToPPid(ppid int) ([]process.Process, error) {
	var (
		pList []process.Process
		pid   []process.Process
		err   error
	)

	if pList, err = GetPidList(); err != nil {
		return nil, fmt.Errorf("GetPidToPPid::GetPidList::pList: %w", err)
	}

	for _, p := range pList {
		if p.PPid == ppid {
			pid = append(pid, p)
		}
	}

	return pid, nil
}

// find pid by using cmd
func GetPidToCmd(cmd string) ([]process.Process, error) {
	var (
		pList []process.Process
		pid   []process.Process
		err   error
	)

	if pList, err = GetPidList(); err != nil {
		return nil, fmt.Errorf("GetPidToCmd::GetPidList::pList: %w", err)
	}

	for _, p := range pList {
		if p.Cmd == cmd {
			pid = append(pid, p)
		}
	}

	if len(pid) == 0 {
		return nil, process.ProcessNotFound
	}

	return pid, nil
}

// SelfPid return self pid
func SelfPid() (process.Process, error) {
	return process.Process{Pid: syscall.Getpid()}, nil
}
