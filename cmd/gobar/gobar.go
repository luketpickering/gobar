package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/luketpickering/gobar/pangoutils"
)

func chomp(s string) string {
	return strings.TrimRight(s, "\n\t ")
}

func chompb(s []byte) string {
	return chomp(string(s))
}


func PrintBlockLine(b Block) {
	if b.err {
		return
	}

	fmt.Printf("\t%v,\n", b.JSONString())
}

func main() {

	tick := time.NewTicker(time.Second)

	blocks := []Blocklike{}

	//We will just ignore blocks that fail the check
	blocks, _ = AddBlock(blocks, &CPUUsageBlock{})
	//Set the polling frequency
	blocks[0].(*CPUUsageBlock).PollFreq = 2

	blocks, _ = AddBlock(blocks, &DiskFreeBlock{})
	blocks, _ = AddBlock(blocks, &TimeBlock{"", ""})
	blocks, _ = AddBlock(blocks, &TimeBlock{"Europe/London", ""})
	blocks, _ = AddBlock(blocks, &TimeBlock{"Asia/Tokyo", ""})
	blocks, _ = AddBlock(blocks, &BatteryBlock{})

	var block_chans = []chan Block{}
	for i := 0; i < len(blocks); i++ {
		block_chans = append(block_chans, make(chan Block))
	}

	fmt.Println("{\"version\":1}")
	fmt.Println("[")
	fmt.Println("\t[],")

	for range tick.C {
		fmt.Println("[")

		for i, bl := range blocks {
			go Poll(bl, block_chans[i])
		}

		for i := 0; i < len(blocks); i++ {
			PrintBlockLine(<-block_chans[i])
		}

		fmt.Println("],")
	}

}
