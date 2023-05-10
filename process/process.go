package process

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"syscall"
)

type Process struct {
	Pid   int
	PPid  int
	Cmd   string
	State string
	Pgrp  int
}

// Process state
var (
	ProcessRunnig        = "R"
	ProcessSleeping      = "S"
	ProcessDiskSleeping  = "D"
	ProcessZombie        = "Z"
	ProcessTracedStopped = "T"
	ProcessPaging        = "W"
	ProcessIdle          = "I"
)

// Error Type
var (
	ProcessNotFound error = errors.New("process : process not found")
)

// Process Signal
const (
	Kill = syscall.SIGKILL
)

// convert pid to process object
func NewProcess(id int) (p Process, err error) {
	// read process status

	statByte, err := ioutil.ReadFile(fmt.Sprintf("/proc/%d/stat", id))
	if err != nil {
		return Process{}, ProcessNotFound
	}
	stat := strings.Split(string(statByte), " ")

	// pid init
	pid, err := strconv.Atoi(stat[0])
	if err != nil {
		return Process{}, err
	}
	p.Pid = pid

	// ppid init
	ppid, err := strconv.Atoi(stat[3])
	if err != nil {
		return Process{}, err
	}
	p.PPid = ppid

	p.Cmd = strings.TrimRight(strings.TrimLeft(stat[1], "("), ")")

	// process state
	state := stat[2]
	switch state {
	case ProcessRunnig:
		p.State = ProcessRunnig
	case ProcessSleeping:
		p.State = ProcessSleeping
	case ProcessDiskSleeping:
		p.State = ProcessDiskSleeping
	case ProcessZombie:
		p.State = ProcessZombie
	case ProcessTracedStopped:
		p.State = ProcessTracedStopped
	case ProcessPaging:
		p.State = ProcessPaging
	case ProcessIdle:
		p.State = ProcessIdle
	}

	// process group id
	pgrp, err := strconv.Atoi(stat[4])
	if err != nil {
		return Process{}, err
	}
	p.Pgrp = pgrp

	return
}

// process Kill
func (p Process) Kill() (err error) {
	sig := os.Kill
	signal, able := sig.(syscall.Signal)
	if able == false {
		return errors.New("process : unsupported signal type")
	}
	if err = syscall.Kill(p.Pid, signal); err != nil {
		if err == syscall.ESRCH {
			return errors.New("process : process is already dead")
		}
		return err
	}
	return
}
