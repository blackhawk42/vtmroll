package vtmroll

import (
	"math/rand/v2"
	"slices"
	"testing"
)

func TestVTMRollerBasic(t *testing.T) {
	roller := NewVTMRoller()
	roller.Rand = rand.New(rand.NewPCG(0, 0))

	// Test basic roll
	result := roller.Roll(5, 2)
	if len(result.rolls) != 5 {
		t.Errorf("Expected 5 rolls, got %d", len(result.rolls))
	}
	if result.hungerDice != 2 {
		t.Errorf("Expected hungerDice = 2, got %d", result.hungerDice)
	}
}

func TestSuccessesCalculation(t *testing.T) {
	roller := NewVTMRoller()
	roller.Rand = rand.New(rand.NewPCG(0, 0))

	// Test with no successes
	roller.SuccessThreshold = 11 // Set threshold high to ensure no successes
	result := roller.Roll(5, 0)
	if result.Successes() != 0 {
		t.Errorf("Expected 0 successes with threshold 11, got %d", result.Successes())
	}

	// Reset threshold
	roller.SuccessThreshold = 6
}

func TestReRollValidation(t *testing.T) {
	roller := NewVTMRoller()
	roller.Rand = rand.New(rand.NewPCG(0, 0))

	// Create initial roll with 2 hunger dice
	result := roller.Roll(8, 2)

	// Test rerolling a hunger die (should fail)
	_, err := roller.ReRoll(result, 0) // Index 0 is a hunger die (first of 2)
	if err == nil {
		t.Error("Expected error when rerolling hunger die")
	}

	// Test rerolling a non-hunger die (should succeed)
	newResult, err := roller.ReRoll(result, 3) // Index 3 is not a hunger die
	if err != nil {
		t.Errorf("Unexpected error when rerolling non-hunger die: %v", err)
	}
	if newResult == nil {
		t.Error("Expected new result from reroll")
	}

	// Test rerolling with invalid index (should fail)
	_, err = roller.ReRoll(result, 10) // Index out of bounds
	if err == nil {
		t.Error("Expected error when rerolling with out-of-bounds index")
	}
}

func TestHungerDiceValidation(t *testing.T) {
	roller := NewVTMRoller()
	roller.Rand = rand.New(rand.NewPCG(0, 0))

	// Test with hungerDice > pool
	result := roller.Roll(5, 10) // hungerDice > pool
	if result.hungerDice != 5 {
		t.Errorf("Expected hungerDice clamped to 5, got %d", result.hungerDice)
	}

	// Test with negative hungerDice
	result2 := roller.Roll(5, -3)
	if result2.hungerDice != 0 {
		t.Errorf("Expected hungerDice clamped to 0, got %d", result2.hungerDice)
	}
}

func TestCriticalSuccessCounting(t *testing.T) {
	// Test specific cases for critical success counting
	tests := []struct {
		name              string
		rolls             []int
		threshold         int
		expectedSuccesses int
	}{
		{
			name:              "No criticals, some successes",
			rolls:             []int{6, 7, 8, 3, 4},
			threshold:         6,
			expectedSuccesses: 3, // 6, 7, 8
		},
		{
			name:              "One half critical (10)",
			rolls:             []int{10, 6, 3, 4, 2},
			threshold:         6,
			expectedSuccesses: 2, // 10 (1) + 6 (1) = 2 (no pair for extra)
		},
		{
			name:              "One critical (pair)",
			rolls:             []int{10, 10, 3, 4, 2},
			threshold:         6,
			expectedSuccesses: 4, // 10+10 = 2 raw + 2 extra for pair
		},
		{
			name:              "Three succeses, on critical",
			rolls:             []int{10, 10, 10, 4, 2},
			threshold:         6,
			expectedSuccesses: 5, // 3 raw + 2 extra for one pair
		},
		{
			name:              "Four criticals (two pairs)",
			rolls:             []int{10, 10, 10, 10, 2},
			threshold:         6,
			expectedSuccesses: 8, // 4 raw + 4 extra for two pairs
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock roller
			roller := NewVTMRoller()
			roller.Rand = rand.New(rand.NewPCG(0, 0))
			roller.SuccessThreshold = tt.threshold

			// Create result manually with the specified rolls
			result := NewVTMRollerResult(tt.rolls, roller, 0)

			if result.Successes() != tt.expectedSuccesses {
				t.Errorf("%s: expected %d successes, got %d", tt.name, tt.expectedSuccesses, result.Successes())
			}
		})
	}
}

