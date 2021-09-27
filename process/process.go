package process

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"sync"
	"syscall"
)

type Process struct {
	Pid   int
	PPid  int
	Cmd   string
	sigMu sync.RWMutex
}

var (
	ProcessNotFound error = errors.New("process : process not founded")
)

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

	// pid init
	pid, err := strconv.Atoi(strings.Split(string(statByte), " ")[0])
	if err != nil {
		return Process{}, err
	}
	p.Pid = pid

	// ppid init
	ppid, err := strconv.Atoi(strings.Split(string(statByte), " ")[3])
	if err != nil {
		return Process{}, err
	}
	p.PPid = ppid

	p.Cmd = strings.TrimRight(strings.TrimLeft(strings.Split(string(statByte), " ")[1], "("), ")")
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
