package internal

import (
	"encoding/hex"
	"errors"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/crypto"
	"github.com/labs3/filecoin-wallet-signing/signer"
)

type AnyMessage struct {
	From address.Address
	Msg  string // original message
	Sig  string // signature
}

func (msg *AnyMessage) SignMsg(privateKey []byte) (string, error) {
	sigType := signer.AddressSigType(msg.From)
	signed, err := signer.Sign(sigType, privateKey, []byte(msg.Msg))
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(signed.Data), nil
}

func (msg *AnyMessage) VerifyMsg() error {
	sigType := signer.AddressSigType(msg.From)
	var signObj signer.SigShim
	switch sigType {
	case crypto.SigTypeSecp256k1:
		signObj = signer.NewSecp256k1Singer()
	case crypto.SigTypeBLS:
		signObj = signer.NewBLSSinger()
	default:
		return errors.New("SigTypeUnknown")
	}

	sig, err := hex.DecodeString(msg.Sig)
	if err != nil {
		return err
	}

	return signObj.Verify(sig, msg.From, []byte(msg.Msg))
}
