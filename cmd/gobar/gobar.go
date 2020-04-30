package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/luketpickering/gobar/internal"
)

func chomp(s string) string {
	return strings.TrimRight(s, "\n\t ")
}

func chompb(s []byte) string {
	return chomp(string(s))
}

type Block struct {
	markup               string
	separator            bool
	separate_block_width int
	full_text            string
	border_px            [4]int
	min_width string

	err bool
}

type Blocklike interface {
	ToBlock() Block
	Update()
	Check() bool
}

func Poll(b Blocklike, ch chan Block) {
	b.Update()
	ch <- b.ToBlock()
}

func AddBlock(blocks []Blocklike, b Blocklike) ([]Blocklike, error) {
	if b.Check() {
		blocks = append(blocks, b)
		return blocks, nil
	}
	return blocks, errors.New("Bad block")
}

func (b Block) JSONString() string {
	return fmt.Sprintf("{\"markup\": \"%v\", \"separator\": %v,\"separate_block_width\": \"%v\", \"border_top\": %v,\"border_right\": %v,\"border_bottom\": %v,\"border_left\": %v, \"min_width\": \"%v\", \"full_text\": \"%v\"}",
		b.markup, b.separator, b.separate_block_width, b.border_px[0], b.border_px[1], b.border_px[2], b.border_px[3], b.min_width, b.full_text)
}

func DefaultBlock() Block {
	return Block{"default", true, 50, "", [4]int{5, 5, 5, 5}, "", false}
}
func PangoBlock() Block {
	blk := DefaultBlock()
	blk.markup = "pango"
	return blk
}
func ErrorBlock() Block {
	blk := DefaultBlock()
	blk.err = true
	return blk
}

type BatteryBlock struct {
	status string
	cap_pc int

	tick             int
	have_notify_send bool
}

func (b *BatteryBlock) Update() {
	b.tick += 1

	b.tick = b.tick % (5 * 60) // Every five minutes

	stat_bytes, err := ioutil.ReadFile("/sys/class/power_supply/BAT0/status")

	if err == nil {
		b.status = chompb(stat_bytes)
	} else {
		b.status = "Unknown"
	}

	cap_bytes, err := ioutil.ReadFile("/sys/class/power_supply/BAT0/capacity")

	if err == nil {
		b.cap_pc, err = strconv.Atoi(chompb(cap_bytes))
	}

	if err != nil {
		b.cap_pc = -1
	}

	// fire off a check of battery %/charging state for notification
	if b.have_notify_send {
		go func() {

			if (b.cap_pc > 0) && (b.cap_pc < 20) && (b.status != "Charging") {

				//Want frequency to depend on %
				send := false
				if b.cap_pc < 5 && (b.tick%30 == 0) {
					send = true
				} else if b.cap_pc < 10 && (b.tick%60 == 0) {
					send = true
				} else {
					send = (b.tick == 0)
				}

				if send {
					exec.Command("notify-send", "Low Battery",
						internal.RedText(fmt.Sprintf("%v%% remaining", b.cap_pc))).Start()
				}
			}
		}()
	}
}

func (b *BatteryBlock) ToBlock() Block {
	out_b := PangoBlock()

	var cap_str, status_str string

	if b.cap_pc == -1 {
		cap_str = "---%"
	} else {
		cap_str = strconv.Itoa(b.cap_pc) + "%"
		if b.cap_pc < 20 {
			cap_str = internal.RedText(cap_str + " \uf243")
		} else if b.cap_pc < 80 {
			cap_str = internal.OrangeText(cap_str + " \uf241")
		} else if b.cap_pc < 60 {
			cap_str = internal.OrangeText(cap_str + " \uf242")
		} else if b.cap_pc < 30 {
			cap_str = internal.OrangeText(cap_str + " \uf243")
		} else {
			cap_str = internal.GreenText(cap_str + " \uf240")
		}
	}

	if b.status == "Unknown" {
		status_str = internal.WhiteText("\uf059")
	} else if b.status == "Full" {
		status_str = internal.GreenText("\uf1e6")
	} else if b.status == "Charging" {
		status_str = internal.OrangeText("\uf1e6\uf0e7")
	} else if b.status == "Discharging" {
		status_str = internal.OrangeText("(Discharging)")
	}

	out_b.full_text = fmt.Sprintf("%v %v", cap_str, status_str)

	return out_b
}

