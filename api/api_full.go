package api

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/consensus-shipyard/go-ipc-types/gateway"
	"github.com/consensus-shipyard/go-ipc-types/sdk"
	"github.com/consensus-shipyard/go-ipc-types/subnetactor"
	"github.com/google/uuid"
	"github.com/ipfs/go-cid"
	blocks "github.com/ipfs/go-libipfs/blocks"
	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-bitfield"
	datatransfer "github.com/filecoin-project/go-data-transfer/v2"
	"github.com/filecoin-project/go-fil-markets/retrievalmarket"
	"github.com/filecoin-project/go-fil-markets/storagemarket"
	"github.com/filecoin-project/go-jsonrpc"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/filecoin-project/go-state-types/builtin/v8/paych"
	"github.com/filecoin-project/go-state-types/builtin/v9/market"
	"github.com/filecoin-project/go-state-types/builtin/v9/miner"
	verifregtypes "github.com/filecoin-project/go-state-types/builtin/v9/verifreg"
	"github.com/filecoin-project/go-state-types/crypto"
	"github.com/filecoin-project/go-state-types/dline"
	abinetwork "github.com/filecoin-project/go-state-types/network"

	apitypes "github.com/filecoin-project/lotus/api/types"
	"github.com/filecoin-project/lotus/chain/actors/builtin"
	lminer "github.com/filecoin-project/lotus/chain/actors/builtin/miner"
	"github.com/filecoin-project/lotus/chain/actors/builtin/power"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/filecoin-project/lotus/chain/types/ethtypes"
	"github.com/filecoin-project/lotus/node/modules/dtypes"
	"github.com/filecoin-project/lotus/node/repo/imports"
)

//go:generate go run github.com/golang/mock/mockgen -destination=mocks/mock_full.go -package=mocks . FullNode

// ChainIO abstracts operations for accessing raw IPLD objects.
type ChainIO interface {
	ChainReadObj(context.Context, cid.Cid) ([]byte, error)
	ChainHasObj(context.Context, cid.Cid) (bool, error)
	ChainPutObj(context.Context, blocks.Block) error
}

const LookbackNoLimit = abi.ChainEpoch(-1)

//                       MODIFYING THE API INTERFACE
//
// NOTE: This is the V1 (Unstable) API - to add methods to the V0 (Stable) API
// you'll have to add those methods to interfaces in `api/v0api`
//
// When adding / changing methods in this file:
// * Do the change here
// * Adjust implementation in `node/impl/`
// * Run `make gen` - this will:
//  * Generate proxy structs
//  * Generate mocks
//  * Generate markdown docs
//  * Generate openrpc blobs

