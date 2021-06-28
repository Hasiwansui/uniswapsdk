package tickmath

import (
	"math"
	"math/big"

	"github.com/nakaochi/espoon/arbitrage/uniswapsdk/utils/constant"
)

/*
var MIN_TICK *big.Int = big.NewInt(-887272)
var MAX_TICK *big.Int = big.NewInt(887272)
var MIN_SQRT_RATIO, _ = new(big.Int).SetString("4295128739", 0)
var MAX_SQRT_RATIO, _ = new(big.Int).SetString("1461446703485210103287273052203988822378723970342", 0)

var ZERO = big.NewInt(0)
var ONE = big.NewInt(1)
var NEGATIVE_ONE = big.NewInt(-1)
var TWO = big.NewInt(2)
var Q32 = new(big.Int).Exp(TWO, big.NewInt(32), nil)
var Q96 = new(big.Int).Exp(TWO, big.NewInt(96), nil)
var Q192 = new(big.Int).Exp(Q96, TWO, nil)
var MaxUint256, _ = new(big.Int).SetString("0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", 0)
*/

var init0 = "0xfffcb933bd6fad37aa2d162d1a594001"
var init1 = "0x100000000000000000000000000000000"

var shift1, _ = new(big.Int).SetString("0xfff97272373d413259a46990580e213a", 0)
var shift2, _ = new(big.Int).SetString("0xfff2e50f5f656932ef12357cf3c7fdcc", 0)
var shift3, _ = new(big.Int).SetString("0xffe5caca7e10e4e61c3624eaa0941cd0", 0)
var shift4, _ = new(big.Int).SetString("0xffcb9843d60f6159c9db58835c926644", 0)
var shift5, _ = new(big.Int).SetString("0xff973b41fa98c081472e6896dfb254c0", 0)
var shift6, _ = new(big.Int).SetString("0xff2ea16466c96a3843ec78b326b52861", 0)
var shift7, _ = new(big.Int).SetString("0xfe5dee046a99a2a811c461f1969c3053", 0)
var shift8, _ = new(big.Int).SetString("0xfcbe86c7900a88aedcffc83b479aa3a4", 0)
var shift9, _ = new(big.Int).SetString("0xf987a7253ac413176f2b074cf7815e54", 0)
var shift10, _ = new(big.Int).SetString("0xf3392b0822b70005940c7a398e4b70f3", 0)
var shift11, _ = new(big.Int).SetString("0xe7159475a2c29b7443b29c7fa6e889d9", 0)
var shift12, _ = new(big.Int).SetString("0xd097f3bdfd2022b8845ad8f792aa5825", 0)
var shift13, _ = new(big.Int).SetString("0xa9f746462d870fdf8a65dc1f90e061e5", 0)
var shift14, _ = new(big.Int).SetString("0x70d869a156d2a1b890bb3df62baf32f7", 0)
var shift15, _ = new(big.Int).SetString("0x31be135f97d08fd981231505542fcfa6", 0)
var shift16, _ = new(big.Int).SetString("0x9aa508b5b7a84e1c677de54f3e99bc9", 0)
var shift17, _ = new(big.Int).SetString("0x5d6af8dedb81196699c329225ee604", 0)
var shift18, _ = new(big.Int).SetString("0x2216e584f5fa1ea926041bedfe98", 0)
var shift19, _ = new(big.Int).SetString("0x48a170391f7dc42444e8fa2", 0)

var mask1 = big.NewInt(0x2)
var mask2 = big.NewInt(0x4)
var mask3 = big.NewInt(0x8)
var mask4 = big.NewInt(0x10)
var mask5 = big.NewInt(0x20)
var mask6 = big.NewInt(0x40)
var mask7 = big.NewInt(0x80)
var mask8 = big.NewInt(0x100)
var mask9 = big.NewInt(0x200)
var mask10 = big.NewInt(0x400)
var mask11 = big.NewInt(0x800)
var mask12 = big.NewInt(0x1000)
var mask13 = big.NewInt(0x2000)
var mask14 = big.NewInt(0x4000)
var mask15 = big.NewInt(0x8000)
var mask16 = big.NewInt(0x10000)
var mask17 = big.NewInt(0x20000)
var mask18 = big.NewInt(0x40000)
var mask19 = big.NewInt(0x80000)
var maskArray = [19]*big.Int{
	mask1,
	mask2,
	mask3,
	mask4,
	mask5,
	mask6,
	mask7,
	mask8,
	mask9,
	mask10,
	mask11,
	mask12,
	mask13,
	mask14,
	mask15,
	mask16,
	mask17,
	mask18,
	mask19}
var shiftArray = [19]*big.Int{
	shift1,
	shift2,
	shift3,
	shift4,
	shift5,
	shift6,
	shift7,
	shift8,
	shift9,
	shift10,
	shift11,
	shift12,
	shift13,
	shift14,
	shift15,
	shift16,
	shift17,
	shift18,
	shift19}
