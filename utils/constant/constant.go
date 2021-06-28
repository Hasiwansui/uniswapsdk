package constant

import (
	"math/big"
)

var MIN_TICK *big.Int = big.NewInt(-887272)
var MAX_TICK *big.Int = big.NewInt(887272)
var MIN_SQRT_RATIO, _ = new(big.Int).SetString("4295128739", 0)
var MAX_SQRT_RATIO, _ = new(big.Int).SetString("1461446703485210103287273052203988822378723970342", 0)
var ETHER, _ = new(big.Int).SetString("1000000000000000000", 0)
var ZERO = big.NewInt(0)
var ONE = big.NewInt(1)
var NEGATIVE_ONE = big.NewInt(-1)
var TWO = big.NewInt(2)
var Q32 = new(big.Int).Exp(TWO, big.NewInt(32), nil)
var Q96 = new(big.Int).Exp(TWO, big.NewInt(96), nil)
var Q192 = new(big.Int).Exp(Q96, TWO, nil)
var MaxUint256, _ = new(big.Int).SetString("0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", 0)

func Pos() {}
