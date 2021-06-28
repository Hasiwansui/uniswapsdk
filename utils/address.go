package utils

import (
	"bytes"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var POOL_INIT_CODE_HASH = "e34f199b19b2b4f47f68442619d555527d244f78a3297ea89325f843f87b8b54"
var FACTORY_ADDRESS = "1F98431c8aD98523631AE4a59f267346ea31F984"

var FEE_LOW = 500
var FEE_MEDIUM = 3000
var FEE_HIGH = 10000
var TICK_LOW = 10
var TICK_MEDIUM = 3000
var TICK_HIGH = 10000

func FeeTick(fee int) int {
	if fee == FEE_LOW {
		return TICK_LOW
	}
	if fee == FEE_MEDIUM {
		return TICK_MEDIUM
	}
	if fee == FEE_HIGH {
		return TICK_HIGH
	}
	return 0
}

//TODO
func ComputePoolAddress(
	factoryAddress string,
	tokenA string,
	tokenB string,
	fee int) string {
	var token0 common.Address
	var token1 common.Address
	cp0, _ := new(big.Int).SetString(tokenA, 0)
	cp1, _ := new(big.Int).SetString(tokenB, 0)
	if cp0.Cmp(cp1) < 0 {
		token0, token1 = common.HexToAddress(tokenA), common.HexToAddress(tokenB)
	} else {
		token0, token1 = common.HexToAddress(tokenB), common.HexToAddress(tokenA)
	}
	var addressType, _ = abi.NewType("address", "address", nil)
	var uint24Type, _ = abi.NewType("uint24", "uint24", nil)
	encode := abi.Arguments{
		{
			Type: addressType,
		},
		{
			Type: addressType,
		},
		{
			Type: uint24Type,
		}}
	head := common.Hex2Bytes("ff")
	btes, _ := encode.Pack(token0, token1, big.NewInt(int64(fee)))

	salt := crypto.Keccak256(btes)
	factory := common.Hex2Bytes(factoryAddress)
	init := common.Hex2Bytes(POOL_INIT_CODE_HASH)
	var buffer bytes.Buffer
	buffer.Write(head)
	buffer.Write(factory)
	buffer.Write(salt)
	buffer.Write(init)
	input := buffer.Bytes()
	output := crypto.Keccak256(input)
	output = output[12:]
	return "0x" + common.Bytes2Hex(output)

}
