package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)


func main() {
	counts := make([]int, 6)
	for i := 1; i <= 6; i++ {
		flag.IntVar(&counts[i-1], strconv.Itoa(i), 0, fmt.Sprintf(`Number of known %d dice`, i))
	}

	flag.Parse()

	known := make(map[Die]int)
	knownDice := 0
	for i, v := range counts {
		known[Die(i+1)] = v
		knownDice += v
	}

	var unknownDice int
	if flag.NArg() != 1 {
		fmt.Fprintf(os.Stderr, "specify dice count\n")
		os.Exit(1)
	} else if dice, err := strconv.ParseUint(flag.Arg(0), 10, 64); err != nil {
		fmt.Fprintf(os.Stderr, "invalid dice count(%v)\n", flag.Arg(0))
		os.Exit(1)
	} else {
		unknownDice = int(dice)
	}

	for bid := range EachBid(unknownDice+knownDice) {
		odds := AtLeast(bid, unknownDice, known)
		fmt.Printf("%.5f - %s\n", odds, bid)
	}
}
