package vtmroll

import (
	"fmt"
	"math/rand/v2"
)

// VTMRollerResult holds the outcome of a single dice roll for Vampire: The Masquerade 5th Edition.
// It contains the individual die results and computed locations of successes, criticals, and special failure states.
type VTMRollerResult struct {
	rolls                           []int
	hungerDice                      int
	rawSuccessesLocations           []int
	failuresLocations               []int
	possibleBestialFailureLocations []int
	halfCriticalsLocations          []int
	halfMessyCriticalLocations      []int
}

func lastPair(n int) int {
	if n%2 == 0 {
		return n
	} else {
		return n - 1
	}
}

// NewVTMRollerResult constructs a VTMRollerResult by analyzing the given rolls against the roller's configuration.
// It computes all success locations, critical pairs, and special failure conditions based on the threshold and upper limit.
func NewVTMRollerResult(rolls []int, vtmr *VTMRoller, hungerDice int) *VTMRollerResult {
	result := &VTMRollerResult{
		rolls:      rolls,
		hungerDice: hungerDice,
	}

	for i, roll := range rolls {
		if vtmr.isSuccess(roll) {
			result.rawSuccessesLocations = append(result.rawSuccessesLocations, i)
		} else {
			result.failuresLocations = append(result.failuresLocations, i)

			if i < hungerDice {
				result.possibleBestialFailureLocations = append(result.possibleBestialFailureLocations, i)
			}
		}

		if roll == vtmr.RollUpperLimit {
			result.halfCriticalsLocations = append(result.halfCriticalsLocations, i)

			if i < hungerDice {
				result.halfMessyCriticalLocations = append(result.halfMessyCriticalLocations, i)
			}
		}
	}

	return result
}

// Successes returns the total number of successes, including raw successes and bonus successes from critical pairs.
// Each complete pair of half-criticals (rolls at RollUpperLimit) grants 2 bonus successes.
func (vtmres *VTMRollerResult) Successes() int {
	return len(vtmres.rawSuccessesLocations) + lastPair(len(vtmres.halfCriticalsLocations))
}

// IsCritical reports whether the roll achieved a critical success (two or more half-criticals).
func (vtmres *VTMRollerResult) IsCritical() bool {
	return len(vtmres.halfCriticalsLocations) >= 2
}

// IsTotalFailure reports whether the roll resulted in zero successes.
func (vtmres *VTMRollerResult) IsTotalFailure() bool {
	return len(vtmres.rawSuccessesLocations) == 0
}

// IsBestialFailure reports whether the roll is a total failure with at least one hunger die involved.
// This represents a catastrophic failure in VTM 5e that triggers the Beast.
func (vtmres *VTMRollerResult) IsBestialFailure() bool {
	return vtmres.IsTotalFailure() && len(vtmres.possibleBestialFailureLocations) > 0
}

// IsMessyCritical reports whether the roll achieved a critical success where at least one half-critical came from a hunger die.
// This represents an unpredictable critical in VTM 5e that deals collateral damage.
func (vtmres *VTMRollerResult) IsMessyCritical() bool {
	if !vtmres.IsCritical() {
		return false
	}

	return len(vtmres.halfMessyCriticalLocations) > 0
}

// GetRolls returns a copy of all individual die rolls in order.
func (vtmres *VTMRollerResult) GetRolls() []int {
	rolls := make([]int, len(vtmres.rolls))
	copy(rolls, vtmres.rolls)
	return rolls
}

// GetHungerDice returns the count of hunger dice in this roll.
func (vtmres *VTMRollerResult) GetHungerDice() int {
	return vtmres.hungerDice
}

// GetRawSuccessesLocations returns a copy of indices where rolls met or exceeded the success threshold.
func (vtmres *VTMRollerResult) GetRawSuccessesLocations() []int {
	locations := make([]int, len(vtmres.rawSuccessesLocations))
	copy(locations, vtmres.rawSuccessesLocations)
	return locations
}

// GetHalfCriticalsLocations returns a copy of indices where rolls matched the upper limit (half-criticals).
func (vtmres *VTMRollerResult) GetHalfCriticalsLocations() []int {
	locations := make([]int, len(vtmres.halfCriticalsLocations))
	copy(locations, vtmres.halfCriticalsLocations)
	return locations
}

// GetHalfMessyCriticalLocations returns a copy of indices where hunger dice rolled the upper limit.
func (vtmres *VTMRollerResult) GetHalfMessyCriticalLocations() []int {
	locations := make([]int, len(vtmres.halfMessyCriticalLocations))
	copy(locations, vtmres.halfMessyCriticalLocations)
	return locations
}

