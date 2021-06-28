package entities

import (
	"fmt"
	"math/big"

	"github.com/nakaochi/espoon/arbitrage/uniswapsdk/utils"
	"github.com/nakaochi/espoon/arbitrage/uniswapsdk/utils/constant"
	"github.com/nakaochi/espoon/arbitrage/uniswapsdk/utils/liquiditymath"
	"github.com/nakaochi/espoon/arbitrage/uniswapsdk/utils/swapmath"
	"github.com/nakaochi/espoon/arbitrage/uniswapsdk/utils/tickmath"
)

type StepComputations struct {
	sqrtPriceStartX96 *big.Int
	tickNext          *big.Int
	initialized       bool
	sqrtPriceNextX96  *big.Int
	amountIn          *big.Int
	amountOut         *big.Int
	feeAmount         *big.Int
}

type State struct {
	AmountSpecifiedRemaining *big.Int
	AmountCalculated         *big.Int
	SqrtPriceX96             *big.Int
	Tick                     *big.Int
	Liquidity                *big.Int
}

type Pool struct {
	token0           Token
	token1           Token
	fee              int
	sqrtRatioX96     *big.Int
	liquidity        *big.Int
	tickCurrent      *big.Int
	tickDataProvider *TickDataProvider

	//_token0Price
	//_token1Price
}
type TokenAmount struct {
	Token  Token
	Amount *big.Int
}

type Token struct {
	Address string
}

func (P *Token) equals(token Token) bool {
	return P.Address == token.Address
}

func (P *Token) sortsBefore(token Token) bool {
	a, _ := new(big.Int).SetString(P.Address, 0)
	b, _ := new(big.Int).SetString(token.Address, 0)
	return a.Cmp(b) < 0
}

func (P *Pool) involvesToken(token Token) bool {
	return token.equals(P.token0) || token.equals(P.token1)
}

func GetAddress(tokenA string, tokenB string, fee int) string {
	return utils.ComputePoolAddress(utils.FACTORY_ADDRESS, tokenA, tokenB, fee)
}

func (P *Pool) swap(
	zeroForOne bool,
	amountSpecified *big.Int,
	sqrtPriceLimitX96 *big.Int) (*big.Int, *big.Int, *big.Int, *big.Int) {
	if sqrtPriceLimitX96 == nil {
		if zeroForOne {
			sqrtPriceLimitX96 = new(big.Int).Add(constant.MIN_SQRT_RATIO, constant.ONE)
		} else {
			sqrtPriceLimitX96 = new(big.Int).Sub(constant.MAX_SQRT_RATIO, constant.ONE)
		}
	}

	if zeroForOne {
		tickmath.Check(sqrtPriceLimitX96.Cmp(constant.MIN_SQRT_RATIO) > 0, "RATIO_MIN")
		tickmath.Check(sqrtPriceLimitX96.Cmp(P.sqrtRatioX96) < 0, "RATIO_CURRENT")
	} else {
		tickmath.Check(sqrtPriceLimitX96.Cmp(constant.MAX_SQRT_RATIO) < 0, "RATIO_MAX")
		tickmath.Check(sqrtPriceLimitX96.Cmp(P.sqrtRatioX96) > 0, "RATIO_CURRENT")
	}
	exactInput := amountSpecified.Cmp(constant.ZERO) >= 0

	state := &State{
		AmountSpecifiedRemaining: new(big.Int).Add(amountSpecified, constant.ZERO), //safe but not efficiency
		AmountCalculated:         big.NewInt(0),
		SqrtPriceX96:             new(big.Int).Add(P.sqrtRatioX96, constant.ZERO),
		Tick:                     new(big.Int).Add(P.tickCurrent, constant.ZERO),
		Liquidity:                new(big.Int).Add(P.liquidity, constant.ZERO)}

	for state.AmountSpecifiedRemaining.Cmp(constant.ZERO) != 0 && state.SqrtPriceX96.Cmp(sqrtPriceLimitX96) != 0 {
		step := &StepComputations{}
		step.sqrtPriceStartX96 = new(big.Int).Add(state.SqrtPriceX96, constant.ZERO)

		step.tickNext, step.initialized =
			(*P.tickDataProvider).NextInitializedTickWithinOneWord(
				state.Tick,
				zeroForOne,
				P.tickSpacing())

		if step.tickNext.Cmp(constant.MIN_TICK) < 0 {
			step.tickNext.Add(constant.MIN_TICK, constant.ZERO)
		} else {
			if step.tickNext.Cmp(constant.MAX_TICK) > 0 {
				step.tickNext.Add(constant.MAX_TICK, constant.ZERO)
			}
		}
		step.sqrtPriceNextX96 = tickmath.GetSqrtRatioAtTick(step.tickNext)
		//TODO SWAP

		var param *big.Int
		var con bool
		if zeroForOne {
			con = step.sqrtPriceNextX96.Cmp(sqrtPriceLimitX96) < 0
		} else {
			con = step.sqrtPriceNextX96.Cmp(sqrtPriceLimitX96) > 0
		}
		if con {
			param = sqrtPriceLimitX96
		} else {
			param = step.sqrtPriceNextX96
		}
		fmt.Println("befor", state.SqrtPriceX96,
			param,
			state.Liquidity,
			state.AmountSpecifiedRemaining,
			P.fee)
		step.sqrtPriceNextX96, step.amountIn, step.amountOut, step.feeAmount = swapmath.ComputeSwapStep(
			state.SqrtPriceX96,
			param,
			state.Liquidity,
			state.AmountSpecifiedRemaining,
			P.fee)

		if exactInput {
			state.AmountSpecifiedRemaining.Sub(
				state.AmountSpecifiedRemaining,
				new(big.Int).Add(step.amountIn, step.feeAmount))
			state.AmountCalculated.Sub(state.AmountCalculated, step.amountOut)
		} else {
			state.AmountSpecifiedRemaining.Add(
				state.AmountSpecifiedRemaining,
				step.amountOut)
			state.AmountCalculated.Add(
				state.AmountCalculated,
				new(big.Int).Add(step.amountIn, step.feeAmount))
		}

		if state.SqrtPriceX96.Cmp(step.sqrtPriceNextX96) == 0 {
			if step.initialized {
				liquidityNet := new(big.Int).Add(
					(*P.tickDataProvider).GetTick(step.tickNext).LiquidityNet,
					constant.ZERO)
				if zeroForOne {
					liquidityNet.Mul(liquidityNet, constant.NEGATIVE_ONE)
				}
				state.Liquidity = liquiditymath.AddDelta(state.Liquidity, liquidityNet)
			}
			if zeroForOne {
				state.Tick.Sub(step.tickNext, constant.ONE)
			} else {
				state.Tick = step.tickNext
			}
		} else {
			if state.SqrtPriceX96.Cmp(step.sqrtPriceStartX96) != 0 {
				state.Tick = tickmath.GetTickAtSqrtRatio(state.SqrtPriceX96)
			}
		}

	}
	return state.AmountCalculated, state.SqrtPriceX96, state.Liquidity, state.Tick

}

