package blocks

import (
	"os/exec"
	"fmt"
	"strconv"
	"bufio"
	"bytes"
	"regexp"

	pgu "github.com/luketpickering/gobar/pangoutils"

)

type CPUTempBlock struct {
	temps []int
	core_regex *regexp.Regexp
	nticks, PollFreq int
}

func (b *CPUTempBlock) Update() {

	//Poll once every PollFreq ticks, and poll when starting, if PollFreq is 1,
	//just poll every tick
	if (b.PollFreq > 1) && (b.nticks != 1) {
		return
	}

	sensors_bytes, sensors_err := exec.Command("sensors").Output()
	if sensors_err != nil {
		return
	}

	scanner := bufio.NewScanner(bytes.NewBuffer(sensors_bytes))
	b.temps = b.temps[:0]

	for scanner.Scan() {
		matches := b.core_regex.FindStringSubmatch(scanner.Text())
		if len(matches) == 0 {
			continue
		}

		fl,_ := strconv.ParseFloat(matches[2],32)
		b.temps = append(b.temps, int(fl))
	}

}

func (b *CPUTempBlock) ToBlock() Block {
	out_b := NewPangoBlock()

	str := ""
	avg := 0
	for _,t := range b.temps {
		avg += t
		if t > 90 {
			str += pgu.NewPangoStrU(fmt.Sprintf("%v/",t)).SetFGColor(pgu.Red).String()
		} else {
			str += fmt.Sprintf("%v/",t)
		}
	}
	avg /= len(b.temps)
	str = str[:len(str)-1]

	ch := "\uf76b" //templow
	if avg > 70 {
		ch = pgu.NewPangoStrU("\uf769").SetFGColor(pgu.Orange).String() // temphigh
	} else if avg > 80 {
		ch = pgu.NewPangoStrU("\uf7e4").SetFGColor(pgu.Red).String() // fire
	}

	out_b.full_text = fmt.Sprintf(" %v %v°C ", ch, str)

	return out_b
}

func (b *CPUTempBlock) Check() bool {
	_, err := exec.LookPath("sensors")
	if err != nil {
		return false
	}

	if b.PollFreq == 0 {
		b.PollFreq = 1
	}


	b.core_regex = regexp.MustCompile(`Core ([0-9]+):\s+\+([0-9]+\.?[0-9]*)°C`)

	return true
}