func TestMessyCriticalManualTests(t *testing.T) {
	// Test specific cases for messy criticals
	tests := []struct {
		name               string
		rolls              []int
		hungerDice         int
		expectedIsCritical bool
		expectedIsMessy    bool
	}{
		{
			name:               "Regular critical (no hunger dice)",
			rolls:              []int{10, 10, 6, 7, 3},
			hungerDice:         0,
			expectedIsCritical: true,
			expectedIsMessy:    false,
		},
		{
			name:               "Messy critical with one hunger half-critical",
			rolls:              []int{10, 10, 6, 7, 3},
			hungerDice:         2,
			expectedIsCritical: true,
			expectedIsMessy:    true,
		},
		{
			name:               "Critical but half-criticals are non-hunger",
			rolls:              []int{6, 7, 10, 10, 3},
			hungerDice:         2,
			expectedIsCritical: true,
			expectedIsMessy:    false,
		},
		{
			name:               "Three half-criticals with one hunger",
			rolls:              []int{10, 10, 10, 7, 3},
			hungerDice:         1,
			expectedIsCritical: true,
			expectedIsMessy:    true,
		},
		{
			name:               "Not a critical (only one half-critical)",
			rolls:              []int{10, 6, 7, 8, 3},
			hungerDice:         2,
			expectedIsCritical: false,
			expectedIsMessy:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			roller := NewVTMRoller()
			roller.Rand = rand.New(rand.NewPCG(0, 0))

			result := NewVTMRollerResult(tt.rolls, roller, tt.hungerDice)

			if result.IsCritical() != tt.expectedIsCritical {
				t.Errorf("%s: IsCritical() = %v, expected %v", tt.name, result.IsCritical(), tt.expectedIsCritical)
			}
			if result.IsMessyCritical() != tt.expectedIsMessy {
				t.Errorf("%s: IsMessyCritical() = %v, expected %v", tt.name, result.IsMessyCritical(), tt.expectedIsMessy)
			}
		})
	}
}

func TestBestialFailureManualTests(t *testing.T) {
	// Test specific cases for bestial failures
	tests := []struct {
		name                   string
		rolls                  []int
		hungerDice             int
		threshold              int
		expectedIsTotalFailure bool
		expectedIsBestial      bool
	}{
		{
			name:                   "Total failure with hunger dice",
			rolls:                  []int{1, 2, 3, 4, 5},
			hungerDice:             2,
			threshold:              6,
			expectedIsTotalFailure: true,
			expectedIsBestial:      true,
		},
		{
			name:                   "Total failure without hunger dice",
			rolls:                  []int{1, 2, 3, 4, 5},
			hungerDice:             0,
			threshold:              6,
			expectedIsTotalFailure: true,
			expectedIsBestial:      false,
		},
		{
			name:                   "Success with hunger dice (not bestial)",
			rolls:                  []int{6, 2, 3, 4, 5},
			hungerDice:             2,
			threshold:              6,
			expectedIsTotalFailure: false,
			expectedIsBestial:      false,
		},
		{
			name:                   "Total failure but hunger dice succeed",
			rolls:                  []int{6, 7, 3, 4, 5},
			hungerDice:             2,
			threshold:              6,
			expectedIsTotalFailure: false, // Has successes
			expectedIsBestial:      false,
		},
		{
			name:                   "Edge case: threshold 10 with all 9s",
			rolls:                  []int{9, 9, 9, 9, 9},
			hungerDice:             2,
			threshold:              10,
			expectedIsTotalFailure: true,
			expectedIsBestial:      false,
		},
		{
			name:                   "Edge case: threshold 10 with all 9s except first is 1",
			rolls:                  []int{1, 9, 9, 9, 9},
			hungerDice:             2,
			threshold:              10,
			expectedIsTotalFailure: true,
			expectedIsBestial:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			roller := NewVTMRoller()
			roller.Rand = rand.New(rand.NewPCG(0, 0))
			roller.SuccessThreshold = tt.threshold

			result := NewVTMRollerResult(tt.rolls, roller, tt.hungerDice)

			if result.IsTotalFailure() != tt.expectedIsTotalFailure {
				t.Errorf("%s: IsTotalFailure() = %v, expected %v", tt.name, result.IsTotalFailure(), tt.expectedIsTotalFailure)
			}
			if result.IsBestialFailure() != tt.expectedIsBestial {
				t.Errorf("%s: IsBestialFailure() = %v, expected %v", tt.name, result.IsBestialFailure(), tt.expectedIsBestial)
			}
		})
	}
}

