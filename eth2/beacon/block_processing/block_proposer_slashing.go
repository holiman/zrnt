package block_processing

import (
	"errors"
	"github.com/protolambda/zrnt/eth2/beacon"
	"github.com/protolambda/zrnt/eth2/util/bls"
	"github.com/protolambda/zrnt/eth2/util/ssz"
)

func ProcessBlockProposerSlashings(state *beacon.BeaconState, block *beacon.BeaconBlock) error {
	if len(block.Body.ProposerSlashings) > beacon.MAX_PROPOSER_SLASHINGS {
		return errors.New("too many proposer slashings")
	}
	for _, ps := range block.Body.ProposerSlashings {
		if err := ProcessProposerSlashing(state, &ps); err != nil {
			return err
		}
	}
	return nil
}

func ProcessProposerSlashing(state *beacon.BeaconState, ps *beacon.ProposerSlashing) error {
	if !state.ValidatorRegistry.IsValidatorIndex(ps.ProposerIndex) {
		return errors.New("invalid proposer index")
	}
	proposer := &state.ValidatorRegistry[ps.ProposerIndex]
	// Verify that the epoch is the same
	if ps.Header1.Slot.ToEpoch() != ps.Header2.Slot.ToEpoch() {
		return errors.New("proposer slashing requires slashing headers to be in same epoch")
	}
	// Check proposer is slashable
	if !proposer.IsSlashable(state.Epoch()) {
		return errors.New("proposer slashing requires proposer to be slashable")
	}
	// But the headers are different
	if ps.Header1.BlockBodyRoot == ps.Header2.BlockBodyRoot {
		return errors.New("proposer slashing requires two different headers")
	}
	if !(
		bls.BlsVerify(proposer.Pubkey, ssz.SignedRoot(ps.Header1), ps.Header1.Signature, beacon.GetDomain(state.Fork, ps.Header1.Slot.ToEpoch(), beacon.DOMAIN_BEACON_BLOCK)) &&
		bls.BlsVerify(proposer.Pubkey, ssz.SignedRoot(ps.Header2), ps.Header2.Signature, beacon.GetDomain(state.Fork, ps.Header2.Slot.ToEpoch(), beacon.DOMAIN_BEACON_BLOCK))) {
		return errors.New("proposer slashing has header with invalid BLS signature")
	}
	return state.SlashValidator(ps.ProposerIndex)
}
