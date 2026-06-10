// Package vtmrollfmt provides formatting and display utilities for Vampire: The
// Masquerade 5th Edition dice roll results.
package vtmrollfmt

import "github.com/blackhawk42/vtmroll/pkg/vtmroll"

// VTMRollResultDiceFormatter formats a roll result into a human-readable
// representation of each die.
type VTMRollResultDiceFormatter interface {
	FormatDice(vtmres vtmroll.VTMRollerResult) []string
}

// VTMRollResultDiceFormatterFunc is an adapter that allows a plain function to
// be used as a VTMRollResultDiceFormatter.
type VTMRollResultDiceFormatterFunc func(vtmres vtmroll.VTMRollerResult) []string

func (f VTMRollResultDiceFormatterFunc) FormatDice(vtmres vtmroll.VTMRollerResult) []string {
	return f(vtmres)
}

// VTMRollResultDiceParser parses string representations of dice rolls back
// into a VTMRollerResult.
//
// Tokenization of the individual rolls is left to the caller.
type VTMRollResultDiceParser interface {
	Parse(rolls []string, roller *vtmroll.VTMRoller, hungerDice int) (vtmroll.VTMRollerResult, error)
}

// VTMRollResultDiceParserFunc is an adapter that allows a plain function to be
// used as a VTMRollResultDiceParser.
type VTMRollResultDiceParserFunc func(rolls []string, roller *vtmroll.VTMRoller, hungerDice int) (vtmroll.VTMRollerResult, error)

func (f VTMRollResultDiceParserFunc) Parse(rolls []string, roller *vtmroll.VTMRoller, hungerDice int) (vtmroll.VTMRollerResult, error) {
	return f(rolls, roller, hungerDice)
}

// VTMRollResultSummaryMessages holds the textual, humaan-readable messages that
// describe the outcome of a dice roll, such as the number of successes and any special
// conditions (critical, bestial failure, etc.).
type VTMRollResultSummaryMessages struct {
	SuccessesMessage        string
	IsCriticalMessage       string
	IsTotalFailureMessage   string
	IsBestialFailureMessage string
	IsMessyCriticalMessage  string
}

// VTMRollResultSummarizer generates human-readable summary messages from a
// roll result.
type VTMRollResultSummarizer interface {
	Summarize(vtmres vtmroll.VTMRollerResult) VTMRollResultSummaryMessages
}

// VTMRollResultSummarizerFunc is an adapter that allows a plain function to be
// used as a VTMRollResultSummarizer.
type VTMRollResultSummarizerFunc func(vtmres vtmroll.VTMRollerResult) VTMRollResultSummaryMessages

func (f VTMRollResultSummarizerFunc) Summarize(vtmres vtmroll.VTMRollerResult) VTMRollResultSummaryMessages {
	return f(vtmres)
}
