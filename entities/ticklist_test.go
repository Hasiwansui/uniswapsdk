package entities

import (
	"math/big"
	"testing"

	"github.com/nakaochi/espoon/arbitrage/uniswapsdk/utils/constant"
	"github.com/stretchr/testify/assert"
)

var lowTick = NewTick(
	new(big.Int).Add(constant.MIN_TICK, constant.ONE),
	big.NewInt(10),
	big.NewInt(10))
var midTick = NewTick(
	big.NewInt(0),
	big.NewInt(5),
	big.NewInt(-5))
var highTick = NewTick(
	new(big.Int).Sub(constant.MAX_TICK, constant.ONE),
	big.NewInt(5),
	big.NewInt(-5))

func Test_Validate(t *testing.T) {
	assert.PanicsWithValue(
		t,
		"LEN",
		func() {
			ValidateList([]Tick{*lowTick}, 1)
		})
	assert.PanicsWithValue(
		t,
		"SORTED",
		func() {
			ValidateList([]Tick{*highTick, *lowTick, *midTick}, 1)
		})
	assert.PanicsWithValue(
		t,
		"TICK_SPACING",
		func() {
			ValidateList([]Tick{*highTick, *lowTick, *midTick}, 1337)
		})
}
func Test_Below(t *testing.T) {
	ticks := []Tick{*lowTick, *midTick, *highTick}
	assert.Equal(t, isBelowSmallest(ticks, constant.MIN_TICK), true)
	assert.Equal(t, isBelowSmallest(ticks, new(big.Int).Add(constant.MIN_TICK, constant.ONE)), false)
}
func Test_Above(t *testing.T) {
	result := []Tick{*lowTick, *midTick, *highTick}
	assert.Equal(t, isAtOrAboveLargest(result, new(big.Int).Sub(constant.MAX_TICK, constant.TWO)), false)
	assert.Equal(t, isAtOrAboveLargest(result, new(big.Int).Sub(constant.MAX_TICK, constant.ONE)), true)
}

func Test_NextInitializedTick(t *testing.T) {
	ticks := []Tick{*lowTick, *midTick, *highTick}
	//LOW-LTE=TRUE
	assert.PanicsWithValue(
		t,
		"BELOW_SMALLEST",
		func() {
			NextInitializedTick(ticks, constant.MIN_TICK, true)
		})
	assert.Equal(
		t,
		NextInitializedTick(
			ticks,
			new(big.Int).Add(constant.MIN_TICK, constant.ONE),
			true).Index.Cmp(lowTick.Index) == 0,
		true)
	assert.Equal(
		t,
		NextInitializedTick(
			ticks,
			new(big.Int).Add(constant.MIN_TICK, constant.TWO),
			true).Index.Cmp(lowTick.Index) == 0,
		true)
	//LOW-LTE=FALSE
	assert.Equal(
		t,
		NextInitializedTick(
			ticks,
			constant.MIN_TICK,
			false).Index.Cmp(lowTick.Index) == 0,
		true)
	assert.Equal(
		t,
		NextInitializedTick(
			ticks,
			new(big.Int).Add(constant.MIN_TICK, constant.ONE),
			false).Index.Cmp(midTick.Index) == 0,
		true)
	//MID-LTE=TRUE
	assert.Equal(
		t,
		NextInitializedTick(
			ticks,
			constant.ZERO,
			true).Index.Cmp(midTick.Index) == 0,
		true)
	assert.Equal(
		t,
		NextInitializedTick(
			ticks,
			constant.ONE,
			true).Index.Cmp(midTick.Index) == 0,
		true)
	//MID-LTE=FALSE
	assert.Equal(
		t,
		NextInitializedTick(
			ticks,
			constant.NEGATIVE_ONE,
			false).Index.Cmp(midTick.Index) == 0,
		true)
	assert.Equal(
		t,
		NextInitializedTick(
			ticks,
			constant.ONE,
			false).Index.Cmp(highTick.Index) == 0,
		true)
	//HIGH-LTE=TRUE
	assert.Equal(
		t,
		NextInitializedTick(
			ticks,
			new(big.Int).Sub(constant.MAX_TICK, constant.ONE),
			true).Index.Cmp(highTick.Index) == 0,
		true)
	assert.Equal(
		t,
		NextInitializedTick(
			ticks,
			constant.MAX_TICK,
			true).Index.Cmp(highTick.Index) == 0,
		true)
	//HIGH-LTE=FALSE
	assert.PanicsWithValue(
		t,
		"AT_OR_ABOVE_LARGEST",
		func() {
			NextInitializedTick(
				ticks,
				new(big.Int).Sub(constant.MAX_TICK, constant.ONE),
				false)
		})
	assert.Equal(
		t,
		NextInitializedTick(
			ticks,
			new(big.Int).Sub(constant.MAX_TICK, constant.TWO),
			false).Index.Cmp(highTick.Index) == 0,
		true)
	assert.Equal(
		t,
		NextInitializedTick(
			ticks,
			new(big.Int).Sub(constant.MAX_TICK, big.NewInt(3)),
			false).Index.Cmp(highTick.Index) == 0,
		true)
}