func TestDifferentThresholdsManualTests(t *testing.T) {
	// Test with different success thresholds
	tests := []struct {
		name              string
		rolls             []int
		threshold         int
		expectedSuccesses int
	}{
		{
			name:              "Threshold 6 (default)",
			rolls:             []int{5, 6, 7, 8, 9},
			threshold:         6,
			expectedSuccesses: 4, // 6, 7, 8, 9
		},
		{
			name:              "Threshold 7",
			rolls:             []int{5, 6, 7, 8, 9},
			threshold:         7,
			expectedSuccesses: 3, // 7, 8, 9
		},
		{
			name:              "Threshold 8",
			rolls:             []int{5, 6, 7, 8, 9},
			threshold:         8,
			expectedSuccesses: 2, // 8, 9
		},
		{
			name:              "Threshold 10 (only 10s succeed)",
			rolls:             []int{9, 9, 10, 10, 1},
			threshold:         10,
			expectedSuccesses: 4, // 2 raw successes + 2 extra for pair
		},
		{
			name:              "Threshold 2 (almost everything succeeds)",
			rolls:             []int{1, 2, 3, 4, 5},
			threshold:         2,
			expectedSuccesses: 4, // 2, 3, 4, 5
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			roller := NewVTMRoller()
			roller.Rand = rand.New(rand.NewPCG(0, 0))
			roller.SuccessThreshold = tt.threshold

			result := NewVTMRollerResult(tt.rolls, roller, 0)

			if result.Successes() != tt.expectedSuccesses {
				t.Errorf("%s: expected %d successes, got %d", tt.name, tt.expectedSuccesses, result.Successes())
			}
		})
	}
}

func TestDifferentLimitsManualTests(t *testing.T) {
	// Test with different roll upper limits (e.g., for different game systems)
	tests := []struct {
		name              string
		rolls             []int
		upperLimit        int
		threshold         int
		expectedSuccesses int
	}{
		{
			name:              "Standard VTM (limit 10, threshold 6)",
			rolls:             []int{9, 10, 10, 6, 3},
			upperLimit:        10,
			threshold:         6,
			expectedSuccesses: 6, // 9, 10, 10, 6 = 4 raw + 2 extra for pair of 10s
		},
		{
			name:              "Limit 8 with threshold 6",
			rolls:             []int{7, 8, 8, 6, 3},
			upperLimit:        8,
			threshold:         6,
			expectedSuccesses: 6, // 7, 8, 8, 6 = 4 raw + 2 extra for pair of 8s (half-criticals)
		},
		{
			name:              "Limit 6 with threshold 6",
			rolls:             []int{5, 6, 6, 4, 3},
			upperLimit:        6,
			threshold:         6,
			expectedSuccesses: 4, // 6, 6 = 2 raw + 2 extra for pair of 6s (half-criticals)
		},
		{
			name:              "Limit 12 with threshold 6",
			rolls:             []int{11, 12, 12, 7, 3},
			upperLimit:        12,
			threshold:         6,
			expectedSuccesses: 6, // 11, 12, 12, 7 = 4 raw + 2 extra for pair of 12s (half-criticals)
		},
		{
			name:              "Limit 8 with threshold 7",
			rolls:             []int{7, 8, 8, 6, 3},
			upperLimit:        8,
			threshold:         7,
			expectedSuccesses: 5, // 7, 8, 8 = 3 raw + 2 extra for pair of 8 (half-criticals)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			roller := NewVTMRoller()
			roller.Rand = rand.New(rand.NewPCG(0, 0))
			roller.RollUpperLimit = tt.upperLimit
			roller.SuccessThreshold = tt.threshold

			result := NewVTMRollerResult(tt.rolls, roller, 0)

			if result.Successes() != tt.expectedSuccesses {
				t.Errorf("%s: expected %d successes, got %d", tt.name, tt.expectedSuccesses, result.Successes())
			}
		})
	}
}

func TestRollsIterator(t *testing.T) {
	roller := NewVTMRoller()
	roller.SuccessThreshold = 6

	tests := []struct {
		name           string
		rolls          []int
		hungerDice     int
		expectedRolls  []int
		expectedTypes  []RollType
	}{
		{
			name:          "Mixed hunger and normal dice",
			rolls:         []int{10, 1, 8, 3, 6},
			hungerDice:    2,
			expectedRolls: []int{10, 1, 8, 3, 6},
			expectedTypes: []RollType{HalfMessyCritical, PossibleBestialFailure, NormalSuccess, NormalFailure, NormalSuccess},
		},
		{
			name:          "No hunger dice",
			rolls:         []int{10, 10, 6, 3},
			hungerDice:    0,
			expectedRolls: []int{10, 10, 6, 3},
			expectedTypes: []RollType{HalfCritical, HalfCritical, NormalSuccess, NormalFailure},
		},
		{
			name:          "All hunger dice",
			rolls:         []int{6, 7, 1},
			hungerDice:    3,
			expectedRolls: []int{6, 7, 1},
			expectedTypes: []RollType{HungerSuccess, HungerSuccess, PossibleBestialFailure},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewVTMRollerResult(tt.rolls, roller, tt.hungerDice)

			var gotRolls []int
			var gotTypes []RollType
			for roll, rt := range result.Rolls() {
				gotRolls = append(gotRolls, roll)
				gotTypes = append(gotTypes, rt)
			}

			if !slices.Equal(gotRolls, tt.expectedRolls) {
				t.Errorf("rolls: got %v, want %v", gotRolls, tt.expectedRolls)
			}
			if !slices.Equal(gotTypes, tt.expectedTypes) {
				t.Errorf("rollTypes: got %v, want %v", gotTypes, tt.expectedTypes)
			}
		})
	}
}

