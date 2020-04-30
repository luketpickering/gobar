package blocks

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