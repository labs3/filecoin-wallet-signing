package signer

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/crypto"
)

const (
	KTBLS       string = "bls"
	KTSecp256k1 string = "secp256k1"
	// KTSecp256k1Ledger string = "secp256k1-ledger"
)

type KeyInfo struct {
	Type       string // secp256k1 or bls
	PrivateKey []byte
}

type Key struct {
	KeyInfo

	PublicKey []byte
	Address   address.Address
}

func ActSigType(typ string) crypto.SigType {
	switch typ {
	case KTBLS:
		return crypto.SigTypeBLS
	case KTSecp256k1:
		return crypto.SigTypeSecp256k1
	default:
		return crypto.SigTypeUnknown
	}
}

// DecodePricateKey Decode the private key exported: eg:7b22...227d
func DecodePricateKey(k string) (*Key, error) {
	pv, err := hex.DecodeString(k)
	if err != nil {
		return nil, err
	}

	key := new(Key)

	err = json.Unmarshal(pv, key)
	if err != nil {
		return nil, err
	}

	var signer SigShim
	if key.Type == KTBLS {
		signer = new(blsSigner)
	} else if key.Type == KTSecp256k1 {
		signer = new(secpSigner)
	} else {
		return nil, fmt.Errorf("unknow key type : %v", key.Type)
	}

	pk, err := signer.ToPublic(key.PrivateKey)
	if err != nil {
		return nil, err
	}
	var addr address.Address
	if key.Type == KTSecp256k1 {
		addr, err = address.NewSecp256k1Address(pk)
		if err != nil {
			return nil, err
		}
	}

	if key.Type == KTBLS {
		addr, err = address.NewBLSAddress(pk)
		if err != nil {
			return nil, err
		}
	}

	return &Key{
		KeyInfo: KeyInfo{
			Type:       key.Type,
			PrivateKey: key.PrivateKey,
		},
		PublicKey: pk,
		Address:   addr,
	}, nil
}
