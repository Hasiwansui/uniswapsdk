package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ComputePoolAddress(t *testing.T) {
	res := ComputePoolAddress(
		FACTORY_ADDRESS,
		"0x6b175474e89094c44da98b954eedeac495271d0f",
		"0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
		500)
	fmt.Println("Pair Address:", res)
	assert.Equal(t, res, "0x6c6bc977e13df9b0de53b251522280bb72383700")
}
