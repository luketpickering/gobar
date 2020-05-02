package blocks

import (
	"os"
	"os/exec"
	"strconv"
	"io/ioutil"
	"fmt"

	pgu "github.com/luketpickering/gobar/pangoutils"
)

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
		b.status = pgu.Chompb(stat_bytes)
	} else {
		b.status = "Unknown"
	}

	cap_bytes, err := ioutil.ReadFile("/sys/class/power_supply/BAT0/capacity")

	if err == nil {
		b.cap_pc, err = strconv.Atoi(pgu.Chompb(cap_bytes))
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
						pgu.NewPangoStrU(fmt.Sprintf("%v%% remaining", b.cap_pc)).SetFGColor(pgu.Red).String()).Start()
				}
			}
		}()
	}
}

func (b *BatteryBlock) ToBlock() Block {
	out_b := NewPangoBlock()

	var cap_str string
	if b.cap_pc == -1 {
		cap_str = ""
	} else {
		// cap_str = strconv.Itoa(b.cap_pc) + "%"
		if b.cap_pc < 20 {
			cap_str = pgu.NewPangoStrU("\uf243").SetFGColor(pgu.Red).String()
		} else if b.cap_pc < 30 {
			cap_str = pgu.NewPangoStrU("\uf243").SetFGColor(pgu.Orange).String()
		} else if b.cap_pc < 60 {
			cap_str = pgu.NewPangoStrU("\uf242").SetFGColor(pgu.Orange).String()
		} else if b.cap_pc < 80 {
			cap_str = pgu.NewPangoStrU("\uf241").SetFGColor(pgu.Green).String()
		} else {
			cap_str = pgu.NewPangoStrU("\uf240").SetFGColor(pgu.Green).String()
		}
	}

	var status_str string

	if b.status == "Unknown" {
		status_str = "\uf059"
	} else if b.status == "Full" {
		cap_str = ""
		status_str = pgu.NewPangoStrU("\uf1e6").SetFGColor(pgu.Green).String()
	} else if b.status == "Charging" {
		status_str = pgu.NewPangoStrU("\uf1e6\uf0e7").SetFGColor(pgu.Green).String()
	} else if b.status == "Discharging" {
		status_str = pgu.NewPangoStrU("\uf0e7").SetFGColor(pgu.Orange).String()
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