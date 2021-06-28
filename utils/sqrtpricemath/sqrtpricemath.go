package sqrtpricemath

import (
	"math/big"

	"github.com/nakaochi/espoon/arbitrage/uniswapsdk/utils/constant"
	"github.com/nakaochi/espoon/arbitrage/uniswapsdk/utils/tickmath"
)

var MaxUint160 = new(big.Int).Sub(new(big.Int).Exp(constant.TWO, big.NewInt(160), nil), constant.ONE)

func multiplyIn256(x *big.Int, y *big.Int) *big.Int {
	product := new(big.Int).Mul(x, y)
	return product.And(product, constant.MaxUint256)
}

func addIn256(x *big.Int, y *big.Int) *big.Int {
	sum := new(big.Int).Add(x, y)
	return sum.And(sum, constant.MaxUint256)
}

func MulDivRoundingUp(a *big.Int, b *big.Int, denominator *big.Int) *big.Int {
	product := new(big.Int).Mul(a, b)
	result := new(big.Int).Div(product, denominator)
	if product.Mod(product, denominator).Cmp(constant.ZERO) != 0 {
		result.Add(result, constant.ONE)
	}
	return result
}

func GetAmount0Delta(
	sqrtRatioAX96 *big.Int,
	sqrtRatioBX96 *big.Int,
	liquidity *big.Int,
	roundUp bool) *big.Int {
	if sqrtRatioAX96.Cmp(sqrtRatioBX96) > 0 {
		sqrtRatioAX96, sqrtRatioBX96 = sqrtRatioBX96, sqrtRatioAX96
	}
	numerator1 := new(big.Int).Lsh(liquidity, 96)
	numerator2 := new(big.Int).Sub(sqrtRatioBX96, sqrtRatioAX96)
	if roundUp {
		res := MulDivRoundingUp(
			MulDivRoundingUp(numerator1, numerator2, sqrtRatioBX96),
			constant.ONE,
			sqrtRatioAX96)
		return res
	} else {
		res := new(big.Int).Mul(numerator1, numerator2)
		res.Div(res, sqrtRatioBX96).Div(res, sqrtRatioAX96)
		return res
	}
}

func GetAmount1Delta(
	sqrtRatioAX96 *big.Int,
	sqrtRatioBX96 *big.Int,
	liquidity *big.Int,
	roundUp bool) *big.Int {
	if sqrtRatioAX96.Cmp(sqrtRatioBX96) > 0 {
		sqrtRatioAX96, sqrtRatioBX96 = sqrtRatioBX96, sqrtRatioAX96
	}
	if roundUp {
		res := new(big.Int).Sub(sqrtRatioBX96, sqrtRatioAX96)
		res = MulDivRoundingUp(
			liquidity,
			res,
			constant.Q96)
		return res
	} else {
		res := new(big.Int).Sub(sqrtRatioBX96, sqrtRatioAX96)
		res.Mul(res, liquidity).Div(res, constant.Q96)
		return res
	}
}

func GetNextSqrtPriceFromInput(
	sqrtPX96 *big.Int,
	liquidity *big.Int,
	amountIn *big.Int,
	zeroForOne bool) *big.Int {
	tickmath.Check(sqrtPX96.Cmp(constant.ZERO) > 0, " ")
	tickmath.Check(liquidity.Cmp(constant.ZERO) > 0, " ")
	if zeroForOne {
		//todo
		return getNextSqrtPriceFromAmount0RoundingUp(
			sqrtPX96,
			liquidity,
			amountIn,
			true)
	} else {
		return getNextSqrtPriceFromAmount1RoundingDown(
			sqrtPX96,
			liquidity,
			amountIn,
			true)
	}
}

func GetNextSqrtPriceFromOutput(
	sqrtPX96 *big.Int,
	liquidity *big.Int,
	amountOut *big.Int,
	zeroForOne bool) *big.Int {
	tickmath.Check(sqrtPX96.Cmp(constant.ZERO) > 0, " ")
	tickmath.Check(liquidity.Cmp(constant.ZERO) > 0, " ")
	if zeroForOne {
		return getNextSqrtPriceFromAmount1RoundingDown(
			sqrtPX96,
			liquidity,
			amountOut,
			false)
	} else {
		return getNextSqrtPriceFromAmount0RoundingUp(
			sqrtPX96,
			liquidity,
			amountOut,
			false)

	}

}

func getNextSqrtPriceFromAmount0RoundingUp(
	sqrtPX96 *big.Int,
	liquidity *big.Int,
	amount *big.Int,
	add bool) *big.Int {
	if amount.Cmp(constant.ZERO) == 0 {
		return sqrtPX96
	}
	numerator1 := new(big.Int).Lsh(liquidity, 96)
	if add {
		product := multiplyIn256(amount, sqrtPX96)
		if new(big.Int).Div(product, amount).Cmp(sqrtPX96) == 0 {
			denominator := addIn256(numerator1, product)
			if denominator.Cmp(numerator1) >= 0 {
				return MulDivRoundingUp(numerator1, sqrtPX96, denominator)
			}
		}
		mid := new(big.Int).Div(numerator1, sqrtPX96)
		mid.Add(mid, amount)
		return MulDivRoundingUp(numerator1, constant.ONE, mid)
	} else {
		product := multiplyIn256(amount, sqrtPX96)
		tickmath.Check(new(big.Int).Div(product, amount).Cmp(sqrtPX96) == 0, " ")
		tickmath.Check(numerator1.Cmp(product) > 0, " ")
		denominator := new(big.Int).Sub(numerator1, product)
		return MulDivRoundingUp(numerator1, sqrtPX96, denominator)
	}
}

func getNextSqrtPriceFromAmount1RoundingDown(
	sqrtPX96 *big.Int,
	liquidity *big.Int,
	amount *big.Int,
	add bool) *big.Int {
	if add {
		var quotient *big.Int
		if amount.Cmp(MaxUint160) <= 0 {
			quotient = new(big.Int).Lsh(amount, 96)
			quotient.Div(quotient, liquidity)
		} else {
			quotient = new(big.Int).Mul(amount, constant.Q96)
			quotient.Div(quotient, liquidity)
		}
		return quotient.Add(quotient, sqrtPX96)
	} else {
		quotient := MulDivRoundingUp(amount, constant.Q96, liquidity)
		tickmath.Check(sqrtPX96.Cmp(quotient) > 0, " ")
		return quotient.Sub(sqrtPX96, quotient)
	}

}
