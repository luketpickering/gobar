package blocks

import (
	"os"
	"os/exec"
	"fmt"

	pgu "github.com/luketpickering/gobar/pangoutils"
)

type TimeBlock struct {
	TZ   string
	Time string
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
		b.Time = pgu.Chompb(date_out)
	} else {

		day_out, day_err := exec.Command("date", "+%a %b %_d").Output()
		if day_err != nil {
			return
		}

		b.Time = pgu.NewPangoStrU(" \uf073 " + pgu.Chompb(day_out) + "  \uf017 " + pgu.Chompb(date_out) + " ").SetFGColor(pgu.DarkGrey).SetBGColor(pgu.Blue).SetFontWeight(pgu.Ultrabold).String()
	}
}

func (b *TimeBlock) ToBlock() Block {
	out_b := NewPangoBlock()
	out_b.full_text = b.Time
	return out_b
}

func (b *TimeBlock) Check() bool {
	_, err := exec.LookPath("date")
	if err != nil {
		return false
	}
	return true
}