package rand

import (
	"context"
	"encoding/binary"

	"github.com/ipfs/go-cid"
	logging "github.com/ipfs/go-log/v2"
	"github.com/minio/blake2b-simd"
	"go.opencensus.io/trace"
	"golang.org/x/xerrors"

	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/crypto"
	"github.com/filecoin-project/go-state-types/network"

	"github.com/filecoin-project/lotus/chain/beacon"
	"github.com/filecoin-project/lotus/chain/store"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/filecoin-project/lotus/chain/vm"
	"github.com/filecoin-project/lotus/eudico-core/global"
)

var log = logging.Logger("rand")

func DrawRandomness(rbase []byte, pers crypto.DomainSeparationTag, round abi.ChainEpoch, entropy []byte) ([]byte, error) {
	h := blake2b.New256()
	if err := binary.Write(h, binary.BigEndian, int64(pers)); err != nil {
		return nil, xerrors.Errorf("deriving randomness: %w", err)
	}
	VRFDigest := blake2b.Sum256(rbase)
	_, err := h.Write(VRFDigest[:])
	if err != nil {
		return nil, xerrors.Errorf("hashing VRFDigest: %w", err)
	}
	if err := binary.Write(h, binary.BigEndian, round); err != nil {
		return nil, xerrors.Errorf("deriving randomness: %w", err)
	}
	_, err = h.Write(entropy)
	if err != nil {
		return nil, xerrors.Errorf("hashing entropy: %w", err)
	}

	return h.Sum(nil), nil
}

func (sr *stateRand) GetBeaconRandomnessTipset(ctx context.Context, round abi.ChainEpoch, lookback bool) (*types.TipSet, error) {
	_, span := trace.StartSpan(ctx, "store.GetBeaconRandomness")
	defer span.End()
	span.AddAttributes(trace.Int64Attribute("round", int64(round)))

	ts, err := sr.cs.LoadTipSet(ctx, types.NewTipSetKey(sr.blks...))
	if err != nil {
		return nil, err
	}

	if round > ts.Height() {
		return nil, xerrors.Errorf("cannot draw randomness from the future")
	}

	searchHeight := round
	if searchHeight < 0 {
		searchHeight = 0
	}

	randTs, err := sr.cs.GetTipsetByHeight(ctx, searchHeight, ts, lookback)
	if err != nil {
		return nil, err
	}

	return randTs, nil
}

func (sr *stateRand) getChainRandomness(ctx context.Context, pers crypto.DomainSeparationTag, round abi.ChainEpoch, entropy []byte, lookback bool) ([]byte, error) {
	_, span := trace.StartSpan(ctx, "store.GetChainRandomness")
	defer span.End()
	span.AddAttributes(trace.Int64Attribute("round", int64(round)))

	ts, err := sr.cs.LoadTipSet(ctx, types.NewTipSetKey(sr.blks...))
	if err != nil {
		return nil, err
	}

	if round > ts.Height() {
		return nil, xerrors.Errorf("cannot draw randomness from the future")
	}

	searchHeight := round
	if searchHeight < 0 {
		searchHeight = 0
	}

	randTs, err := sr.cs.GetTipsetByHeight(ctx, searchHeight, ts, lookback)
	if err != nil {
		return nil, err
	}

	mtb := randTs.MinTicketBlock()

	// if at (or just past -- for null epochs) appropriate epoch
	// or at genesis (works for negative epochs)
	return DrawRandomness(mtb.Ticket.VRFProof, pers, round, entropy)
}

type NetworkVersionGetter func(context.Context, abi.ChainEpoch) network.Version

type stateRand struct {
	cs                   *store.ChainStore
	blks                 []cid.Cid
	beacon               beacon.Schedule
	networkVersionGetter NetworkVersionGetter
}

func NewStateRand(cs *store.ChainStore, blks []cid.Cid, b beacon.Schedule, networkVersionGetter NetworkVersionGetter) vm.Rand {
	// we don't use winningPoSt and we use fake windowPoSt when running
	// mir consensus (at least for now)
	if global.IsConsensusAlgorithm(global.MirConsensus) {
		log.Debug("=================================================================================")
		log.Debug("DANGER ZONE! YOU ARE USING FAKE RANDOMNESS FOR YOUR VM AND PROOFS. USE WITH CARE!")
		log.Debug("=================================================================================")
		return &fakeRand{}
	}

	return &stateRand{
		cs:                   cs,
		blks:                 blks,
		beacon:               b,
		networkVersionGetter: networkVersionGetter,
	}
}

