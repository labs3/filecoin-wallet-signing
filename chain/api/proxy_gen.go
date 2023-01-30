package api

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-bitfield"
	datatransfer "github.com/filecoin-project/go-data-transfer"
	"github.com/filecoin-project/go-fil-markets/retrievalmarket"
	"github.com/filecoin-project/go-fil-markets/storagemarket"
	"github.com/filecoin-project/go-jsonrpc/auth"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/filecoin-project/go-state-types/builtin/v9/market"
	"github.com/filecoin-project/go-state-types/builtin/v9/miner"
	"github.com/filecoin-project/go-state-types/builtin/v9/paych"
	verifregtypes "github.com/filecoin-project/go-state-types/builtin/v9/verifreg"
	"github.com/filecoin-project/go-state-types/crypto"
	"github.com/filecoin-project/go-state-types/dline"
	abinetwork "github.com/filecoin-project/go-state-types/network"
	"github.com/filecoin-project/go-state-types/proof"
	"github.com/google/uuid"
	blocks "github.com/ipfs/go-block-format"
	"github.com/ipfs/go-cid"
	"github.com/libp2p/go-libp2p/core/metrics"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	"golang.org/x/xerrors"

	"github.com/labs3/filecoin-wallet-signing/chain/types"
)

var ErrNotSupported = xerrors.New("method not supported")

type CommonStruct struct {
	Internal struct {
		AuthNew func(p0 context.Context, p1 []auth.Permission) ([]byte, error) `perm:"admin"`

		AuthVerify func(p0 context.Context, p1 string) ([]auth.Permission, error) `perm:"read"`

		Closing func(p0 context.Context) (<-chan struct{}, error) `perm:"read"`

		Discover func(p0 context.Context) (OpenRPCDocument, error) `perm:"read"`

		LogAlerts func(p0 context.Context) ([]Alert, error) `perm:"admin"`

		LogList func(p0 context.Context) ([]string, error) `perm:"write"`

		LogSetLevel func(p0 context.Context, p1 string, p2 string) error `perm:"write"`

		Session func(p0 context.Context) (uuid.UUID, error) `perm:"read"`

		Shutdown func(p0 context.Context) error `perm:"admin"`

		StartTime func(p0 context.Context) (time.Time, error) `perm:"read"`

		Version func(p0 context.Context) (APIVersion, error) `perm:"read"`
	}
}