func (P *Pool) tickSpacing() int {
	return utils.FeeTick(P.fee)
}

func (P *Pool) GetOutputAmount(
	inputAmount TokenAmount,
	sqrtPriceLimitX96 *big.Int) (*TokenAmount, *Pool) {
	tickmath.Check(P.involvesToken(inputAmount.Token), "TOKEN")
	zeroForOne := inputAmount.Token.equals(P.token0)
	outputAmount, sqrtRatioX96, liquidity, tickCurrent :=
		P.swap(
			zeroForOne,
			inputAmount.Amount,
			sqrtPriceLimitX96)

	var outputToken Token
	if zeroForOne {
		outputToken = P.token1
	} else {
		outputToken = P.token0
	}
	outputAmount.Mul(outputAmount, constant.NEGATIVE_ONE)
	return &TokenAmount{outputToken, outputAmount}, NewPool(P.token0, P.token1, P.fee, sqrtRatioX96, liquidity, tickCurrent, *P.tickDataProvider)

}

func (P *Pool) GetInputAmount(
	outputAmount TokenAmount,
	sqrtPriceLimitX96 *big.Int) (*TokenAmount, *Pool) {
	tickmath.Check(P.involvesToken(outputAmount.Token), "TOKEN")
	zeroForOne := outputAmount.Token.equals(P.token1)
	inputAmount, sqrtRatioX96, liquidity, tickCurrent :=
		P.swap(
			zeroForOne,
			new(big.Int).Mul(outputAmount.Amount, constant.NEGATIVE_ONE),
			sqrtPriceLimitX96)
	var inputToken Token
	if zeroForOne {
		inputToken = P.token0
	} else {
		inputToken = P.token1
	}
	return &TokenAmount{inputToken, inputAmount}, NewPool(P.token0, P.token1, P.fee, sqrtRatioX96, liquidity, tickCurrent, *P.tickDataProvider)
}

func NewPool(
	tokenA Token,
	tokenB Token,
	fee int,
	sqrtRatioX96 *big.Int,
	liquidity *big.Int,
	tickCurrent *big.Int,
	ticks TickDataProvider) *Pool {
	tickmath.Check(fee < 1000000, "FEE")
	tickCurrentSqrtRatioX96 := tickmath.GetSqrtRatioAtTick(tickCurrent)
	nextTickSqrtRatioX96 := tickmath.GetSqrtRatioAtTick(new(big.Int).Add(tickCurrent, constant.ONE))
	tickmath.Check(
		sqrtRatioX96.Cmp(tickCurrentSqrtRatioX96) >= 0 && sqrtRatioX96.Cmp(nextTickSqrtRatioX96) <= 0,
		"PRICE_BOUNDS")
	pool := &Pool{}
	if tokenA.sortsBefore(tokenB) {
		pool.token0, pool.token1 = tokenA, tokenB
	} else {
		pool.token0, pool.token1 = tokenB, tokenA
	}
	pool.fee = fee
	pool.sqrtRatioX96 = sqrtRatioX96
	pool.liquidity = liquidity
	pool.tickCurrent = tickCurrent
	pool.tickDataProvider = &ticks
	return pool
}