func (b *BatteryBlock) Check() bool {
	_, stat_err := os.Stat("/sys/class/power_supply/BAT0/status")
	_, cap_err := os.Stat("/sys/class/power_supply/BAT0/capacity")

	_, notify_send_err := exec.LookPath("notify-send")
	if notify_send_err != nil {
		b.have_notify_send = false
	}

	b.have_notify_send = true

	return !(os.IsNotExist(stat_err) || os.IsNotExist(cap_err))
}

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

	procstat_cmd := exec.Command("cat", "/proc/stat")
	procstat_pipe, procstat_err := procstat_cmd.StdoutPipe()

	if procstat_err != nil {
		return
	}
	scanner := bufio.NewScanner(procstat_pipe)
	procstat_cmd.Start()

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

	//I don't care about the rest of the output
	go procstat_cmd.Wait()

}
func (b *CPUUsageBlock) ToBlock() Block {
	out_b := DefaultBlock()
	out_b.min_width = "\uf2db 100%%"

	if b.readerr {
		return ErrorBlock()
	}

	if (b.Usage[0][0].User == 0) || (b.Usage[0][0].User == 0) {
		out_b.full_text = fmt.Sprintf("\uf2db -%%")
		return out_b
	}

	last_tick := b.now_idx
	prev_tick := (b.now_idx + 1) % 2

	CPUUsage := Usage(b.Usage[last_tick], b.Usage[prev_tick])

	out_b.full_text = fmt.Sprintf("\uf2db %v%%", CPUUsage)
	out_b.min_width = "\uf2db 100%%"

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
	nproc, err := strconv.Atoi(chompb(nproc_bytes))
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

type TimeBlock struct {
	TZ   string
	time string
}

func (b *TimeBlock) Update() {

	date_cmd := exec.Command("date", "+%H:%M (%Z)")

	if len(b.TZ) != 0 {
		date_cmd.Env = append(os.Environ(),
			fmt.Sprintf("TZ=%v", b.TZ),
		)
	}

	date_out, date_err := date_cmd.Output()
	if date_err != nil {
		return
	}

	if len(b.TZ) != 0 {
		b.time = chompb(date_out)
	} else {

		day_out, day_err := exec.Command("date", "+%a %b %_d").Output()
		if day_err != nil {
			return
		}

		b.time = internal.BoldText(internal.OrangeBackground(internal.DarkGreyText(" \uf073 " + chompb(day_out) + "  \uf017 " + chompb(date_out) + " ")))
	}
}

func (b *TimeBlock) ToBlock() Block {
	out_b := PangoBlock()
	out_b.full_text = b.time
	return out_b
}

func (b *TimeBlock) Check() bool {
	_, err := exec.LookPath("date")
	if err != nil {
		return false
	}
	return true
}

type DiskFreeBlock struct {
	homedir string
	df_str  string
}

func (b *DiskFreeBlock) Update() {

	df_out, df_err := exec.Command("df", "--output=pcent", b.homedir).Output()
	if df_err != nil {
		return
	}

	splits := strings.Fields(string(df_out))

	if len(splits) < 2 {
		return
	}

	b.df_str = "\uf0a0 " + splits[1]
}

func (b *DiskFreeBlock) ToBlock() Block {
	out_b := PangoBlock()
	out_b.full_text = b.df_str
	return out_b
}

func (b *DiskFreeBlock) Check() bool {
	_, err := exec.LookPath("df")
	if err != nil {
		return false
	}

	b.homedir = os.Getenv("HOME")

	return true
}

func PrintBlockLine(b Block) {
	if b.err {
		return
	}

	fmt.Printf("\t%v,\n", b.JSONString())
}

func main() {

	tick := time.NewTicker(time.Second)

	blocks := []Blocklike{}

	//We will just ignore blocks that fail the check
	blocks, _ = AddBlock(blocks, &CPUUsageBlock{})
	//Set the polling frequency
	blocks[0].(*CPUUsageBlock).PollFreq = 2

	blocks, _ = AddBlock(blocks, &DiskFreeBlock{})
	blocks, _ = AddBlock(blocks, &TimeBlock{"", ""})
	blocks, _ = AddBlock(blocks, &TimeBlock{"Europe/London", ""})
	blocks, _ = AddBlock(blocks, &TimeBlock{"Asia/Tokyo", ""})
	blocks, _ = AddBlock(blocks, &BatteryBlock{})

	var block_chans = []chan Block{}
	for i := 0; i < len(blocks); i++ {
		block_chans = append(block_chans, make(chan Block))
	}

	fmt.Println("{\"version\":1}")
	fmt.Println("[")
	fmt.Println("\t[],")

	for range tick.C {
		fmt.Println("[")

		for i, bl := range blocks {
			go Poll(bl, block_chans[i])
		}

		for i := 0; i < len(blocks); i++ {
			PrintBlockLine(<-block_chans[i])
		}

		fmt.Println("],")
	}

}