type FullNodeStruct struct {
	CommonStruct

	NetStruct

	Internal struct {
		ChainBlockstoreInfo func(p0 context.Context) (map[string]interface{}, error) `perm:"read"`

		ChainCheckBlockstore func(p0 context.Context) error `perm:"admin"`

		ChainDeleteObj func(p0 context.Context, p1 cid.Cid) error `perm:"admin"`

		ChainExport func(p0 context.Context, p1 abi.ChainEpoch, p2 bool, p3 types.TipSetKey) (<-chan []byte, error) `perm:"read"`

		ChainGetBlock func(p0 context.Context, p1 cid.Cid) (*types.BlockHeader, error) `perm:"read"`

		ChainGetBlockMessages func(p0 context.Context, p1 cid.Cid) (*BlockMessages, error) `perm:"read"`

		ChainGetGenesis func(p0 context.Context) (*types.TipSet, error) `perm:"read"`

		ChainGetMessage func(p0 context.Context, p1 cid.Cid) (*types.Message, error) `perm:"read"`

		ChainGetMessagesInTipset func(p0 context.Context, p1 types.TipSetKey) ([]Message, error) `perm:"read"`

		ChainGetNode func(p0 context.Context, p1 string) (*IpldObject, error) `perm:"read"`

		ChainGetParentMessages func(p0 context.Context, p1 cid.Cid) ([]Message, error) `perm:"read"`

		ChainGetParentReceipts func(p0 context.Context, p1 cid.Cid) ([]*types.MessageReceipt, error) `perm:"read"`

		ChainGetPath func(p0 context.Context, p1 types.TipSetKey, p2 types.TipSetKey) ([]*HeadChange, error) `perm:"read"`

		ChainGetTipSet func(p0 context.Context, p1 types.TipSetKey) (*types.TipSet, error) `perm:"read"`

		ChainGetTipSetAfterHeight func(p0 context.Context, p1 abi.ChainEpoch, p2 types.TipSetKey) (*types.TipSet, error) `perm:"read"`

		ChainGetTipSetByHeight func(p0 context.Context, p1 abi.ChainEpoch, p2 types.TipSetKey) (*types.TipSet, error) `perm:"read"`

		ChainHasObj func(p0 context.Context, p1 cid.Cid) (bool, error) `perm:"read"`

		ChainHead func(p0 context.Context) (*types.TipSet, error) `perm:"read"`

		ChainNotify func(p0 context.Context) (<-chan []*HeadChange, error) `perm:"read"`

		ChainPrune func(p0 context.Context, p1 PruneOpts) error `perm:"admin"`

		ChainPutObj func(p0 context.Context, p1 blocks.Block) error `perm:"admin"`

		ChainReadObj func(p0 context.Context, p1 cid.Cid) ([]byte, error) `perm:"read"`

		ChainSetHead func(p0 context.Context, p1 types.TipSetKey) error `perm:"admin"`

		ChainStatObj func(p0 context.Context, p1 cid.Cid, p2 cid.Cid) (ObjStat, error) `perm:"read"`

		ChainTipSetWeight func(p0 context.Context, p1 types.TipSetKey) (types.BigInt, error) `perm:"read"`

		ClientCalcCommP func(p0 context.Context, p1 string) (*CommPRet, error) `perm:"write"`

		ClientCancelDataTransfer func(p0 context.Context, p1 datatransfer.TransferID, p2 peer.ID, p3 bool) error `perm:"write"`

		ClientCancelRetrievalDeal func(p0 context.Context, p1 retrievalmarket.DealID) error `perm:"write"`

		ClientDataTransferUpdates func(p0 context.Context) (<-chan DataTransferChannel, error) `perm:"write"`

		ClientDealPieceCID func(p0 context.Context, p1 cid.Cid) (DataCIDSize, error) `perm:"read"`

		ClientDealSize func(p0 context.Context, p1 cid.Cid) (DataSize, error) `perm:"read"`

		ClientExport func(p0 context.Context, p1 ExportRef, p2 FileRef) error `perm:"admin"`

		ClientFindData func(p0 context.Context, p1 cid.Cid, p2 *cid.Cid) ([]QueryOffer, error) `perm:"read"`

		ClientGenCar func(p0 context.Context, p1 FileRef, p2 string) error `perm:"write"`

		ClientGetDealInfo func(p0 context.Context, p1 cid.Cid) (*DealInfo, error) `perm:"read"`

		ClientGetDealStatus func(p0 context.Context, p1 uint64) (string, error) `perm:"read"`

		ClientGetDealUpdates func(p0 context.Context) (<-chan DealInfo, error) `perm:"write"`

		ClientGetRetrievalUpdates func(p0 context.Context) (<-chan RetrievalInfo, error) `perm:"write"`

		ClientHasLocal func(p0 context.Context, p1 cid.Cid) (bool, error) `perm:"write"`

		ClientImport func(p0 context.Context, p1 FileRef) (*ImportRes, error) `perm:"admin"`

		ClientListDataTransfers func(p0 context.Context) ([]DataTransferChannel, error) `perm:"write"`

		ClientListDeals func(p0 context.Context) ([]DealInfo, error) `perm:"write"`

		ClientListImports func(p0 context.Context) ([]Import, error) `perm:"write"`

		ClientListRetrievals func(p0 context.Context) ([]RetrievalInfo, error) `perm:"write"`

		ClientMinerQueryOffer func(p0 context.Context, p1 address.Address, p2 cid.Cid, p3 *cid.Cid) (QueryOffer, error) `perm:"read"`

		ClientQueryAsk func(p0 context.Context, p1 peer.ID, p2 address.Address) (*StorageAsk, error) `perm:"read"`

		ClientRemoveImport func(p0 context.Context, p1 ID) error `perm:"admin"`

		ClientRestartDataTransfer func(p0 context.Context, p1 datatransfer.TransferID, p2 peer.ID, p3 bool) error `perm:"write"`

		ClientRetrieve func(p0 context.Context, p1 RetrievalOrder) (*RestrievalRes, error) `perm:"admin"`

		ClientRetrieveTryRestartInsufficientFunds func(p0 context.Context, p1 address.Address) error `perm:"write"`

		ClientRetrieveWait func(p0 context.Context, p1 retrievalmarket.DealID) error `perm:"admin"`

		ClientStartDeal func(p0 context.Context, p1 *StartDealParams) (*cid.Cid, error) `perm:"admin"`

		ClientStatelessDeal func(p0 context.Context, p1 *StartDealParams) (*cid.Cid, error) `perm:"write"`

		CreateBackup func(p0 context.Context, p1 string) error `perm:"admin"`

		GasEstimateFeeCap func(p0 context.Context, p1 *types.Message, p2 int64, p3 types.TipSetKey) (types.BigInt, error) `perm:"read"`

		GasEstimateGasLimit func(p0 context.Context, p1 *types.Message, p2 types.TipSetKey) (int64, error) `perm:"read"`

		GasEstimateGasPremium func(p0 context.Context, p1 uint64, p2 address.Address, p3 int64, p4 types.TipSetKey) (types.BigInt, error) `perm:"read"`

		GasEstimateMessageGas func(p0 context.Context, p1 *types.Message, p2 *MessageSendSpec, p3 types.TipSetKey) (*types.Message, error) `perm:"read"`

		MarketAddBalance func(p0 context.Context, p1 address.Address, p2 address.Address, p3 types.BigInt) (cid.Cid, error) `perm:"sign"`

		MarketGetReserved func(p0 context.Context, p1 address.Address) (types.BigInt, error) `perm:"sign"`

		MarketReleaseFunds func(p0 context.Context, p1 address.Address, p2 types.BigInt) error `perm:"sign"`

		MarketReserveFunds func(p0 context.Context, p1 address.Address, p2 address.Address, p3 types.BigInt) (cid.Cid, error) `perm:"sign"`

		MarketWithdraw func(p0 context.Context, p1 address.Address, p2 address.Address, p3 types.BigInt) (cid.Cid, error) `perm:"sign"`

		MinerCreateBlock func(p0 context.Context, p1 *BlockTemplate) (*types.BlockMsg, error) `perm:"write"`

		MinerGetBaseInfo func(p0 context.Context, p1 address.Address, p2 abi.ChainEpoch, p3 types.TipSetKey) (*MiningBaseInfo, error) `perm:"read"`

		MpoolBatchPush func(p0 context.Context, p1 []*types.SignedMessage) ([]cid.Cid, error) `perm:"write"`

		MpoolBatchPushMessage func(p0 context.Context, p1 []*types.Message, p2 *MessageSendSpec) ([]*types.SignedMessage, error) `perm:"sign"`

		MpoolBatchPushUntrusted func(p0 context.Context, p1 []*types.SignedMessage) ([]cid.Cid, error) `perm:"write"`

		MpoolCheckMessages func(p0 context.Context, p1 []*MessagePrototype) ([][]MessageCheckStatus, error) `perm:"read"`

		MpoolCheckPendingMessages func(p0 context.Context, p1 address.Address) ([][]MessageCheckStatus, error) `perm:"read"`

		MpoolCheckReplaceMessages func(p0 context.Context, p1 []*types.Message) ([][]MessageCheckStatus, error) `perm:"read"`

		MpoolClear func(p0 context.Context, p1 bool) error `perm:"write"`

		MpoolGetConfig func(p0 context.Context) (*types.MpoolConfig, error) `perm:"read"`

		MpoolGetNonce func(p0 context.Context, p1 address.Address) (uint64, error) `perm:"read"`

		MpoolPending func(p0 context.Context, p1 types.TipSetKey) ([]*types.SignedMessage, error) `perm:"read"`

		MpoolPush func(p0 context.Context, p1 *types.SignedMessage) (cid.Cid, error) `perm:"write"`

		MpoolPushMessage func(p0 context.Context, p1 *types.Message, p2 *MessageSendSpec) (*types.SignedMessage, error) `perm:"sign"`

		MpoolPushUntrusted func(p0 context.Context, p1 *types.SignedMessage) (cid.Cid, error) `perm:"write"`

		MpoolSelect func(p0 context.Context, p1 types.TipSetKey, p2 float64) ([]*types.SignedMessage, error) `perm:"read"`

		MpoolSetConfig func(p0 context.Context, p1 *types.MpoolConfig) error `perm:"admin"`

		MpoolSub func(p0 context.Context) (<-chan MpoolUpdate, error) `perm:"read"`

		MsigAddApprove func(p0 context.Context, p1 address.Address, p2 address.Address, p3 uint64, p4 address.Address, p5 address.Address, p6 bool) (*MessagePrototype, error) `perm:"sign"`

		MsigAddCancel func(p0 context.Context, p1 address.Address, p2 address.Address, p3 uint64, p4 address.Address, p5 bool) (*MessagePrototype, error) `perm:"sign"`

		MsigAddPropose func(p0 context.Context, p1 address.Address, p2 address.Address, p3 address.Address, p4 bool) (*MessagePrototype, error) `perm:"sign"`

		MsigApprove func(p0 context.Context, p1 address.Address, p2 uint64, p3 address.Address) (*MessagePrototype, error) `perm:"sign"`

		MsigApproveTxnHash func(p0 context.Context, p1 address.Address, p2 uint64, p3 address.Address, p4 address.Address, p5 types.BigInt, p6 address.Address, p7 uint64, p8 []byte) (*MessagePrototype, error) `perm:"sign"`

		MsigCancel func(p0 context.Context, p1 address.Address, p2 uint64, p3 address.Address) (*MessagePrototype, error) `perm:"sign"`

		MsigCancelTxnHash func(p0 context.Context, p1 address.Address, p2 uint64, p3 address.Address, p4 types.BigInt, p5 address.Address, p6 uint64, p7 []byte) (*MessagePrototype, error) `perm:"sign"`

		MsigCreate func(p0 context.Context, p1 uint64, p2 []address.Address, p3 abi.ChainEpoch, p4 types.BigInt, p5 address.Address, p6 types.BigInt) (*MessagePrototype, error) `perm:"sign"`

		MsigGetAvailableBalance func(p0 context.Context, p1 address.Address, p2 types.TipSetKey) (types.BigInt, error) `perm:"read"`

		MsigGetPending func(p0 context.Context, p1 address.Address, p2 types.TipSetKey) ([]*MsigTransaction, error) `perm:"read"`

		MsigGetVested func(p0 context.Context, p1 address.Address, p2 types.TipSetKey, p3 types.TipSetKey) (types.BigInt, error) `perm:"read"`

		MsigGetVestingSchedule func(p0 context.Context, p1 address.Address, p2 types.TipSetKey) (MsigVesting, error) `perm:"read"`

		MsigPropose func(p0 context.Context, p1 address.Address, p2 address.Address, p3 types.BigInt, p4 address.Address, p5 uint64, p6 []byte) (*MessagePrototype, error) `perm:"sign"`

		MsigRemoveSigner func(p0 context.Context, p1 address.Address, p2 address.Address, p3 address.Address, p4 bool) (*MessagePrototype, error) `perm:"sign"`

		MsigSwapApprove func(p0 context.Context, p1 address.Address, p2 address.Address, p3 uint64, p4 address.Address, p5 address.Address, p6 address.Address) (*MessagePrototype, error) `perm:"sign"`

		MsigSwapCancel func(p0 context.Context, p1 address.Address, p2 address.Address, p3 uint64, p4 address.Address, p5 address.Address) (*MessagePrototype, error) `perm:"sign"`

		MsigSwapPropose func(p0 context.Context, p1 address.Address, p2 address.Address, p3 address.Address, p4 address.Address) (*MessagePrototype, error) `perm:"sign"`

		NodeStatus func(p0 context.Context, p1 bool) (NodeStatus, error) `perm:"read"`

		PaychAllocateLane func(p0 context.Context, p1 address.Address) (uint64, error) `perm:"sign"`

		PaychAvailableFunds func(p0 context.Context, p1 address.Address) (*ChannelAvailableFunds, error) `perm:"sign"`

		PaychAvailableFundsByFromTo func(p0 context.Context, p1 address.Address, p2 address.Address) (*ChannelAvailableFunds, error) `perm:"sign"`

		PaychCollect func(p0 context.Context, p1 address.Address) (cid.Cid, error) `perm:"sign"`

		PaychFund func(p0 context.Context, p1 address.Address, p2 address.Address, p3 types.BigInt) (*ChannelInfo, error) `perm:"sign"`

		PaychGet func(p0 context.Context, p1 address.Address, p2 address.Address, p3 types.BigInt, p4 PaychGetOpts) (*ChannelInfo, error) `perm:"sign"`

		PaychGetWaitReady func(p0 context.Context, p1 cid.Cid) (address.Address, error) `perm:"sign"`

		PaychList func(p0 context.Context) ([]address.Address, error) `perm:"read"`

		PaychNewPayment func(p0 context.Context, p1 address.Address, p2 address.Address, p3 []VoucherSpec) (*PaymentInfo, error) `perm:"sign"`

		PaychSettle func(p0 context.Context, p1 address.Address) (cid.Cid, error) `perm:"sign"`

		PaychStatus func(p0 context.Context, p1 address.Address) (*PaychStatus, error) `perm:"read"`

		PaychVoucherAdd func(p0 context.Context, p1 address.Address, p2 *paych.SignedVoucher, p3 []byte, p4 types.BigInt) (types.BigInt, error) `perm:"write"`

		PaychVoucherCheckSpendable func(p0 context.Context, p1 address.Address, p2 *paych.SignedVoucher, p3 []byte, p4 []byte) (bool, error) `perm:"read"`

		PaychVoucherCheckValid func(p0 context.Context, p1 address.Address, p2 *paych.SignedVoucher) error `perm:"read"`

		PaychVoucherCreate func(p0 context.Context, p1 address.Address, p2 types.BigInt, p3 uint64) (*VoucherCreateResult, error) `perm:"sign"`

		PaychVoucherList func(p0 context.Context, p1 address.Address) ([]*paych.SignedVoucher, error) `perm:"write"`

		PaychVoucherSubmit func(p0 context.Context, p1 address.Address, p2 *paych.SignedVoucher, p3 []byte, p4 []byte) (cid.Cid, error) `perm:"sign"`

		RaftLeader func(p0 context.Context) (peer.ID, error) `perm:"read"`

		RaftState func(p0 context.Context) (*RaftStateData, error) `perm:"read"`

		StateAccountKey func(p0 context.Context, p1 address.Address, p2 types.TipSetKey) (address.Address, error) `perm:"read"`

		StateActorCodeCIDs func(p0 context.Context, p1 abinetwork.Version) (map[string]cid.Cid, error) `perm:"read"`

		StateActorManifestCID func(p0 context.Context, p1 abinetwork.Version) (cid.Cid, error) `perm:"read"`

		StateAllMinerFaults func(p0 context.Context, p1 abi.ChainEpoch, p2 types.TipSetKey) ([]*Fault, error) `perm:"read"`

		StateCall func(p0 context.Context, p1 *types.Message, p2 types.TipSetKey) (*InvocResult, error) `perm:"read"`

		StateChangedActors func(p0 context.Context, p1 cid.Cid, p2 cid.Cid) (map[string]types.Actor, error) `perm:"read"`

		StateCirculatingSupply func(p0 context.Context, p1 types.TipSetKey) (abi.TokenAmount, error) `perm:"read"`

		StateCompute func(p0 context.Context, p1 abi.ChainEpoch, p2 []*types.Message, p3 types.TipSetKey) (*ComputeStateOutput, error) `perm:"read"`

		StateComputeDataCID func(p0 context.Context, p1 address.Address, p2 abi.RegisteredSealProof, p3 []abi.DealID, p4 types.TipSetKey) (cid.Cid, error) `perm:"read"`

		StateDealProviderCollateralBounds func(p0 context.Context, p1 abi.PaddedPieceSize, p2 bool, p3 types.TipSetKey) (DealCollateralBounds, error) `perm:"read"`

		StateDecodeParams func(p0 context.Context, p1 address.Address, p2 abi.MethodNum, p3 []byte, p4 types.TipSetKey) (interface{}, error) `perm:"read"`

		StateEncodeParams func(p0 context.Context, p1 cid.Cid, p2 abi.MethodNum, p3 json.RawMessage) ([]byte, error) `perm:"read"`

		StateGetActor func(p0 context.Context, p1 address.Address, p2 types.TipSetKey) (*types.Actor, error) `perm:"read"`

		StateGetAllocation func(p0 context.Context, p1 address.Address, p2 verifregtypes.AllocationId, p3 types.TipSetKey) (*verifregtypes.Allocation, error) `perm:"read"`

		StateGetAllocationForPendingDeal func(p0 context.Context, p1 abi.DealID, p2 types.TipSetKey) (*verifregtypes.Allocation, error) `perm:"read"`

		StateGetAllocations func(p0 context.Context, p1 address.Address, p2 types.TipSetKey) (map[verifregtypes.AllocationId]verifregtypes.Allocation, error) `perm:"read"`

		StateGetBeaconEntry func(p0 context.Context, p1 abi.ChainEpoch) (*types.BeaconEntry, error) `perm:"read"`

		StateGetClaim func(p0 context.Context, p1 address.Address, p2 verifregtypes.ClaimId, p3 types.TipSetKey) (*verifregtypes.Claim, error) `perm:"read"`

		StateGetClaims func(p0 context.Context, p1 address.Address, p2 types.TipSetKey) (map[verifregtypes.ClaimId]verifregtypes.Claim, error) `perm:"read"`

		StateGetNetworkParams func(p0 context.Context) (*NetworkParams, error) `perm:"read"`

		StateGetRandomnessFromBeacon func(p0 context.Context, p1 crypto.DomainSeparationTag, p2 abi.ChainEpoch, p3 []byte, p4 types.TipSetKey) (abi.Randomness, error) `perm:"read"`

		StateGetRandomnessFromTickets func(p0 context.Context, p1 crypto.DomainSeparationTag, p2 abi.ChainEpoch, p3 []byte, p4 types.TipSetKey) (abi.Randomness, error) `perm:"read"`

		StateListActors func(p0 context.Context, p1 types.TipSetKey) ([]address.Address, error) `perm:"read"`

		StateListMessages func(p0 context.Context, p1 *MessageMatch, p2 types.TipSetKey, p3 abi.ChainEpoch) ([]cid.Cid, error) `perm:"read"`

		StateListMiners func(p0 context.Context, p1 types.TipSetKey) ([]address.Address, error) `perm:"read"`

		StateLookupID func(p0 context.Context, p1 address.Address, p2 types.TipSetKey) (address.Address, error) `perm:"read"`

		StateLookupRobustAddress func(p0 context.Context, p1 address.Address, p2 types.TipSetKey) (address.Address, error) `perm:"read"`

		StateMarketBalance func(p0 context.Context, p1 address.Address, p2 types.TipSetKey) (MarketBalance, error) `perm:"read"`

		StateMarketDeals func(p0 context.Context, p1 types.TipSetKey) (map[string]*MarketDeal, error) `perm:"read"`

		StateMarketParticipants func(p0 context.Context, p1 types.TipSetKey) (map[string]MarketBalance, error) `perm:"read"`

		StateMarketStorageDeal func(p0 context.Context, p1 abi.DealID, p2 types.TipSetKey) (*MarketDeal, error) `perm:"read"`

		StateMinerActiveSectors func(p0 context.Context, p1 address.Address, p2 types.TipSetKey) ([]*miner.SectorOnChainInfo, error) `perm:"read"`

		StateMinerAllocated func(p0 context.Context, p1 address.Address, p2 types.TipSetKey) (*bitfield.BitField, error) `perm:"read"`

		StateMinerAvailableBalance func(p0 context.Context, p1 address.Address, p2 types.TipSetKey) (types.BigInt, error) `perm:"read"`

		StateMinerDeadlines func(p0 context.Context, p1 address.Address, p2 types.TipSetKey) ([]Deadline, error) `perm:"read"`

		StateMinerFaults func(p0 context.Context, p1 address.Address, p2 types.TipSetKey) (bitfield.BitField, error) `perm:"read"`

		StateMinerInfo func(p0 context.Context, p1 address.Address, p2 types.TipSetKey) (MinerInfo, error) `perm:"read"`

		StateMinerInitialPledgeCollateral func(p0 context.Context, p1 address.Address, p2 miner.SectorPreCommitInfo, p3 types.TipSetKey) (types.BigInt, error) `perm:"read"`

		StateMinerPartitions func(p0 context.Context, p1 address.Address, p2 uint64, p3 types.TipSetKey) ([]Partition, error) `perm:"read"`

		StateMinerPower func(p0 context.Context, p1 address.Address, p2 types.TipSetKey) (*MinerPower, error) `perm:"read"`

		StateMinerPreCommitDepositForPower func(p0 context.Context, p1 address.Address, p2 miner.SectorPreCommitInfo, p3 types.TipSetKey) (types.BigInt, error) `perm:"read"`

		StateMinerProvingDeadline func(p0 context.Context, p1 address.Address, p2 types.TipSetKey) (*dline.Info, error) `perm:"read"`

		StateMinerRecoveries func(p0 context.Context, p1 address.Address, p2 types.TipSetKey) (bitfield.BitField, error) `perm:"read"`

		StateMinerSectorAllocated func(p0 context.Context, p1 address.Address, p2 abi.SectorNumber, p3 types.TipSetKey) (bool, error) `perm:"read"`

		StateMinerSectorCount func(p0 context.Context, p1 address.Address, p2 types.TipSetKey) (MinerSectors, error) `perm:"read"`

		StateMinerSectors func(p0 context.Context, p1 address.Address, p2 *bitfield.BitField, p3 types.TipSetKey) ([]*miner.SectorOnChainInfo, error) `perm:"read"`

		StateNetworkName func(p0 context.Context) (NetworkName, error) `perm:"read"`

		StateNetworkVersion func(p0 context.Context, p1 types.TipSetKey) (NetworkVersion, error) `perm:"read"`

		StateReadState func(p0 context.Context, p1 address.Address, p2 types.TipSetKey) (*ActorState, error) `perm:"read"`

		StateReplay func(p0 context.Context, p1 types.TipSetKey, p2 cid.Cid) (*InvocResult, error) `perm:"read"`

		StateSearchMsg func(p0 context.Context, p1 types.TipSetKey, p2 cid.Cid, p3 abi.ChainEpoch, p4 bool) (*MsgLookup, error) `perm:"read"`

		StateSectorExpiration func(p0 context.Context, p1 address.Address, p2 abi.SectorNumber, p3 types.TipSetKey) (*SectorExpiration, error) `perm:"read"`

		StateSectorGetInfo func(p0 context.Context, p1 address.Address, p2 abi.SectorNumber, p3 types.TipSetKey) (*miner.SectorOnChainInfo, error) `perm:"read"`

		StateSectorPartition func(p0 context.Context, p1 address.Address, p2 abi.SectorNumber, p3 types.TipSetKey) (*SectorLocation, error) `perm:"read"`

		StateSectorPreCommitInfo func(p0 context.Context, p1 address.Address, p2 abi.SectorNumber, p3 types.TipSetKey) (*miner.SectorPreCommitOnChainInfo, error) `perm:"read"`

		StateVMCirculatingSupplyInternal func(p0 context.Context, p1 types.TipSetKey) (CirculatingSupply, error) `perm:"read"`

		StateVerifiedClientStatus func(p0 context.Context, p1 address.Address, p2 types.TipSetKey) (*abi.StoragePower, error) `perm:"read"`

		StateVerifiedRegistryRootKey func(p0 context.Context, p1 types.TipSetKey) (address.Address, error) `perm:"read"`

		StateVerifierStatus func(p0 context.Context, p1 address.Address, p2 types.TipSetKey) (*abi.StoragePower, error) `perm:"read"`

		StateWaitMsg func(p0 context.Context, p1 cid.Cid, p2 uint64, p3 abi.ChainEpoch, p4 bool) (*MsgLookup, error) `perm:"read"`

		SyncCheckBad func(p0 context.Context, p1 cid.Cid) (string, error) `perm:"read"`

		SyncCheckpoint func(p0 context.Context, p1 types.TipSetKey) error `perm:"admin"`

		SyncIncomingBlocks func(p0 context.Context) (<-chan *types.BlockHeader, error) `perm:"read"`

		SyncMarkBad func(p0 context.Context, p1 cid.Cid) error `perm:"admin"`

		SyncState func(p0 context.Context) (*SyncState, error) `perm:"read"`

		SyncSubmitBlock func(p0 context.Context, p1 *types.BlockMsg) error `perm:"write"`

		SyncUnmarkAllBad func(p0 context.Context) error `perm:"admin"`

		SyncUnmarkBad func(p0 context.Context, p1 cid.Cid) error `perm:"admin"`

		SyncValidateTipset func(p0 context.Context, p1 types.TipSetKey) (bool, error) `perm:"read"`

		WalletBalance func(p0 context.Context, p1 address.Address) (types.BigInt, error) `perm:"read"`

		WalletDefaultAddress func(p0 context.Context) (address.Address, error) `perm:"write"`

		WalletDelete func(p0 context.Context, p1 address.Address) error `perm:"admin"`

		WalletExport func(p0 context.Context, p1 address.Address) (*types.KeyInfo, error) `perm:"admin"`

		WalletHas func(p0 context.Context, p1 address.Address) (bool, error) `perm:"write"`

		WalletImport func(p0 context.Context, p1 *types.KeyInfo) (address.Address, error) `perm:"admin"`

		WalletList func(p0 context.Context) ([]address.Address, error) `perm:"write"`

		WalletNew func(p0 context.Context, p1 types.KeyType) (address.Address, error) `perm:"write"`

		WalletSetDefault func(p0 context.Context, p1 address.Address) error `perm:"write"`

		WalletSign func(p0 context.Context, p1 address.Address, p2 []byte) (*crypto.Signature, error) `perm:"sign"`

		WalletSignMessage func(p0 context.Context, p1 address.Address, p2 *types.Message) (*types.SignedMessage, error) `perm:"sign"`

		WalletValidateAddress func(p0 context.Context, p1 string) (address.Address, error) `perm:"read"`

		WalletVerify func(p0 context.Context, p1 address.Address, p2 []byte, p3 *crypto.Signature) (bool, error) `perm:"read"`
	}
}

