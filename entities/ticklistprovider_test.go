package entities

import (
	"math/big"
	"testing"

	"github.com/nakaochi/espoon/arbitrage/uniswapsdk/utils/constant"
	"github.com/stretchr/testify/assert"
)

func Test_constructor(t *testing.T) {
	assert.PanicsWithValue( //throws for 0 tick spacing
		t,
		"TICK_SPACING_NONZERO",
		func() {
			ticks := TickListProvider{}
			ticks.Initialize(
				[]Tick{
					Tick{big.NewInt(-1), big.NewInt(1), big.NewInt(-1)},
					Tick{big.NewInt(1), big.NewInt(1), big.NewInt(2)}},
				0)
		})
}
func Test_GetTick(t *testing.T) {
	assert.PanicsWithValue( //throws if tick not in list
		t,
		"NOT_CONTAINED",
		func() {
			ticks := TickListProvider{}
			ticks.Initialize(
				[]Tick{
					Tick{big.NewInt(-1), big.NewInt(1), big.NewInt(-1)},
					Tick{big.NewInt(1), big.NewInt(1), big.NewInt(1)}},
				1)
			ticks.GetTick(constant.ZERO)
		})
	assert.Equal( //gets the smallest tick from the list
		t,
		func() bool {
			ticks := []Tick{
				Tick{big.NewInt(-1), big.NewInt(1), big.NewInt(-1)},
				Tick{big.NewInt(1), big.NewInt(1), big.NewInt(1)}}
			provider := TickListProvider{}
			provider.Initialize(ticks, 1)
			tick := provider.GetTick(constant.NEGATIVE_ONE)
			return tick.LiquidityNet.Cmp(constant.NEGATIVE_ONE) == 0 &&
				tick.LiquidityGross.Cmp(constant.ONE) == 0
		}(),
		true)
	assert.Equal( //gets the largest tick from the list
		t,
		func() bool {
			ticks := []Tick{
				Tick{big.NewInt(-1), big.NewInt(1), big.NewInt(-1)},
				Tick{big.NewInt(1), big.NewInt(1), big.NewInt(1)}}
			provider := TickListProvider{}
			provider.Initialize(ticks, 1)
			tick := provider.GetTick(constant.ONE)
			return tick.LiquidityNet.Cmp(constant.ONE) == 0 &&
				tick.LiquidityGross.Cmp(constant.ONE) == 0
		}(),
		true)
}