// GetPossibleBestialFailureLocations returns a copy of indices where hunger dice failed to meet the threshold.
func (vtmres *VTMRollerResult) GetPossibleBestialFailureLocations() []int {
	locations := make([]int, len(vtmres.possibleBestialFailureLocations))
	copy(locations, vtmres.possibleBestialFailureLocations)
	return locations
}

// VTMRoller executes dice rolls for Vampire: The Masquerade 5th Edition using a configured threshold and die range.
//
// VTMRoller is not thread-safe; do not modify fields or call methods concurrently from multiple goroutines.
type VTMRoller struct {
	// Rand is the random source for the roller. It is initialized with a unique seed by default,
	// but can be overridden with a custom *rand.Rand for reproducible results (useful in tests).
	Rand *rand.Rand

	// RollLowerLimit is the minimum value a die can roll (inclusive). Default for standard VTM5 is 1.
	RollLowerLimit int

	// RollUpperLimit is the maximum value a die can roll (inclusive). Default for standard VTM5 is 10.
	// Rolls matching this value are half-criticals and can form critical pairs.
	RollUpperLimit int

	// SuccessThreshold is the minimum roll value (inclusive) needed for a success. Default for standard VTM5 is 6.
	SuccessThreshold int
}

// NewVTMRoller creates a new roller for VTM dice with default configuration and a unique random source.
//
// Default values are the standard ones for VTM: RollLowerLimit = 1, RollUpperLimit = 10, SuccessThreshold = 6.
// The Rand field is initialized with a unique seed from the system RNG.
//
// For reproducible results in tests, override the Rand field after construction:
//
//	roller := NewVTMRoller()
//	roller.Rand = rand.New(rand.NewPCG(0, 0))
func NewVTMRoller() *VTMRoller {
	return &VTMRoller{
		Rand:             rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64())),
		RollLowerLimit:   1,
		RollUpperLimit:   10,
		SuccessThreshold: 6,
	}
}

func (vtmr *VTMRoller) isSuccess(roll int) bool {
	return roll >= vtmr.SuccessThreshold
}

func (vtmr *VTMRoller) isUpperLimit(roll int) bool {
	return roll == vtmr.RollUpperLimit
}

func (vtmr *VTMRoller) isLowerLimit(roll int) bool {
	return roll == vtmr.RollLowerLimit
}

// RollDie generates a single die roll within the configured range [RollLowerLimit, RollUpperLimit].
func (vtmr *VTMRoller) RollDie() int {
	return vtmr.Rand.IntN(vtmr.RollUpperLimit+1-vtmr.RollLowerLimit) + vtmr.RollLowerLimit
}

// Roll performs a dice roll with the given pool size and hunger dice count.
// Hunger dice count is clamped to [0, pool]. The returned result contains all roll details,
// computed successes, and special outcome states.
func (vtmr *VTMRoller) Roll(pool int, hungerDice int) *VTMRollerResult {
	// Validate hungerDice parameter
	if hungerDice < 0 {
		hungerDice = 0
	}
	if hungerDice > pool {
		hungerDice = pool
	}

	rolls := make([]int, pool)
	for i := range rolls {
		rolls[i] = vtmr.RollDie()
	}

	result := NewVTMRollerResult(rolls, vtmr, hungerDice)

	return result
}

// ReRoll re-rolls specified dice from a prior roll, creating a new result with updated outcomes.
// Hunger dice cannot be rerolled. Indices must be valid (0-based, within pool size, and not hunger dice).
// Returns an error if any index is invalid or refers to a hunger die.
func (vtmr *VTMRoller) ReRoll(oldResult *VTMRollerResult, rerollPlaces ...int) (*VTMRollerResult, error) {
	newRolls := make([]int, len(oldResult.rolls))
	copy(newRolls, oldResult.rolls)

	for _, rp := range rerollPlaces {
		if rp < oldResult.hungerDice {
			return nil, fmt.Errorf("cannot reroll hunger die (hunger: %d, offset %d)", oldResult.hungerDice, rp)
		}

		if rp >= len(newRolls) {
			return nil, fmt.Errorf("offset greater than or equal to rolls (len: %d, offset: %d)", len(newRolls), rp)
		}

		newRolls[rp] = vtmr.RollDie()
	}

	newResult := NewVTMRollerResult(newRolls, vtmr, oldResult.hungerDice)

	return newResult, nil
}