type NetStruct struct {
	Internal struct {
		ID func(p0 context.Context) (peer.ID, error) `perm:"read"`

		NetAddrsListen func(p0 context.Context) (peer.AddrInfo, error) `perm:"read"`

		NetAgentVersion func(p0 context.Context, p1 peer.ID) (string, error) `perm:"read"`

		NetAutoNatStatus func(p0 context.Context) (NatInfo, error) `perm:"read"`

		NetBandwidthStats func(p0 context.Context) (metrics.Stats, error) `perm:"read"`

		NetBandwidthStatsByPeer func(p0 context.Context) (map[string]metrics.Stats, error) `perm:"read"`

		NetBandwidthStatsByProtocol func(p0 context.Context) (map[protocol.ID]metrics.Stats, error) `perm:"read"`

		NetBlockAdd func(p0 context.Context, p1 NetBlockList) error `perm:"admin"`

		NetBlockList func(p0 context.Context) (NetBlockList, error) `perm:"read"`

		NetBlockRemove func(p0 context.Context, p1 NetBlockList) error `perm:"admin"`

		NetConnect func(p0 context.Context, p1 peer.AddrInfo) error `perm:"write"`

		NetConnectedness func(p0 context.Context, p1 peer.ID) (network.Connectedness, error) `perm:"read"`

		NetDisconnect func(p0 context.Context, p1 peer.ID) error `perm:"write"`

		NetFindPeer func(p0 context.Context, p1 peer.ID) (peer.AddrInfo, error) `perm:"read"`

		NetLimit func(p0 context.Context, p1 string) (NetLimit, error) `perm:"read"`

		NetPeerInfo func(p0 context.Context, p1 peer.ID) (*ExtendedPeerInfo, error) `perm:"read"`

		NetPeers func(p0 context.Context) ([]peer.AddrInfo, error) `perm:"read"`

		NetPing func(p0 context.Context, p1 peer.ID) (time.Duration, error) `perm:"read"`

		NetProtectAdd func(p0 context.Context, p1 []peer.ID) error `perm:"admin"`

		NetProtectList func(p0 context.Context) ([]peer.ID, error) `perm:"read"`

		NetProtectRemove func(p0 context.Context, p1 []peer.ID) error `perm:"admin"`

		NetPubsubScores func(p0 context.Context) ([]PubsubScore, error) `perm:"read"`

		NetSetLimit func(p0 context.Context, p1 string, p2 NetLimit) error `perm:"admin"`

		NetStat func(p0 context.Context, p1 string) (NetStat, error) `perm:"read"`
	}
}

type SignableStruct struct {
	Internal struct {
		Sign func(p0 context.Context, p1 SignFunc) error ``
	}
}

func (s *CommonStruct) AuthVerify(p0 context.Context, p1 string) ([]auth.Permission, error) {
	if s.Internal.AuthVerify == nil {
		return *new([]auth.Permission), ErrNotSupported
	}
	return s.Internal.AuthVerify(p0, p1)
}

func (s *CommonStruct) Closing(p0 context.Context) (<-chan struct{}, error) {
	if s.Internal.Closing == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.Closing(p0)
}

func (s *CommonStruct) Discover(p0 context.Context) (OpenRPCDocument, error) {
	if s.Internal.Discover == nil {
		return *new(OpenRPCDocument), ErrNotSupported
	}
	return s.Internal.Discover(p0)
}

func (s *CommonStruct) LogAlerts(p0 context.Context) ([]Alert, error) {
	if s.Internal.LogAlerts == nil {
		return *new([]Alert), ErrNotSupported
	}
	return s.Internal.LogAlerts(p0)
}

func (s *CommonStruct) LogList(p0 context.Context) ([]string, error) {
	if s.Internal.LogList == nil {
		return *new([]string), ErrNotSupported
	}
	return s.Internal.LogList(p0)
}

func (s *CommonStruct) LogSetLevel(p0 context.Context, p1 string, p2 string) error {
	if s.Internal.LogSetLevel == nil {
		return ErrNotSupported
	}
	return s.Internal.LogSetLevel(p0, p1, p2)
}

func (s *CommonStruct) Session(p0 context.Context) (uuid.UUID, error) {
	if s.Internal.Session == nil {
		return *new(uuid.UUID), ErrNotSupported
	}
	return s.Internal.Session(p0)
}

func (s *CommonStruct) Shutdown(p0 context.Context) error {
	if s.Internal.Shutdown == nil {
		return ErrNotSupported
	}
	return s.Internal.Shutdown(p0)
}

func (s *CommonStruct) StartTime(p0 context.Context) (time.Time, error) {
	if s.Internal.StartTime == nil {
		return *new(time.Time), ErrNotSupported
	}
	return s.Internal.StartTime(p0)
}

func (s *CommonStruct) Version(p0 context.Context) (APIVersion, error) {
	if s.Internal.Version == nil {
		return *new(APIVersion), ErrNotSupported
	}
	return s.Internal.Version(p0)
}

func (s *FullNodeStruct) ChainBlockstoreInfo(p0 context.Context) (map[string]interface{}, error) {
	if s.Internal.ChainBlockstoreInfo == nil {
		return *new(map[string]interface{}), ErrNotSupported
	}
	return s.Internal.ChainBlockstoreInfo(p0)
}

func (s *FullNodeStruct) ChainCheckBlockstore(p0 context.Context) error {
	if s.Internal.ChainCheckBlockstore == nil {
		return ErrNotSupported
	}
	return s.Internal.ChainCheckBlockstore(p0)
}

func (s *FullNodeStruct) ChainDeleteObj(p0 context.Context, p1 cid.Cid) error {
	if s.Internal.ChainDeleteObj == nil {
		return ErrNotSupported
	}
	return s.Internal.ChainDeleteObj(p0, p1)
}

func (s *FullNodeStruct) ChainExport(p0 context.Context, p1 abi.ChainEpoch, p2 bool, p3 types.TipSetKey) (<-chan []byte, error) {
	if s.Internal.ChainExport == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.ChainExport(p0, p1, p2, p3)
}

func (s *FullNodeStruct) ChainGetBlock(p0 context.Context, p1 cid.Cid) (*types.BlockHeader, error) {
	if s.Internal.ChainGetBlock == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.ChainGetBlock(p0, p1)
}

func (s *FullNodeStruct) ChainGetBlockMessages(p0 context.Context, p1 cid.Cid) (*BlockMessages, error) {
	if s.Internal.ChainGetBlockMessages == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.ChainGetBlockMessages(p0, p1)
}

func (s *FullNodeStruct) ChainGetGenesis(p0 context.Context) (*types.TipSet, error) {
	if s.Internal.ChainGetGenesis == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.ChainGetGenesis(p0)
}

func (s *FullNodeStruct) ChainGetMessage(p0 context.Context, p1 cid.Cid) (*types.Message, error) {
	if s.Internal.ChainGetMessage == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.ChainGetMessage(p0, p1)
}

func (s *FullNodeStruct) ChainGetMessagesInTipset(p0 context.Context, p1 types.TipSetKey) ([]Message, error) {
	if s.Internal.ChainGetMessagesInTipset == nil {
		return *new([]Message), ErrNotSupported
	}
	return s.Internal.ChainGetMessagesInTipset(p0, p1)
}

func (s *FullNodeStruct) ChainGetNode(p0 context.Context, p1 string) (*IpldObject, error) {
	if s.Internal.ChainGetNode == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.ChainGetNode(p0, p1)
}

func (s *FullNodeStruct) ChainGetParentMessages(p0 context.Context, p1 cid.Cid) ([]Message, error) {
	if s.Internal.ChainGetParentMessages == nil {
		return *new([]Message), ErrNotSupported
	}
	return s.Internal.ChainGetParentMessages(p0, p1)
}

func (s *FullNodeStruct) ChainGetParentReceipts(p0 context.Context, p1 cid.Cid) ([]*types.MessageReceipt, error) {
	if s.Internal.ChainGetParentReceipts == nil {
		return *new([]*types.MessageReceipt), ErrNotSupported
	}
	return s.Internal.ChainGetParentReceipts(p0, p1)
}

func (s *FullNodeStruct) ChainGetPath(p0 context.Context, p1 types.TipSetKey, p2 types.TipSetKey) ([]*HeadChange, error) {
	if s.Internal.ChainGetPath == nil {
		return *new([]*HeadChange), ErrNotSupported
	}
	return s.Internal.ChainGetPath(p0, p1, p2)
}

func (s *FullNodeStruct) ChainGetTipSet(p0 context.Context, p1 types.TipSetKey) (*types.TipSet, error) {
	if s.Internal.ChainGetTipSet == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.ChainGetTipSet(p0, p1)
}

func (s *FullNodeStruct) ChainGetTipSetAfterHeight(p0 context.Context, p1 abi.ChainEpoch, p2 types.TipSetKey) (*types.TipSet, error) {
	if s.Internal.ChainGetTipSetAfterHeight == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.ChainGetTipSetAfterHeight(p0, p1, p2)
}

func (s *FullNodeStruct) ChainGetTipSetByHeight(p0 context.Context, p1 abi.ChainEpoch, p2 types.TipSetKey) (*types.TipSet, error) {
	if s.Internal.ChainGetTipSetByHeight == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.ChainGetTipSetByHeight(p0, p1, p2)
}

func (s *FullNodeStruct) ChainHasObj(p0 context.Context, p1 cid.Cid) (bool, error) {
	if s.Internal.ChainHasObj == nil {
		return false, ErrNotSupported
	}
	return s.Internal.ChainHasObj(p0, p1)
}

func (s *FullNodeStruct) ChainHead(p0 context.Context) (*types.TipSet, error) {
	if s.Internal.ChainHead == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.ChainHead(p0)
}

func (s *FullNodeStruct) ChainNotify(p0 context.Context) (<-chan []*HeadChange, error) {
	if s.Internal.ChainNotify == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.ChainNotify(p0)
}

func (s *FullNodeStruct) ChainPrune(p0 context.Context, p1 PruneOpts) error {
	if s.Internal.ChainPrune == nil {
		return ErrNotSupported
	}
	return s.Internal.ChainPrune(p0, p1)
}

func (s *FullNodeStruct) ChainPutObj(p0 context.Context, p1 blocks.Block) error {
	if s.Internal.ChainPutObj == nil {
		return ErrNotSupported
	}
	return s.Internal.ChainPutObj(p0, p1)
}

func (s *FullNodeStruct) ChainReadObj(p0 context.Context, p1 cid.Cid) ([]byte, error) {
	if s.Internal.ChainReadObj == nil {
		return *new([]byte), ErrNotSupported
	}
	return s.Internal.ChainReadObj(p0, p1)
}

func (s *FullNodeStruct) ChainSetHead(p0 context.Context, p1 types.TipSetKey) error {
	if s.Internal.ChainSetHead == nil {
		return ErrNotSupported
	}
	return s.Internal.ChainSetHead(p0, p1)
}

type ObjStat struct {
	Size  uint64
	Links uint64
}

func (s *FullNodeStruct) ChainStatObj(p0 context.Context, p1 cid.Cid, p2 cid.Cid) (ObjStat, error) {
	if s.Internal.ChainStatObj == nil {
		return *new(ObjStat), ErrNotSupported
	}
	return s.Internal.ChainStatObj(p0, p1, p2)
}

func (s *FullNodeStruct) ChainTipSetWeight(p0 context.Context, p1 types.TipSetKey) (types.BigInt, error) {
	if s.Internal.ChainTipSetWeight == nil {
		return *new(types.BigInt), ErrNotSupported
	}
	return s.Internal.ChainTipSetWeight(p0, p1)
}

func (s *FullNodeStruct) ClientCalcCommP(p0 context.Context, p1 string) (*CommPRet, error) {
	if s.Internal.ClientCalcCommP == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.ClientCalcCommP(p0, p1)
}

func (s *FullNodeStruct) ClientCancelDataTransfer(p0 context.Context, p1 datatransfer.TransferID, p2 peer.ID, p3 bool) error {
	if s.Internal.ClientCancelDataTransfer == nil {
		return ErrNotSupported
	}
	return s.Internal.ClientCancelDataTransfer(p0, p1, p2, p3)
}

func (s *FullNodeStruct) ClientCancelRetrievalDeal(p0 context.Context, p1 retrievalmarket.DealID) error {
	if s.Internal.ClientCancelRetrievalDeal == nil {
		return ErrNotSupported
	}
	return s.Internal.ClientCancelRetrievalDeal(p0, p1)
}

type DataTransferChannel struct {
	TransferID  datatransfer.TransferID
	Status      datatransfer.Status
	BaseCID     cid.Cid
	IsInitiator bool
	IsSender    bool
	Voucher     string
	Message     string
	OtherPeer   peer.ID
	Transferred uint64
	Stages      *datatransfer.ChannelStages
}

func (s *FullNodeStruct) ClientDataTransferUpdates(p0 context.Context) (<-chan DataTransferChannel, error) {
	if s.Internal.ClientDataTransferUpdates == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.ClientDataTransferUpdates(p0)
}

func (s *FullNodeStruct) ClientDealPieceCID(p0 context.Context, p1 cid.Cid) (DataCIDSize, error) {
	if s.Internal.ClientDealPieceCID == nil {
		return *new(DataCIDSize), ErrNotSupported
	}
	return s.Internal.ClientDealPieceCID(p0, p1)
}

func (s *FullNodeStruct) ClientDealSize(p0 context.Context, p1 cid.Cid) (DataSize, error) {
	if s.Internal.ClientDealSize == nil {
		return *new(DataSize), ErrNotSupported
	}
	return s.Internal.ClientDealSize(p0, p1)
}

type ExportRef struct {
	Root cid.Cid

	// DAGs array specifies a list of DAGs to export
	// - If exporting into unixfs files, only one DAG is supported, DataSelector is only used to find the targeted root node
	// - If exporting into a car file
	//   - When exactly one text-path DataSelector is specified exports the subgraph and its full merkle-path from the original root
	//   - Otherwise ( multiple paths and/or JSON selector specs) determines each individual subroot and exports the subtrees as a multi-root car
	// - When not specified defaults to a single DAG:
	//   - Data - the entire DAG: `{"R":{"l":{"none":{}},":>":{"a":{">":{"@":{}}}}}}`
	DAGs []DagSpec

	FromLocalCAR string // if specified, get data from a local CARv2 file.
	DealID       retrievalmarket.DealID
}

