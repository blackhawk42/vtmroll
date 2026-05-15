package vtmroll

import (
	"fmt"
	"iter"
	"math/rand/v2"
)

// RollType classifies each die in a VTM roll pool.
//
// Each die has exactly one RollType, determined by its position (hunger vs. normal) and its rolled value.
// The types are mutually exclusive: a die is one of NormalSuccess, NormalFailure,
// HungerSuccess, HungerFailure, HalfCritical, HalfMessyCritical, or PossibleBestialFailure.
type RollType int

const (
	// NormalSuccess is a non-hunger die that met the success threshold but is not a half-critical.
	NormalSuccess RollType = iota

	// NormalFailure is a non-hunger die that failed to meet the success threshold.
	NormalFailure

	// HungerSuccess is a hunger die that met the success threshold but is not a half-critical.
	HungerSuccess

	// HungerFailure is a hunger die that failed to meet the success threshold and did not roll the lower limit.
	HungerFailure

	// HalfCritical is a non-hunger die that rolled the upper limit (e.g., 10 on a d10).
	//
	// Two or more half-criticals form a critical pair, granting bonus successes.
	HalfCritical

	// HalfMessyCritical is a hunger die that rolled the upper limit.
	//
	// If a critical includes at least one half-messy-critical, the result is a messy critical.
	HalfMessyCritical

	// PossibleBestialFailure is a hunger die that rolled the lower limit (e.g., 1 on a d10).
	//
	// In a total failure, any possible bestial failure makes the result a bestial failure.
	PossibleBestialFailure
)

// String returns a human-readable name for the RollType.
func (rt RollType) String() string {
	switch rt {
	case NormalSuccess:
		return "NormalSuccess"
	case NormalFailure:
		return "NormalFailure"
	case HungerSuccess:
		return "HungerSuccess"
	case HungerFailure:
		return "HungerFailure"
	case HalfCritical:
		return "HalfCritical"
	case HalfMessyCritical:
		return "HalfMessyCritical"
	case PossibleBestialFailure:
		return "PossibleBestialFailure"
	default:
		return fmt.Sprintf("RollType(%d)", int(rt))
	}
}

// VTMRollerResult holds the outcome of a single dice roll for Vampire: The Masquerade 5th Edition.
//
// It contains the individual die results and computed locations of successes, criticals, and special failure states.
type VTMRollerResult struct {
	rolls            []int
	hungerDice       int
	rollTypes        []RollType
	successes        int
	failures         int
	isCritical       bool
	isTotalFailure   bool
	isBestialFailure bool
	isMessyCritical  bool
}

func lastPair(n int) int {
	if n%2 == 0 {
		return n
	} else {
		return n - 1
	}
}

// NewVTMRollerResult constructs a VTMRollerResult by analyzing the given rolls against the roller's configuration.
//
// It computes all success locations, critical pairs, and special failure conditions based on the threshold and upper limit.
func NewVTMRollerResult(rolls []int, vtmr *VTMRoller, hungerDice int) VTMRollerResult {
	result := VTMRollerResult{
		rolls:      rolls,
		hungerDice: hungerDice,
		rollTypes:  make([]RollType, len(rolls)),
	}

	var normalSuccesses int
	var normalFailures int
	var hungerSuccesses int
	var hungerFailures int
	var halfCriticals int
	var halfMessyCriticals int
	var possibleBestialFailures int

	for i, roll := range rolls {
		if i < hungerDice {
			switch {
			case vtmr.IsUpperLimit(roll):
				result.rollTypes[i] = HalfMessyCritical
				halfMessyCriticals++
			case vtmr.IsLowerLimit(roll):
				result.rollTypes[i] = PossibleBestialFailure
				possibleBestialFailures++
			case vtmr.IsSuccess(roll):
				result.rollTypes[i] = HungerSuccess
				hungerSuccesses++
			default:
				result.rollTypes[i] = HungerFailure
				hungerFailures++
			}
		} else {
			switch {
			case vtmr.IsUpperLimit(roll):
				result.rollTypes[i] = HalfCritical
				halfCriticals++
			case vtmr.IsSuccess(roll):
				result.rollTypes[i] = NormalSuccess
				normalSuccesses++
			default:
				result.rollTypes[i] = NormalFailure
				normalFailures++
			}
		}
	}

	result.successes = normalSuccesses +
		halfCriticals +
		hungerSuccesses +
		halfMessyCriticals +
		lastPair(halfCriticals+halfMessyCriticals)

	result.failures = normalFailures + hungerFailures + possibleBestialFailures

	result.isCritical = (halfCriticals + halfMessyCriticals) >= 2

	result.isTotalFailure = result.successes == 0

	result.isBestialFailure = result.isTotalFailure && (possibleBestialFailures >= 1)

	result.isMessyCritical = result.isCritical && (halfMessyCriticals >= 1)

	return result
}

