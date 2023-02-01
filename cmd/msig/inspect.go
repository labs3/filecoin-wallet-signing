package msig

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/specs-actors/v8/actors/builtin/miner"
	"github.com/filecoin-project/specs-actors/v8/actors/builtin/multisig"
	"github.com/filecoin-project/specs-actors/v8/actors/util/adt"
	cbor "github.com/ipfs/go-ipld-cbor"
	"github.com/spf13/cobra"
	"golang.org/x/xerrors"

	"github.com/labs3/filecoin-wallet-signing/chain/blockstore"
	"github.com/labs3/filecoin-wallet-signing/internal"
	"github.com/labs3/filecoin-wallet-signing/pkg"
)

// inspectCmd represents the msiginspect command
var inspectCmd = &cobra.Command{
	Use:   "inspect <multisigAddress> ",
	Short: "inspect multisigAddress ",

	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 1 {
			_ = cmd.Help()
			return
		}

		//mtsaddr, err := address.NewFromString("t2i35vaqpkqpx3rcmqpttayaa3k4b7qm2fgrqiq3q")
		mtsaddr, err := address.NewFromString(args[0])
		if err != nil {
			fmt.Println("decode multisigAddress failed:", err.Error())
			return
		}

		if mtsaddr.Protocol() != address.Actor && mtsaddr.Protocol() != address.ID {
			fmt.Println("please input a correct multisigAddress")
			return
		}

		multisigID, err := internal.Lapi.StateLookupID(internal.Ctx, mtsaddr, *internal.CurrentTsk)
		if err != nil {
			fmt.Println("get address ID failed:", err.Error())
			return
		}

		fmt.Printf("Address: %s, ID: %s \n", mtsaddr.String(), multisigID.String())

		a, err := internal.Lapi.StateGetActor(internal.Ctx, mtsaddr, *internal.CurrentTsk)
		if err != nil {
			fmt.Println("Failed to get the address information:", err.Error())
			return
		}

		hd, err := internal.Lapi.ChainReadObj(internal.Ctx, a.Head)
		if err != nil {
			fmt.Println("Failed to get the address HEAD:", err.Error())
			return
		}

		var mstate multisig.State

		err = mstate.UnmarshalCBOR(bytes.NewReader(hd))
		if err != nil {
			fmt.Println("unmarshal address state failed:", err.Error())
			return
		}

		fmt.Printf("Number of signatories %v threshold  %v \n", len(mstate.Signers), mstate.NumApprovalsThreshold)
		for _, signer := range mstate.Signers {
			signerAddr, err := internal.Lapi.StateAccountKey(internal.Ctx, signer, *internal.CurrentTsk)
			if err != nil {
				fmt.Println("get singer of multisigAddress failed : ", err.Error())
				return
			}
			fmt.Printf("%s : %s \n", signer.String(), signerAddr.String())
		}

		store := adt.WrapStore(internal.Ctx, cbor.NewCborStore(blockstore.NewAPIBlockstore(internal.Lapi)))

		arr, err := adt.AsMap(store, mstate.PendingTxns, 5)
		if err != nil {
			fmt.Println("map address pending transaction failed:", err.Error())
			return
		}
		ks, err := arr.CollectKeys()
		if err != nil {
			fmt.Println("Collect address pending transaction failed:", err.Error())
			return
		}
		if len(ks) == 0 {
			fmt.Println("No pending transactions")
			return
		}
		fmt.Println("Pending transaction: ")
		var out multisig.Transaction
		err = arr.ForEach(&out, func(key string) error {
			txid, n := binary.Varint([]byte(key))
			if n <= 0 {
				return xerrors.Errorf("invalid pending transaction key: %v", key)
			}
			p := ""
			msg := ""
			var mwdp miner.WithdrawBalanceParams
			msg = "send out"
			if out.Method == 16 {
				err = mwdp.UnmarshalCBOR(bytes.NewReader(out.Params))
				if err != nil {
					fmt.Println("Parameter parsing failed:", err.Error())
					return nil
				}
				b, _ := json.Marshal(mwdp)
				p = string(b)
				msg = fmt.Sprintf("withdraw from miner  %v FIL", pkg.ToFloat64(mwdp.AmountRequested))
			}
			if out.Method == 23 {
				addr := address.Address{}
				err = addr.UnmarshalCBOR(bytes.NewReader(out.Params))
				if err != nil {
					fmt.Println("Parameter parsing failed:", err.Error())
					return nil
				}

				msg = fmt.Sprintf("change miner %v owner is %v ", out.To.String(), addr.String())
			}
			fmt.Printf("pending id: %v , to : %v , method: %v , amount: %v FIL, Params: %s, approved %v, ps: %s \n",
				txid, out.To, out.Method, pkg.ToFloat64(out.Value), p, out.Approved, msg)
			return nil
		})
		if err != nil {
			fmt.Println("get address pinding transation failed:", err.Error())
			return
		}
	},
}