type DagSpec struct {
	// DataSelector matches data to be retrieved
	// - when using textselector, the path specifies subtree
	// - the matched graph must have a single root
	DataSelector *Selector

	// ExportMerkleProof is applicable only when exporting to a CAR file via a path textselector
	// When true, in addition to the selection target, the resulting CAR will contain every block along the
	// path back to, and including the original root
	// When false the resulting CAR contains only the blocks of the target subdag
	ExportMerkleProof bool
}

type FileRef struct {
	Path  string
	IsCAR bool
}

func (s *FullNodeStruct) ClientExport(p0 context.Context, p1 ExportRef, p2 FileRef) error {
	if s.Internal.ClientExport == nil {
		return ErrNotSupported
	}
	return s.Internal.ClientExport(p0, p1, p2)
}

func (s *FullNodeStruct) ClientFindData(p0 context.Context, p1 cid.Cid, p2 *cid.Cid) ([]QueryOffer, error) {
	if s.Internal.ClientFindData == nil {
		return *new([]QueryOffer), ErrNotSupported
	}
	return s.Internal.ClientFindData(p0, p1, p2)
}

func (s *FullNodeStruct) ClientGenCar(p0 context.Context, p1 FileRef, p2 string) error {
	if s.Internal.ClientGenCar == nil {
		return ErrNotSupported
	}
	return s.Internal.ClientGenCar(p0, p1, p2)
}

func (s *FullNodeStruct) ClientGetDealInfo(p0 context.Context, p1 cid.Cid) (*DealInfo, error) {
	if s.Internal.ClientGetDealInfo == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.ClientGetDealInfo(p0, p1)
}

func (s *FullNodeStruct) ClientGetDealStatus(p0 context.Context, p1 uint64) (string, error) {
	if s.Internal.ClientGetDealStatus == nil {
		return "", ErrNotSupported
	}
	return s.Internal.ClientGetDealStatus(p0, p1)
}

func (s *FullNodeStruct) ClientGetDealUpdates(p0 context.Context) (<-chan DealInfo, error) {
	if s.Internal.ClientGetDealUpdates == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.ClientGetDealUpdates(p0)
}

func (s *FullNodeStruct) ClientGetRetrievalUpdates(p0 context.Context) (<-chan RetrievalInfo, error) {
	if s.Internal.ClientGetRetrievalUpdates == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.ClientGetRetrievalUpdates(p0)
}

func (s *FullNodeStruct) ClientHasLocal(p0 context.Context, p1 cid.Cid) (bool, error) {
	if s.Internal.ClientHasLocal == nil {
		return false, ErrNotSupported
	}
	return s.Internal.ClientHasLocal(p0, p1)
}

type ImportRes struct {
	Root     cid.Cid
	ImportID ID
}

func (s *FullNodeStruct) ClientImport(p0 context.Context, p1 FileRef) (*ImportRes, error) {
	if s.Internal.ClientImport == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.ClientImport(p0, p1)
}

func (s *FullNodeStruct) ClientListDataTransfers(p0 context.Context) ([]DataTransferChannel, error) {
	if s.Internal.ClientListDataTransfers == nil {
		return *new([]DataTransferChannel), ErrNotSupported
	}
	return s.Internal.ClientListDataTransfers(p0)
}

func (s *FullNodeStruct) ClientListDeals(p0 context.Context) ([]DealInfo, error) {
	if s.Internal.ClientListDeals == nil {
		return *new([]DealInfo), ErrNotSupported
	}
	return s.Internal.ClientListDeals(p0)
}

type Import struct {
	Key ID
	Err string

	Root *cid.Cid

	// Source is the provenance of the import, e.g. "import", "unknown", else.
	// Currently useless but may be used in the future.
	Source string

	// FilePath is the path of the original file. It is important that the file
	// is retained at this path, because it will be referenced during
	// the transfer (when we do the UnixFS chunking, we don't duplicate the
	// leaves, but rather point to chunks of the original data through
	// positional references).
	FilePath string

	// CARPath is the path of the CAR file containing the DAG for this import.
	CARPath string
}

func (s *FullNodeStruct) ClientListImports(p0 context.Context) ([]Import, error) {
	if s.Internal.ClientListImports == nil {
		return *new([]Import), ErrNotSupported
	}
	return s.Internal.ClientListImports(p0)
}

func (s *FullNodeStruct) ClientListRetrievals(p0 context.Context) ([]RetrievalInfo, error) {
	if s.Internal.ClientListRetrievals == nil {
		return *new([]RetrievalInfo), ErrNotSupported
	}
	return s.Internal.ClientListRetrievals(p0)
}

func (s *FullNodeStruct) ClientMinerQueryOffer(p0 context.Context, p1 address.Address, p2 cid.Cid, p3 *cid.Cid) (QueryOffer, error) {
	if s.Internal.ClientMinerQueryOffer == nil {
		return *new(QueryOffer), ErrNotSupported
	}
	return s.Internal.ClientMinerQueryOffer(p0, p1, p2, p3)
}

type StorageAsk struct {
	Response *storagemarket.StorageAsk

	DealProtocols []string
}

func (s *FullNodeStruct) ClientQueryAsk(p0 context.Context, p1 peer.ID, p2 address.Address) (*StorageAsk, error) {
	if s.Internal.ClientQueryAsk == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.ClientQueryAsk(p0, p1, p2)
}

func (s *FullNodeStruct) ClientRemoveImport(p0 context.Context, p1 ID) error {
	if s.Internal.ClientRemoveImport == nil {
		return ErrNotSupported
	}
	return s.Internal.ClientRemoveImport(p0, p1)
}

func (s *FullNodeStruct) ClientRestartDataTransfer(p0 context.Context, p1 datatransfer.TransferID, p2 peer.ID, p3 bool) error {
	if s.Internal.ClientRestartDataTransfer == nil {
		return ErrNotSupported
	}
	return s.Internal.ClientRestartDataTransfer(p0, p1, p2, p3)
}

func (s *FullNodeStruct) ClientRetrieve(p0 context.Context, p1 RetrievalOrder) (*RestrievalRes, error) {
	if s.Internal.ClientRetrieve == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.ClientRetrieve(p0, p1)
}

func (s *FullNodeStruct) ClientRetrieveTryRestartInsufficientFunds(p0 context.Context, p1 address.Address) error {
	if s.Internal.ClientRetrieveTryRestartInsufficientFunds == nil {
		return ErrNotSupported
	}
	return s.Internal.ClientRetrieveTryRestartInsufficientFunds(p0, p1)
}

func (s *FullNodeStruct) ClientRetrieveWait(p0 context.Context, p1 retrievalmarket.DealID) error {
	if s.Internal.ClientRetrieveWait == nil {
		return ErrNotSupported
	}
	return s.Internal.ClientRetrieveWait(p0, p1)
}

func (s *FullNodeStruct) ClientStartDeal(p0 context.Context, p1 *StartDealParams) (*cid.Cid, error) {
	if s.Internal.ClientStartDeal == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.ClientStartDeal(p0, p1)
}

func (s *FullNodeStruct) ClientStatelessDeal(p0 context.Context, p1 *StartDealParams) (*cid.Cid, error) {
	if s.Internal.ClientStatelessDeal == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.ClientStatelessDeal(p0, p1)
}

func (s *FullNodeStruct) CreateBackup(p0 context.Context, p1 string) error {
	if s.Internal.CreateBackup == nil {
		return ErrNotSupported
	}
	return s.Internal.CreateBackup(p0, p1)
}

func (s *FullNodeStruct) GasEstimateFeeCap(p0 context.Context, p1 *types.Message, p2 int64, p3 types.TipSetKey) (types.BigInt, error) {
	if s.Internal.GasEstimateFeeCap == nil {
		return *new(types.BigInt), ErrNotSupported
	}
	return s.Internal.GasEstimateFeeCap(p0, p1, p2, p3)
}

func (s *FullNodeStruct) GasEstimateGasLimit(p0 context.Context, p1 *types.Message, p2 types.TipSetKey) (int64, error) {
	if s.Internal.GasEstimateGasLimit == nil {
		return 0, ErrNotSupported
	}
	return s.Internal.GasEstimateGasLimit(p0, p1, p2)
}

func (s *FullNodeStruct) GasEstimateGasPremium(p0 context.Context, p1 uint64, p2 address.Address, p3 int64, p4 types.TipSetKey) (types.BigInt, error) {
	if s.Internal.GasEstimateGasPremium == nil {
		return *new(types.BigInt), ErrNotSupported
	}
	return s.Internal.GasEstimateGasPremium(p0, p1, p2, p3, p4)
}

func (s *FullNodeStruct) GasEstimateMessageGas(p0 context.Context, p1 *types.Message, p2 *MessageSendSpec, p3 types.TipSetKey) (*types.Message, error) {
	if s.Internal.GasEstimateMessageGas == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.GasEstimateMessageGas(p0, p1, p2, p3)
}

func (s *FullNodeStruct) MarketAddBalance(p0 context.Context, p1 address.Address, p2 address.Address, p3 types.BigInt) (cid.Cid, error) {
	if s.Internal.MarketAddBalance == nil {
		return *new(cid.Cid), ErrNotSupported
	}
	return s.Internal.MarketAddBalance(p0, p1, p2, p3)
}

func (s *FullNodeStruct) MarketGetReserved(p0 context.Context, p1 address.Address) (types.BigInt, error) {
	if s.Internal.MarketGetReserved == nil {
		return *new(types.BigInt), ErrNotSupported
	}
	return s.Internal.MarketGetReserved(p0, p1)
}

func (s *FullNodeStruct) MarketReleaseFunds(p0 context.Context, p1 address.Address, p2 types.BigInt) error {
	if s.Internal.MarketReleaseFunds == nil {
		return ErrNotSupported
	}
	return s.Internal.MarketReleaseFunds(p0, p1, p2)
}

func (s *FullNodeStruct) MarketReserveFunds(p0 context.Context, p1 address.Address, p2 address.Address, p3 types.BigInt) (cid.Cid, error) {
	if s.Internal.MarketReserveFunds == nil {
		return *new(cid.Cid), ErrNotSupported
	}
	return s.Internal.MarketReserveFunds(p0, p1, p2, p3)
}

func (s *FullNodeStruct) MarketWithdraw(p0 context.Context, p1 address.Address, p2 address.Address, p3 types.BigInt) (cid.Cid, error) {
	if s.Internal.MarketWithdraw == nil {
		return *new(cid.Cid), ErrNotSupported
	}
	return s.Internal.MarketWithdraw(p0, p1, p2, p3)
}

func (s *FullNodeStruct) MinerCreateBlock(p0 context.Context, p1 *BlockTemplate) (*types.BlockMsg, error) {
	if s.Internal.MinerCreateBlock == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.MinerCreateBlock(p0, p1)
}

func (s *FullNodeStruct) MinerGetBaseInfo(p0 context.Context, p1 address.Address, p2 abi.ChainEpoch, p3 types.TipSetKey) (*MiningBaseInfo, error) {
	if s.Internal.MinerGetBaseInfo == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.MinerGetBaseInfo(p0, p1, p2, p3)
}

func (s *FullNodeStruct) MpoolBatchPush(p0 context.Context, p1 []*types.SignedMessage) ([]cid.Cid, error) {
	if s.Internal.MpoolBatchPush == nil {
		return *new([]cid.Cid), ErrNotSupported
	}
	return s.Internal.MpoolBatchPush(p0, p1)
}

func (s *FullNodeStruct) MpoolBatchPushMessage(p0 context.Context, p1 []*types.Message, p2 *MessageSendSpec) ([]*types.SignedMessage, error) {
	if s.Internal.MpoolBatchPushMessage == nil {
		return *new([]*types.SignedMessage), ErrNotSupported
	}
	return s.Internal.MpoolBatchPushMessage(p0, p1, p2)
}

func (s *FullNodeStruct) MpoolBatchPushUntrusted(p0 context.Context, p1 []*types.SignedMessage) ([]cid.Cid, error) {
	if s.Internal.MpoolBatchPushUntrusted == nil {
		return *new([]cid.Cid), ErrNotSupported
	}
	return s.Internal.MpoolBatchPushUntrusted(p0, p1)
}

func (s *FullNodeStruct) MpoolCheckMessages(p0 context.Context, p1 []*MessagePrototype) ([][]MessageCheckStatus, error) {
	if s.Internal.MpoolCheckMessages == nil {
		return *new([][]MessageCheckStatus), ErrNotSupported
	}
	return s.Internal.MpoolCheckMessages(p0, p1)
}

func (s *FullNodeStruct) MpoolCheckPendingMessages(p0 context.Context, p1 address.Address) ([][]MessageCheckStatus, error) {
	if s.Internal.MpoolCheckPendingMessages == nil {
		return *new([][]MessageCheckStatus), ErrNotSupported
	}
	return s.Internal.MpoolCheckPendingMessages(p0, p1)
}

func (s *FullNodeStruct) MpoolCheckReplaceMessages(p0 context.Context, p1 []*types.Message) ([][]MessageCheckStatus, error) {
	if s.Internal.MpoolCheckReplaceMessages == nil {
		return *new([][]MessageCheckStatus), ErrNotSupported
	}
	return s.Internal.MpoolCheckReplaceMessages(p0, p1)
}

func (s *FullNodeStruct) MpoolClear(p0 context.Context, p1 bool) error {
	if s.Internal.MpoolClear == nil {
		return ErrNotSupported
	}
	return s.Internal.MpoolClear(p0, p1)
}

func (s *FullNodeStruct) MpoolGetConfig(p0 context.Context) (*types.MpoolConfig, error) {
	if s.Internal.MpoolGetConfig == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.MpoolGetConfig(p0)
}

func (s *FullNodeStruct) MpoolGetNonce(p0 context.Context, p1 address.Address) (uint64, error) {
	if s.Internal.MpoolGetNonce == nil {
		return 0, ErrNotSupported
	}
	return s.Internal.MpoolGetNonce(p0, p1)
}

func (s *FullNodeStruct) MpoolPending(p0 context.Context, p1 types.TipSetKey) ([]*types.SignedMessage, error) {
	if s.Internal.MpoolPending == nil {
		return *new([]*types.SignedMessage), ErrNotSupported
	}
	return s.Internal.MpoolPending(p0, p1)
}

func (s *FullNodeStruct) MpoolPush(p0 context.Context, p1 *types.SignedMessage) (cid.Cid, error) {
	if s.Internal.MpoolPush == nil {
		return *new(cid.Cid), ErrNotSupported
	}
	return s.Internal.MpoolPush(p0, p1)
}

func (s *FullNodeStruct) MpoolPushMessage(p0 context.Context, p1 *types.Message, p2 *MessageSendSpec) (*types.SignedMessage, error) {
	if s.Internal.MpoolPushMessage == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.MpoolPushMessage(p0, p1, p2)
}

