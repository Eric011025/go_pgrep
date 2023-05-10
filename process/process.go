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
)

// Error Tyep
var (
	ProcessNotFound error = errors.New("process : process not founded")
)

// Process Signal
const (
	Kill = syscall.SIGKILL
)

// convert pid to process object
func NewProcess(id int) (Process, error) {
	var (
		p        Process
		pid      int
		ppid     int
		pgrp     int
		statByte []byte
		encap    bool
		err      error
	)

	// read process status
	statByte, err = ioutil.ReadFile(fmt.Sprintf("/proc/%d/stat", id))
	if err != nil {
		return Process{}, ProcessNotFound
	}

	// split stat to fields. We split by space, but not when it's encapsulated by '(' and ')'
	stat := strings.FieldsFunc(string(statByte), func(r rune) bool {
		if r == '(' || r == ')' {
			encap = !encap
		}
		return !encap && r == ' '
	})

	// pid init
	if pid, err = strconv.Atoi(stat[0]); err != nil {
		return Process{}, err
	}
	p.Pid = pid

	// ppid init
	if ppid, err = strconv.Atoi(stat[3]); err != nil {
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
	}

	// process group id
	if pgrp, err = strconv.Atoi(stat[4]); err != nil {
		return Process{}, err
	}
	p.Pgrp = pgrp

	return p, nil
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
