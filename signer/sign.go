package signer

import (
	"fmt"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/crypto"
)

// SigShim is used for introducing signature functions
type SigShim interface {
	GenPrivate() ([]byte, error)
	ToPublic(pk []byte) ([]byte, error)
	Sign(pk []byte, msg []byte) ([]byte, error)
	Verify(sig []byte, a address.Address, msg []byte) error
}

var sigs map[crypto.SigType]SigShim

func init() {
	sigs = make(map[crypto.SigType]SigShim, 2)
	sigs[crypto.SigTypeBLS] = new(blsSigner)
	sigs[crypto.SigTypeSecp256k1] = new(secpSigner)
}

func NewSecp256k1Singer() SigShim {
	return new(secpSigner)
}

func NewBLSSinger() SigShim {
	return new(blsSigner)
}

// Sign takes in signature type, private key and message. Returns a signature for that message.
// Valid sigTypes are: "secp256k1" and "bls"
func Sign(sigType crypto.SigType, privkey []byte, msg []byte) (*crypto.Signature, error) {
	sv, ok := sigs[sigType]
	if !ok {
		return nil, fmt.Errorf("cannot sign message with signature of unsupported type: %v", sigType)
	}

	sb, err := sv.Sign(privkey, msg)
	if err != nil {
		return nil, err
	}
	return &crypto.Signature{
		Type: sigType,
		Data: sb,
	}, nil
}

func AddressSigType(addr address.Address) crypto.SigType {
	if addr.Protocol() == address.SECP256K1 {
		return crypto.SigTypeSecp256k1
	}

	if addr.Protocol() == address.BLS {
		return crypto.SigTypeBLS
	}

	return crypto.SigTypeUnknown
}
