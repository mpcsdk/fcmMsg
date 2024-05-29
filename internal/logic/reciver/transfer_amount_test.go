package reciver

import (
	"fmt"
	"math"
	"math/big"
	"testing"
)

func TestAmount(t *testing.T) {
	fbalance := big.NewFloat(0)
	fbalance.SetString("1000000")
	fval := fbalance.Quo(fbalance, big.NewFloat(math.Pow10(18)))
	fmt.Println(fval.Text('f', -1))
	////
}
