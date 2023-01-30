package pkg

import (
	"fmt"
	"math/big"

	"github.com/filecoin-project/go-state-types/abi"

	"github.com/labs3/filecoin-wallet-signing/signer"
)

// FilecoinPrecision 10^18
const FilecoinPrecision = uint64(1_000_000_000_000_000_000)

// ToFloat64 convert TokenAmount to FIL(float64)
func ToFloat64(f abi.TokenAmount) float64 {
	var zero float64 = 0

	if f.Int == nil {
		return zero
	}

	fp := big.NewInt(int64(FilecoinPrecision))
	fil, _ := new(big.Rat).SetFrac(f.Int, fp).Float64()

	return fil
}

// ReadPrivteKey read and decode private key from console input
func ReadPrivteKey() (*signer.Key, error) {
	ps := ""
	fmt.Print("Please enter the private key: ")
	_, err := fmt.Scanln(&ps)
	if err != nil {
		return nil, err
	}

	return signer.DecodePricateKey(ps)
}
