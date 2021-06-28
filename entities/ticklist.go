package entities

import (
	"math"
	"math/big"

	"github.com/nakaochi/espoon/arbitrage/uniswapsdk/utils/constant"
	"github.com/nakaochi/espoon/arbitrage/uniswapsdk/utils/tickmath"
)

func tickComparatoe(a Tick, b Tick) int {
	return a.Index.Cmp(b.Index)
}

func supportValidation(ticks []Tick, tickSpacing int) bool {
	spacing := big.NewInt(int64(tickSpacing))
	tmp := new(big.Int)
	for i := 0; i < len(ticks); i++ {
		if tmp.Mod(ticks[i].Index, spacing).Cmp(constant.ZERO) != 0 {
			return false
		}
	}
	return true
}

func isSorted(ticks []Tick) bool {
	tickmath.Check(len(ticks) > 1, "LEN")
	for i := 0; i < len(ticks)-1; i++ {
		if tickComparatoe(ticks[i], ticks[i+1]) == 1 {
			return false
		}
	}
	return true
}

func ValidateList(ticks []Tick, tickSpacing int) {
	tickmath.Check(tickSpacing > 0, "TICK_SPACING_NONZERO")
	tickmath.Check(supportValidation(ticks, tickSpacing), "TICK_SPACING")
	tickmath.Check(isSorted(ticks), "SORTED")
}

func isBelowSmallest(ticks []Tick, tick *big.Int) bool {
	tickmath.Check(len(ticks) > 0, "LENGTH")
	return tick.Cmp(ticks[0].Index) < 0
}

func isAtOrAboveLargest(ticks []Tick, tick *big.Int) bool {
	tickmath.Check(len(ticks) > 0, "LENGTH")
	return tick.Cmp(ticks[len(ticks)-1].Index) >= 0
}

func GetTick(ticks []Tick, index *big.Int) *Tick {
	tick := &ticks[BinarySearch(ticks, index)]
	tickmath.Check(tick.Index.Cmp(index) == 0, "NOT_CONTAINED")
	return tick
}

func BinarySearch(ticks []Tick, tick *big.Int) int {
	tickmath.Check(!isBelowSmallest(ticks, tick), "BELOW_SMALLEST")
	var l float64 = 0
	var r float64 = float64(len(ticks) - 1)
	var i int
	for {
		i = int(math.Floor((l + r) / 2))
		if ticks[i].Index.Cmp(tick) <= 0 && (i == len(ticks)-1 || ticks[i+1].Index.Cmp(tick) > 0) {
			return i
		}
		if ticks[i].Index.Cmp(tick) < 0 {
			l = float64(i + 1)
		} else {
			r = float64(i - 1)
		}
	}

}
func NextInitializedTick(ticks []Tick, tick *big.Int, lte bool) Tick {
	if lte {
		tickmath.Check(!isBelowSmallest(ticks, tick), "BELOW_SMALLEST")
		if isAtOrAboveLargest(ticks, tick) {
			return ticks[len(ticks)-1]
		}
		index := BinarySearch(ticks, tick)
		return ticks[index]
	} else {
		tickmath.Check(!isAtOrAboveLargest(ticks, tick), "AT_OR_ABOVE_LARGEST")
		if isBelowSmallest(ticks, tick) {
			return ticks[0]
		}
		index := BinarySearch(ticks, tick)
		return ticks[index+1]
	}
}

func max(a *big.Int, b *big.Int) *big.Int {
	if a.Cmp(b) > 0 {
		return a
	}
	return b
}

func min(a *big.Int, b *big.Int) *big.Int {
	if a.Cmp(b) < 0 {
		return a
	}
	return b
}

func NextInitializedTickWithinOneWord(ticks []Tick, tick *big.Int, lte bool, tickSpacing int) (*big.Int, bool) {
	spacing := big.NewInt(int64(tickSpacing))
	compressed := new(big.Int).Div(tick, spacing)
	if lte {
		wordPos := new(big.Int).Rsh(compressed, 8)
		minimum := wordPos.Lsh(wordPos, 8).Mul(wordPos, spacing)
		if isBelowSmallest(ticks, tick) {
			return minimum, false
		}
		index := NextInitializedTick(ticks, tick, lte).Index
		nextInitializedTick := max(minimum, index)
		return nextInitializedTick, nextInitializedTick.Cmp(index) == 0
	} else {
		wordPos := new(big.Int).Rsh(compressed.Add(compressed, constant.ONE), 8)
		maximum := wordPos.Add(wordPos, constant.ONE).Lsh(wordPos, 8).Mul(wordPos, spacing).Sub(wordPos, constant.ONE)

		if isAtOrAboveLargest(ticks, tick) {
			return maximum, false
		}
		index := NextInitializedTick(ticks, tick, lte).Index
		nextInitializedTick := min(maximum, index)
		return nextInitializedTick, nextInitializedTick.Cmp(index) == 0
	}
}
