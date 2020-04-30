package blocks

type Blocklike interface {
	ToBlock() Block
	Update()
	Check() bool
}

type Blocklist []Blocklike

func AppendBlocklike(blocks Blocklist, b Blocklike) (Blocklist, error) {
	if b.Check() {
		blocks = append(blocks, b)
		return blocks, nil
	}
	return blocks, errors.New("Bad block")
}

func Poll(b Blocklike, ch chan Block) {
	b.Update()
	ch <- b.ToBlock()
}

func PollAll(bs Blocklist, chs []chan Block) {
	for i, b := bs {
		b.Update()
		chs[i] <- b.ToBlock()
	}
}