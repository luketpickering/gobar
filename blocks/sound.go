package blocks

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
	"bytes"
	"regexp"
	"strconv"

	pgu "github.com/luketpickering/gobar/pangoutils"
)



type SoundBlock struct {
	nticks, PollFreq int

	micjack_id, vol_id, cap_id string
	vol_min, vol_max, vol_pc int

	micjack_on bool
}

func (b *SoundBlock) Update() {

	b.nticks += 1
	b.nticks = b.nticks % b.PollFreq

	//Poll once every PollFreq ticks, and poll when starting, if PollFreq is 1,
	//just poll every tick
	if (b.PollFreq > 1) && (b.nticks != 1) {
		return
	}

	//Check the mic jack
	amixjack_bytes, amixjack_err := exec.Command("amixer", "-c", "0", "cget", b.micjack_id).Output()
	if amixjack_err != nil {
		return
	}

	scanner := bufio.NewScanner(bytes.NewBuffer(amixjack_bytes))
	var jack_regex = regexp.MustCompile(`values=on`)
	b.micjack_on = false

	for scanner.Scan() {
		ln := strings.TrimLeft(scanner.Text(), " ")
		if ln[0] != ':' {
			continue
		}
		if jack_regex.MatchString(ln) {
			b.micjack_on = true
		break		}
		
	}



	amixvol_bytes, amixvol_err := exec.Command("amixer", "-c", "0", "cget", b.vol_id).Output()
	if amixvol_err != nil {
		return
	}

	scanner = bufio.NewScanner(bytes.NewBuffer(amixvol_bytes))

	var vol_regex = regexp.MustCompile(`values=([0-9]+)`)

	for scanner.Scan() {
		ln := strings.TrimLeft(scanner.Text(), " ")
		if ln[0] != ':' {
			continue
		}
		matches := vol_regex.FindStringSubmatch(ln)
		if len(matches) == 0 {
			continue
		}
		vol_raw, _ := strconv.Atoi(matches[1])
		b.vol_pc = ((vol_raw - b.vol_min)*100)/(b.vol_max - b.vol_min)
		
	}
}

func (b *SoundBlock) ToBlock() Block {
	out_b := NewPangoBlock()

	vol_str := fmt.Sprintf("%v%%", b.vol_pc)
	if b.vol_pc > 50 {
		vol_str = "\uf028 " + vol_str
	} else if b.vol_pc > 15 {
		vol_str = "\uf027 " + vol_str
	} else if b.vol_pc > 0 {
		vol_str = "\uf026 " + vol_str
	} else {
		vol_str = "\uf6a9 " + vol_str
	}

	if b.micjack_on {
		vol_str += ", \uf025"
	}

	out_b.full_text = pgu.NewPangoStrU(" " + vol_str + " ").String()
	return out_b
}

func (b *SoundBlock) Check() bool {
	_, err := exec.LookPath("amixer")
	if err != nil {
		return false
	}

	if b.PollFreq == 0 {
		b.PollFreq = 1
	}

	amixcontrols_bytes, amixcontrols_err := exec.Command("amixer", "-c", "0", "controls").Output()
	if amixcontrols_err != nil {
		return false
	}

	scanner := bufio.NewScanner(bytes.NewBuffer(amixcontrols_bytes))

	var numid_regex = regexp.MustCompile(`numid=([0-9]+)`)

	for scanner.Scan() {
	
		s := scanner.Text()

		if strings.Contains(s,"Headphone Mic Jack"){
			id_s := numid_regex.FindStringSubmatch(s)
			b.micjack_id = "numid=" + id_s[1]
		} else if strings.Contains(s,"Master Playback Volume"){
			id_s := numid_regex.FindStringSubmatch(s)
			b.vol_id = "numid=" + id_s[1]
		} else if strings.Contains(s,"Capture Source"){
			id_s := numid_regex.FindStringSubmatch(s)
			b.cap_id = "numid=" + id_s[1]
		}
	}


	amixvol_bytes, amixvol_err := exec.Command("amixer", "-c", "0", "cget", b.vol_id).Output()
	if amixvol_err != nil {
		return false
	}

	scanner = bufio.NewScanner(bytes.NewBuffer(amixvol_bytes))

	var volrange_regex = regexp.MustCompile(`min=([0-9]+),max=([0-9]+)`)

	for scanner.Scan() {
		matches := volrange_regex.FindStringSubmatch(scanner.Text())
		if len(matches) == 0 {
			continue
		}
		b.vol_min, _ = strconv.Atoi(matches[1])
		b.vol_max, _ = strconv.Atoi(matches[2])
	}

	return true
}