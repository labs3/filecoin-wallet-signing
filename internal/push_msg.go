package internal

import (
	"encoding/json"
	"fmt"

	"github.com/labs3/filecoin-wallet-signing/chain/types"
	"github.com/labs3/filecoin-wallet-signing/signer"
)

func PushSignedMsg(msg *types.Message, privateKey []byte) error {
	nonce, err := Lapi.MpoolGetNonce(Ctx, msg.From)
	if err != nil {
		fmt.Println("Mpool GetNonce failed: ", err)
		return err
	}
	msg.Nonce = nonce
	msgWithGas, err := Lapi.GasEstimateMessageGas(Ctx, msg, nil, *CurrentTsk)
	if err != nil {
		fmt.Println("GasEstimateMessageGas failed: ", err)
		return err
	}

	blk, err := msgWithGas.ToStorageBlock()
	if err != nil {
		fmt.Println("msg.ToStorageBlock() failed: ", err)
		return err
	}

	sigType := signer.AddressSigType(msg.From)
	signed, err := signer.Sign(sigType, privateKey, blk.Cid().Bytes())
	if err != nil {
		fmt.Println("sign failed: ", err.Error())
		return err
	}
	signedMsg := types.SignedMessage{
		Message:   *msgWithGas,
		Signature: *signed,
	}
	b, _ := json.MarshalIndent(signedMsg, " ", " ")
	fmt.Println("Signed message: ", string(b))

	msgCid, err := Lapi.MpoolPush(Ctx, &signedMsg)
	if err != nil {
		fmt.Println("push message failed: ", err.Error())
		return err
	}

	fmt.Println("message CID:", msgCid.String())
	return nil
}