func TestRollsIteratorEarlyBreak(t *testing.T) {
	roller := NewVTMRoller()
	roller.SuccessThreshold = 6
	result := NewVTMRollerResult([]int{10, 1, 8, 3, 6}, roller, 2)

	count := 0
	for range result.Rolls() {
		count++
		if count == 2 {
			break
		}
	}

	if count != 2 {
		t.Errorf("expected 2 iterations after break, got %d", count)
	}
}

func TestReRollManualTests(t *testing.T) {
	tests := []struct {
		name                  string
		initialRolls          []int
		hungerDice            int
		indicesToReroll       []int
		expectError           bool
		expectedNewRollsCount int // We can't predict exact rolls, but we can check count
	}{
		{
			name:                  "Reroll single non-hunger die",
			initialRolls:          []int{1, 2, 3, 4, 5},
			hungerDice:            2,
			indicesToReroll:       []int{3},
			expectError:           false,
			expectedNewRollsCount: 5,
		},
		{
			name:                  "Reroll multiple non-hunger dice",
			initialRolls:          []int{1, 2, 3, 4, 5},
			hungerDice:            2,
			indicesToReroll:       []int{3, 4},
			expectError:           false,
			expectedNewRollsCount: 5,
		},
		{
			name:                  "Try to reroll hunger die (should error)",
			initialRolls:          []int{1, 2, 3, 4, 5},
			hungerDice:            2,
			indicesToReroll:       []int{0},
			expectError:           true,
			expectedNewRollsCount: 5,
		},
		{
			name:                  "Try to reroll out of bounds index",
			initialRolls:          []int{1, 2, 3, 4, 5},
			hungerDice:            2,
			indicesToReroll:       []int{10},
			expectError:           true,
			expectedNewRollsCount: 5,
		},
		{
			name:                  "Reroll same index twice",
			initialRolls:          []int{1, 2, 3, 4, 5},
			hungerDice:            2,
			indicesToReroll:       []int{3, 3},
			expectError:           false,
			expectedNewRollsCount: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			roller := NewVTMRoller()
			roller.Rand = rand.New(rand.NewPCG(0, 0))

			// Create initial result
			initialResult := NewVTMRollerResult(tt.initialRolls, roller, tt.hungerDice)

			// Try to reroll
			newResult, err := roller.ReRoll(initialResult, tt.indicesToReroll...)

			if tt.expectError {
				if err == nil {
					t.Errorf("%s: expected error but got none", tt.name)
				}
				return
			}

			if err != nil {
				t.Errorf("%s: unexpected error: %v", tt.name, err)
				return
			}

			if newResult == nil {
				t.Errorf("%s: newResult is nil", tt.name)
				return
			}

			if len(newResult.rolls) != tt.expectedNewRollsCount {
				t.Errorf("%s: expected %d rolls, got %d", tt.name, tt.expectedNewRollsCount, len(newResult.rolls))
			}

			// Verify that the rerolled indices have changed (they will be different due to RNG)
			for _, idx := range tt.indicesToReroll {
				if idx < len(initialResult.rolls) && idx < len(newResult.rolls) {
					if initialResult.rolls[idx] == newResult.rolls[idx] {
						// This could happen by chance, but it's unlikely
						t.Logf("%s: index %d has same value before and after reroll (%d)", tt.name, idx, initialResult.rolls[idx])
					}
				}
			}

			// Verify non-rerolled indices are unchanged
			for i := 0; i < len(initialResult.rolls); i++ {
				isRerolled := false
				for _, idx := range tt.indicesToReroll {
					if i == idx {
						isRerolled = true
						break
					}
				}

				if !isRerolled && i < len(newResult.rolls) {
					if initialResult.rolls[i] != newResult.rolls[i] {
						t.Errorf("%s: non-rerolled index %d changed from %d to %d", tt.name, i, initialResult.rolls[i], newResult.rolls[i])
					}
				}
			}
		})
	}
}