// Successes returns the total number of successes, including raw successes and bonus successes from critical pairs.
//
// Each complete pair of half-criticals (rolls at RollUpperLimit) grants 2 bonus successes.
func (vtmres VTMRollerResult) Successes() int {
	return vtmres.successes
}

// Failures returns the total number of dice that did not meet the success threshold,
// including hunger failures and possible bestial failures.
func (vtmres VTMRollerResult) Failures() int {
	return vtmres.failures
}

// IsCritical reports whether the roll achieved a critical success (two or more half-criticals).
func (vtmres VTMRollerResult) IsCritical() bool {
	return vtmres.isCritical
}

// IsTotalFailure reports whether the roll resulted in zero successes.
func (vtmres VTMRollerResult) IsTotalFailure() bool {
	return vtmres.isTotalFailure
}

// IsBestialFailure reports whether the roll is a total failure with at least one hunger die involved.
//
// This represents a catastrophic failure in VTM 5e that triggers the Beast.
func (vtmres VTMRollerResult) IsBestialFailure() bool {
	return vtmres.isBestialFailure
}

// IsMessyCritical reports whether the roll achieved a critical success where at least one half-critical came from a hunger die.
//
// This represents an unpredictable critical in VTM 5e that deals collateral damage.
func (vtmres VTMRollerResult) IsMessyCritical() bool {
	return vtmres.isMessyCritical
}

// GetRolls returns a copy of all individual die rolls in order.
func (vtmres VTMRollerResult) GetRolls() []int {
	rolls := make([]int, len(vtmres.rolls))
	copy(rolls, vtmres.rolls)
	return rolls
}

// GetHungerDice returns the count of hunger dice in this roll.
func (vtmres VTMRollerResult) HungerDice() int {
	return vtmres.hungerDice
}

// GetRollTypes returns a copy of the RollType classification for each die in order.
func (vtmres VTMRollerResult) GetRollTypes() []RollType {
	rollTypes := make([]RollType, len(vtmres.rollTypes))
	copy(rollTypes, vtmres.rollTypes)
	return rollTypes
}

// Rolls returns an iterator over each die's roll value and RollType classification.
//
// The iterator yields (rollValue, rollType) pairs in die order. Hunger dice come first.
func (vtmres VTMRollerResult) Rolls() iter.Seq2[int, RollType] {
	return func(yield func(int, RollType) bool) {
		for i := range vtmres.rolls {
			if !yield(vtmres.rolls[i], vtmres.rollTypes[i]) {
				return
			}
		}
	}
}

// Len returns the amount of rolls in the result (the length).
func (vtmres VTMRollerResult) Len() int {
	return len(vtmres.rolls)
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

// IsSuccess reports whether the given roll meets or exceeds the success threshold.
func (vtmr *VTMRoller) IsSuccess(roll int) bool {
	return roll >= vtmr.SuccessThreshold
}

// IsUpperLimit reports whether the given roll equals the upper limit of the die range,
// indicating a potential critical result.
func (vtmr *VTMRoller) IsUpperLimit(roll int) bool {
	return roll == vtmr.RollUpperLimit
}

// IsLowerLimit reports whether the given roll equals the lower limit of the die range,
// indicating a potential bestial failure for hunger dice.
func (vtmr *VTMRoller) IsLowerLimit(roll int) bool {
	return roll == vtmr.RollLowerLimit
}

// RollDie generates a single die roll within the configured range [RollLowerLimit, RollUpperLimit].
func (vtmr *VTMRoller) RollDie() int {
	return vtmr.Rand.IntN(vtmr.RollUpperLimit+1-vtmr.RollLowerLimit) + vtmr.RollLowerLimit
}

// Roll performs a dice roll with the given pool size and hunger dice count.
//
// Hunger dice count is clamped to [0, pool]. The returned result contains all roll details,
// computed successes, and special outcome states.
func (vtmr *VTMRoller) Roll(pool int, hungerDice int) VTMRollerResult {
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
//
// Hunger dice cannot be rerolled. Indices must be valid (0-based, within pool size, and not hunger dice).
// Returns an error if any index is invalid or refers to a hunger die.
func (vtmr *VTMRoller) ReRoll(oldResult VTMRollerResult, rerollPlaces ...int) (VTMRollerResult, error) {
	newRolls := make([]int, len(oldResult.rolls))
	copy(newRolls, oldResult.rolls)

	for _, rp := range rerollPlaces {
		if rp < oldResult.hungerDice {
			return VTMRollerResult{}, fmt.Errorf("cannot reroll hunger die (hunger: %d, offset %d)", oldResult.hungerDice, rp)
		}

		if rp >= len(newRolls) {
			return VTMRollerResult{}, fmt.Errorf("offset greater than or equal to rolls (len: %d, offset: %d)", len(newRolls), rp)
		}

		newRolls[rp] = vtmr.RollDie()
	}

	newResult := NewVTMRollerResult(newRolls, vtmr, oldResult.hungerDice)

	return newResult, nil
}
