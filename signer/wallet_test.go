package signer

import (
	"testing"

	"github.com/filecoin-project/go-address"
	"github.com/stretchr/testify/assert"
)

func init() {
	address.CurrentNetwork = address.Testnet
}

func Test_DecodePrivateKey(t *testing.T) {
	blsPrivateKeys := map[string]string{
		"t3rmgdcyj4utrz4mryuveoj5uycnb3t5ggjzeo5e5iyxmiwjg42amvmbioke3epn55aruu3hifirz7gkw5qg3q": "7b2254797065223a22626c73222c22507269766174654b6579223a2275534649664173674d58367a387a5139727276752b67786943316871697441585a79524e68362b79306d493d227d",
		"t3qcqbu5ktaqikhad3t6o4ertiecr4lvduoxq3corpfqa5wato3v4aouxf5pbrwkny2sikgxz6crbmrl4sblaa": "7b2254797065223a22626c73222c22507269766174654b6579223a22695a663446484b7a472f3939446e313136395a74596b69546257504139467a46476b3837315172645432553d227d",
	}

	secp256k1PrivateKeys := map[string]string{
		"t162t6wgg7svrfafrvm4uwpyfklmtz3ydhkqpsaay": "7b2254797065223a22736563703235366b31222c22507269766174654b6579223a2234372f63794d4e68384f757a3450486f3144626972663453376c5a4f59355236775031616f5771763834733d227d",
		"t1na3ksu4k2g4ofurzufjhn5vcfd5625pops3umsy": "7b2254797065223a22736563703235366b31222c22507269766174654b6579223a22414d4b31727a52536b537070526b31566157583561784e34734270307a2f4955505a546a44386b4b6365493d227d",
		"t1esyzcpgfc3komlbgvjozjonjthmvkjnybovzgaq": "7b2254797065223a22736563703235366b31222c22507269766174654b6579223a2244374267505a6b4e42586c7364766b6c45355862546d534f416342556c672f6c72515a55705a2f496435413d227d",
		"t1huy4zdjez3optamvo4ccllutjsqzzhqti6tpf5i": "7b2254797065223a22736563703235366b31222c22507269766174654b6579223a224a6a4e3766496e33577145473270434955524d65596c4359547274366c325430764762526b54426d35754d3d227d",
		"t1jnudagll3jmldid5eg2nxtwjmpxts7q5zf4drni": "7b2254797065223a22736563703235366b31222c22507269766174654b6579223a226f6d6b4279516d57626753663952554135356b78715962395a30634a5857427961576b4a4f4e6c67536e413d227d",
		"t1shartfibmku6dnmtfa4x6263g5jm4wpxy2ilw2y": "7b2254797065223a22736563703235366b31222c22507269766174654b6579223a225761356e35707136315a5263515653367437596d5635303033465644563973376766682f68373455634e773d227d",
	}

	for k, v := range blsPrivateKeys {
		privateKey, err := DecodePricateKey(v)
		assert.Nil(t, err)
		assert.Equal(t, privateKey.Address.String(), k)
	}

	for k, v := range secp256k1PrivateKeys {
		privateKey, err := DecodePricateKey(v)
		assert.Nil(t, err)
		assert.Equal(t, privateKey.Address.String(), k)
	}
}

func Test_secp256k1_(t *testing.T) {

	sger := NewSecp256k1Singer()
	priv, err := sger.GenPrivate()
	assert.Nil(t, err)
	pubk, err := sger.ToPublic(priv)
	assert.Nil(t, err)
	addr, err := address.NewSecp256k1Address(pubk)
	assert.Nil(t, err)
	t.Log("address: ", addr.String())

	msgs := []string{
		"",
		"a",
		"aaaczxcasdasda",
		"f6as0d+6 a+sd8+-as/d/ -0as+d5 +a9sd08 +a0sd+ +asd ",
		"asdasa-s*f341 -+0 d9+as0 d10GFA+-*AS07/7* +/-70[ *] +A	GG	GFA	*+/faASDsdasd-as/f-*-asGdAas-dSFAaGs+d-*as.0]-0=GSDFAG983940ir1/'.a';ksdDSFSAFAA",
		"ASG三大发😊😋😎😍😘生地0+a苏打水：> 😀😁😂😃😄😅😆😉😗😙😚😇😐😑😶😏😣😥😮😯😪😫😴😌😛😜😝😒😓😔😕😲😷😖😞😟😤😢😭😦😧😨😬",
	}

	for _, v := range msgs {
		msgSign, err := sger.Sign(priv, []byte(v))
		assert.Nil(t, err)
		err = sger.Verify(msgSign, addr, []byte(v))
		assert.Nil(t, err)
	}

}

func Test_BLS(t *testing.T) {
	sger := NewBLSSinger()
	priv, err := sger.GenPrivate()
	assert.Nil(t, err)
	pubk, err := sger.ToPublic(priv)
	assert.Nil(t, err)
	addr, err := address.NewBLSAddress(pubk)
	assert.Nil(t, err)
	t.Log("address: ", addr.String())

	msgs := []string{
		"",
		"a",
		"aaaczxcasdasda",
		"f6as0d+6 a+sd8+-as/d/ -0as+d5 +a9sd08 +a0sd+ +asd ",
		"asdasa-s*f341 -+0 d9+as0 d10GFA+-*AS07/7* +/-70[ *] +A	GG	GFA	*+/😣😥😮😯😪dasd-as/f-*-asGdAas-dSFAaGs+d-*as.0]-0=GSDFAG983940ir1/'.a';ksdDSFSAFAA41 -+0 d9+as0 d10GFA+-*AS07/7* +/-70[ *] +A\tGG\tGFA\t*+/😣😥😮😯😪dasd-as/f-*-asGdAas-dSFAaGs+d-*as.0]-0=GSDFAG",
		"ASG三大发😊😋😎😍😘生地0+a苏打水：> 😀😁😂😃😄😅😆😉😗😙😚😇😐😑😶😏😫😴😌😛😜😝😒😓😔😕😲😷😖😞😟😤😢😭😦😧😨😬",
	}

	for _, v := range msgs {
		msgSign, err := sger.Sign(priv, []byte(v))
		assert.Nil(t, err)
		err = sger.Verify(msgSign, addr, []byte(v))
		assert.Nil(t, err)
	}
}
