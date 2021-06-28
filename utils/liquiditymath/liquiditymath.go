package liquiditymath

import (
	"math/big"

	"github.com/nakaochi/espoon/arbitrage/uniswapsdk/utils/constant"
)

func AddDelta(x *big.Int, y *big.Int) *big.Int {
	if y.Cmp(constant.ZERO) < 0 {
		result := new(big.Int).Mul(y, constant.NEGATIVE_ONE)
		result.Sub(x, result)
		return result
	}
	return new(big.Int).Add(x, y)
}
