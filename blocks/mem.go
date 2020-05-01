package blocks

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type MemAvailBlock struct {
	mem_str  string
	nticks, PollFreq int
}

func (b *MemAvailBlock) Update() {

	b.nticks += 1
	b.nticks = b.nticks % b.PollFreq

	//Poll once every PollFreq ticks, and poll when starting, if PollFreq is 1,
	//just poll every tick
	if (b.PollFreq > 1) && (b.nticks != 1) {
		return
	}


	memi_cmd := exec.Command("cat", "/proc/meminfo")
	memi_pipe, memi_err := memi_cmd.StdoutPipe()

	if memi_err != nil {
		return
	}
	scanner := bufio.NewScanner(memi_pipe)
	memi_cmd.Start()

	memt, mema := 0,0

	for scanner.Scan() {
		splits := strings.Fields(scanner.Text())

		switch splits[0] {
		case "MemTotal:": {
			memt,_ = strconv.Atoi(splits[1])
		}
		case "MemAvailable:": {
			mema,_ = strconv.Atoi(splits[1])			
		}
	default: {
		continue
		}
	}

	}

	if (memt > 0) && (mema > 0) {
		b.mem_str = fmt.Sprintf("\uf538 %v%%",int((1.0 - float32(mema)/float32(memt))*100))
	}
}

func (b *MemAvailBlock) ToBlock() Block {
	out_b := NewPangoBlock()
	out_b.full_text = b.mem_str
	return out_b
}

func (b *MemAvailBlock) Check() bool {
	_, memi_err := os.Stat("/proc/meminfo")

	if memi_err != nil {
		return false
	}

	if b.PollFreq == 0 {
		b.PollFreq = 1
	}

	return true
}