// FullNode API is a low-level interface to the Filecoin network full node
type FullNode interface {
	Common
	Net

	// MethodGroup: Chain
	// The Chain method group contains methods for interacting with the
	// blockchain, but that do not require any form of state computation.

	// ChainNotify returns channel with chain head updates.
	// First message is guaranteed to be of len == 1, and type == 'current'.
	ChainNotify(context.Context) (<-chan []*HeadChange, error) //perm:read

	// ChainHead returns the current head of the chain.
	ChainHead(context.Context) (*types.TipSet, error) //perm:read

	// ChainGetBlock returns the block specified by the given CID.
	ChainGetBlock(context.Context, cid.Cid) (*types.BlockHeader, error) //perm:read
	// ChainGetTipSet returns the tipset specified by the given TipSetKey.
	ChainGetTipSet(context.Context, types.TipSetKey) (*types.TipSet, error) //perm:read

	// ChainGetBlockMessages returns messages stored in the specified block.
	//
	// Note: If there are multiple blocks in a tipset, it's likely that some
	// messages will be duplicated. It's also possible for blocks in a tipset to have
	// different messages from the same sender at the same nonce. When that happens,
	// only the first message (in a block with lowest ticket) will be considered
	// for execution
	//
	// NOTE: THIS METHOD SHOULD ONLY BE USED FOR GETTING MESSAGES IN A SPECIFIC BLOCK
	//
	// DO NOT USE THIS METHOD TO GET MESSAGES INCLUDED IN A TIPSET
	// Use ChainGetParentMessages, which will perform correct message deduplication
	ChainGetBlockMessages(ctx context.Context, blockCid cid.Cid) (*BlockMessages, error) //perm:read

	// ChainGetParentReceipts returns receipts for messages in parent tipset of
	// the specified block. The receipts in the list returned is one-to-one with the
	// messages returned by a call to ChainGetParentMessages with the same blockCid.
	ChainGetParentReceipts(ctx context.Context, blockCid cid.Cid) ([]*types.MessageReceipt, error) //perm:read

	// ChainGetParentMessages returns messages stored in parent tipset of the
	// specified block.
	ChainGetParentMessages(ctx context.Context, blockCid cid.Cid) ([]Message, error) //perm:read

	// ChainGetMessagesInTipset returns message stores in current tipset
	ChainGetMessagesInTipset(ctx context.Context, tsk types.TipSetKey) ([]Message, error) //perm:read

	// ChainGetTipSetByHeight looks back for a tipset at the specified epoch.
	// If there are no blocks at the specified epoch, a tipset at an earlier epoch
	// will be returned.
	ChainGetTipSetByHeight(context.Context, abi.ChainEpoch, types.TipSetKey) (*types.TipSet, error) //perm:read

	// ChainGetTipSetAfterHeight looks back for a tipset at the specified epoch.
	// If there are no blocks at the specified epoch, the first non-nil tipset at a later epoch
	// will be returned.
	ChainGetTipSetAfterHeight(context.Context, abi.ChainEpoch, types.TipSetKey) (*types.TipSet, error) //perm:read

	// ChainReadObj reads ipld nodes referenced by the specified CID from chain
	// blockstore and returns raw bytes.
	ChainReadObj(context.Context, cid.Cid) ([]byte, error) //perm:read

	// ChainDeleteObj deletes node referenced by the given CID
	ChainDeleteObj(context.Context, cid.Cid) error //perm:admin

	// ChainHasObj checks if a given CID exists in the chain blockstore.
	ChainHasObj(context.Context, cid.Cid) (bool, error) //perm:read

	// ChainPutObj puts a given object into the block store
	ChainPutObj(context.Context, blocks.Block) error //perm:admin

	// ChainStatObj returns statistics about the graph referenced by 'obj'.
	// If 'base' is also specified, then the returned stat will be a diff
	// between the two objects.
	ChainStatObj(ctx context.Context, obj cid.Cid, base cid.Cid) (ObjStat, error) //perm:read

	// ChainSetHead forcefully sets current chain head. Use with caution.
	ChainSetHead(context.Context, types.TipSetKey) error //perm:admin

	// ChainGetGenesis returns the genesis tipset.
	ChainGetGenesis(context.Context) (*types.TipSet, error) //perm:read

	// ChainTipSetWeight computes weight for the specified tipset.
	ChainTipSetWeight(context.Context, types.TipSetKey) (types.BigInt, error) //perm:read
	ChainGetNode(ctx context.Context, p string) (*IpldObject, error)          //perm:read

	// ChainGetMessage reads a message referenced by the specified CID from the
	// chain blockstore.
	ChainGetMessage(context.Context, cid.Cid) (*types.Message, error) //perm:read

	// ChainGetPath returns a set of revert/apply operations needed to get from
	// one tipset to another, for example:
	// ```
	//        to
	//         ^
	// from   tAA
	//   ^     ^
	// tBA    tAB
	//  ^---*--^
	//      ^
	//     tRR
	// ```
	// Would return `[revert(tBA), apply(tAB), apply(tAA)]`
	ChainGetPath(ctx context.Context, from types.TipSetKey, to types.TipSetKey) ([]*HeadChange, error) //perm:read

	// ChainExport returns a stream of bytes with CAR dump of chain data.
	// The exported chain data includes the header chain from the given tipset
	// back to genesis, the entire genesis state, and the most recent 'nroots'
	// state trees.
	// If oldmsgskip is set, messages from before the requested roots are also not included.
	ChainExport(ctx context.Context, nroots abi.ChainEpoch, oldmsgskip bool, tsk types.TipSetKey) (<-chan []byte, error) //perm:read

	// ChainExportRangeInternal triggers the export of a chain
	// CAR-snapshot directly to disk. It is similar to ChainExport,
	// except, depending on options, the snapshot can include receipts,
	// messages and stateroots for the length between the specified head
	// and tail, thus producing "archival-grade" snapshots that include
	// all the on-chain data.  The header chain is included back to
	// genesis and these snapshots can be used to initialize Filecoin
	// nodes.
	ChainExportRangeInternal(ctx context.Context, head, tail types.TipSetKey, cfg ChainExportConfig) error //perm:admin

	// ChainPrune forces compaction on cold store and garbage collects; only supported if you
	// are using the splitstore
	ChainPrune(ctx context.Context, opts PruneOpts) error //perm:admin

	// ChainHotGC does online (badger) GC on the hot store; only supported if you are using
	// the splitstore
	ChainHotGC(ctx context.Context, opts HotGCOpts) error //perm:admin

	// ChainCheckBlockstore performs an (asynchronous) health check on the chain/state blockstore
	// if supported by the underlying implementation.
	ChainCheckBlockstore(context.Context) error //perm:admin

	// ChainBlockstoreInfo returns some basic information about the blockstore
	ChainBlockstoreInfo(context.Context) (map[string]interface{}, error) //perm:read

	// ChainGetEvents returns the events under an event AMT root CID.
	ChainGetEvents(context.Context, cid.Cid) ([]types.Event, error) //perm:read

	// GasEstimateFeeCap estimates gas fee cap
	GasEstimateFeeCap(context.Context, *types.Message, int64, types.TipSetKey) (types.BigInt, error) //perm:read

	// GasEstimateGasLimit estimates gas used by the message and returns it.
	// It fails if message fails to execute.
	GasEstimateGasLimit(context.Context, *types.Message, types.TipSetKey) (int64, error) //perm:read

	// GasEstimateGasPremium estimates what gas price should be used for a
	// message to have high likelihood of inclusion in `nblocksincl` epochs.

	GasEstimateGasPremium(_ context.Context, nblocksincl uint64,
		sender address.Address, gaslimit int64, tsk types.TipSetKey) (types.BigInt, error) //perm:read

	// GasEstimateMessageGas estimates gas values for unset message gas fields
	GasEstimateMessageGas(context.Context, *types.Message, *MessageSendSpec, types.TipSetKey) (*types.Message, error) //perm:read

	// MethodGroup: Sync
	// The Sync method group contains methods for interacting with and
	// observing the lotus sync service.

	// SyncState returns the current status of the lotus sync system.
	SyncState(context.Context) (*SyncState, error) //perm:read

	// SyncSubmitBlock can be used to submit a newly created block to the.
	// network through this node
	SyncSubmitBlock(ctx context.Context, blk *types.BlockMsg) error //perm:write

	// SyncPurgeForRecovery "forgets" all state after a Mir checkpoint to create
	// a clean slate from which the daemon can sync according to the
	// checkpoint provided by Mir. This functions moves the chain head to
	// the height pointed by the checkpoint.
	SyncPurgeForRecovery(ctx context.Context, height abi.ChainEpoch) error //perm:write

	// SyncIncomingBlocks returns a channel streaming incoming, potentially not
	// yet synced block headers.
	SyncIncomingBlocks(ctx context.Context) (<-chan *types.BlockHeader, error) //perm:read

	// SyncCheckpoint marks a blocks as checkpointed, meaning that it won't ever fork away from it.
	SyncCheckpoint(ctx context.Context, tsk types.TipSetKey) error //perm:admin

	// SyncMarkBad marks a blocks as bad, meaning that it won't ever by synced.
	// Use with extreme caution.
	SyncMarkBad(ctx context.Context, bcid cid.Cid) error //perm:admin

	// SyncUnmarkBad unmarks a blocks as bad, making it possible to be validated and synced again.
	SyncUnmarkBad(ctx context.Context, bcid cid.Cid) error //perm:admin

	// SyncUnmarkAllBad purges bad block cache, making it possible to sync to chains previously marked as bad
	SyncUnmarkAllBad(ctx context.Context) error //perm:admin

	// SyncCheckBad checks if a block was marked as bad, and if it was, returns
	// the reason.
	SyncCheckBad(ctx context.Context, bcid cid.Cid) (string, error) //perm:read

	// SyncValidateTipset indicates whether the provided tipset is valid or not
	SyncValidateTipset(ctx context.Context, tsk types.TipSetKey) (bool, error) //perm:read

	// SyncFetchTipSet fetches a tipSet from the specific peer (the same way the hello
	// protocol would do) and informs the syncer if it advances our view of the chain.
	SyncFetchTipSetFromPeer(ctx context.Context, p peer.ID, tsk types.TipSetKey) (*types.TipSet, error) //perm:read

	// MethodGroup: Mpool
	// The Mpool methods are for interacting with the message pool. The message pool
	// manages all incoming and outgoing 'messages' going over the network.

	// MpoolPending returns pending mempool messages.
	MpoolPending(context.Context, types.TipSetKey) ([]*types.SignedMessage, error) //perm:read

	// MpoolSelect returns a list of pending messages for inclusion in the next block
	MpoolSelect(context.Context, types.TipSetKey, float64) ([]*types.SignedMessage, error) //perm:read

	// MpoolPush pushes a signed message to mempool.
	MpoolPush(context.Context, *types.SignedMessage) (cid.Cid, error) //perm:write

	// MpoolPushUntrusted pushes a signed message to mempool from untrusted sources.
	MpoolPushUntrusted(context.Context, *types.SignedMessage) (cid.Cid, error) //perm:write

	// MpoolPushMessage atomically assigns a nonce, signs, and pushes a message
	// to mempool.
	// maxFee is only used when GasFeeCap/GasPremium fields aren't specified
	//
	// When maxFee is set to 0, MpoolPushMessage will guess appropriate fee
	// based on current chain conditions
	MpoolPushMessage(ctx context.Context, msg *types.Message, spec *MessageSendSpec) (*types.SignedMessage, error) //perm:sign

	// MpoolBatchPush batch pushes a signed message to mempool.
	MpoolBatchPush(context.Context, []*types.SignedMessage) ([]cid.Cid, error) //perm:write

	// MpoolBatchPushUntrusted batch pushes a signed message to mempool from untrusted sources.
	MpoolBatchPushUntrusted(context.Context, []*types.SignedMessage) ([]cid.Cid, error) //perm:write

	// MpoolBatchPushMessage batch pushes a unsigned message to mempool.
	MpoolBatchPushMessage(context.Context, []*types.Message, *MessageSendSpec) ([]*types.SignedMessage, error) //perm:sign

	// MpoolCheckMessages performs logical checks on a batch of messages
	MpoolCheckMessages(context.Context, []*MessagePrototype) ([][]MessageCheckStatus, error) //perm:read
	// MpoolCheckPendingMessages performs logical checks for all pending messages from a given address
	MpoolCheckPendingMessages(context.Context, address.Address) ([][]MessageCheckStatus, error) //perm:read
	// MpoolCheckReplaceMessages performs logical checks on pending messages with replacement
	MpoolCheckReplaceMessages(context.Context, []*types.Message) ([][]MessageCheckStatus, error) //perm:read

	// MpoolGetNonce gets next nonce for the specified sender.
	// Note that this method may not be atomic. Use MpoolPushMessage instead.
	MpoolGetNonce(context.Context, address.Address) (uint64, error) //perm:read
	MpoolSub(context.Context) (<-chan MpoolUpdate, error)           //perm:read

	// MpoolClear clears pending messages from the mpool
	MpoolClear(context.Context, bool) error //perm:write

	// MpoolGetConfig returns (a copy of) the current mpool config
	MpoolGetConfig(context.Context) (*types.MpoolConfig, error) //perm:read
	// MpoolSetConfig sets the mpool config to (a copy of) the supplied config
	MpoolSetConfig(context.Context, *types.MpoolConfig) error //perm:admin

	// MethodGroup: Miner

	MinerGetBaseInfo(context.Context, address.Address, abi.ChainEpoch, types.TipSetKey) (*MiningBaseInfo, error) //perm:read
	MinerCreateBlock(context.Context, *BlockTemplate) (*types.BlockMsg, error)                                   //perm:write

	// // UX ?

	// MethodGroup: WalletF

	// WalletNew creates a new address in the wallet with the given sigType.
	// Available key types: bls, secp256k1, secp256k1-ledger
	// Support for numerical types: 1 - secp256k1, 2 - BLS is deprecated
	WalletNew(context.Context, types.KeyType) (address.Address, error) //perm:write
	// WalletHas indicates whether the given address is in the wallet.
	WalletHas(context.Context, address.Address) (bool, error) //perm:write
	// WalletList lists all the addresses in the wallet.
	WalletList(context.Context) ([]address.Address, error) //perm:write
	// WalletBalance returns the balance of the given address at the current head of the chain.
	WalletBalance(context.Context, address.Address) (types.BigInt, error) //perm:read
	// WalletSign signs the given bytes using the given address.
	WalletSign(context.Context, address.Address, []byte) (*crypto.Signature, error) //perm:sign
	// WalletSignMessage signs the given message using the given address.
	WalletSignMessage(context.Context, address.Address, *types.Message) (*types.SignedMessage, error) //perm:sign
	// WalletVerify takes an address, a signature, and some bytes, and indicates whether the signature is valid.
	// The address does not have to be in the wallet.
	WalletVerify(context.Context, address.Address, []byte, *crypto.Signature) (bool, error) //perm:read
	// WalletDefaultAddress returns the address marked as default in the wallet.
	WalletDefaultAddress(context.Context) (address.Address, error) //perm:write
	// WalletSetDefault marks the given address as as the default one.
	WalletSetDefault(context.Context, address.Address) error //perm:write
	// WalletExport returns the private key of an address in the wallet.
	WalletExport(context.Context, address.Address) (*types.KeyInfo, error) //perm:admin
	// WalletImport receives a KeyInfo, which includes a private key, and imports it into the wallet.
	WalletImport(context.Context, *types.KeyInfo) (address.Address, error) //perm:admin
	// WalletDelete deletes an address from the wallet.
	WalletDelete(context.Context, address.Address) error //perm:admin
	// WalletValidateAddress validates whether a given string can be decoded as a well-formed address
	WalletValidateAddress(context.Context, string) (address.Address, error) //perm:read

	// Other

	// MethodGroup: Client
	// The Client methods all have to do with interacting with the storage and
	// retrieval markets as a client

	// ClientImport imports file under the specified path into filestore.
	ClientImport(ctx context.Context, ref FileRef) (*ImportRes, error) //perm:admin
	// ClientRemoveImport removes file import
	ClientRemoveImport(ctx context.Context, importID imports.ID) error //perm:admin
	// ClientStartDeal proposes a deal with a miner.
	ClientStartDeal(ctx context.Context, params *StartDealParams) (*cid.Cid, error) //perm:admin
	// ClientStatelessDeal fire-and-forget-proposes an offline deal to a miner without subsequent tracking.
	ClientStatelessDeal(ctx context.Context, params *StartDealParams) (*cid.Cid, error) //perm:write
	// ClientGetDealInfo returns the latest information about a given deal.
	ClientGetDealInfo(context.Context, cid.Cid) (*DealInfo, error) //perm:read
	// ClientListDeals returns information about the deals made by the local client.
	ClientListDeals(ctx context.Context) ([]DealInfo, error) //perm:write
	// ClientGetDealUpdates returns the status of updated deals
	ClientGetDealUpdates(ctx context.Context) (<-chan DealInfo, error) //perm:write
	// ClientGetDealStatus returns status given a code
	ClientGetDealStatus(ctx context.Context, statusCode uint64) (string, error) //perm:read
	// ClientHasLocal indicates whether a certain CID is locally stored.
	ClientHasLocal(ctx context.Context, root cid.Cid) (bool, error) //perm:write
	// ClientFindData identifies peers that have a certain file, and returns QueryOffers (one per peer).
	ClientFindData(ctx context.Context, root cid.Cid, piece *cid.Cid) ([]QueryOffer, error) //perm:read
	// ClientMinerQueryOffer returns a QueryOffer for the specific miner and file.
	ClientMinerQueryOffer(ctx context.Context, miner address.Address, root cid.Cid, piece *cid.Cid) (QueryOffer, error) //perm:read
	// ClientRetrieve initiates the retrieval of a file, as specified in the order.
	ClientRetrieve(ctx context.Context, params RetrievalOrder) (*RestrievalRes, error) //perm:admin
	// ClientRetrieveWait waits for retrieval to be complete
	ClientRetrieveWait(ctx context.Context, deal retrievalmarket.DealID) error //perm:admin
	// ClientExport exports a file stored in the local filestore to a system file
	ClientExport(ctx context.Context, exportRef ExportRef, fileRef FileRef) error //perm:admin
	// ClientListRetrievals returns information about retrievals made by the local client
	ClientListRetrievals(ctx context.Context) ([]RetrievalInfo, error) //perm:write
	// ClientGetRetrievalUpdates returns status of updated retrieval deals
	ClientGetRetrievalUpdates(ctx context.Context) (<-chan RetrievalInfo, error) //perm:write
	// ClientQueryAsk returns a signed StorageAsk from the specified miner.
	ClientQueryAsk(ctx context.Context, p peer.ID, miner address.Address) (*StorageAsk, error) //perm:read
	// ClientCalcCommP calculates the CommP and data size of the specified CID
	ClientDealPieceCID(ctx context.Context, root cid.Cid) (DataCIDSize, error) //perm:read
	// ClientCalcCommP calculates the CommP for a specified file
	ClientCalcCommP(ctx context.Context, inpath string) (*CommPRet, error) //perm:write
	// ClientGenCar generates a CAR file for the specified file.
	ClientGenCar(ctx context.Context, ref FileRef, outpath string) error //perm:write
	// ClientDealSize calculates real deal data size
	ClientDealSize(ctx context.Context, root cid.Cid) (DataSize, error) //perm:read
	// ClientListTransfers returns the status of all ongoing transfers of data
	ClientListDataTransfers(ctx context.Context) ([]DataTransferChannel, error)        //perm:write
	ClientDataTransferUpdates(ctx context.Context) (<-chan DataTransferChannel, error) //perm:write
	// ClientRestartDataTransfer attempts to restart a data transfer with the given transfer ID and other peer
	ClientRestartDataTransfer(ctx context.Context, transferID datatransfer.TransferID, otherPeer peer.ID, isInitiator bool) error //perm:write
	// ClientCancelDataTransfer cancels a data transfer with the given transfer ID and other peer
	ClientCancelDataTransfer(ctx context.Context, transferID datatransfer.TransferID, otherPeer peer.ID, isInitiator bool) error //perm:write
	// ClientRetrieveTryRestartInsufficientFunds attempts to restart stalled retrievals on a given payment channel
	// which are stuck due to insufficient funds
	ClientRetrieveTryRestartInsufficientFunds(ctx context.Context, paymentChannel address.Address) error //perm:write

	// ClientCancelRetrievalDeal cancels an ongoing retrieval deal based on DealID
	ClientCancelRetrievalDeal(ctx context.Context, dealid retrievalmarket.DealID) error //perm:write

	// ClientUnimport removes references to the specified file from filestore
	// ClientUnimport(path string)

	// ClientListImports lists imported files and their root CIDs
	ClientListImports(ctx context.Context) ([]Import, error) //perm:write

	// ClientListAsks() []Ask

	// MethodGroup: State
	// The State methods are used to query, inspect, and interact with chain state.
	// Most methods take a TipSetKey as a parameter. The state looked up is the parent state of the tipset.
	// A nil TipSetKey can be provided as a param, this will cause the heaviest tipset in the chain to be used.

	// StateCall runs the given message and returns its result without any persisted changes.
	//
	// StateCall applies the message to the tipset's parent state. The
	// message is not applied on-top-of the messages in the passed-in
	// tipset.
	StateCall(context.Context, *types.Message, types.TipSetKey) (*InvocResult, error) //perm:read
	// StateReplay replays a given message, assuming it was included in a block in the specified tipset.
	//
	// If a tipset key is provided, and a replacing message is not found on chain,
	// the method will return an error saying that the message wasn't found
	//
	// If no tipset key is provided, the appropriate tipset is looked up, and if
	// the message was gas-repriced, the on-chain message will be replayed - in
	// that case the returned InvocResult.MsgCid will not match the Cid param
	//
	// If the caller wants to ensure that exactly the requested message was executed,
	// they MUST check that InvocResult.MsgCid is equal to the provided Cid.
	// Without this check both the requested and original message may appear as
	// successfully executed on-chain, which may look like a double-spend.
	//
	// A replacing message is a message with a different CID, any of Gas values, and
	// different signature, but with all other parameters matching (source/destination,
	// nonce, params, etc.)
	StateReplay(context.Context, types.TipSetKey, cid.Cid) (*InvocResult, error) //perm:read
	// StateGetActor returns the indicated actor's nonce and balance.
	StateGetActor(ctx context.Context, actor address.Address, tsk types.TipSetKey) (*types.Actor, error) //perm:read
	// StateReadState returns the indicated actor's state.
	StateReadState(ctx context.Context, actor address.Address, tsk types.TipSetKey) (*ActorState, error) //perm:read
	// StateListMessages looks back and returns all messages with a matching to or from address, stopping at the given height.
	StateListMessages(ctx context.Context, match *MessageMatch, tsk types.TipSetKey, toht abi.ChainEpoch) ([]cid.Cid, error) //perm:read
	// StateDecodeParams attempts to decode the provided params, based on the recipient actor address and method number.
	StateDecodeParams(ctx context.Context, toAddr address.Address, method abi.MethodNum, params []byte, tsk types.TipSetKey) (interface{}, error) //perm:read
	// StateEncodeParams attempts to encode the provided json params to the binary from
	StateEncodeParams(ctx context.Context, toActCode cid.Cid, method abi.MethodNum, params json.RawMessage) ([]byte, error) //perm:read

	// StateNetworkName returns the name of the network the node is synced to
	StateNetworkName(context.Context) (dtypes.NetworkName, error) //perm:read
	// StateMinerSectors returns info about the given miner's sectors. If the filter bitfield is nil, all sectors are included.
	StateMinerSectors(context.Context, address.Address, *bitfield.BitField, types.TipSetKey) ([]*miner.SectorOnChainInfo, error) //perm:read
	// StateMinerActiveSectors returns info about sectors that a given miner is actively proving.
	StateMinerActiveSectors(context.Context, address.Address, types.TipSetKey) ([]*miner.SectorOnChainInfo, error) //perm:read
	// StateMinerProvingDeadline calculates the deadline at some epoch for a proving period
	// and returns the deadline-related calculations.
	StateMinerProvingDeadline(context.Context, address.Address, types.TipSetKey) (*dline.Info, error) //perm:read
	// StateMinerPower returns the power of the indicated miner
	StateMinerPower(context.Context, address.Address, types.TipSetKey) (*MinerPower, error) //perm:read
	// StateMinerInfo returns info about the indicated miner
	StateMinerInfo(context.Context, address.Address, types.TipSetKey) (MinerInfo, error) //perm:read
	// StateMinerDeadlines returns all the proving deadlines for the given miner
	StateMinerDeadlines(context.Context, address.Address, types.TipSetKey) ([]Deadline, error) //perm:read
	// StateMinerPartitions returns all partitions in the specified deadline
	StateMinerPartitions(ctx context.Context, m address.Address, dlIdx uint64, tsk types.TipSetKey) ([]Partition, error) //perm:read
	// StateMinerFaults returns a bitfield indicating the faulty sectors of the given miner
	StateMinerFaults(context.Context, address.Address, types.TipSetKey) (bitfield.BitField, error) //perm:read
	// StateAllMinerFaults returns all non-expired Faults that occur within lookback epochs of the given tipset
	StateAllMinerFaults(ctx context.Context, lookback abi.ChainEpoch, ts types.TipSetKey) ([]*Fault, error) //perm:read
	// StateMinerRecoveries returns a bitfield indicating the recovering sectors of the given miner
	StateMinerRecoveries(context.Context, address.Address, types.TipSetKey) (bitfield.BitField, error) //perm:read
	// StateMinerInitialPledgeCollateral returns the precommit deposit for the specified miner's sector
	StateMinerPreCommitDepositForPower(context.Context, address.Address, miner.SectorPreCommitInfo, types.TipSetKey) (types.BigInt, error) //perm:read
	// StateMinerInitialPledgeCollateral returns the initial pledge collateral for the specified miner's sector
	StateMinerInitialPledgeCollateral(context.Context, address.Address, miner.SectorPreCommitInfo, types.TipSetKey) (types.BigInt, error) //perm:read
	// StateMinerAvailableBalance returns the portion of a miner's balance that can be withdrawn or spent
	StateMinerAvailableBalance(context.Context, address.Address, types.TipSetKey) (types.BigInt, error) //perm:read
	// StateMinerSectorAllocated checks if a sector number is marked as allocated.
	StateMinerSectorAllocated(context.Context, address.Address, abi.SectorNumber, types.TipSetKey) (bool, error) //perm:read
	// StateSectorPreCommitInfo returns the PreCommit info for the specified miner's sector.
	// Returns nil and no error if the sector isn't precommitted.
	//
	// Note that the sector number may be allocated while PreCommitInfo is nil. This means that either allocated sector
	// numbers were compacted, and the sector number was marked as allocated in order to reduce size of the allocated
	// sectors bitfield, or that the sector was precommitted, but the precommit has expired.
	StateSectorPreCommitInfo(context.Context, address.Address, abi.SectorNumber, types.TipSetKey) (*miner.SectorPreCommitOnChainInfo, error) //perm:read
	// StateSectorGetInfo returns the on-chain info for the specified miner's sector. Returns null in case the sector info isn't found
	// NOTE: returned info.Expiration may not be accurate in some cases, use StateSectorExpiration to get accurate
	// expiration epoch
	StateSectorGetInfo(context.Context, address.Address, abi.SectorNumber, types.TipSetKey) (*miner.SectorOnChainInfo, error) //perm:read
	// StateSectorExpiration returns epoch at which given sector will expire
	StateSectorExpiration(context.Context, address.Address, abi.SectorNumber, types.TipSetKey) (*lminer.SectorExpiration, error) //perm:read
	// StateSectorPartition finds deadline/partition with the specified sector
	StateSectorPartition(ctx context.Context, maddr address.Address, sectorNumber abi.SectorNumber, tok types.TipSetKey) (*lminer.SectorLocation, error) //perm:read
	// StateSearchMsg looks back up to limit epochs in the chain for a message, and returns its receipt and the tipset where it was executed
	//
	// NOTE: If a replacing message is found on chain, this method will return
	// a MsgLookup for the replacing message - the MsgLookup.Message will be a different
	// CID than the one provided in the 'cid' param, MsgLookup.Receipt will contain the
	// result of the execution of the replacing message.
	//
	// If the caller wants to ensure that exactly the requested message was executed,
	// they must check that MsgLookup.Message is equal to the provided 'cid', or set the
	// `allowReplaced` parameter to false. Without this check, and with `allowReplaced`
	// set to true, both the requested and original message may appear as
	// successfully executed on-chain, which may look like a double-spend.
	//
	// A replacing message is a message with a different CID, any of Gas values, and
	// different signature, but with all other parameters matching (source/destination,
	// nonce, params, etc.)
	StateSearchMsg(ctx context.Context, from types.TipSetKey, msg cid.Cid, limit abi.ChainEpoch, allowReplaced bool) (*MsgLookup, error) //perm:read
	// StateWaitMsg looks back up to limit epochs in the chain for a message.
	// If not found, it blocks until the message arrives on chain, and gets to the
	// indicated confidence depth.
	//
	// NOTE: If a replacing message is found on chain, this method will return
	// a MsgLookup for the replacing message - the MsgLookup.Message will be a different
	// CID than the one provided in the 'cid' param, MsgLookup.Receipt will contain the
	// result of the execution of the replacing message.
	//
	// If the caller wants to ensure that exactly the requested message was executed,
	// they must check that MsgLookup.Message is equal to the provided 'cid', or set the
	// `allowReplaced` parameter to false. Without this check, and with `allowReplaced`
	// set to true, both the requested and original message may appear as
	// successfully executed on-chain, which may look like a double-spend.
	//
	// A replacing message is a message with a different CID, any of Gas values, and
	// different signature, but with all other parameters matching (source/destination,
	// nonce, params, etc.)
	StateWaitMsg(ctx context.Context, cid cid.Cid, confidence uint64, limit abi.ChainEpoch, allowReplaced bool) (*MsgLookup, error) //perm:read
	// StateListMiners returns the addresses of every miner that has claimed power in the Power Actor
	StateListMiners(context.Context, types.TipSetKey) ([]address.Address, error) //perm:read
	// StateListActors returns the addresses of every actor in the state
	StateListActors(context.Context, types.TipSetKey) ([]address.Address, error) //perm:read
	// StateMarketBalance looks up the Escrow and Locked balances of the given address in the Storage Market
	StateMarketBalance(context.Context, address.Address, types.TipSetKey) (MarketBalance, error) //perm:read
	// StateMarketParticipants returns the Escrow and Locked balances of every participant in the Storage Market
	StateMarketParticipants(context.Context, types.TipSetKey) (map[string]MarketBalance, error) //perm:read
	// StateMarketDeals returns information about every deal in the Storage Market
	StateMarketDeals(context.Context, types.TipSetKey) (map[string]*MarketDeal, error) //perm:read
	// StateMarketStorageDeal returns information about the indicated deal
	StateMarketStorageDeal(context.Context, abi.DealID, types.TipSetKey) (*MarketDeal, error) //perm:read
	// StateGetAllocationForPendingDeal returns the allocation for a given deal ID of a pending deal. Returns nil if
	// pending allocation is not found.
	StateGetAllocationForPendingDeal(ctx context.Context, dealId abi.DealID, tsk types.TipSetKey) (*verifregtypes.Allocation, error) //perm:read
	// StateGetAllocation returns the allocation for a given address and allocation ID.
	StateGetAllocation(ctx context.Context, clientAddr address.Address, allocationId verifregtypes.AllocationId, tsk types.TipSetKey) (*verifregtypes.Allocation, error) //perm:read
	// StateGetAllocations returns the all the allocations for a given client.
	StateGetAllocations(ctx context.Context, clientAddr address.Address, tsk types.TipSetKey) (map[verifregtypes.AllocationId]verifregtypes.Allocation, error) //perm:read
	// StateGetClaim returns the claim for a given address and claim ID.
	StateGetClaim(ctx context.Context, providerAddr address.Address, claimId verifregtypes.ClaimId, tsk types.TipSetKey) (*verifregtypes.Claim, error) //perm:read
	// StateGetClaims returns the all the claims for a given provider.
	StateGetClaims(ctx context.Context, providerAddr address.Address, tsk types.TipSetKey) (map[verifregtypes.ClaimId]verifregtypes.Claim, error) //perm:read
	// StateComputeDataCID computes DataCID from a set of on-chain deals
	StateComputeDataCID(ctx context.Context, maddr address.Address, sectorType abi.RegisteredSealProof, deals []abi.DealID, tsk types.TipSetKey) (cid.Cid, error) //perm:read
	// StateLookupID retrieves the ID address of the given address
	StateLookupID(context.Context, address.Address, types.TipSetKey) (address.Address, error) //perm:read
	// StateAccountKey returns the public key address of the given ID address for secp and bls accounts
	StateAccountKey(context.Context, address.Address, types.TipSetKey) (address.Address, error) //perm:read
	// StateLookupRobustAddress returns the public key address of the given ID address for non-account addresses (multisig, miners etc)
	StateLookupRobustAddress(context.Context, address.Address, types.TipSetKey) (address.Address, error) //perm:read
	// StateChangedActors returns all the actors whose states change between the two given state CIDs
	// TODO: Should this take tipset keys instead?
	StateChangedActors(context.Context, cid.Cid, cid.Cid) (map[string]types.Actor, error) //perm:read
	// StateMinerSectorCount returns the number of sectors in a miner's sector set and proving set
	StateMinerSectorCount(context.Context, address.Address, types.TipSetKey) (MinerSectors, error) //perm:read
	// StateMinerAllocated returns a bitfield containing all sector numbers marked as allocated in miner state
	StateMinerAllocated(context.Context, address.Address, types.TipSetKey) (*bitfield.BitField, error) //perm:read
	// StateCompute is a flexible command that applies the given messages on the given tipset.
	// The messages are run as though the VM were at the provided height.
	//
	// When called, StateCompute will:
	// - Load the provided tipset, or use the current chain head if not provided
	// - Compute the tipset state of the provided tipset on top of the parent state
	//   - (note that this step runs before vmheight is applied to the execution)
	//   - Execute state upgrade if any were scheduled at the epoch, or in null
	//     blocks preceding the tipset
	//   - Call the cron actor on null blocks preceding the tipset
	//   - For each block in the tipset
	//     - Apply messages in blocks in the specified
	//     - Award block reward by calling the reward actor
	//   - Call the cron actor for the current epoch
	// - If the specified vmheight is higher than the current epoch, apply any
	//   needed state upgrades to the state
	// - Apply the specified messages to the state
	//
	// The vmheight parameter sets VM execution epoch, and can be used to simulate
	// message execution in different network versions. If the specified vmheight
	// epoch is higher than the epoch of the specified tipset, any state upgrades
	// until the vmheight will be executed on the state before applying messages
	// specified by the user.
	//
	// Note that the initial tipset state computation is not affected by the
	// vmheight parameter - only the messages in the `apply` set are
	//
	// If the caller wants to simply compute the state, vmheight should be set to
	// the epoch of the specified tipset.
	//
	// Messages in the `apply` parameter must have the correct nonces, and gas
	// values set.
	StateCompute(context.Context, abi.ChainEpoch, []*types.Message, types.TipSetKey) (*ComputeStateOutput, error) //perm:read
	// StateVerifierStatus returns the data cap for the given address.
	// Returns nil if there is no entry in the data cap table for the
	// address.
	StateVerifierStatus(ctx context.Context, addr address.Address, tsk types.TipSetKey) (*abi.StoragePower, error) //perm:read
	// StateVerifiedClientStatus returns the data cap for the given address.
	// Returns nil if there is no entry in the data cap table for the
	// address.
	StateVerifiedClientStatus(ctx context.Context, addr address.Address, tsk types.TipSetKey) (*abi.StoragePower, error) //perm:read
	// StateVerifiedRegistryRootKey returns the address of the Verified Registry's root key
	StateVerifiedRegistryRootKey(ctx context.Context, tsk types.TipSetKey) (address.Address, error) //perm:read
	// StateDealProviderCollateralBounds returns the min and max collateral a storage provider
	// can issue. It takes the deal size and verified status as parameters.
	StateDealProviderCollateralBounds(context.Context, abi.PaddedPieceSize, bool, types.TipSetKey) (DealCollateralBounds, error) //perm:read

	// StateCirculatingSupply returns the exact circulating supply of Filecoin at the given tipset.
	// This is not used anywhere in the protocol itself, and is only for external consumption.
	StateCirculatingSupply(context.Context, types.TipSetKey) (abi.TokenAmount, error) //perm:read
	// StateVMCirculatingSupplyInternal returns an approximation of the circulating supply of Filecoin at the given tipset.
	// This is the value reported by the runtime interface to actors code.
	StateVMCirculatingSupplyInternal(context.Context, types.TipSetKey) (CirculatingSupply, error) //perm:read
	// StateNetworkVersion returns the network version at the given tipset
	StateNetworkVersion(context.Context, types.TipSetKey) (apitypes.NetworkVersion, error) //perm:read
	// StateActorCodeCIDs returns the CIDs of all the builtin actors for the given network version
	StateActorCodeCIDs(context.Context, abinetwork.Version) (map[string]cid.Cid, error) //perm:read
	// StateActorManifestCID returns the CID of the builtin actors manifest for the given network version
	StateActorManifestCID(context.Context, abinetwork.Version) (cid.Cid, error) //perm:read

	// StateGetRandomnessFromTickets is used to sample the chain for randomness.
	StateGetRandomnessFromTickets(ctx context.Context, personalization crypto.DomainSeparationTag, randEpoch abi.ChainEpoch, entropy []byte, tsk types.TipSetKey) (abi.Randomness, error) //perm:read
	// StateGetRandomnessFromBeacon is used to sample the beacon for randomness.
	StateGetRandomnessFromBeacon(ctx context.Context, personalization crypto.DomainSeparationTag, randEpoch abi.ChainEpoch, entropy []byte, tsk types.TipSetKey) (abi.Randomness, error) //perm:read

	// StateGetBeaconEntry returns the beacon entry for the given filecoin epoch. If
	// the entry has not yet been produced, the call will block until the entry
	// becomes available
	StateGetBeaconEntry(ctx context.Context, epoch abi.ChainEpoch) (*types.BeaconEntry, error) //perm:read

	// StateGetNetworkParams return current network params
	StateGetNetworkParams(ctx context.Context) (*NetworkParams, error) //perm:read

	// MethodGroup: Msig
	// The Msig methods are used to interact with multisig wallets on the
	// filecoin network

	// MsigGetAvailableBalance returns the portion of a multisig's balance that can be withdrawn or spent
	MsigGetAvailableBalance(context.Context, address.Address, types.TipSetKey) (types.BigInt, error) //perm:read
	// MsigGetVestingSchedule returns the vesting details of a given multisig.
	MsigGetVestingSchedule(context.Context, address.Address, types.TipSetKey) (MsigVesting, error) //perm:read
	// MsigGetVested returns the amount of FIL that vested in a multisig in a certain period.
	// It takes the following params: <multisig address>, <start epoch>, <end epoch>
	MsigGetVested(context.Context, address.Address, types.TipSetKey, types.TipSetKey) (types.BigInt, error) //perm:read

	// MsigGetPending returns pending transactions for the given multisig
	// wallet. Once pending transactions are fully approved, they will no longer
	// appear here.
	MsigGetPending(context.Context, address.Address, types.TipSetKey) ([]*MsigTransaction, error) //perm:read

	// MsigCreate creates a multisig wallet
	// It takes the following params: <required number of senders>, <approving addresses>, <unlock duration>
	// <initial balance>, <sender address of the create msg>, <gas price>
	MsigCreate(context.Context, uint64, []address.Address, abi.ChainEpoch, types.BigInt, address.Address, types.BigInt) (*MessagePrototype, error) //perm:sign

	// MsigPropose proposes a multisig message
	// It takes the following params: <multisig address>, <recipient address>, <value to transfer>,
	// <sender address of the propose msg>, <method to call in the proposed message>, <params to include in the proposed message>
	MsigPropose(context.Context, address.Address, address.Address, types.BigInt, address.Address, uint64, []byte) (*MessagePrototype, error) //perm:sign

	// MsigApprove approves a previously-proposed multisig message by transaction ID
	// It takes the following params: <multisig address>, <proposed transaction ID> <signer address>
	MsigApprove(context.Context, address.Address, uint64, address.Address) (*MessagePrototype, error) //perm:sign

	// MsigApproveTxnHash approves a previously-proposed multisig message, specified
	// using both transaction ID and a hash of the parameters used in the
	// proposal. This method of approval can be used to ensure you only approve
	// exactly the transaction you think you are.
	// It takes the following params: <multisig address>, <proposed message ID>, <proposer address>, <recipient address>, <value to transfer>,
	// <sender address of the approve msg>, <method to call in the proposed message>, <params to include in the proposed message>
	MsigApproveTxnHash(context.Context, address.Address, uint64, address.Address, address.Address, types.BigInt, address.Address, uint64, []byte) (*MessagePrototype, error) //perm:sign

	// MsigCancel cancels a previously-proposed multisig message
	// It takes the following params: <multisig address>, <proposed transaction ID> <signer address>
	MsigCancel(context.Context, address.Address, uint64, address.Address) (*MessagePrototype, error) //perm:sign

	// MsigCancel cancels a previously-proposed multisig message
	// It takes the following params: <multisig address>, <proposed transaction ID>, <recipient address>, <value to transfer>,
	// <sender address of the cancel msg>, <method to call in the proposed message>, <params to include in the proposed message>
	MsigCancelTxnHash(context.Context, address.Address, uint64, address.Address, types.BigInt, address.Address, uint64, []byte) (*MessagePrototype, error) //perm:sign

	// MsigAddPropose proposes adding a signer in the multisig
	// It takes the following params: <multisig address>, <sender address of the propose msg>,
	// <new signer>, <whether the number of required signers should be increased>
	MsigAddPropose(context.Context, address.Address, address.Address, address.Address, bool) (*MessagePrototype, error) //perm:sign

	// MsigAddApprove approves a previously proposed AddSigner message
	// It takes the following params: <multisig address>, <sender address of the approve msg>, <proposed message ID>,
	// <proposer address>, <new signer>, <whether the number of required signers should be increased>
	MsigAddApprove(context.Context, address.Address, address.Address, uint64, address.Address, address.Address, bool) (*MessagePrototype, error) //perm:sign

	// MsigAddCancel cancels a previously proposed AddSigner message
	// It takes the following params: <multisig address>, <sender address of the cancel msg>, <proposed message ID>,
	// <new signer>, <whether the number of required signers should be increased>
	MsigAddCancel(context.Context, address.Address, address.Address, uint64, address.Address, bool) (*MessagePrototype, error) //perm:sign

	// MsigSwapPropose proposes swapping 2 signers in the multisig
	// It takes the following params: <multisig address>, <sender address of the propose msg>,
	// <old signer>, <new signer>
	MsigSwapPropose(context.Context, address.Address, address.Address, address.Address, address.Address) (*MessagePrototype, error) //perm:sign

	// MsigSwapApprove approves a previously proposed SwapSigner
	// It takes the following params: <multisig address>, <sender address of the approve msg>, <proposed message ID>,
	// <proposer address>, <old signer>, <new signer>
	MsigSwapApprove(context.Context, address.Address, address.Address, uint64, address.Address, address.Address, address.Address) (*MessagePrototype, error) //perm:sign

	// MsigSwapCancel cancels a previously proposed SwapSigner message
	// It takes the following params: <multisig address>, <sender address of the cancel msg>, <proposed message ID>,
	// <old signer>, <new signer>
	MsigSwapCancel(context.Context, address.Address, address.Address, uint64, address.Address, address.Address) (*MessagePrototype, error) //perm:sign

	// MsigRemoveSigner proposes the removal of a signer from the multisig.
	// It accepts the multisig to make the change on, the proposer address to
	// send the message from, the address to be removed, and a boolean
	// indicating whether or not the signing threshold should be lowered by one
	// along with the address removal.
	MsigRemoveSigner(ctx context.Context, msig address.Address, proposer address.Address, toRemove address.Address, decrease bool) (*MessagePrototype, error) //perm:sign

	// MarketAddBalance adds funds to the market actor
	MarketAddBalance(ctx context.Context, wallet, addr address.Address, amt types.BigInt) (cid.Cid, error) //perm:sign
	// MarketGetReserved gets the amount of funds that are currently reserved for the address
	MarketGetReserved(ctx context.Context, addr address.Address) (types.BigInt, error) //perm:sign
	// MarketReserveFunds reserves funds for a deal
	MarketReserveFunds(ctx context.Context, wallet address.Address, addr address.Address, amt types.BigInt) (cid.Cid, error) //perm:sign
	// MarketReleaseFunds releases funds reserved by MarketReserveFunds
	MarketReleaseFunds(ctx context.Context, addr address.Address, amt types.BigInt) error //perm:sign
	// MarketWithdraw withdraws unlocked funds from the market actor
	MarketWithdraw(ctx context.Context, wallet, addr address.Address, amt types.BigInt) (cid.Cid, error) //perm:sign

	// MethodGroup: Paych
	// The Paych methods are for interacting with and managing payment channels

	// PaychGet gets or creates a payment channel between address pair
	//  The specified amount will be reserved for use. If there aren't enough non-reserved funds
	//    available, funds will be added through an on-chain message.
	//  - When opts.OffChain is true, this call will not cause any messages to be sent to the chain (no automatic
	//    channel creation/funds adding). If the operation can't be performed without sending a message an error will be
	//    returned. Note that even when this option is specified, this call can be blocked by previous operations on the
	//    channel waiting for on-chain operations.
	PaychGet(ctx context.Context, from, to address.Address, amt types.BigInt, opts PaychGetOpts) (*ChannelInfo, error) //perm:sign
	// PaychFund gets or creates a payment channel between address pair.
	// The specified amount will be added to the channel through on-chain send for future use
	PaychFund(ctx context.Context, from, to address.Address, amt types.BigInt) (*ChannelInfo, error)                    //perm:sign
	PaychGetWaitReady(context.Context, cid.Cid) (address.Address, error)                                                //perm:sign
	PaychAvailableFunds(ctx context.Context, ch address.Address) (*ChannelAvailableFunds, error)                        //perm:sign
	PaychAvailableFundsByFromTo(ctx context.Context, from, to address.Address) (*ChannelAvailableFunds, error)          //perm:sign
	PaychList(context.Context) ([]address.Address, error)                                                               //perm:read
	PaychStatus(context.Context, address.Address) (*PaychStatus, error)                                                 //perm:read
	PaychSettle(context.Context, address.Address) (cid.Cid, error)                                                      //perm:sign
	PaychCollect(context.Context, address.Address) (cid.Cid, error)                                                     //perm:sign
	PaychAllocateLane(ctx context.Context, ch address.Address) (uint64, error)                                          //perm:sign
	PaychNewPayment(ctx context.Context, from, to address.Address, vouchers []VoucherSpec) (*PaymentInfo, error)        //perm:sign
	PaychVoucherCheckValid(context.Context, address.Address, *paych.SignedVoucher) error                                //perm:read
	PaychVoucherCheckSpendable(context.Context, address.Address, *paych.SignedVoucher, []byte, []byte) (bool, error)    //perm:read
	PaychVoucherCreate(context.Context, address.Address, types.BigInt, uint64) (*VoucherCreateResult, error)            //perm:sign
	PaychVoucherAdd(context.Context, address.Address, *paych.SignedVoucher, []byte, types.BigInt) (types.BigInt, error) //perm:write
	PaychVoucherList(context.Context, address.Address) ([]*paych.SignedVoucher, error)                                  //perm:write
	PaychVoucherSubmit(context.Context, address.Address, *paych.SignedVoucher, []byte, []byte) (cid.Cid, error)         //perm:sign

	// MethodGroup: Node
	// These methods are general node management and status commands

	NodeStatus(ctx context.Context, inclChainStatus bool) (NodeStatus, error) //perm:read

	// MethodGroup: Eth
	// These methods are used for Ethereum-compatible JSON-RPC calls
	//
	// EthAccounts will always return [] since we don't expect Lotus to manage private keys
	EthAccounts(ctx context.Context) ([]ethtypes.EthAddress, error) //perm:read
	// EthAddressToFilecoinAddress converts an EthAddress into an f410 Filecoin Address
	EthAddressToFilecoinAddress(ctx context.Context, ethAddress ethtypes.EthAddress) (address.Address, error) //perm:read
	// FilecoinAddressToEthAddress converts an f410 or f0 Filecoin Address to an EthAddress
	FilecoinAddressToEthAddress(ctx context.Context, filecoinAddress address.Address) (ethtypes.EthAddress, error) //perm:read
	// EthBlockNumber returns the height of the latest (heaviest) TipSet
	EthBlockNumber(ctx context.Context) (ethtypes.EthUint64, error) //perm:read
	// EthGetBlockTransactionCountByNumber returns the number of messages in the TipSet
	EthGetBlockTransactionCountByNumber(ctx context.Context, blkNum ethtypes.EthUint64) (ethtypes.EthUint64, error) //perm:read
	// EthGetBlockTransactionCountByHash returns the number of messages in the TipSet
	EthGetBlockTransactionCountByHash(ctx context.Context, blkHash ethtypes.EthHash) (ethtypes.EthUint64, error) //perm:read

	EthGetBlockByHash(ctx context.Context, blkHash ethtypes.EthHash, fullTxInfo bool) (ethtypes.EthBlock, error)                               //perm:read
	EthGetBlockByNumber(ctx context.Context, blkNum string, fullTxInfo bool) (ethtypes.EthBlock, error)                                        //perm:read
	EthGetTransactionByHash(ctx context.Context, txHash *ethtypes.EthHash) (*ethtypes.EthTx, error)                                            //perm:read
	EthGetTransactionByHashLimited(ctx context.Context, txHash *ethtypes.EthHash, limit abi.ChainEpoch) (*ethtypes.EthTx, error)               //perm:read
	EthGetTransactionHashByCid(ctx context.Context, cid cid.Cid) (*ethtypes.EthHash, error)                                                    //perm:read
	EthGetMessageCidByTransactionHash(ctx context.Context, txHash *ethtypes.EthHash) (*cid.Cid, error)                                         //perm:read
	EthGetTransactionCount(ctx context.Context, sender ethtypes.EthAddress, blkOpt string) (ethtypes.EthUint64, error)                         //perm:read
	EthGetTransactionReceipt(ctx context.Context, txHash ethtypes.EthHash) (*EthTxReceipt, error)                                              //perm:read
	EthGetTransactionReceiptLimited(ctx context.Context, txHash ethtypes.EthHash, limit abi.ChainEpoch) (*EthTxReceipt, error)                 //perm:read
	EthGetTransactionByBlockHashAndIndex(ctx context.Context, blkHash ethtypes.EthHash, txIndex ethtypes.EthUint64) (ethtypes.EthTx, error)    //perm:read
	EthGetTransactionByBlockNumberAndIndex(ctx context.Context, blkNum ethtypes.EthUint64, txIndex ethtypes.EthUint64) (ethtypes.EthTx, error) //perm:read

	EthGetCode(ctx context.Context, address ethtypes.EthAddress, blkOpt string) (ethtypes.EthBytes, error)                                    //perm:read
	EthGetStorageAt(ctx context.Context, address ethtypes.EthAddress, position ethtypes.EthBytes, blkParam string) (ethtypes.EthBytes, error) //perm:read
	EthGetBalance(ctx context.Context, address ethtypes.EthAddress, blkParam string) (ethtypes.EthBigInt, error)                              //perm:read
	EthChainId(ctx context.Context) (ethtypes.EthUint64, error)                                                                               //perm:read
	NetVersion(ctx context.Context) (string, error)                                                                                           //perm:read
	NetListening(ctx context.Context) (bool, error)                                                                                           //perm:read
	EthProtocolVersion(ctx context.Context) (ethtypes.EthUint64, error)                                                                       //perm:read
	EthGasPrice(ctx context.Context) (ethtypes.EthBigInt, error)                                                                              //perm:read
	EthFeeHistory(ctx context.Context, p jsonrpc.RawParams) (ethtypes.EthFeeHistory, error)                                                   //perm:read

	EthMaxPriorityFeePerGas(ctx context.Context) (ethtypes.EthBigInt, error)                      //perm:read
	EthEstimateGas(ctx context.Context, tx ethtypes.EthCall) (ethtypes.EthUint64, error)          //perm:read
	EthCall(ctx context.Context, tx ethtypes.EthCall, blkParam string) (ethtypes.EthBytes, error) //perm:read

	EthSendRawTransaction(ctx context.Context, rawTx ethtypes.EthBytes) (ethtypes.EthHash, error) //perm:read

	// Returns event logs matching given filter spec.
	EthGetLogs(ctx context.Context, filter *ethtypes.EthFilterSpec) (*ethtypes.EthFilterResult, error) //perm:read

	// Polling method for a filter, returns event logs which occurred since last poll.
	// (requires write perm since timestamp of last filter execution will be written)
	EthGetFilterChanges(ctx context.Context, id ethtypes.EthFilterID) (*ethtypes.EthFilterResult, error) //perm:write

	// Returns event logs matching filter with given id.
	// (requires write perm since timestamp of last filter execution will be written)
	EthGetFilterLogs(ctx context.Context, id ethtypes.EthFilterID) (*ethtypes.EthFilterResult, error) //perm:write

	// Installs a persistent filter based on given filter spec.
	EthNewFilter(ctx context.Context, filter *ethtypes.EthFilterSpec) (ethtypes.EthFilterID, error) //perm:write

	// Installs a persistent filter to notify when a new block arrives.
	EthNewBlockFilter(ctx context.Context) (ethtypes.EthFilterID, error) //perm:write

	// Installs a persistent filter to notify when new messages arrive in the message pool.
	EthNewPendingTransactionFilter(ctx context.Context) (ethtypes.EthFilterID, error) //perm:write

	// Uninstalls a filter with given id.
	EthUninstallFilter(ctx context.Context, id ethtypes.EthFilterID) (bool, error) //perm:write

	// Subscribe to different event types using websockets
	// eventTypes is one or more of:
	//  - newHeads: notify when new blocks arrive.
	//  - pendingTransactions: notify when new messages arrive in the message pool.
	//  - logs: notify new event logs that match a criteria
	// params contains additional parameters used with the log event type
	// The client will receive a stream of EthSubscriptionResponse values until EthUnsubscribe is called.
	EthSubscribe(ctx context.Context, params jsonrpc.RawParams) (ethtypes.EthSubscriptionID, error) //perm:write

	// Unsubscribe from a websocket subscription
	EthUnsubscribe(ctx context.Context, id ethtypes.EthSubscriptionID) (bool, error) //perm:write

	// Returns the client version
	Web3ClientVersion(ctx context.Context) (string, error) //perm:read

	// CreateBackup creates node backup onder the specified file name. The
	// method requires that the lotus daemon is running with the
	// LOTUS_BACKUP_BASE_PATH environment variable set to some path, and that
	// the path specified when calling CreateBackup is within the base path
	CreateBackup(ctx context.Context, fpath string) error //perm:admin

	RaftState(ctx context.Context) (*RaftStateData, error) //perm:read
	RaftLeader(ctx context.Context) (peer.ID, error)       //perm:read

	// IPC-specific methods //

	// IPCAddSubnetActor deploys a new subnet actor.
	IPCAddSubnetActor(ctx context.Context, wallet address.Address, params subnetactor.ConstructParams) (address.Address, error)                          //perm:write
	IPCReadGatewayState(ctx context.Context, actor address.Address, tsk types.TipSetKey) (*gateway.State, error)                                         //perm:read
	IPCReadSubnetActorState(ctx context.Context, sn sdk.SubnetID, tsk types.TipSetKey) (*subnetactor.State, error)                                       //perm:read
	IPCGetPrevCheckpointForChild(ctx context.Context, gatewayAddr address.Address, subnet sdk.SubnetID) (cid.Cid, error)                                 //perm:read
	IPCGetCheckpointTemplate(ctx context.Context, gatewayAddr address.Address, epoch abi.ChainEpoch) (*gateway.BottomUpCheckpoint, error)                //perm:read
	IPCListChildSubnets(ctx context.Context, gatewayAddr address.Address) ([]gateway.Subnet, error)                                                      //perm:read
	IPCHasVotedBottomUpCheckpoint(ctx context.Context, sn sdk.SubnetID, e abi.ChainEpoch, v address.Address) (bool, error)                               //perm:read
	IPCHasVotedTopDownCheckpoint(ctx context.Context, gw address.Address, e abi.ChainEpoch, v address.Address) (bool, error)                             //perm:read
	IPCListCheckpoints(ctx context.Context, sn sdk.SubnetID, from, to abi.ChainEpoch) ([]*gateway.BottomUpCheckpoint, error)                             //perm:read
	IPCGetCheckpoint(ctx context.Context, sn sdk.SubnetID, epoch abi.ChainEpoch) (*gateway.BottomUpCheckpoint, error)                                    //perm:read
	IPCGetTopDownMsgs(ctx context.Context, gatewayAddr address.Address, sn sdk.SubnetID, tsk types.TipSetKey, nonce uint64) ([]*gateway.CrossMsg, error) //perm:read
	IPCGetGenesisEpochForSubnet(ctx context.Context, gatewayAddr address.Address, sn sdk.SubnetID) (abi.ChainEpoch, error)                               //perm:read

	// Serialized representation of IPC calls.
	// This calls are serialized version of some of the IPC calls. They return directly the CBOR IPCGetCheckpointSerialized
	// version of the output of the call. These are really convenient to use the same type of serialization used
	// in actor's state, removing the need of then intermediate serialization introduced by the Lotus API.
	IPCGetCheckpointSerialized(ctx context.Context, sn sdk.SubnetID, epoch abi.ChainEpoch) ([]byte, error)                                              //perm:read
	IPCListCheckpointsSerialized(ctx context.Context, sn sdk.SubnetID, from, to abi.ChainEpoch) ([][]byte, error)                                       //perm:read
	IPCGetCheckpointTemplateSerialized(ctx context.Context, gatewayAddr address.Address, epoch abi.ChainEpoch) ([]byte, error)                          //perm:read
	IPCGetTopDownMsgsSerialized(ctx context.Context, gatewayAddr address.Address, sn sdk.SubnetID, tsk types.TipSetKey, nonce uint64) ([][]byte, error) //perm:read

}

