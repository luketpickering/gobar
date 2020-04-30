package blocks

type Block struct {
	markup               string
	separator            bool
	separate_block_width int
	full_text            string
	border_px            [4]int
	min_width string

	err bool
}

func NewBlock() Block {
	return Block{"default", true, 50, "", [4]int{5, 5, 5, 5}, "", false}
}
func NewPangoBlock() Block {
	blk := DefaultBlock()
	blk.markup = "pango"
	return blk
}
func NewErrorBlock() Block {
	blk := DefaultBlock()
	blk.err = true
	return blk
}


func (b Block) JSONString() string {
	return fmt.Sprintf("{\"markup\": \"%v\", \"separator\": %v,\"separate_block_width\": \"%v\", \"border_top\": %v,\"border_right\": %v,\"border_bottom\": %v,\"border_left\": %v, \"min_width\": \"%v\", \"full_text\": \"%v\"}",
		b.markup, b.separator, b.separate_block_width, b.border_px[0], b.border_px[1], b.border_px[2], b.border_px[3], b.min_width, b.full_text)
}