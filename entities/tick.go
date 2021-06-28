package entities

import (
	"math/big"

	"github.com/nakaochi/espoon/arbitrage/uniswapsdk/utils/constant"
	"github.com/nakaochi/espoon/arbitrage/uniswapsdk/utils/tickmath"
)

type Tick struct {
	Index          *big.Int
	LiquidityGross *big.Int
	LiquidityNet   *big.Int
}

func NewTick(index *big.Int, liquidityGross *big.Int, liquidityNet *big.Int) *Tick {
	tickmath.Check(index.Cmp(constant.MIN_TICK) >= 0 && index.Cmp(constant.MAX_TICK) <= 0, "TICK")
	tick := &Tick{}
	tick.Index = index
	tick.LiquidityGross = liquidityGross
	tick.LiquidityNet = liquidityNet
	return tick
}