var POWERS_OF_2 = [8]*big.Int{
	new(big.Int).Exp(constant.TWO, big.NewInt(128), nil),
	new(big.Int).Exp(constant.TWO, big.NewInt(64), nil),
	new(big.Int).Exp(constant.TWO, big.NewInt(32), nil),
	new(big.Int).Exp(constant.TWO, big.NewInt(16), nil),
	new(big.Int).Exp(constant.TWO, big.NewInt(8), nil),
	new(big.Int).Exp(constant.TWO, big.NewInt(4), nil),
	new(big.Int).Exp(constant.TWO, big.NewInt(2), nil),
	new(big.Int).Exp(constant.TWO, big.NewInt(1), nil)}

func mulShift(val *big.Int, mulBy *big.Int) { //NOT CONST
	val.Mul(val, mulBy)
	val.Rsh(val, 128)

}

func Check(b bool, msg string) {
	if b {
		return
	}

	if msg == "" {
		panic("failing invariant")
	}
	panic(msg)
}

func GetSqrtRatioAtTick(tick *big.Int) *big.Int {
	tickBig := new(big.Int).Add(tick, constant.ZERO)
	Check(tickBig.Cmp(constant.MIN_TICK) >= 0 && tickBig.Cmp(constant.MAX_TICK) <= 0, "TICK")
	tickAbs := new(big.Int).Abs(tickBig)

	var ratio *big.Int
	var tmp = new(big.Int)
	if tmp.And(tickAbs, constant.ONE).Cmp(constant.ZERO) != 0 {
		ratio, _ = new(big.Int).SetString(init0, 0)
	} else {
		ratio, _ = new(big.Int).SetString(init1, 0)
	}
	for i := 0; i < 19; i++ {
		if tmp.And(tickAbs, maskArray[i]).Cmp(constant.ZERO) != 0 {
			mulShift(ratio, shiftArray[i])
		}
	}
	if tick.Cmp(constant.ZERO) > 0 {
		ratio = ratio.Div(constant.MaxUint256, ratio)
	}

	if tmp.Mod(ratio, constant.Q32).Cmp(constant.ZERO) > 0 {
		ratio.Div(ratio, constant.Q32).Add(ratio, constant.ONE)
		return ratio
	}
	ratio.Div(ratio, constant.Q32)
	return ratio
}

func mostSignificantBit(x *big.Int) int { //CONST
	X := new(big.Int).Add(x, constant.ZERO)
	Check(X.Cmp(constant.ZERO) > 0, "ZERO")
	Check(X.Cmp(constant.MaxUint256) <= 0, "MAX")
	msb := 0
	for i := 0; i != 8; i++ {
		if X.Cmp(POWERS_OF_2[i]) >= 0 {
			power := uint(math.Pow(2, float64(7-i)))
			X.Rsh(X, power)
			msb += int(power)
		}
	}
	return msb
}

func GetTickAtSqrtRatio(sqrtRatioX96 *big.Int) *big.Int {
	Check(sqrtRatioX96.Cmp(constant.MIN_SQRT_RATIO) >= 0 && sqrtRatioX96.Cmp(constant.MAX_SQRT_RATIO) < 0, "SQRT_RATIO")
	sqrtRatioX128 := new(big.Int).Lsh(sqrtRatioX96, 32)
	msb := mostSignificantBit(sqrtRatioX128)
	Bmsb := big.NewInt(int64(msb))
	Int128 := big.NewInt(128)
	var r *big.Int = new(big.Int)
	if Bmsb.Cmp(Int128) >= 0 {
		r.Rsh(sqrtRatioX128, uint(msb-127))
	} else {
		r.Lsh(sqrtRatioX128, uint(127-msb))
	}
	log_2 := new(big.Int).Sub(Bmsb, Int128)
	log_2.Lsh(log_2, 64)
	tmp := new(big.Int)
	for i := 0; i < 14; i++ {
		r.Mul(r, r).Rsh(r, 127)
		f := new(big.Int).Rsh(r, 128)
		log_2 = log_2.Or(log_2, tmp.Lsh(f, uint(63-i)))
		r.Rsh(r, uint(f.Int64()))
	}
	log_sqrt10001 := tmp.Mul(log_2, mul10001)

	tickLow := new(big.Int).Sub(log_sqrt10001, tlow)
	tickLow.Rsh(tickLow, 128)

	tickHigh := new(big.Int).Add(log_sqrt10001, thigh)
	tickHigh.Rsh(tickHigh, 128)

	if tickLow.Cmp(tickHigh) == 0 {
		return tickLow
	} else {
		if GetSqrtRatioAtTick(tickHigh).Cmp(sqrtRatioX96) <= 0 {
			return tickHigh
		}
	}
	return tickLow

}

var mul10001, _ = new(big.Int).SetString("255738958999603826347141", 0)
var tlow, _ = new(big.Int).SetString("3402992956809132418596140100660247210", 0)
var thigh, _ = new(big.Int).SetString("291339464771989622907027621153398088495", 0)
