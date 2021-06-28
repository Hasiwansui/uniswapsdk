package entities

import (
	"math"
	"math/big"
	"testing"

	"github.com/nakaochi/espoon/arbitrage/uniswapsdk/utils"
	"github.com/nakaochi/espoon/arbitrage/uniswapsdk/utils/constant"
	"github.com/stretchr/testify/assert"
)

var USDC = Token{"0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"}
var DAI = Token{"0x6B175474E89094C44Da98b954EedeAC495271d0F"}

func encodeSqrtRatioX96(amount1 *big.Int, amount0 *big.Int) *big.Int {
	numerator := new(big.Int).Lsh(amount1, 192)
	ratioX192 := new(big.Int).Div(numerator, amount0)
	return ratioX192.Sqrt(ratioX192)
}

func nearestUsableTick(tick *big.Int, tickSpacing int) *big.Int {
	rounded := math.Round(float64(tick.Int64())/float64(tickSpacing)) * float64(tickSpacing)
	rd := big.NewInt(int64(rounded))
	if rd.Cmp(constant.MIN_TICK) < 0 {
		return rd.Add(rd, big.NewInt(int64(tickSpacing)))
	} else {
		if rd.Cmp(constant.MAX_TICK) > 0 {
			return rd.Sub(rd, big.NewInt(int64(tickSpacing)))
		}
	}
	return rd
}

var pool = NewPool(
	USDC,
	DAI,
	utils.FEE_LOW,
	encodeSqrtRatioX96(
		big.NewInt(1),
		big.NewInt(1)),
	constant.ETHER,
	big.NewInt(0),
	TickListProvider{
		[]Tick{
			Tick{
				nearestUsableTick(constant.MIN_TICK, utils.FeeTick(utils.FEE_LOW)),
				constant.ETHER,
				constant.ETHER},
			Tick{
				nearestUsableTick(constant.MAX_TICK, utils.FeeTick(utils.FEE_LOW)),
				constant.ETHER,
				new(big.Int).Mul(constant.ETHER, constant.NEGATIVE_ONE),
			},
		},
	})

func Test_swap(t *testing.T) {
	res, _ := pool.GetOutputAmount(TokenAmount{USDC, big.NewInt(100)}, nil)
	assert.Equal(
		t,
		res.Token.equals(DAI) && res.Amount.Cmp(big.NewInt(98)) == 0,
		true)
	res, _ = pool.GetOutputAmount(TokenAmount{DAI, big.NewInt(100)}, nil)
	assert.Equal(
		t,
		res.Token.equals(USDC) && res.Amount.Cmp(big.NewInt(98)) == 0,
		true)
	res, _ = pool.GetInputAmount(TokenAmount{DAI, big.NewInt(98)}, nil)
	assert.Equal(
		t,
		res.Token.equals(USDC) && res.Amount.Cmp(big.NewInt(100)) == 0,
		true)
	res, _ = pool.GetInputAmount(TokenAmount{USDC, big.NewInt(98)}, nil)
	assert.Equal(
		t,
		res.Token.equals(DAI) && res.Amount.Cmp(big.NewInt(100)) == 0,
		true)
}
