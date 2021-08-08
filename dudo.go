package main

import (
	"fmt"
	"math"
	"math/big"
	"sort"
)

type Die uint8

const (
	Ace Die = iota + 1
	Two
	Three
	Four
	Five
	Six
)

func (d Die) String() string {
	switch d {
	case Ace:
		return "Ace"
	case Two:
		return "Two"
	case Three:
		return "Three"
	case Four:
		return "Four"
	case Five:
		return "Five"
	case Six:
		return "Six"
	}
	return ""
}

type Bid struct {
	Count int
	Value Die
}

func (b Bid) Less(c Bid) bool {
	return b.strength() < c.strength()
}

func (b Bid) strength() int {
	s := b.Count * 10
	if b.Value == Ace {
		s *= 2
	}
	s += int(b.Value)
	return s
}

func (b Bid) String() string {
	pluralize := ""
	if b.Count > 1 {
		if b.Value == Six {
			pluralize = "es"
		} else {
			pluralize = "s"
		}
	}
	return fmt.Sprintf("%d %s%s", b.Count, b.Value, pluralize)
}

func EachBid(dice int) <-chan Bid {
	ch := make(chan Bid)
	go func() {
		bids := make([]Bid, 0, 6*dice)
		for c := 1; c <= dice; c++ {
			for v := Ace; v <= Six; v++ {
				bids = append(bids, Bid{c, v})
			}
		}
		sort.Slice(bids, func(i, j int) bool {
			return bids[i].Less(bids[j])
		})
		for _, b := range bids {
			ch <- b
		}
		close(ch)
	}()
	return ch
}

func AtLeast(bid Bid, unknownDice int, known map[Die]int) float64 {
	knownDice := 0
	for _, c := range known {
		knownDice += c
	}

	var p float64
	for c := bid.Count; c <= unknownDice+knownDice; c++ {
		p += Exactly(Bid{c, bid.Value}, unknownDice, known)
		if p >= 1.0 {
			return 1.0
		}
	}
	return p
}

func Exactly(bid Bid, unknownDice int, known map[Die]int) float64 {
	need := bid
	need.Count -= known[need.Value]
	if need.Value != Ace {
		need.Count -= known[Ace]
	}

	if need.Count <= 0 {
		return 1.0
	}

	p := float64(2) / float64(6)
	if need.Value == Ace {
		p = float64(1) / float64(6)
	}
	q := float64(1) - p

	return math.Pow(p, float64(need.Count)) * math.Pow(q, float64(unknownDice-need.Count)) * float64(choose(unknownDice, need.Count).Uint64())
}

func choose(n, k int) *big.Int {
	return new(big.Int).Div(factorial(n, n-k+1), factorial(k, 1))
}

func factorial(n, k int) *big.Int {
	if n > k {
		return new(big.Int).Mul(big.NewInt(int64(n)), factorial(n-1, k))
	}
	return big.NewInt(int64(k))
}
