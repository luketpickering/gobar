package blocks

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