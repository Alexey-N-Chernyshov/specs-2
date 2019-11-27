package storage_mining

import (
	"math"

	block "github.com/filecoin-project/specs/systems/filecoin_blockchain/struct/block"
)

// update pointer to most recent challenge
func (cs *ChallengeStatus_I) OnNewChallenge(currEpoch block.ChainEpoch) {
	cs.LastChallengeEpoch_ = currEpoch
}

// Update pointer to most recent successful challenge response (both ePoSt and sPoSt)
// Call by _onSuccessfulPoSt
func (cs *ChallengeStatus_I) OnPoStSuccess(currEpoch block.ChainEpoch) {
	cs._lastPoStSuccessEpoch_ = currEpoch
}

// Update pointer to most recent challenge response failure
// Call by  _onMissedSurprisePoSt
func (cs *ChallengeStatus_I) OnPoStFailure(currEpoch block.ChainEpoch) {
	cs._lastPoStFailureEpoch_ = currEpoch
}

func (cs *ChallengeStatus_I) LastPoStResponseEpoch() block.ChainEpoch {
	return block.ChainEpoch(math.Max(float64(cs._lastPoStSuccessEpoch()), float64(cs._lastPoStFailureEpoch())))
}

func (cs *ChallengeStatus_I) IsChallenged() bool {
	// true if most recent challenge has gone unanswered
	return cs.LastChallengeEpoch() > cs.LastPoStResponseEpoch()
}

func (cs *ChallengeStatus_I) CanBeElected(currEpoch block.ChainEpoch) bool {
	// true if most recent successful post (surprise or election) was recent enough
	// and not currently getting challenged

	// pull in from consts
	PROVING_PERIOD := block.ChainEpoch(0)
	return !cs.IsChallenged() && currEpoch < cs._lastPoStSuccessEpoch()+PROVING_PERIOD
}

func (cs *ChallengeStatus_I) ShouldChallenge(currEpoch block.ChainEpoch, challengeFreePeriod block.ChainEpoch) bool {
	return currEpoch > (cs.LastChallengeEpoch()+challengeFreePeriod) && !cs.IsChallenged()
}