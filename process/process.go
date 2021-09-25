package process

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

type Process struct {
	Pid  int
	PPid int
	Cmd  string
}

//
func NewProcess(id int) (p Process, err error) {
	statByte, err := ioutil.ReadFile(fmt.Sprintf("/proc/%d/stat", id))
	if err != nil {
		return Process{}, err
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
	return
}

// kill function
