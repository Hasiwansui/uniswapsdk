package tickmath

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/nakaochi/espoon/arbitrage/uniswapsdk/utils/constant"
	"github.com/stretchr/testify/assert"
)

func Test_GetSqrtRatioAtTick(t *testing.T) {
	fmt.Println("GetSqrtRatioAtTick")
	assert.Panics(
		t,
		func() {
			GetSqrtRatioAtTick(new(big.Int).Sub(constant.MIN_TICK, constant.ONE))
		})
	assert.Panics(
		t,
		func() {
			GetSqrtRatioAtTick(new(big.Int).Add(constant.MAX_TICK, constant.ONE))
		})
	assert.True(t, GetSqrtRatioAtTick(constant.MIN_TICK).Cmp(constant.MIN_SQRT_RATIO) == 0)
	assert.True(t, GetSqrtRatioAtTick(constant.MAX_TICK).Cmp(constant.MAX_SQRT_RATIO) == 0)
	assert.True(t, GetSqrtRatioAtTick(constant.ZERO).Cmp(new(big.Int).Lsh(constant.ONE, 96)) == 0)
}

func Test_GetTickAtSqrtRatio(t *testing.T) {
	assert.True(t, GetTickAtSqrtRatio(constant.MIN_SQRT_RATIO).Cmp(constant.MIN_TICK) == 0)
	assert.True(t, GetTickAtSqrtRatio(new(big.Int).Sub(constant.MAX_SQRT_RATIO, constant.ONE)).Cmp(new(big.Int).Sub(constant.MAX_TICK, constant.ONE)) == 0)
}
