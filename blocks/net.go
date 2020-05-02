package blocks

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
	"bytes"

	pgu "github.com/luketpickering/gobar/pangoutils"
)

const (
  kWifi = iota
  kEthernet
)

type Network struct {
	name string
	ltype int
}

type NetworkBlock struct {
	nticks, PollFreq int
	links []Network
}

func (b *NetworkBlock) Update() {

	b.nticks += 1
	b.nticks = b.nticks % b.PollFreq

	//Poll once every PollFreq ticks, and poll when starting, if PollFreq is 1,
	//just poll every tick
	if (b.PollFreq > 1) && (b.nticks != 1) {
		return
	}

	procstat_bytes, procstat_err := exec.Command("nmcli", "-t", "dev").Output()
	if procstat_err != nil {
		return
	}

	b.links = b.links[:0]

	scanner := bufio.NewScanner(bytes.NewBuffer(procstat_bytes))

	for scanner.Scan() {
		splits := strings.Split(scanner.Text(),":")
		if splits[2] != "connected" {
			continue
		}
		switch splits[1] {
	case "ethernet": {
		b.links = append(b.links, Network{splits[3], kEthernet})
		
	}
	case "wifi": {
		b.links = append(b.links, Network{splits[3], kWifi})		
	}
}


	}
}

func (b *NetworkBlock) ToBlock() Block {
	out_b := NewPangoBlock()

	if len(b.links) == 0 {
		return NewErrorBlock()
	}

	var strb []string

	for _,v := range b.links {
		if v.ltype == kWifi {
			strb = append(strb, fmt.Sprintf("\uf1eb %v",v.name))
		} else {
			strb = append(strb, "\uf796")
		}
	}

	out_b.full_text = pgu.NewPangoStrU( " " + strings.Join(strb, ", ") + " ").SetFGColor(pgu.Green).String()

	return out_b
}

func (b *NetworkBlock) Check() bool {
	_, err := exec.LookPath("nmcli")
	if err != nil {
		return false
	}

	if b.PollFreq == 0 {
		b.PollFreq = 1
	}

	return true
}