// network v0-12
func (sr *stateRand) getBeaconRandomnessV1(ctx context.Context, pers crypto.DomainSeparationTag, round abi.ChainEpoch, entropy []byte) ([]byte, error) {
	randTs, err := sr.GetBeaconRandomnessTipset(ctx, round, true)
	if err != nil {
		return nil, err
	}

	be, err := sr.cs.GetLatestBeaconEntry(ctx, randTs)
	if err != nil {
		return nil, err
	}

	// if at (or just past -- for null epochs) appropriate epoch
	// or at genesis (works for negative epochs)
	return DrawRandomness(be.Data, pers, round, entropy)
}

// network v13
func (sr *stateRand) getBeaconRandomnessV2(ctx context.Context, pers crypto.DomainSeparationTag, round abi.ChainEpoch, entropy []byte) ([]byte, error) {
	randTs, err := sr.GetBeaconRandomnessTipset(ctx, round, false)
	if err != nil {
		return nil, err
	}

	be, err := sr.cs.GetLatestBeaconEntry(ctx, randTs)
	if err != nil {
		return nil, err
	}

	// if at (or just past -- for null epochs) appropriate epoch
	// or at genesis (works for negative epochs)
	return DrawRandomness(be.Data, pers, round, entropy)
}

// network v14 and on
func (sr *stateRand) getBeaconRandomnessV3(ctx context.Context, pers crypto.DomainSeparationTag, filecoinEpoch abi.ChainEpoch, entropy []byte) ([]byte, error) {
	if filecoinEpoch < 0 {
		return sr.getBeaconRandomnessV2(ctx, pers, filecoinEpoch, entropy)
	}

	be, err := sr.extractBeaconEntryForEpoch(ctx, filecoinEpoch)
	if err != nil {
		log.Errorf("failed to get beacon entry as expected: %s", err)
		return nil, err
	}

	return DrawRandomness(be.Data, pers, filecoinEpoch, entropy)
}

func (sr *stateRand) GetChainRandomness(ctx context.Context, pers crypto.DomainSeparationTag, filecoinEpoch abi.ChainEpoch, entropy []byte) ([]byte, error) {
	nv := sr.networkVersionGetter(ctx, filecoinEpoch)

	if nv >= network.Version13 {
		return sr.getChainRandomness(ctx, pers, filecoinEpoch, entropy, false)
	}

	return sr.getChainRandomness(ctx, pers, filecoinEpoch, entropy, true)
}

func (sr *stateRand) GetBeaconRandomness(ctx context.Context, pers crypto.DomainSeparationTag, filecoinEpoch abi.ChainEpoch, entropy []byte) ([]byte, error) {
	nv := sr.networkVersionGetter(ctx, filecoinEpoch)

	if nv >= network.Version14 {
		return sr.getBeaconRandomnessV3(ctx, pers, filecoinEpoch, entropy)
	} else if nv == network.Version13 {
		return sr.getBeaconRandomnessV2(ctx, pers, filecoinEpoch, entropy)
	} else {
		return sr.getBeaconRandomnessV1(ctx, pers, filecoinEpoch, entropy)
	}
}

func (sr *stateRand) extractBeaconEntryForEpoch(ctx context.Context, filecoinEpoch abi.ChainEpoch) (*types.BeaconEntry, error) {
	randTs, err := sr.GetBeaconRandomnessTipset(ctx, filecoinEpoch, false)
	if err != nil {
		return nil, err
	}

	nv := sr.networkVersionGetter(ctx, filecoinEpoch)

	round := sr.beacon.BeaconForEpoch(filecoinEpoch).MaxBeaconRoundForEpoch(nv, filecoinEpoch)

	for i := 0; i < 20; i++ {
		cbe := randTs.Blocks()[0].BeaconEntries
		for _, v := range cbe {
			if v.Round == round {
				return &v, nil
			}
		}

		next, err := sr.cs.LoadTipSet(ctx, randTs.Parents())
		if err != nil {
			return nil, xerrors.Errorf("failed to load parents when searching back for beacon entry: %w", err)
		}

		randTs = next
	}

	return nil, xerrors.Errorf("didn't find beacon for round %d (epoch %d)", round, filecoinEpoch)
}
