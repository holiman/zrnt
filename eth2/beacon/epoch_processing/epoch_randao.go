package epoch_processing

import (
	"github.com/protolambda/zrnt/eth2/beacon"
)

func ProcessEpochRandao(state *beacon.BeaconState) {
	state.LatestRandaoMixes[(state.Epoch()+1)%beacon.LATEST_RANDAO_MIXES_LENGTH] = state.GetRandaoMix(state.Epoch())
}
