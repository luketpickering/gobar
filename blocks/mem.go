package blocks

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"bytes"

	pgu "github.com/luketpickering/gobar/pangoutils"
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


	memi_bytes, memi_err := exec.Command("cat", "/proc/meminfo").Output()
	if memi_err != nil {
		return
	}
	scanner := bufio.NewScanner(bytes.NewBuffer(memi_bytes))

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
		mem_pc := int((1.0 - float32(mema)/float32(memt))*100)
		if mem_pc > 90 {
			b.mem_str = pgu.NewPangoStrU(fmt.Sprintf(" \uf538 %v%% ",mem_pc)).SetBGColor(pgu.Red).SetFGColor(pgu.DarkGrey).String()
		} else {
			b.mem_str = fmt.Sprintf(" \uf538 %v%% ",mem_pc)
		}
	}
}

func (b *MemAvailBlock) ToBlock() Block {
	out_b := NewPangoBlock()
	out_b.full_text = b.mem_str
	out_b.min_width = " \uf538 100% "
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