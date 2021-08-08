package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)


func main() {
	flag.Parse()

	if flag.NArg() != 1 {
		fmt.Fprintf(os.Stderr, "specify dice count\n")
		os.Exit(1)
	} else if dice, err := strconv.ParseUint(flag.Arg(0), 10, 64); err != nil {
		fmt.Fprintf(os.Stderr, "invalid dice count(%v)\n", flag.Arg(0))
		os.Exit(1)
	} else {
		for bid := range EachBid(int(dice)) {
			odds := AtLeast(bid, int(dice), nil)
			fmt.Printf("%.5f - %s\n", odds, bid)
		}
	}
}