func Test_NextInitializedTickWithinOneWord(t *testing.T) {
	ticks := []Tick{*lowTick, *midTick, *highTick}
	//words around 0 ,lte=true
	assert.Equal(
		t,
		func() bool {
			tick, bl := NextInitializedTickWithinOneWord(
				ticks,
				big.NewInt(-257),
				true,
				1)
			return tick.Cmp(big.NewInt(-512)) == 0 && bl == false
		}(),
		true)
	assert.Equal(
		t,
		func() bool {
			tick, bl := NextInitializedTickWithinOneWord(
				ticks,
				big.NewInt(-256),
				true,
				1)
			return tick.Cmp(big.NewInt(-256)) == 0 && bl == false
		}(),
		true)
	assert.Equal(
		t,
		func() bool {
			tick, bl := NextInitializedTickWithinOneWord(
				ticks,
				big.NewInt(-1),
				true,
				1)
			return tick.Cmp(big.NewInt(-256)) == 0 && bl == false
		}(),
		true)
	assert.Equal(
		t,
		func() bool {
			tick, bl := NextInitializedTickWithinOneWord(
				ticks,
				big.NewInt(0),
				true,
				1)
			return tick.Cmp(big.NewInt(0)) == 0 && bl == true
		}(),
		true)
	assert.Equal(
		t,
		func() bool {
			tick, bl := NextInitializedTickWithinOneWord(
				ticks,
				big.NewInt(1),
				true,
				1)
			return tick.Cmp(big.NewInt(0)) == 0 && bl == true
		}(),
		true)
	assert.Equal(
		t,
		func() bool {
			tick, bl := NextInitializedTickWithinOneWord(
				ticks,
				big.NewInt(255),
				true,
				1)
			return tick.Cmp(big.NewInt(0)) == 0 && bl == true
		}(),
		true)
	assert.Equal(
		t,
		func() bool {
			tick, bl := NextInitializedTickWithinOneWord(
				ticks,
				big.NewInt(256),
				true,
				1)
			return tick.Cmp(big.NewInt(256)) == 0 && bl == false
		}(),
		true)
	assert.Equal(
		t,
		func() bool {
			tick, bl := NextInitializedTickWithinOneWord(
				ticks,
				big.NewInt(257),
				true,
				1)
			return tick.Cmp(big.NewInt(256)) == 0 && bl == false
		}(),
		true)
	//words around 0 ,lte=false
	assert.Equal(
		t,
		func() bool {
			tick, bl := NextInitializedTickWithinOneWord(
				ticks,
				big.NewInt(-258),
				false,
				1)
			return tick.Cmp(big.NewInt(-257)) == 0 && bl == false
		}(),
		true)
	assert.Equal(
		t,
		func() bool {
			tick, bl := NextInitializedTickWithinOneWord(
				ticks,
				big.NewInt(-257),
				false,
				1)
			return tick.Cmp(big.NewInt(-1)) == 0 && bl == false
		}(),
		true)
	assert.Equal(
		t,
		func() bool {
			tick, bl := NextInitializedTickWithinOneWord(
				ticks,
				big.NewInt(-256),
				false,
				1)
			return tick.Cmp(big.NewInt(-1)) == 0 && bl == false
		}(),
		true)
	assert.Equal(
		t,
		func() bool {
			tick, bl := NextInitializedTickWithinOneWord(
				ticks,
				big.NewInt(-2),
				false,
				1)
			return tick.Cmp(big.NewInt(-1)) == 0 && bl == false
		}(),
		true)
	assert.Equal(
		t,
		func() bool {
			tick, bl := NextInitializedTickWithinOneWord(
				ticks,
				big.NewInt(-1),
				false,
				1)
			return tick.Cmp(big.NewInt(0)) == 0 && bl == true
		}(),
		true)
	assert.Equal(
		t,
		func() bool {
			tick, bl := NextInitializedTickWithinOneWord(
				ticks,
				big.NewInt(0),
				false,
				1)
			return tick.Cmp(big.NewInt(255)) == 0 && bl == false
		}(),
		true)
	assert.Equal(
		t,
		func() bool {
			tick, bl := NextInitializedTickWithinOneWord(
				ticks,
				big.NewInt(1),
				false,
				1)
			return tick.Cmp(big.NewInt(255)) == 0 && bl == false
		}(),
		true)
	assert.Equal(
		t,
		func() bool {
			tick, bl := NextInitializedTickWithinOneWord(
				ticks,
				big.NewInt(254),
				false,
				1)
			return tick.Cmp(big.NewInt(255)) == 0 && bl == false
		}(),
		true)
	assert.Equal(
		t,
		func() bool {
			tick, bl := NextInitializedTickWithinOneWord(
				ticks,
				big.NewInt(255),
				false,
				1)
			return tick.Cmp(big.NewInt(511)) == 0 && bl == false
		}(),
		true)
	assert.Equal(
		t,
		func() bool {
			tick, bl := NextInitializedTickWithinOneWord(
				ticks,
				big.NewInt(256),
				false,
				1)
			return tick.Cmp(big.NewInt(511)) == 0 && bl == false
		}(),
		true)
}
