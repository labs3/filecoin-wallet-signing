package api

import (
	"encoding/json"
	"time"

	datatransfer "github.com/filecoin-project/go-data-transfer"
	"github.com/filecoin-project/go-fil-markets/retrievalmarket"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/google/uuid"
	"github.com/ipfs/go-cid"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
)

type NetStat struct {
	System    *network.ScopeStat           `json:",omitempty"`
	Transient *network.ScopeStat           `json:",omitempty"`
	Services  map[string]network.ScopeStat `json:",omitempty"`
	Protocols map[string]network.ScopeStat `json:",omitempty"`
	Peers     map[string]network.ScopeStat `json:",omitempty"`
}

type NetLimit struct {
	Memory int64 `json:",omitempty"`

	Streams, StreamsInbound, StreamsOutbound int
	Conns, ConnsInbound, ConnsOutbound       int
	FD                                       int
}

type PubsubScore struct {
	ID    peer.ID
	Score *pubsub.PeerScoreSnapshot
}

type ExtendedPeerInfo struct {
	ID          peer.ID
	Agent       string
	Addrs       []string
	Protocols   []string
	ConnMgrMeta *ConnMgrInfo
}

type ConnMgrInfo struct {
	FirstSeen time.Time
	Value     int
	Tags      map[string]int
	Conns     map[string]time.Time
}

type NetBlockList struct {
	Peers     []peer.ID
	IPAddrs   []string
	IPSubnets []string
}

type MessageCheckStatus struct {
	Cid cid.Cid
	CheckStatus
}

type CheckStatus struct {
	Code CheckStatusCode
	OK   bool
	Err  string
	Hint map[string]interface{}
}

type CheckStatusCode int

const (
	_ CheckStatusCode = iota
	// Message Checks
	CheckStatusMessageSerialize
	CheckStatusMessageSize
	CheckStatusMessageValidity
	CheckStatusMessageMinGas
	CheckStatusMessageMinBaseFee
	CheckStatusMessageBaseFee
	CheckStatusMessageBaseFeeLowerBound
	CheckStatusMessageBaseFeeUpperBound
	CheckStatusMessageGetStateNonce
	CheckStatusMessageNonce
	CheckStatusMessageGetStateBalance
	CheckStatusMessageBalance
)

type MessageSendSpec struct {
	MaxFee  abi.TokenAmount
	MsgUuid uuid.UUID
}

type RestrievalRes struct {
	DealID retrievalmarket.DealID
}

type RetrievalInfo struct {
	PayloadCID   cid.Cid
	ID           retrievalmarket.DealID
	PieceCID     *cid.Cid
	PricePerByte abi.TokenAmount
	UnsealPrice  abi.TokenAmount

	Status        retrievalmarket.DealStatus
	Message       string // more information about deal state, particularly errors
	Provider      peer.ID
	BytesReceived uint64
	BytesPaidFor  uint64
	TotalPaid     abi.TokenAmount

	TransferChannelID *datatransfer.ChannelID
	DataTransfer      *DataTransferChannel

	// optional event if part of ClientGetRetrievalUpdates
	Event *retrievalmarket.ClientEvent
}

type NetworkName string

// AlertType is a unique alert identifier
type AlertType struct {
	System, Subsystem string
}

// AlertEvent contains information about alert state transition
type AlertEvent struct {
	Type    string // either 'raised' or 'resolved'
	Message json.RawMessage
	Time    time.Time
}

type Alert struct {
	Type   AlertType
	Active bool

	LastActive   *AlertEvent // NOTE: pointer for nullability, don't mutate the referenced object!
	LastResolved *AlertEvent

	journalType EventType
}

// EventType represents the signature of an event.
type EventType struct {
	System string
	Event  string

	// enabled stores whether this event type is enabled.
	enabled bool

	// safe is a sentinel marker that's set to true if this EventType was
	// constructed correctly (via Journal#RegisterEventType).
	safe bool
}

func (et EventType) String() string {
	return et.System + ":" + et.Event
}

type ID uint64

type Claim struct {
	// Sum of raw byte power for a miner's sectors.
	RawBytePower abi.StoragePower

	// Sum of quality adjusted power for a miner's sectors.
	QualityAdjPower abi.StoragePower
}

type SectorExpiration struct {
	OnTime abi.ChainEpoch

	// non-zero if sector is faulty, epoch at which it will be permanently
	// removed if it doesn't recover
	Early abi.ChainEpoch
}

type SectorLocation struct {
	Deadline  uint64
	Partition uint64
}
