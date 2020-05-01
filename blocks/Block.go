package blocks

import "fmt"

type Block struct {
	markup               string
	separator            bool
	separate_block_width int
	full_text            string
	min_width string

	err bool
}

func NewBlock() Block {
	return Block{"default", true, 50, "", "", false}
}

func NewPangoBlock() Block {
	blk := NewBlock()
	blk.markup = "pango"
	return blk
}

func NewErrorBlock() Block {
	blk := NewBlock()
	blk.err = true
	return blk
}


func (b Block) JSONString() string {
	return fmt.Sprintf("{\"markup\": \"%v\", \"separator\": %v,\"separate_block_width\": \"%v\", \"min_width\": \"%v\", \"align\": \"center\", \"full_text\": \"%v\"}",
		b.markup, b.separator, b.separate_block_width, b.min_width, b.full_text)
}

func (b Block) PrintBlockLine() string {
	if b.err {
		return ""
	}

	return fmt.Sprintf("\t%v,\n", b.JSONString())
}
