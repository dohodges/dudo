package main

import (
	"fmt"
	"math"
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
				bids = append(bids, Bid{c,v})
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

func AtLeast(bid Bid, dice int, known []Die) float64 {
	var p float64
	for c := bid.Count; c <= dice; c++ {
		p += Exactly(Bid{c,bid.Value}, dice, known)
	}
	return p
}

func Exactly(bid Bid, dice int, known []Die) float64 {
	p := float64(2)/float64(6)
	if bid.Value == Ace {
		p = float64(1)/float64(6)
	}
	q := float64(1) - p

	return math.Pow(p, float64(bid.Count)) * math.Pow(q, float64(dice - bid.Count)) * float64(choose(dice, bid.Count))
}

func choose(n, k int) int {
	return int(factorial(n)/(factorial(k)*factorial(n-k)))
}

func factorial(n int) uint64 {
	if n > 0 {
		return uint64(n) * factorial(n-1)
	}
	return 1
}
