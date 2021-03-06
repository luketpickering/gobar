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

type CPUUsage struct {
	User, System, Idle int
}

func (c CPUUsage) Usage() int {
	return c.User + c.System
}

func (c CPUUsage) Total() int {
	return c.Usage() + c.Idle
}

func Usage(now, prev []CPUUsage) int {

	sum_usage := 0

	for cpu_it, _ := range now {
		delta_usage := now[cpu_it].Usage() - prev[cpu_it].Usage()
		delta_total := now[cpu_it].Total() - prev[cpu_it].Total()
		sum_usage += (delta_usage * 100) / delta_total
	}

	return sum_usage
}

type CPUUsageBlock struct {
	nproc            int
	Usage            [2][]CPUUsage
	nticks, PollFreq int
	now_idx          int
	readerr          bool
}

func (b *CPUUsageBlock) Update() {

	b.nticks += 1
	b.nticks = b.nticks % b.PollFreq

	//Poll once every PollFreq ticks, and poll when starting, if PollFreq is 1,
	//just poll every tick
	if (b.PollFreq > 1) && (b.nticks != 1) {
		return
	}

	procstat_bytes, procstat_err := exec.Command("cat", "/proc/stat").Output()
	if procstat_err != nil {
		return
	}

	scanner := bufio.NewScanner(bytes.NewBuffer(procstat_bytes))

	//Skip the first line which is the average
	scanner.Scan()

	b.now_idx += 1
	b.now_idx = b.now_idx % 2

	// Read nproc lines
	for i := 0; i < b.nproc; i++ {
		scanner.Scan()
		splits := strings.Fields(scanner.Text())

		User, User_err := strconv.Atoi(splits[1])
		System, System_err := strconv.Atoi(splits[3])
		Idle, Idle_err := strconv.Atoi(splits[4])

		b.readerr = (User_err != nil) || (System_err != nil) || (Idle_err != nil)

		if b.readerr {
			return
		}

		b.Usage[b.now_idx][i] = CPUUsage{User, System, Idle}
	}
}

func (b *CPUUsageBlock) ToBlock() Block {
	out_b := NewPangoBlock()
	out_b.min_width = " \uf2db ---% "

	if b.readerr {
		return NewErrorBlock()
	}

	if (b.Usage[0][0].User == 0) || (b.Usage[0][0].User == 0) {
		out_b.full_text = " \uf2db ---% "
		return out_b
	}

	last_tick := b.now_idx
	prev_tick := (b.now_idx + 1) % 2

	CPUUsage := Usage(b.Usage[last_tick], b.Usage[prev_tick])


	if CPUUsage > ((b.nproc - 1)*100) {
		out_b.full_text = pgu.NewPangoStrU(fmt.Sprintf(" \uf2db %v%% ",CPUUsage)).SetBGColor(pgu.Red).SetFGColor(pgu.DarkGrey).String()
	} else if CPUUsage > (b.nproc*100/2) {
		out_b.full_text = pgu.NewPangoStrU(fmt.Sprintf(" \uf2db %v%% ",CPUUsage)).SetFGColor(pgu.Orange).String()
	} else {
		out_b.full_text = fmt.Sprintf(" \uf2db %v%% ",CPUUsage)
	}

	return out_b
}

func (b *CPUUsageBlock) Check() bool {
	loc, err := exec.LookPath("nproc")
	if err != nil {
		return false
	}

	nproc_cmd := exec.Command(loc)
	nproc_bytes, nproc_err := nproc_cmd.Output()
	if nproc_err != nil {
		return false
	}
	nproc, err := strconv.Atoi(pgu.Chompb(nproc_bytes))
	if err != nil {
		return false
	}

	if nproc == 0 {
		return false
	}

	b.nproc = nproc

	b.Usage[0] = make([]CPUUsage, nproc)
	b.Usage[1] = make([]CPUUsage, nproc)

	if b.PollFreq == 0 {
		b.PollFreq = 1
	}

	_, stat_err := os.Stat("/proc/stat")

	if stat_err != nil {
		return false
	}

	return true
}