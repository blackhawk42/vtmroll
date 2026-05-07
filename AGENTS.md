# AGENTS.md

## Project Overview

**vtmroll** is a Go library (and tools making use of that library) that simulates dice rolls for Vampire: The Masquerade 5th Edition. It implements custom rules around success thresholds, hunger dice, and critical mechanics.

## Build & Test Commands

```bash
# Run tests
go test ./...
```

## Project Structure

- `pkg/vtmroll/`: Core dice-rolling logic
  - `vtmroll.go`: `VTMRoller` (roll executor) and `VTMRollerResult` (result container with public getter methods)
  - `vtmroll_test.go`: Comprehensive test suite covering all game rules
- `pkg/vtmrollfmt/`: Formatting and display utilities for roll results (in development)

## Key Domain Logic to Preserve

### Roll Mechanics
- **Dice pool & hunger dice**: Pool contains N dice, first H are marked as "hunger" (affect special outcomes)
- **Success threshold**: Default 6, configurable via `VTMRoller.SuccessThreshold`
- **Dice range**: Default 1–10, configurable via `VTMRoller.RollLowerLimit` and `VTMRoller.RollUpperLimit`

### Critical & Special Outcomes
- **Critical (pair)**: 2+ rolls at `RollUpperLimit` → each complete pair grants 2 bonus successes
- **Messy critical**: Critical where at least one half-critical is from a hunger die
- **Bestial failure**: Total failure (0 successes) with at least one hunger die
- **Reroll rules**: Cannot reroll hunger dice; validates reroll indices strictly

### Result Calculation
- `Successes()` = raw successes + (half-criticals / 2) * 2
- Hunger dice always occupy the first N positions in the roll slice

## Testing Strategy

Tests in `vtmroll_test.go` use **table-driven tests** with manual roll creation (`NewVTMRollerResult()`) to test logic without RNG unpredictability. Most assertions are on counts and booleans rather than exact roll values.

Key test patterns:
- Manual VTMRollerResult construction for deterministic logic testing
- RNG seeding (`rand.NewPCG(0, 0)`) for reproducible randomness when needed
- Boundary conditions: threshold at upper/lower limits, pool edge cases, reroll validation

## Important Implementation Notes

- **Result immutability**: `Get*()` methods return copies, not references, to prevent external mutations
- **Default RNG initialization**: `NewVTMRoller()` creates a new random source with `rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64()))`, ensuring each roller instance has a unique sequence by default
- **RNG override for reproducible results**: The `Rand` field is public. Tests and code needing repeatable rolls can override it: `roller.Rand = rand.New(rand.NewPCG(0, 0))`. Set this before calling `Roll()`.
- **Hunger clamping**: `Roll(pool, hungerDice)` clamps hunger to `[0, pool]`
- **Reroll validation**: `ReRoll()` prevents rerolling hunger dice; validates reroll indices strictly
