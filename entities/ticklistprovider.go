package entities

import (
	"math/big"
)

type TickDataProvider interface {
	GetTick(tick *big.Int) *Tick
	NextInitializedTickWithinOneWord(tick *big.Int, lte bool, tickSpacing int) (*big.Int, bool)
}

type TickListProvider struct {
	ticks []Tick
}

func (P *TickListProvider) Initialize(ticks []Tick, tickSpacing int) {
	ValidateList(ticks, tickSpacing)
	P.ticks = ticks[:]
}

func (P TickListProvider) GetTick(tick *big.Int) *Tick {
	return GetTick(P.ticks, tick)
}
func (P TickListProvider) NextInitializedTickWithinOneWord(tick *big.Int, lte bool, tickSpacing int) (*big.Int, bool) {
	return NextInitializedTickWithinOneWord(P.ticks, tick, lte, tickSpacing)
}
