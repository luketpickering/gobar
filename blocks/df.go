package blocks

import (
	"os"
	"os/exec"
	"strings"
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

	b.df_str = "\uf0a0 " + splits[1]
}

func (b *DiskFreeBlock) ToBlock() Block {
	out_b := NewPangoBlock()
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