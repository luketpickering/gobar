package blocks

import (
	"os"
	"os/exec"
	"strings"
	"fmt"
	"strconv"

	pgu "github.com/luketpickering/gobar/pangoutils"

)

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

	df_pc, _ := strconv.Atoi(strings.TrimRight(splits[1],"%"))

	if df_pc > 90 {
		b.df_str = pgu.NewPangoStrU(fmt.Sprintf(" \uf0a0 %v%% ",df_pc)).SetBGColor(pgu.Red).SetFGColor(pgu.DarkGrey).String()		
	} else if df_pc > 80 {
		b.df_str = pgu.NewPangoStrU(fmt.Sprintf(" \uf0a0 %v%% ",df_pc)).SetFGColor(pgu.Orange).String()
	} else {
		b.df_str = fmt.Sprintf(" \uf0a0 %v%% ",df_pc)
	
	}
}

func (b *DiskFreeBlock) ToBlock() Block {
	out_b := NewPangoBlock()
	out_b.full_text = b.df_str
	out_b.min_width = " \uf0a0 100% "
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