func (s *FullNodeStruct) MpoolPushUntrusted(p0 context.Context, p1 *types.SignedMessage) (cid.Cid, error) {
	if s.Internal.MpoolPushUntrusted == nil {
		return *new(cid.Cid), ErrNotSupported
	}
	return s.Internal.MpoolPushUntrusted(p0, p1)
}

func (s *FullNodeStruct) MpoolSelect(p0 context.Context, p1 types.TipSetKey, p2 float64) ([]*types.SignedMessage, error) {
	if s.Internal.MpoolSelect == nil {
		return *new([]*types.SignedMessage), ErrNotSupported
	}
	return s.Internal.MpoolSelect(p0, p1, p2)
}

func (s *FullNodeStruct) MpoolSetConfig(p0 context.Context, p1 *types.MpoolConfig) error {
	if s.Internal.MpoolSetConfig == nil {
		return ErrNotSupported
	}
	return s.Internal.MpoolSetConfig(p0, p1)
}

func (s *FullNodeStruct) MpoolSub(p0 context.Context) (<-chan MpoolUpdate, error) {
	if s.Internal.MpoolSub == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.MpoolSub(p0)
}

func (s *FullNodeStruct) MsigAddApprove(p0 context.Context, p1 address.Address, p2 address.Address, p3 uint64, p4 address.Address, p5 address.Address, p6 bool) (*MessagePrototype, error) {
	if s.Internal.MsigAddApprove == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.MsigAddApprove(p0, p1, p2, p3, p4, p5, p6)
}

func (s *FullNodeStruct) MsigAddCancel(p0 context.Context, p1 address.Address, p2 address.Address, p3 uint64, p4 address.Address, p5 bool) (*MessagePrototype, error) {
	if s.Internal.MsigAddCancel == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.MsigAddCancel(p0, p1, p2, p3, p4, p5)
}

func (s *FullNodeStruct) MsigAddPropose(p0 context.Context, p1 address.Address, p2 address.Address, p3 address.Address, p4 bool) (*MessagePrototype, error) {
	if s.Internal.MsigAddPropose == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.MsigAddPropose(p0, p1, p2, p3, p4)
}

func (s *FullNodeStruct) MsigApprove(p0 context.Context, p1 address.Address, p2 uint64, p3 address.Address) (*MessagePrototype, error) {
	if s.Internal.MsigApprove == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.MsigApprove(p0, p1, p2, p3)
}

func (s *FullNodeStruct) MsigApproveTxnHash(p0 context.Context, p1 address.Address, p2 uint64, p3 address.Address, p4 address.Address, p5 types.BigInt, p6 address.Address, p7 uint64, p8 []byte) (*MessagePrototype, error) {
	if s.Internal.MsigApproveTxnHash == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.MsigApproveTxnHash(p0, p1, p2, p3, p4, p5, p6, p7, p8)
}

func (s *FullNodeStruct) MsigCancel(p0 context.Context, p1 address.Address, p2 uint64, p3 address.Address) (*MessagePrototype, error) {
	if s.Internal.MsigCancel == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.MsigCancel(p0, p1, p2, p3)
}

func (s *FullNodeStruct) MsigCancelTxnHash(p0 context.Context, p1 address.Address, p2 uint64, p3 address.Address, p4 types.BigInt, p5 address.Address, p6 uint64, p7 []byte) (*MessagePrototype, error) {
	if s.Internal.MsigCancelTxnHash == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.MsigCancelTxnHash(p0, p1, p2, p3, p4, p5, p6, p7)
}

func (s *FullNodeStruct) MsigCreate(p0 context.Context, p1 uint64, p2 []address.Address, p3 abi.ChainEpoch, p4 types.BigInt, p5 address.Address, p6 types.BigInt) (*MessagePrototype, error) {
	if s.Internal.MsigCreate == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.MsigCreate(p0, p1, p2, p3, p4, p5, p6)
}

func (s *FullNodeStruct) MsigGetAvailableBalance(p0 context.Context, p1 address.Address, p2 types.TipSetKey) (types.BigInt, error) {
	if s.Internal.MsigGetAvailableBalance == nil {
		return *new(types.BigInt), ErrNotSupported
	}
	return s.Internal.MsigGetAvailableBalance(p0, p1, p2)
}

func (s *FullNodeStruct) MsigGetPending(p0 context.Context, p1 address.Address, p2 types.TipSetKey) ([]*MsigTransaction, error) {
	if s.Internal.MsigGetPending == nil {
		return *new([]*MsigTransaction), ErrNotSupported
	}
	return s.Internal.MsigGetPending(p0, p1, p2)
}

func (s *FullNodeStruct) MsigGetVested(p0 context.Context, p1 address.Address, p2 types.TipSetKey, p3 types.TipSetKey) (types.BigInt, error) {
	if s.Internal.MsigGetVested == nil {
		return *new(types.BigInt), ErrNotSupported
	}
	return s.Internal.MsigGetVested(p0, p1, p2, p3)
}

func (s *FullNodeStruct) MsigGetVestingSchedule(p0 context.Context, p1 address.Address, p2 types.TipSetKey) (MsigVesting, error) {
	if s.Internal.MsigGetVestingSchedule == nil {
		return *new(MsigVesting), ErrNotSupported
	}
	return s.Internal.MsigGetVestingSchedule(p0, p1, p2)
}

func (s *FullNodeStruct) MsigPropose(p0 context.Context, p1 address.Address, p2 address.Address, p3 types.BigInt, p4 address.Address, p5 uint64, p6 []byte) (*MessagePrototype, error) {
	if s.Internal.MsigPropose == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.MsigPropose(p0, p1, p2, p3, p4, p5, p6)
}

func (s *FullNodeStruct) MsigRemoveSigner(p0 context.Context, p1 address.Address, p2 address.Address, p3 address.Address, p4 bool) (*MessagePrototype, error) {
	if s.Internal.MsigRemoveSigner == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.MsigRemoveSigner(p0, p1, p2, p3, p4)
}

func (s *FullNodeStruct) MsigSwapApprove(p0 context.Context, p1 address.Address, p2 address.Address, p3 uint64, p4 address.Address, p5 address.Address, p6 address.Address) (*MessagePrototype, error) {
	if s.Internal.MsigSwapApprove == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.MsigSwapApprove(p0, p1, p2, p3, p4, p5, p6)
}

func (s *FullNodeStruct) MsigSwapCancel(p0 context.Context, p1 address.Address, p2 address.Address, p3 uint64, p4 address.Address, p5 address.Address) (*MessagePrototype, error) {
	if s.Internal.MsigSwapCancel == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.MsigSwapCancel(p0, p1, p2, p3, p4, p5)
}

func (s *FullNodeStruct) MsigSwapPropose(p0 context.Context, p1 address.Address, p2 address.Address, p3 address.Address, p4 address.Address) (*MessagePrototype, error) {
	if s.Internal.MsigSwapPropose == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.MsigSwapPropose(p0, p1, p2, p3, p4)
}

type NodeStatus struct {
	SyncStatus  NodeSyncStatus
	PeerStatus  NodePeerStatus
	ChainStatus NodeChainStatus
}

type NodeSyncStatus struct {
	Epoch  uint64
	Behind uint64
}

type NodePeerStatus struct {
	PeersToPublishMsgs   int
	PeersToPublishBlocks int
}

type NodeChainStatus struct {
	BlocksPerTipsetLast100      float64
	BlocksPerTipsetLastFinality float64
}

func (s *FullNodeStruct) NodeStatus(p0 context.Context, p1 bool) (NodeStatus, error) {
	if s.Internal.NodeStatus == nil {
		return *new(NodeStatus), ErrNotSupported
	}
	return s.Internal.NodeStatus(p0, p1)
}

func (s *FullNodeStruct) PaychAllocateLane(p0 context.Context, p1 address.Address) (uint64, error) {
	if s.Internal.PaychAllocateLane == nil {
		return 0, ErrNotSupported
	}
	return s.Internal.PaychAllocateLane(p0, p1)
}

func (s *FullNodeStruct) PaychAvailableFunds(p0 context.Context, p1 address.Address) (*ChannelAvailableFunds, error) {
	if s.Internal.PaychAvailableFunds == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.PaychAvailableFunds(p0, p1)
}

func (s *FullNodeStruct) PaychAvailableFundsByFromTo(p0 context.Context, p1 address.Address, p2 address.Address) (*ChannelAvailableFunds, error) {
	if s.Internal.PaychAvailableFundsByFromTo == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.PaychAvailableFundsByFromTo(p0, p1, p2)
}

func (s *FullNodeStruct) PaychCollect(p0 context.Context, p1 address.Address) (cid.Cid, error) {
	if s.Internal.PaychCollect == nil {
		return *new(cid.Cid), ErrNotSupported
	}
	return s.Internal.PaychCollect(p0, p1)
}

func (s *FullNodeStruct) PaychFund(p0 context.Context, p1 address.Address, p2 address.Address, p3 types.BigInt) (*ChannelInfo, error) {
	if s.Internal.PaychFund == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.PaychFund(p0, p1, p2, p3)
}

func (s *FullNodeStruct) PaychGet(p0 context.Context, p1 address.Address, p2 address.Address, p3 types.BigInt, p4 PaychGetOpts) (*ChannelInfo, error) {
	if s.Internal.PaychGet == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.PaychGet(p0, p1, p2, p3, p4)
}

func (s *FullNodeStruct) PaychGetWaitReady(p0 context.Context, p1 cid.Cid) (address.Address, error) {
	if s.Internal.PaychGetWaitReady == nil {
		return *new(address.Address), ErrNotSupported
	}
	return s.Internal.PaychGetWaitReady(p0, p1)
}

func (s *FullNodeStruct) PaychList(p0 context.Context) ([]address.Address, error) {
	if s.Internal.PaychList == nil {
		return *new([]address.Address), ErrNotSupported
	}
	return s.Internal.PaychList(p0)
}

func (s *FullNodeStruct) PaychNewPayment(p0 context.Context, p1 address.Address, p2 address.Address, p3 []VoucherSpec) (*PaymentInfo, error) {
	if s.Internal.PaychNewPayment == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.PaychNewPayment(p0, p1, p2, p3)
}

func (s *FullNodeStruct) PaychSettle(p0 context.Context, p1 address.Address) (cid.Cid, error) {
	if s.Internal.PaychSettle == nil {
		return *new(cid.Cid), ErrNotSupported
	}
	return s.Internal.PaychSettle(p0, p1)
}

func (s *FullNodeStruct) PaychStatus(p0 context.Context, p1 address.Address) (*PaychStatus, error) {
	if s.Internal.PaychStatus == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.PaychStatus(p0, p1)
}

func (s *FullNodeStruct) PaychVoucherAdd(p0 context.Context, p1 address.Address, p2 *paych.SignedVoucher, p3 []byte, p4 types.BigInt) (types.BigInt, error) {
	if s.Internal.PaychVoucherAdd == nil {
		return *new(types.BigInt), ErrNotSupported
	}
	return s.Internal.PaychVoucherAdd(p0, p1, p2, p3, p4)
}

func (s *FullNodeStruct) PaychVoucherCheckSpendable(p0 context.Context, p1 address.Address, p2 *paych.SignedVoucher, p3 []byte, p4 []byte) (bool, error) {
	if s.Internal.PaychVoucherCheckSpendable == nil {
		return false, ErrNotSupported
	}
	return s.Internal.PaychVoucherCheckSpendable(p0, p1, p2, p3, p4)
}

func (s *FullNodeStruct) PaychVoucherCheckValid(p0 context.Context, p1 address.Address, p2 *paych.SignedVoucher) error {
	if s.Internal.PaychVoucherCheckValid == nil {
		return ErrNotSupported
	}
	return s.Internal.PaychVoucherCheckValid(p0, p1, p2)
}

func (s *FullNodeStruct) PaychVoucherCreate(p0 context.Context, p1 address.Address, p2 types.BigInt, p3 uint64) (*VoucherCreateResult, error) {
	if s.Internal.PaychVoucherCreate == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.PaychVoucherCreate(p0, p1, p2, p3)
}

func (s *FullNodeStruct) PaychVoucherList(p0 context.Context, p1 address.Address) ([]*paych.SignedVoucher, error) {
	if s.Internal.PaychVoucherList == nil {
		return *new([]*paych.SignedVoucher), ErrNotSupported
	}
	return s.Internal.PaychVoucherList(p0, p1)
}

func (s *FullNodeStruct) PaychVoucherSubmit(p0 context.Context, p1 address.Address, p2 *paych.SignedVoucher, p3 []byte, p4 []byte) (cid.Cid, error) {
	if s.Internal.PaychVoucherSubmit == nil {
		return *new(cid.Cid), ErrNotSupported
	}
	return s.Internal.PaychVoucherSubmit(p0, p1, p2, p3, p4)
}

func (s *FullNodeStruct) RaftLeader(p0 context.Context) (peer.ID, error) {
	if s.Internal.RaftLeader == nil {
		return *new(peer.ID), ErrNotSupported
	}
	return s.Internal.RaftLeader(p0)
}

func (s *FullNodeStruct) RaftState(p0 context.Context) (*RaftStateData, error) {
	if s.Internal.RaftState == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.RaftState(p0)
}

type NonceMapType map[address.Address]uint64
type MsgUuidMapType map[uuid.UUID]*types.SignedMessage

type RaftStateData struct {
	NonceMap NonceMapType
	MsgUuids MsgUuidMapType
}

func (s *FullNodeStruct) StateAccountKey(p0 context.Context, p1 address.Address, p2 types.TipSetKey) (address.Address, error) {
	if s.Internal.StateAccountKey == nil {
		return *new(address.Address), ErrNotSupported
	}
	return s.Internal.StateAccountKey(p0, p1, p2)
}

func (s *FullNodeStruct) StateActorCodeCIDs(p0 context.Context, p1 abinetwork.Version) (map[string]cid.Cid, error) {
	if s.Internal.StateActorCodeCIDs == nil {
		return *new(map[string]cid.Cid), ErrNotSupported
	}
	return s.Internal.StateActorCodeCIDs(p0, p1)
}

func (s *FullNodeStruct) StateActorManifestCID(p0 context.Context, p1 abinetwork.Version) (cid.Cid, error) {
	if s.Internal.StateActorManifestCID == nil {
		return *new(cid.Cid), ErrNotSupported
	}
	return s.Internal.StateActorManifestCID(p0, p1)
}

func (s *FullNodeStruct) StateAllMinerFaults(p0 context.Context, p1 abi.ChainEpoch, p2 types.TipSetKey) ([]*Fault, error) {
	if s.Internal.StateAllMinerFaults == nil {
		return *new([]*Fault), ErrNotSupported
	}
	return s.Internal.StateAllMinerFaults(p0, p1, p2)
}

func (s *FullNodeStruct) StateCall(p0 context.Context, p1 *types.Message, p2 types.TipSetKey) (*InvocResult, error) {
	if s.Internal.StateCall == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.StateCall(p0, p1, p2)
}

