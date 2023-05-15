package process

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"
)

type Process struct {
	Pid   int    // process id
	PPid  int    // parent process id
	Pgrp  int    // process group id
	Cmd   string // process command
	State string // process state
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
	ProcessUnknown       = "Unknown"
)

// Process Signal
const (
	Kill = syscall.SIGKILL
)

// convert pid to process object
func New(id int) (Process, error) {
	var (
		p        Process
		statByte []byte
		encap    bool
		err      error
	)

	// read process status
	if statByte, err = os.ReadFile(fmt.Sprintf("/proc/%d/stat", id)); err != nil {
		return Process{}, fmt.Errorf("process list read fail : %w", err)
	}

	// split stat to fields. We split by space, but not when it's encapsulated by '(' and ')'
	stat := strings.FieldsFunc(string(statByte), func(r rune) bool {
		if r == '(' || r == ')' {
			encap = !encap
		}
		return !encap && r == ' '
	})

	// pid init
	if p.Pid, err = strconv.Atoi(stat[0]); err != nil {
		return Process{}, fmt.Errorf("process pid init fail : %w", err)
	}

	// ppid init
	if p.PPid, err = strconv.Atoi(stat[3]); err != nil {
		return Process{}, fmt.Errorf("process ppid init fail : %w", err)
	}

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
	default:
		p.State = ProcessUnknown
	}

	// process group id
	if p.Pgrp, err = strconv.Atoi(stat[4]); err != nil {
		return Process{}, fmt.Errorf("process gid init fail : %w", err)
	}

	return p, nil
}

// Kill can be used to kill a process
func (p Process) Kill() error {
	var (
		signal syscall.Signal
		able   bool
		err    error
	)

	sig := os.Kill

	if signal, able = sig.(syscall.Signal); !able {
		return errors.New("process : unsupported signal type")
	}

	if err = syscall.Kill(p.Pid, signal); err != nil {
		if err == syscall.ESRCH {
			return errors.New("process : process is already dead")
		}

		return err
	}

	return nil
}