// reverse interface to the client, called after EthSubscribe
type EthSubscriber interface {
	// note: the parameter is ethtypes.EthSubscriptionResponse serialized as json object
	EthSubscription(ctx context.Context, r jsonrpc.RawParams) error // rpc_method:eth_subscription notify:true
}

type StorageAsk struct {
	Response *storagemarket.StorageAsk

	DealProtocols []string
}

type FileRef struct {
	Path  string
	IsCAR bool
}

type MinerSectors struct {
	// Live sectors that should be proven.
	Live uint64
	// Sectors actively contributing to power.
	Active uint64
	// Sectors with failed proofs.
	Faulty uint64
}

type ImportRes struct {
	Root     cid.Cid
	ImportID imports.ID
}

type Import struct {
	Key imports.ID
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

type MsgLookup struct {
	Message   cid.Cid // Can be different than requested, in case it was replaced, but only gas values changed
	Receipt   types.MessageReceipt
	ReturnDec interface{}
	TipSet    types.TipSetKey
	Height    abi.ChainEpoch
}

type MsgGasCost struct {
	Message            cid.Cid // Can be different than requested, in case it was replaced, but only gas values changed
	GasUsed            abi.TokenAmount
	BaseFeeBurn        abi.TokenAmount
	OverEstimationBurn abi.TokenAmount
	MinerPenalty       abi.TokenAmount
	MinerTip           abi.TokenAmount
	Refund             abi.TokenAmount
	TotalCost          abi.TokenAmount
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
	MinerPower  power.Claim
	TotalPower  power.Claim
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

type IpldObject struct {
	Cid cid.Cid
	Obj interface{}
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
	Sectors           []builtin.ExtendedSectorInfo
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
	WinningPoStProof []builtin.PoStProof
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
type HeadChange struct {
	Type string
	Val  *types.TipSet
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

type PruneOpts struct {
	MovingGC    bool
	RetainState int64
}

type HotGCOpts struct {
	Threshold float64
	Periodic  bool
	Moving    bool
}

type EthTxReceipt struct {
	TransactionHash   ethtypes.EthHash     `json:"transactionHash"`
	TransactionIndex  ethtypes.EthUint64   `json:"transactionIndex"`
	BlockHash         ethtypes.EthHash     `json:"blockHash"`
	BlockNumber       ethtypes.EthUint64   `json:"blockNumber"`
	From              ethtypes.EthAddress  `json:"from"`
	To                *ethtypes.EthAddress `json:"to"`
	StateRoot         ethtypes.EthHash     `json:"root"`
	Status            ethtypes.EthUint64   `json:"status"`
	ContractAddress   *ethtypes.EthAddress `json:"contractAddress"`
	CumulativeGasUsed ethtypes.EthUint64   `json:"cumulativeGasUsed"`
	GasUsed           ethtypes.EthUint64   `json:"gasUsed"`
	EffectiveGasPrice ethtypes.EthBigInt   `json:"effectiveGasPrice"`
	LogsBloom         ethtypes.EthBytes    `json:"logsBloom"`
	Logs              []ethtypes.EthLog    `json:"logs"`
	Type              ethtypes.EthUint64   `json:"type"`
}