func (s *FullNodeStruct) StateChangedActors(p0 context.Context, p1 cid.Cid, p2 cid.Cid) (map[string]types.Actor, error) {
	if s.Internal.StateChangedActors == nil {
		return *new(map[string]types.Actor), ErrNotSupported
	}
	return s.Internal.StateChangedActors(p0, p1, p2)
}

func (s *FullNodeStruct) StateCirculatingSupply(p0 context.Context, p1 types.TipSetKey) (abi.TokenAmount, error) {
	if s.Internal.StateCirculatingSupply == nil {
		return *new(abi.TokenAmount), ErrNotSupported
	}
	return s.Internal.StateCirculatingSupply(p0, p1)
}

func (s *FullNodeStruct) StateCompute(p0 context.Context, p1 abi.ChainEpoch, p2 []*types.Message, p3 types.TipSetKey) (*ComputeStateOutput, error) {
	if s.Internal.StateCompute == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.StateCompute(p0, p1, p2, p3)
}

func (s *FullNodeStruct) StateComputeDataCID(p0 context.Context, p1 address.Address, p2 abi.RegisteredSealProof, p3 []abi.DealID, p4 types.TipSetKey) (cid.Cid, error) {
	if s.Internal.StateComputeDataCID == nil {
		return *new(cid.Cid), ErrNotSupported
	}
	return s.Internal.StateComputeDataCID(p0, p1, p2, p3, p4)
}

func (s *FullNodeStruct) StateDealProviderCollateralBounds(p0 context.Context, p1 abi.PaddedPieceSize, p2 bool, p3 types.TipSetKey) (DealCollateralBounds, error) {
	if s.Internal.StateDealProviderCollateralBounds == nil {
		return *new(DealCollateralBounds), ErrNotSupported
	}
	return s.Internal.StateDealProviderCollateralBounds(p0, p1, p2, p3)
}

func (s *FullNodeStruct) StateDecodeParams(p0 context.Context, p1 address.Address, p2 abi.MethodNum, p3 []byte, p4 types.TipSetKey) (interface{}, error) {
	if s.Internal.StateDecodeParams == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.StateDecodeParams(p0, p1, p2, p3, p4)
}

func (s *FullNodeStruct) StateEncodeParams(p0 context.Context, p1 cid.Cid, p2 abi.MethodNum, p3 json.RawMessage) ([]byte, error) {
	if s.Internal.StateEncodeParams == nil {
		return *new([]byte), ErrNotSupported
	}
	return s.Internal.StateEncodeParams(p0, p1, p2, p3)
}

func (s *FullNodeStruct) StateGetActor(p0 context.Context, p1 address.Address, p2 types.TipSetKey) (*types.Actor, error) {
	if s.Internal.StateGetActor == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.StateGetActor(p0, p1, p2)
}

func (s *FullNodeStruct) StateGetAllocation(p0 context.Context, p1 address.Address, p2 verifregtypes.AllocationId, p3 types.TipSetKey) (*verifregtypes.Allocation, error) {
	if s.Internal.StateGetAllocation == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.StateGetAllocation(p0, p1, p2, p3)
}

func (s *FullNodeStruct) StateGetAllocationForPendingDeal(p0 context.Context, p1 abi.DealID, p2 types.TipSetKey) (*verifregtypes.Allocation, error) {
	if s.Internal.StateGetAllocationForPendingDeal == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.StateGetAllocationForPendingDeal(p0, p1, p2)
}

func (s *FullNodeStruct) StateGetAllocations(p0 context.Context, p1 address.Address, p2 types.TipSetKey) (map[verifregtypes.AllocationId]verifregtypes.Allocation, error) {
	if s.Internal.StateGetAllocations == nil {
		return *new(map[verifregtypes.AllocationId]verifregtypes.Allocation), ErrNotSupported
	}
	return s.Internal.StateGetAllocations(p0, p1, p2)
}

func (s *FullNodeStruct) StateGetBeaconEntry(p0 context.Context, p1 abi.ChainEpoch) (*types.BeaconEntry, error) {
	if s.Internal.StateGetBeaconEntry == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.StateGetBeaconEntry(p0, p1)
}

func (s *FullNodeStruct) StateGetClaim(p0 context.Context, p1 address.Address, p2 verifregtypes.ClaimId, p3 types.TipSetKey) (*verifregtypes.Claim, error) {
	if s.Internal.StateGetClaim == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.StateGetClaim(p0, p1, p2, p3)
}

func (s *FullNodeStruct) StateGetClaims(p0 context.Context, p1 address.Address, p2 types.TipSetKey) (map[verifregtypes.ClaimId]verifregtypes.Claim, error) {
	if s.Internal.StateGetClaims == nil {
		return *new(map[verifregtypes.ClaimId]verifregtypes.Claim), ErrNotSupported
	}
	return s.Internal.StateGetClaims(p0, p1, p2)
}

func (s *FullNodeStruct) StateGetNetworkParams(p0 context.Context) (*NetworkParams, error) {
	if s.Internal.StateGetNetworkParams == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.StateGetNetworkParams(p0)
}

type NetworkParams struct {
	NetworkName             NetworkName
	BlockDelaySecs          uint64
	ConsensusMinerMinPower  abi.StoragePower
	SupportedProofTypes     []abi.RegisteredSealProof
	PreCommitChallengeDelay abi.ChainEpoch
	ForkUpgradeParams       ForkUpgradeParams
}

type ForkUpgradeParams struct {
	UpgradeSmokeHeight         abi.ChainEpoch
	UpgradeBreezeHeight        abi.ChainEpoch
	UpgradeIgnitionHeight      abi.ChainEpoch
	UpgradeLiftoffHeight       abi.ChainEpoch
	UpgradeAssemblyHeight      abi.ChainEpoch
	UpgradeRefuelHeight        abi.ChainEpoch
	UpgradeTapeHeight          abi.ChainEpoch
	UpgradeKumquatHeight       abi.ChainEpoch
	UpgradePriceListOopsHeight abi.ChainEpoch
	BreezeGasTampingDuration   abi.ChainEpoch
	UpgradeCalicoHeight        abi.ChainEpoch
	UpgradePersianHeight       abi.ChainEpoch
	UpgradeOrangeHeight        abi.ChainEpoch
	UpgradeClausHeight         abi.ChainEpoch
	UpgradeTrustHeight         abi.ChainEpoch
	UpgradeNorwegianHeight     abi.ChainEpoch
	UpgradeTurboHeight         abi.ChainEpoch
	UpgradeHyperdriveHeight    abi.ChainEpoch
	UpgradeChocolateHeight     abi.ChainEpoch
	UpgradeOhSnapHeight        abi.ChainEpoch
	UpgradeSkyrHeight          abi.ChainEpoch
	UpgradeSharkHeight         abi.ChainEpoch
}

func (s *FullNodeStruct) StateGetRandomnessFromBeacon(p0 context.Context, p1 crypto.DomainSeparationTag, p2 abi.ChainEpoch, p3 []byte, p4 types.TipSetKey) (abi.Randomness, error) {
	if s.Internal.StateGetRandomnessFromBeacon == nil {
		return *new(abi.Randomness), ErrNotSupported
	}
	return s.Internal.StateGetRandomnessFromBeacon(p0, p1, p2, p3, p4)
}

func (s *FullNodeStruct) StateGetRandomnessFromTickets(p0 context.Context, p1 crypto.DomainSeparationTag, p2 abi.ChainEpoch, p3 []byte, p4 types.TipSetKey) (abi.Randomness, error) {
	if s.Internal.StateGetRandomnessFromTickets == nil {
		return *new(abi.Randomness), ErrNotSupported
	}
	return s.Internal.StateGetRandomnessFromTickets(p0, p1, p2, p3, p4)
}

func (s *FullNodeStruct) StateListActors(p0 context.Context, p1 types.TipSetKey) ([]address.Address, error) {
	if s.Internal.StateListActors == nil {
		return *new([]address.Address), ErrNotSupported
	}
	return s.Internal.StateListActors(p0, p1)
}

func (s *FullNodeStruct) StateListMessages(p0 context.Context, p1 *MessageMatch, p2 types.TipSetKey, p3 abi.ChainEpoch) ([]cid.Cid, error) {
	if s.Internal.StateListMessages == nil {
		return *new([]cid.Cid), ErrNotSupported
	}
	return s.Internal.StateListMessages(p0, p1, p2, p3)
}

func (s *FullNodeStruct) StateListMiners(p0 context.Context, p1 types.TipSetKey) ([]address.Address, error) {
	if s.Internal.StateListMiners == nil {
		return *new([]address.Address), ErrNotSupported
	}
	return s.Internal.StateListMiners(p0, p1)
}

func (s *FullNodeStruct) StateLookupID(p0 context.Context, p1 address.Address, p2 types.TipSetKey) (address.Address, error) {
	if s.Internal.StateLookupID == nil {
		return *new(address.Address), ErrNotSupported
	}
	return s.Internal.StateLookupID(p0, p1, p2)
}

func (s *FullNodeStruct) StateLookupRobustAddress(p0 context.Context, p1 address.Address, p2 types.TipSetKey) (address.Address, error) {
	if s.Internal.StateLookupRobustAddress == nil {
		return *new(address.Address), ErrNotSupported
	}
	return s.Internal.StateLookupRobustAddress(p0, p1, p2)
}

func (s *FullNodeStruct) StateMarketBalance(p0 context.Context, p1 address.Address, p2 types.TipSetKey) (MarketBalance, error) {
	if s.Internal.StateMarketBalance == nil {
		return *new(MarketBalance), ErrNotSupported
	}
	return s.Internal.StateMarketBalance(p0, p1, p2)
}

func (s *FullNodeStruct) StateMarketDeals(p0 context.Context, p1 types.TipSetKey) (map[string]*MarketDeal, error) {
	if s.Internal.StateMarketDeals == nil {
		return *new(map[string]*MarketDeal), ErrNotSupported
	}
	return s.Internal.StateMarketDeals(p0, p1)
}

func (s *FullNodeStruct) StateMarketParticipants(p0 context.Context, p1 types.TipSetKey) (map[string]MarketBalance, error) {
	if s.Internal.StateMarketParticipants == nil {
		return *new(map[string]MarketBalance), ErrNotSupported
	}
	return s.Internal.StateMarketParticipants(p0, p1)
}

func (s *FullNodeStruct) StateMarketStorageDeal(p0 context.Context, p1 abi.DealID, p2 types.TipSetKey) (*MarketDeal, error) {
	if s.Internal.StateMarketStorageDeal == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.StateMarketStorageDeal(p0, p1, p2)
}

func (s *FullNodeStruct) StateMinerActiveSectors(p0 context.Context, p1 address.Address, p2 types.TipSetKey) ([]*miner.SectorOnChainInfo, error) {
	if s.Internal.StateMinerActiveSectors == nil {
		return *new([]*miner.SectorOnChainInfo), ErrNotSupported
	}
	return s.Internal.StateMinerActiveSectors(p0, p1, p2)
}

func (s *FullNodeStruct) StateMinerAllocated(p0 context.Context, p1 address.Address, p2 types.TipSetKey) (*bitfield.BitField, error) {
	if s.Internal.StateMinerAllocated == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.StateMinerAllocated(p0, p1, p2)
}

func (s *FullNodeStruct) StateMinerAvailableBalance(p0 context.Context, p1 address.Address, p2 types.TipSetKey) (types.BigInt, error) {
	if s.Internal.StateMinerAvailableBalance == nil {
		return *new(types.BigInt), ErrNotSupported
	}
	return s.Internal.StateMinerAvailableBalance(p0, p1, p2)
}

func (s *FullNodeStruct) StateMinerDeadlines(p0 context.Context, p1 address.Address, p2 types.TipSetKey) ([]Deadline, error) {
	if s.Internal.StateMinerDeadlines == nil {
		return *new([]Deadline), ErrNotSupported
	}
	return s.Internal.StateMinerDeadlines(p0, p1, p2)
}

func (s *FullNodeStruct) StateMinerFaults(p0 context.Context, p1 address.Address, p2 types.TipSetKey) (bitfield.BitField, error) {
	if s.Internal.StateMinerFaults == nil {
		return *new(bitfield.BitField), ErrNotSupported
	}
	return s.Internal.StateMinerFaults(p0, p1, p2)
}

func (s *FullNodeStruct) StateMinerInfo(p0 context.Context, p1 address.Address, p2 types.TipSetKey) (MinerInfo, error) {
	if s.Internal.StateMinerInfo == nil {
		return *new(MinerInfo), ErrNotSupported
	}
	return s.Internal.StateMinerInfo(p0, p1, p2)
}

func (s *FullNodeStruct) StateMinerInitialPledgeCollateral(p0 context.Context, p1 address.Address, p2 miner.SectorPreCommitInfo, p3 types.TipSetKey) (types.BigInt, error) {
	if s.Internal.StateMinerInitialPledgeCollateral == nil {
		return *new(types.BigInt), ErrNotSupported
	}
	return s.Internal.StateMinerInitialPledgeCollateral(p0, p1, p2, p3)
}

func (s *FullNodeStruct) StateMinerPartitions(p0 context.Context, p1 address.Address, p2 uint64, p3 types.TipSetKey) ([]Partition, error) {
	if s.Internal.StateMinerPartitions == nil {
		return *new([]Partition), ErrNotSupported
	}
	return s.Internal.StateMinerPartitions(p0, p1, p2, p3)
}

func (s *FullNodeStruct) StateMinerPower(p0 context.Context, p1 address.Address, p2 types.TipSetKey) (*MinerPower, error) {
	if s.Internal.StateMinerPower == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.StateMinerPower(p0, p1, p2)
}

func (s *FullNodeStruct) StateMinerPreCommitDepositForPower(p0 context.Context, p1 address.Address, p2 miner.SectorPreCommitInfo, p3 types.TipSetKey) (types.BigInt, error) {
	if s.Internal.StateMinerPreCommitDepositForPower == nil {
		return *new(types.BigInt), ErrNotSupported
	}
	return s.Internal.StateMinerPreCommitDepositForPower(p0, p1, p2, p3)
}

func (s *FullNodeStruct) StateMinerProvingDeadline(p0 context.Context, p1 address.Address, p2 types.TipSetKey) (*dline.Info, error) {
	if s.Internal.StateMinerProvingDeadline == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.StateMinerProvingDeadline(p0, p1, p2)
}

func (s *FullNodeStruct) StateMinerRecoveries(p0 context.Context, p1 address.Address, p2 types.TipSetKey) (bitfield.BitField, error) {
	if s.Internal.StateMinerRecoveries == nil {
		return *new(bitfield.BitField), ErrNotSupported
	}
	return s.Internal.StateMinerRecoveries(p0, p1, p2)
}

