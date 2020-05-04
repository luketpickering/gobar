package main

import (
	"fmt"
	"strings"
	"time"

	. "github.com/luketpickering/gobar/blocks"
)

func main() {

	tick := time.NewTicker(time.Second)

	blocks := []Blocklike{}

	//We will just ignore blocks that fail the check
	blocks, _ = AppendBlocklike(blocks, &SoundBlock{})
	blocks, _ = AppendBlocklike(blocks, &CPUTempBlock{})
	blocks, _ = AppendBlocklike(blocks, &CPUUsageBlock{})
	//Set the polling frequency
	blocks[len(blocks)-1].(*CPUUsageBlock).PollFreq = 2

	blocks, _ = AppendBlocklike(blocks, &MemAvailBlock{})

	blocks, _ = AppendBlocklike(blocks, &DiskFreeBlock{})
	blocks, _ = AppendBlocklike(blocks, &NetworkBlock{})
	
	blocks, _ = AppendBlocklike(blocks, &TimeBlock{"", ""})
	blocks, _ = AppendBlocklike(blocks, &TimeBlock{"Europe/London", ""})
	blocks, _ = AppendBlocklike(blocks, &TimeBlock{"Asia/Tokyo", ""})
	
	blocks, _ = AppendBlocklike(blocks, &BatteryBlock{})

	var block_chans = []chan Block{}
	for i := 0; i < len(blocks); i++ {
		block_chans = append(block_chans, make(chan Block))
	}

	fmt.Println("{\"version\":1}")
	fmt.Println("[")
	fmt.Println("\t[],")

	for range tick.C {
		strbuildr := []string{"["}

		PollAll(blocks, block_chans)

		for i := 0; i < len(blocks); i++ {
			strbuildr = append(strbuildr,(<-block_chans[i]).PrintBlockLine())
		}

		strbuildr = append(strbuildr,"],")

		fmt.Println(strings.Join(strbuildr," "))
	}

}
