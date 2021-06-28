package swapmath

import (
	"math/big"

	"github.com/nakaochi/espoon/arbitrage/uniswapsdk/utils/constant"
	"github.com/nakaochi/espoon/arbitrage/uniswapsdk/utils/sqrtpricemath"
)

var MAX_FEE = new(big.Int).Exp(big.NewInt(10), big.NewInt(6), nil)

type returnValues struct {
	sqrtRatioNextX96 *big.Int
	amountIn         *big.Int
	amountOut        *big.Int
	feeAmount        *big.Int
}

func ComputeSwapStep(
	sqrtRatioCurrentX96 *big.Int,
	sqrtRatioTargetX96 *big.Int,
	liquidity *big.Int,
	amountRemaining *big.Int,
	feePips int) (*big.Int, *big.Int, *big.Int, *big.Int) {
	returnvalues := &returnValues{}
	zeroForeOne := sqrtRatioCurrentX96.Cmp(sqrtRatioTargetX96) >= 0
	exactIn := amountRemaining.Cmp(constant.ZERO) >= 0

	feePipsB := big.NewInt(int64(feePips))

	if exactIn {
		amountRemainingLessFee := new(big.Int).Sub(MAX_FEE, feePipsB)
		amountRemainingLessFee.Mul(amountRemainingLessFee, amountRemaining)
		amountRemainingLessFee.Div(amountRemainingLessFee, MAX_FEE)
		if zeroForeOne {
			returnvalues.amountIn = sqrtpricemath.GetAmount0Delta(
				sqrtRatioTargetX96,
				sqrtRatioCurrentX96,
				liquidity,
				true)
		} else {
			returnvalues.amountIn = sqrtpricemath.GetAmount1Delta(
				sqrtRatioCurrentX96,
				sqrtRatioTargetX96,
				liquidity,
				true)
		}

		if amountRemainingLessFee.Cmp(returnvalues.amountIn) >= 0 {
			returnvalues.sqrtRatioNextX96 = sqrtRatioTargetX96
		} else {
			returnvalues.sqrtRatioNextX96 = sqrtpricemath.GetNextSqrtPriceFromInput(
				sqrtRatioCurrentX96,
				liquidity,
				amountRemainingLessFee,
				zeroForeOne)
		}

	} else {
		if zeroForeOne {
			returnvalues.amountOut = sqrtpricemath.GetAmount1Delta(
				sqrtRatioTargetX96,
				sqrtRatioCurrentX96,
				liquidity,
				false)
		} else {
			returnvalues.amountOut = sqrtpricemath.GetAmount0Delta(
				sqrtRatioCurrentX96,
				sqrtRatioTargetX96,
				liquidity,
				false)
		}
		if new(big.Int).Mul(amountRemaining, constant.NEGATIVE_ONE).Cmp(returnvalues.amountOut) >= 0 {
			returnvalues.sqrtRatioNextX96 = sqrtRatioTargetX96
		} else {
			returnvalues.sqrtRatioNextX96 = sqrtpricemath.GetNextSqrtPriceFromOutput(
				sqrtRatioCurrentX96,
				liquidity,
				new(big.Int).Mul(amountRemaining, constant.NEGATIVE_ONE),
				zeroForeOne)
		}

	}

	max := sqrtRatioTargetX96.Cmp(returnvalues.sqrtRatioNextX96) == 0

	if zeroForeOne {
		if max && exactIn {
			returnvalues.amountIn = returnvalues.amountIn
		} else {
			returnvalues.amountIn = sqrtpricemath.GetAmount0Delta(
				returnvalues.sqrtRatioNextX96,
				sqrtRatioCurrentX96,
				liquidity,
				true)
		}
		if max && !exactIn {
			returnvalues.amountOut = returnvalues.amountOut
		} else {
			returnvalues.amountOut = sqrtpricemath.GetAmount1Delta(
				returnvalues.sqrtRatioNextX96,
				sqrtRatioCurrentX96,
				liquidity,
				false)
		}
		/*
			if max && exactIn {
				returnvalues.amountOut = sqrtpricemath.GetAmount1Delta(
					returnvalues.sqrtRatioNextX96,
					sqrtRatioCurrentX96,
					liquidity,
					false)
			} else {
				returnvalues.amountIn = sqrtpricemath.GetAmount0Delta(
					returnvalues.sqrtRatioNextX96,
					sqrtRatioCurrentX96,
					liquidity,
					true)
			}*/
	} else {
		if max && exactIn {
			returnvalues.amountIn = returnvalues.amountIn
		} else {
			returnvalues.amountIn = sqrtpricemath.GetAmount1Delta(
				sqrtRatioCurrentX96,
				returnvalues.sqrtRatioNextX96,
				liquidity,
				true)
		}
		if max && !exactIn {
			returnvalues.amountOut = returnvalues.amountOut
		} else {
			returnvalues.amountOut = sqrtpricemath.GetAmount0Delta(
				sqrtRatioCurrentX96,
				returnvalues.sqrtRatioNextX96,
				liquidity,
				false)
		}
		/*
			if max && exactIn {
				returnvalues.amountOut = sqrtpricemath.GetAmount0Delta(
					sqrtRatioCurrentX96,
					returnvalues.sqrtRatioNextX96,
					liquidity,
					false)
			} else {
				returnvalues.amountIn = sqrtpricemath.GetAmount1Delta(
					sqrtRatioCurrentX96,
					returnvalues.sqrtRatioNextX96,
					liquidity,
					true)
			}*/
	}

	tmp := new(big.Int).Mul(amountRemaining, constant.NEGATIVE_ONE)
	if !exactIn && returnvalues.amountOut.Cmp(tmp) > 0 {
		returnvalues.amountOut = tmp
	}

	if exactIn && returnvalues.sqrtRatioNextX96.Cmp(sqrtRatioTargetX96) != 0 {
		returnvalues.feeAmount = new(big.Int).Sub(amountRemaining, returnvalues.amountIn)
	} else {
		returnvalues.feeAmount = sqrtpricemath.MulDivRoundingUp(
			returnvalues.amountIn,
			feePipsB,
			new(big.Int).Sub(MAX_FEE, feePipsB))
	}
	return returnvalues.sqrtRatioNextX96, returnvalues.amountIn, returnvalues.amountOut, returnvalues.feeAmount

}