func (s *FullNodeStruct) StateMinerSectorAllocated(p0 context.Context, p1 address.Address, p2 abi.SectorNumber, p3 types.TipSetKey) (bool, error) {
	if s.Internal.StateMinerSectorAllocated == nil {
		return false, ErrNotSupported
	}
	return s.Internal.StateMinerSectorAllocated(p0, p1, p2, p3)
}

func (s *FullNodeStruct) StateMinerSectorCount(p0 context.Context, p1 address.Address, p2 types.TipSetKey) (MinerSectors, error) {
	if s.Internal.StateMinerSectorCount == nil {
		return *new(MinerSectors), ErrNotSupported
	}
	return s.Internal.StateMinerSectorCount(p0, p1, p2)
}

func (s *FullNodeStruct) StateMinerSectors(p0 context.Context, p1 address.Address, p2 *bitfield.BitField, p3 types.TipSetKey) ([]*miner.SectorOnChainInfo, error) {
	if s.Internal.StateMinerSectors == nil {
		return *new([]*miner.SectorOnChainInfo), ErrNotSupported
	}
	return s.Internal.StateMinerSectors(p0, p1, p2, p3)
}

func (s *FullNodeStruct) StateNetworkName(p0 context.Context) (NetworkName, error) {
	if s.Internal.StateNetworkName == nil {
		return *new(NetworkName), ErrNotSupported
	}
	return s.Internal.StateNetworkName(p0)
}

func (s *FullNodeStruct) StateNetworkVersion(p0 context.Context, p1 types.TipSetKey) (NetworkVersion, error) {
	if s.Internal.StateNetworkVersion == nil {
		return *new(NetworkVersion), ErrNotSupported
	}
	return s.Internal.StateNetworkVersion(p0, p1)
}

func (s *FullNodeStruct) StateReadState(p0 context.Context, p1 address.Address, p2 types.TipSetKey) (*ActorState, error) {
	if s.Internal.StateReadState == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.StateReadState(p0, p1, p2)
}

func (s *FullNodeStruct) StateReplay(p0 context.Context, p1 types.TipSetKey, p2 cid.Cid) (*InvocResult, error) {
	if s.Internal.StateReplay == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.StateReplay(p0, p1, p2)
}

func (s *FullNodeStruct) StateSearchMsg(p0 context.Context, p1 types.TipSetKey, p2 cid.Cid, p3 abi.ChainEpoch, p4 bool) (*MsgLookup, error) {
	if s.Internal.StateSearchMsg == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.StateSearchMsg(p0, p1, p2, p3, p4)
}

func (s *FullNodeStruct) StateSectorExpiration(p0 context.Context, p1 address.Address, p2 abi.SectorNumber, p3 types.TipSetKey) (*SectorExpiration, error) {
	if s.Internal.StateSectorExpiration == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.StateSectorExpiration(p0, p1, p2, p3)
}

func (s *FullNodeStruct) StateSectorGetInfo(p0 context.Context, p1 address.Address, p2 abi.SectorNumber, p3 types.TipSetKey) (*miner.SectorOnChainInfo, error) {
	if s.Internal.StateSectorGetInfo == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.StateSectorGetInfo(p0, p1, p2, p3)
}

func (s *FullNodeStruct) StateSectorPartition(p0 context.Context, p1 address.Address, p2 abi.SectorNumber, p3 types.TipSetKey) (*SectorLocation, error) {
	if s.Internal.StateSectorPartition == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.StateSectorPartition(p0, p1, p2, p3)
}

func (s *FullNodeStruct) StateSectorPreCommitInfo(p0 context.Context, p1 address.Address, p2 abi.SectorNumber, p3 types.TipSetKey) (*miner.SectorPreCommitOnChainInfo, error) {
	if s.Internal.StateSectorPreCommitInfo == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.StateSectorPreCommitInfo(p0, p1, p2, p3)
}

func (s *FullNodeStruct) StateVMCirculatingSupplyInternal(p0 context.Context, p1 types.TipSetKey) (CirculatingSupply, error) {
	if s.Internal.StateVMCirculatingSupplyInternal == nil {
		return *new(CirculatingSupply), ErrNotSupported
	}
	return s.Internal.StateVMCirculatingSupplyInternal(p0, p1)
}

func (s *FullNodeStruct) StateVerifiedClientStatus(p0 context.Context, p1 address.Address, p2 types.TipSetKey) (*abi.StoragePower, error) {
	if s.Internal.StateVerifiedClientStatus == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.StateVerifiedClientStatus(p0, p1, p2)
}

func (s *FullNodeStruct) StateVerifiedRegistryRootKey(p0 context.Context, p1 types.TipSetKey) (address.Address, error) {
	if s.Internal.StateVerifiedRegistryRootKey == nil {
		return *new(address.Address), ErrNotSupported
	}
	return s.Internal.StateVerifiedRegistryRootKey(p0, p1)
}

func (s *FullNodeStruct) StateVerifierStatus(p0 context.Context, p1 address.Address, p2 types.TipSetKey) (*abi.StoragePower, error) {
	if s.Internal.StateVerifierStatus == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.StateVerifierStatus(p0, p1, p2)
}

func (s *FullNodeStruct) StateWaitMsg(p0 context.Context, p1 cid.Cid, p2 uint64, p3 abi.ChainEpoch, p4 bool) (*MsgLookup, error) {
	if s.Internal.StateWaitMsg == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.StateWaitMsg(p0, p1, p2, p3, p4)
}

func (s *FullNodeStruct) SyncCheckBad(p0 context.Context, p1 cid.Cid) (string, error) {
	if s.Internal.SyncCheckBad == nil {
		return "", ErrNotSupported
	}
	return s.Internal.SyncCheckBad(p0, p1)
}

func (s *FullNodeStruct) SyncCheckpoint(p0 context.Context, p1 types.TipSetKey) error {
	if s.Internal.SyncCheckpoint == nil {
		return ErrNotSupported
	}
	return s.Internal.SyncCheckpoint(p0, p1)
}

func (s *FullNodeStruct) SyncIncomingBlocks(p0 context.Context) (<-chan *types.BlockHeader, error) {
	if s.Internal.SyncIncomingBlocks == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.SyncIncomingBlocks(p0)
}

func (s *FullNodeStruct) SyncMarkBad(p0 context.Context, p1 cid.Cid) error {
	if s.Internal.SyncMarkBad == nil {
		return ErrNotSupported
	}
	return s.Internal.SyncMarkBad(p0, p1)
}

func (s *FullNodeStruct) SyncState(p0 context.Context) (*SyncState, error) {
	if s.Internal.SyncState == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.SyncState(p0)
}

func (s *FullNodeStruct) SyncSubmitBlock(p0 context.Context, p1 *types.BlockMsg) error {
	if s.Internal.SyncSubmitBlock == nil {
		return ErrNotSupported
	}
	return s.Internal.SyncSubmitBlock(p0, p1)
}

func (s *FullNodeStruct) SyncUnmarkAllBad(p0 context.Context) error {
	if s.Internal.SyncUnmarkAllBad == nil {
		return ErrNotSupported
	}
	return s.Internal.SyncUnmarkAllBad(p0)
}

func (s *FullNodeStruct) SyncUnmarkBad(p0 context.Context, p1 cid.Cid) error {
	if s.Internal.SyncUnmarkBad == nil {
		return ErrNotSupported
	}
	return s.Internal.SyncUnmarkBad(p0, p1)
}

func (s *FullNodeStruct) SyncValidateTipset(p0 context.Context, p1 types.TipSetKey) (bool, error) {
	if s.Internal.SyncValidateTipset == nil {
		return false, ErrNotSupported
	}
	return s.Internal.SyncValidateTipset(p0, p1)
}

func (s *FullNodeStruct) WalletBalance(p0 context.Context, p1 address.Address) (types.BigInt, error) {
	if s.Internal.WalletBalance == nil {
		return *new(types.BigInt), ErrNotSupported
	}
	return s.Internal.WalletBalance(p0, p1)
}

func (s *FullNodeStruct) WalletDefaultAddress(p0 context.Context) (address.Address, error) {
	if s.Internal.WalletDefaultAddress == nil {
		return *new(address.Address), ErrNotSupported
	}
	return s.Internal.WalletDefaultAddress(p0)
}

func (s *FullNodeStruct) WalletDelete(p0 context.Context, p1 address.Address) error {
	if s.Internal.WalletDelete == nil {
		return ErrNotSupported
	}
	return s.Internal.WalletDelete(p0, p1)
}

func (s *FullNodeStruct) WalletExport(p0 context.Context, p1 address.Address) (*types.KeyInfo, error) {
	if s.Internal.WalletExport == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.WalletExport(p0, p1)
}

func (s *FullNodeStruct) WalletHas(p0 context.Context, p1 address.Address) (bool, error) {
	if s.Internal.WalletHas == nil {
		return false, ErrNotSupported
	}
	return s.Internal.WalletHas(p0, p1)
}

func (s *FullNodeStruct) WalletImport(p0 context.Context, p1 *types.KeyInfo) (address.Address, error) {
	if s.Internal.WalletImport == nil {
		return *new(address.Address), ErrNotSupported
	}
	return s.Internal.WalletImport(p0, p1)
}

func (s *FullNodeStruct) WalletList(p0 context.Context) ([]address.Address, error) {
	if s.Internal.WalletList == nil {
		return *new([]address.Address), ErrNotSupported
	}
	return s.Internal.WalletList(p0)
}

func (s *FullNodeStruct) WalletNew(p0 context.Context, p1 types.KeyType) (address.Address, error) {
	if s.Internal.WalletNew == nil {
		return *new(address.Address), ErrNotSupported
	}
	return s.Internal.WalletNew(p0, p1)
}

func (s *FullNodeStruct) WalletSetDefault(p0 context.Context, p1 address.Address) error {
	if s.Internal.WalletSetDefault == nil {
		return ErrNotSupported
	}
	return s.Internal.WalletSetDefault(p0, p1)
}

func (s *FullNodeStruct) WalletSign(p0 context.Context, p1 address.Address, p2 []byte) (*crypto.Signature, error) {
	if s.Internal.WalletSign == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.WalletSign(p0, p1, p2)
}

func (s *FullNodeStruct) WalletSignMessage(p0 context.Context, p1 address.Address, p2 *types.Message) (*types.SignedMessage, error) {
	if s.Internal.WalletSignMessage == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.WalletSignMessage(p0, p1, p2)
}

func (s *FullNodeStruct) WalletValidateAddress(p0 context.Context, p1 string) (address.Address, error) {
	if s.Internal.WalletValidateAddress == nil {
		return *new(address.Address), ErrNotSupported
	}
	return s.Internal.WalletValidateAddress(p0, p1)
}

func (s *FullNodeStruct) WalletVerify(p0 context.Context, p1 address.Address, p2 []byte, p3 *crypto.Signature) (bool, error) {
	if s.Internal.WalletVerify == nil {
		return false, ErrNotSupported
	}
	return s.Internal.WalletVerify(p0, p1, p2, p3)
}

func (s *NetStruct) ID(p0 context.Context) (peer.ID, error) {
	if s.Internal.ID == nil {
		return *new(peer.ID), ErrNotSupported
	}
	return s.Internal.ID(p0)
}

func (s *NetStruct) NetAddrsListen(p0 context.Context) (peer.AddrInfo, error) {
	if s.Internal.NetAddrsListen == nil {
		return *new(peer.AddrInfo), ErrNotSupported
	}
	return s.Internal.NetAddrsListen(p0)
}

func (s *NetStruct) NetAgentVersion(p0 context.Context, p1 peer.ID) (string, error) {
	if s.Internal.NetAgentVersion == nil {
		return "", ErrNotSupported
	}
	return s.Internal.NetAgentVersion(p0, p1)
}

func (s *NetStruct) NetAutoNatStatus(p0 context.Context) (NatInfo, error) {
	if s.Internal.NetAutoNatStatus == nil {
		return *new(NatInfo), ErrNotSupported
	}
	return s.Internal.NetAutoNatStatus(p0)
}

func (s *NetStruct) NetBandwidthStats(p0 context.Context) (metrics.Stats, error) {
	if s.Internal.NetBandwidthStats == nil {
		return *new(metrics.Stats), ErrNotSupported
	}
	return s.Internal.NetBandwidthStats(p0)
}

func (s *NetStruct) NetBandwidthStatsByPeer(p0 context.Context) (map[string]metrics.Stats, error) {
	if s.Internal.NetBandwidthStatsByPeer == nil {
		return *new(map[string]metrics.Stats), ErrNotSupported
	}
	return s.Internal.NetBandwidthStatsByPeer(p0)
}

func (s *NetStruct) NetBandwidthStatsByProtocol(p0 context.Context) (map[protocol.ID]metrics.Stats, error) {
	if s.Internal.NetBandwidthStatsByProtocol == nil {
		return *new(map[protocol.ID]metrics.Stats), ErrNotSupported
	}
	return s.Internal.NetBandwidthStatsByProtocol(p0)
}

func (s *NetStruct) NetBlockAdd(p0 context.Context, p1 NetBlockList) error {
	if s.Internal.NetBlockAdd == nil {
		return ErrNotSupported
	}
	return s.Internal.NetBlockAdd(p0, p1)
}

func (s *NetStruct) NetBlockList(p0 context.Context) (NetBlockList, error) {
	if s.Internal.NetBlockList == nil {
		return *new(NetBlockList), ErrNotSupported
	}
	return s.Internal.NetBlockList(p0)
}

func (s *NetStruct) NetBlockRemove(p0 context.Context, p1 NetBlockList) error {
	if s.Internal.NetBlockRemove == nil {
		return ErrNotSupported
	}
	return s.Internal.NetBlockRemove(p0, p1)
}

func (s *NetStruct) NetConnect(p0 context.Context, p1 peer.AddrInfo) error {
	if s.Internal.NetConnect == nil {
		return ErrNotSupported
	}
	return s.Internal.NetConnect(p0, p1)
}

func (s *NetStruct) NetConnectedness(p0 context.Context, p1 peer.ID) (network.Connectedness, error) {
	if s.Internal.NetConnectedness == nil {
		return *new(network.Connectedness), ErrNotSupported
	}
	return s.Internal.NetConnectedness(p0, p1)
}

func (s *NetStruct) NetDisconnect(p0 context.Context, p1 peer.ID) error {
	if s.Internal.NetDisconnect == nil {
		return ErrNotSupported
	}
	return s.Internal.NetDisconnect(p0, p1)
}

func (s *NetStruct) NetFindPeer(p0 context.Context, p1 peer.ID) (peer.AddrInfo, error) {
	if s.Internal.NetFindPeer == nil {
		return *new(peer.AddrInfo), ErrNotSupported
	}
	return s.Internal.NetFindPeer(p0, p1)
}

func (s *NetStruct) NetLimit(p0 context.Context, p1 string) (NetLimit, error) {
	if s.Internal.NetLimit == nil {
		return *new(NetLimit), ErrNotSupported
	}
	return s.Internal.NetLimit(p0, p1)
}

