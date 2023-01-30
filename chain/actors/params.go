package actors

import (
	"bytes"

	"github.com/filecoin-project/go-state-types/exitcode"
	cbg "github.com/whyrusleeping/cbor-gen"

	"github.com/labs3/filecoin-wallet-signing/chain/types"
)

func SerializeParams(i cbg.CBORMarshaler) ([]byte, types.ActorError) {
	buf := new(bytes.Buffer)
	if err := i.MarshalCBOR(buf); err != nil {
		// TODO: shouldnt this be a fatal error?
		return nil, types.Absorb(err, exitcode.ErrSerialization, "failed to encode parameter")
	}
	return buf.Bytes(), nil
}