func (s *NetStruct) NetPeerInfo(p0 context.Context, p1 peer.ID) (*ExtendedPeerInfo, error) {
	if s.Internal.NetPeerInfo == nil {
		return nil, ErrNotSupported
	}
	return s.Internal.NetPeerInfo(p0, p1)
}

func (s *NetStruct) NetPeers(p0 context.Context) ([]peer.AddrInfo, error) {
	if s.Internal.NetPeers == nil {
		return *new([]peer.AddrInfo), ErrNotSupported
	}
	return s.Internal.NetPeers(p0)
}

func (s *NetStruct) NetPing(p0 context.Context, p1 peer.ID) (time.Duration, error) {
	if s.Internal.NetPing == nil {
		return *new(time.Duration), ErrNotSupported
	}
	return s.Internal.NetPing(p0, p1)
}

func (s *NetStruct) NetProtectAdd(p0 context.Context, p1 []peer.ID) error {
	if s.Internal.NetProtectAdd == nil {
		return ErrNotSupported
	}
	return s.Internal.NetProtectAdd(p0, p1)
}

func (s *NetStruct) NetProtectList(p0 context.Context) ([]peer.ID, error) {
	if s.Internal.NetProtectList == nil {
		return *new([]peer.ID), ErrNotSupported
	}
	return s.Internal.NetProtectList(p0)
}
func (s *NetStruct) NetProtectRemove(p0 context.Context, p1 []peer.ID) error {
	if s.Internal.NetProtectRemove == nil {
		return ErrNotSupported
	}
	return s.Internal.NetProtectRemove(p0, p1)
}

func (s *NetStruct) NetPubsubScores(p0 context.Context) ([]PubsubScore, error) {
	if s.Internal.NetPubsubScores == nil {
		return *new([]PubsubScore), ErrNotSupported
	}
	return s.Internal.NetPubsubScores(p0)
}

func (s *NetStruct) NetSetLimit(p0 context.Context, p1 string, p2 NetLimit) error {
	if s.Internal.NetSetLimit == nil {
		return ErrNotSupported
	}
	return s.Internal.NetSetLimit(p0, p1, p2)
}

func (s *NetStruct) NetStat(p0 context.Context, p1 string) (NetStat, error) {
	if s.Internal.NetStat == nil {
		return *new(NetStat), ErrNotSupported
	}
	return s.Internal.NetStat(p0, p1)
}

type SignFunc = func(context.Context, []byte) (*crypto.Signature, error)

func (s *SignableStruct) Sign(p0 context.Context, p1 SignFunc) error {
	if s.Internal.Sign == nil {
		return ErrNotSupported
	}
	return s.Internal.Sign(p0, p1)
}

// APIVersion provides various build-time information
type APIVersion struct {
	Version string

	// APIVersion is a binary encoded semver version of the remote implementing
	// this api
	//
	// See APIVersion in build/version.go
	APIVersion Version

	// TODO: git commit / os / genesis cid?

	// Seconds
	BlockDelay uint64
}

func (v APIVersion) String() string {
	return fmt.Sprintf("%s+api%s", v.Version, v.APIVersion.String())
}

type MessagePrototype struct {
	Message    types.Message
	ValidNonce bool
}

// BlsMessages[x].cid = Cids[x]
// SecpkMessages[y].cid = Cids[BlsMessages.length + y]
type BlockMessages struct {
	BlsMessages   []*types.Message
	SecpkMessages []*types.SignedMessage

	Cids []cid.Cid
}

type Message struct {
	Cid     cid.Cid
	Message *types.Message
}

type IpldObject struct {
	Cid cid.Cid
	Obj interface{}
}

type HeadChange struct {
	Type string
	Val  *types.TipSet
}

type PruneOpts struct {
	MovingGC    bool
	RetainState int64
}

type ActorState struct {
	Balance types.BigInt
	Code    cid.Cid
	State   interface{}
}

type PCHDir int

const (
	PCHUndef PCHDir = iota
	PCHInbound
	PCHOutbound
)

type PaychGetOpts struct {
	OffChain bool
}

type PaychStatus struct {
	ControlAddr address.Address
	Direction   PCHDir
}

type ChannelInfo struct {
	Channel      address.Address
	WaitSentinel cid.Cid
}

type ChannelAvailableFunds struct {
	// Channel is the address of the channel
	Channel *address.Address
	// From is the from address of the channel (channel creator)
	From address.Address
	// To is the to address of the channel
	To address.Address

	// ConfirmedAmt is the total amount of funds that have been confirmed on-chain for the channel
	ConfirmedAmt types.BigInt
	// PendingAmt is the amount of funds that are pending confirmation on-chain
	PendingAmt types.BigInt

	// NonReservedAmt is part of ConfirmedAmt that is available for use (e.g. when the payment channel was pre-funded)
	NonReservedAmt types.BigInt
	// PendingAvailableAmt is the amount of funds that are pending confirmation on-chain that will become available once confirmed
	PendingAvailableAmt types.BigInt

	// PendingWaitSentinel can be used with PaychGetWaitReady to wait for
	// confirmation of pending funds
	PendingWaitSentinel *cid.Cid
	// QueuedAmt is the amount that is queued up behind a pending request
	QueuedAmt types.BigInt

	// VoucherRedeemedAmt is the amount that is redeemed by vouchers on-chain
	// and in the local datastore
	VoucherReedeemedAmt types.BigInt
}

type PaymentInfo struct {
	Channel      address.Address
	WaitSentinel cid.Cid
	Vouchers     []*paych.SignedVoucher
}

type VoucherSpec struct {
	Amount      types.BigInt
	TimeLockMin abi.ChainEpoch
	TimeLockMax abi.ChainEpoch
	MinSettle   abi.ChainEpoch

	Extra *paych.ModVerifyParams
}

// VoucherCreateResult is the response to calling PaychVoucherCreate
type VoucherCreateResult struct {
	// Voucher that was created, or nil if there was an error or if there
	// were insufficient funds in the channel
	Voucher *paych.SignedVoucher
	// Shortfall is the additional amount that would be needed in the channel
	// in order to be able to create the voucher
	Shortfall types.BigInt
}

type MinerPower struct {
	MinerPower  Claim
	TotalPower  Claim
	HasMinPower bool
}

type QueryOffer struct {
	Err string

	Root  cid.Cid
	Piece *cid.Cid

	Size                    uint64
	MinPrice                types.BigInt
	UnsealPrice             types.BigInt
	PricePerByte            abi.TokenAmount
	PaymentInterval         uint64
	PaymentIntervalIncrease uint64
	Miner                   address.Address
	MinerPeer               retrievalmarket.RetrievalPeer
}

func (o *QueryOffer) Order(client address.Address) RetrievalOrder {
	return RetrievalOrder{
		Root:                    o.Root,
		Piece:                   o.Piece,
		Size:                    o.Size,
		Total:                   o.MinPrice,
		UnsealPrice:             o.UnsealPrice,
		PaymentInterval:         o.PaymentInterval,
		PaymentIntervalIncrease: o.PaymentIntervalIncrease,
		Client:                  client,

		Miner:     o.Miner,
		MinerPeer: &o.MinerPeer,
	}
}

type MarketBalance struct {
	Escrow big.Int
	Locked big.Int
}

type MarketDeal struct {
	Proposal market.DealProposal
	State    market.DealState
}
type Selector string

type RetrievalOrder struct {
	Root         cid.Cid
	Piece        *cid.Cid
	DataSelector *Selector

	// todo: Size/Total are only used for calculating price per byte; we should let users just pass that
	Size  uint64
	Total types.BigInt

	UnsealPrice             types.BigInt
	PaymentInterval         uint64
	PaymentIntervalIncrease uint64
	Client                  address.Address
	Miner                   address.Address
	MinerPeer               *retrievalmarket.RetrievalPeer

	RemoteStore *RemoteStoreID `json:"RemoteStore,omitempty"`
}

type RemoteStoreID = uuid.UUID

type InvocResult struct {
	MsgCid         cid.Cid
	Msg            *types.Message
	MsgRct         *types.MessageReceipt
	GasCost        MsgGasCost
	ExecutionTrace types.ExecutionTrace
	Error          string
	Duration       time.Duration
}

type MsgGasCost struct {
	Message            cid.Cid
	GasUsed            abi.TokenAmount
	BaseFeeBurn        abi.TokenAmount
	OverEstimationBurn abi.TokenAmount
	MinerPenalty       abi.TokenAmount
	MinerTip           abi.TokenAmount
	Refund             abi.TokenAmount
	TotalCost          abi.TokenAmount
}

type MethodCall struct {
	types.MessageReceipt
	Error string
}

type StartDealParams struct {
	Data               *storagemarket.DataRef
	Wallet             address.Address
	Miner              address.Address
	EpochPrice         types.BigInt
	MinBlocksDuration  uint64
	ProviderCollateral big.Int
	DealStartEpoch     abi.ChainEpoch
	FastRetrieval      bool
	VerifiedDeal       bool
}

func (s *StartDealParams) UnmarshalJSON(raw []byte) (err error) {
	type sdpAlias StartDealParams

	sdp := sdpAlias{
		FastRetrieval: true,
	}

	if err := json.Unmarshal(raw, &sdp); err != nil {
		return err
	}

	*s = StartDealParams(sdp)

	return nil
}

type ActiveSync struct {
	WorkerID uint64
	Base     *types.TipSet
	Target   *types.TipSet

	Stage  SyncStateStage
	Height abi.ChainEpoch

	Start   time.Time
	End     time.Time
	Message string
}

type SyncState struct {
	ActiveSyncs []ActiveSync

	VMApplied uint64
}

type SyncStateStage int

const (
	StageIdle = SyncStateStage(iota)
	StageHeaders
	StagePersistHeaders
	StageMessages
	StageSyncComplete
	StageSyncErrored
	StageFetchingMessages
)

func (v SyncStateStage) String() string {
	switch v {
	case StageIdle:
		return "idle"
	case StageHeaders:
		return "header sync"
	case StagePersistHeaders:
		return "persisting headers"
	case StageMessages:
		return "message sync"
	case StageSyncComplete:
		return "complete"
	case StageSyncErrored:
		return "error"
	case StageFetchingMessages:
		return "fetching messages"
	default:
		return fmt.Sprintf("<unknown: %d>", v)
	}
}

type MpoolChange int

const (
	MpoolAdd MpoolChange = iota
	MpoolRemove
)

type MpoolUpdate struct {
	Type    MpoolChange
	Message *types.SignedMessage
}

type ComputeStateOutput struct {
	Root  cid.Cid
	Trace []*InvocResult
}

type DealCollateralBounds struct {
	Min abi.TokenAmount
	Max abi.TokenAmount
}

type CirculatingSupply struct {
	FilVested           abi.TokenAmount
	FilMined            abi.TokenAmount
	FilBurnt            abi.TokenAmount
	FilLocked           abi.TokenAmount
	FilCirculating      abi.TokenAmount
	FilReserveDisbursed abi.TokenAmount
}

type MiningBaseInfo struct {
	MinerPower        types.BigInt
	NetworkPower      types.BigInt
	Sectors           []proof.ExtendedSectorInfo
	WorkerKey         address.Address
	SectorSize        abi.SectorSize
	PrevBeaconEntry   types.BeaconEntry
	BeaconEntries     []types.BeaconEntry
	EligibleForMining bool
}

type BlockTemplate struct {
	Miner            address.Address
	Parents          types.TipSetKey
	Ticket           *types.Ticket
	Eproof           *types.ElectionProof
	BeaconValues     []types.BeaconEntry
	Messages         []*types.SignedMessage
	Epoch            abi.ChainEpoch
	Timestamp        uint64
	WinningPoStProof []proof.PoStProof
}

type DataSize struct {
	PayloadSize int64
	PieceSize   abi.PaddedPieceSize
}

type DataCIDSize struct {
	PayloadSize int64
	PieceSize   abi.PaddedPieceSize
	PieceCID    cid.Cid
}

type CommPRet struct {
	Root cid.Cid
	Size abi.UnpaddedPieceSize
}

type MsigProposeResponse int

const (
	MsigApprove MsigProposeResponse = iota
	MsigCancel
)

type Deadline struct {
	PostSubmissions      bitfield.BitField
	DisputableProofCount uint64
}

type Partition struct {
	AllSectors        bitfield.BitField
	FaultySectors     bitfield.BitField
	RecoveringSectors bitfield.BitField
	LiveSectors       bitfield.BitField
	ActiveSectors     bitfield.BitField
}

type Fault struct {
	Miner address.Address
	Epoch abi.ChainEpoch
}

var EmptyVesting = MsigVesting{
	InitialBalance: types.EmptyInt,
	StartEpoch:     -1,
	UnlockDuration: -1,
}

type MsigVesting struct {
	InitialBalance abi.TokenAmount
	StartEpoch     abi.ChainEpoch
	UnlockDuration abi.ChainEpoch
}

type MessageMatch struct {
	To   address.Address
	From address.Address
}

type MsigTransaction struct {
	ID     int64
	To     address.Address
	Value  abi.TokenAmount
	Method abi.MethodNum
	Params []byte

	Approved []address.Address
}

type DealInfo struct {
	ProposalCid cid.Cid
	State       storagemarket.StorageDealStatus
	Message     string // more information about deal state, particularly errors
	DealStages  *storagemarket.DealStages
	Provider    address.Address

	DataRef  *storagemarket.DataRef
	PieceCID cid.Cid
	Size     uint64

	PricePerEpoch types.BigInt
	Duration      uint64

	DealID abi.DealID

	CreationTime time.Time
	Verified     bool

	TransferChannelID *datatransfer.ChannelID
	DataTransfer      *DataTransferChannel
}

type NatInfo struct {
	Reachability network.Reachability
	PublicAddr   string
}
type MsgLookup struct {
	Message   cid.Cid
	Receipt   types.MessageReceipt
	ReturnDec interface{}
	TipSet    types.TipSetKey
	Height    abi.ChainEpoch
}

type MinerSectors struct {
	// Live sectors that should be proven.
	Live uint64
	// Sectors actively contributing to power.
	Active uint64
	// Sectors with failed proofs.
	Faulty uint64
}

type MinerInfo struct {
	Owner                      address.Address   // Must be an ID-address.
	Worker                     address.Address   // Must be an ID-address.
	NewWorker                  address.Address   // Must be an ID-address.
	ControlAddresses           []address.Address // Must be an ID-addresses.
	WorkerChangeEpoch          abi.ChainEpoch
	PeerId                     *peer.ID
	Multiaddrs                 []abi.Multiaddrs
	WindowPoStProofType        abi.RegisteredPoStProof
	SectorSize                 abi.SectorSize
	WindowPoStPartitionSectors uint64
	ConsensusFaultElapsed      abi.ChainEpoch
	Beneficiary                address.Address
	BeneficiaryTerm            *miner.BeneficiaryTerm
	PendingBeneficiaryTerm     *miner.PendingBeneficiaryChange
}

type NetworkVersion = abinetwork.Version
type OpenRPCDocument map[string]interface{}
