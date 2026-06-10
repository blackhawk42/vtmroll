// Package fmtbuiltins provides built-in implementations of the formatters,
// parsers, and style presets defined in vtmrollfmt.
package fmtbuiltins

import (
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/blackhawk42/vtmroll/pkg/vtmroll"
	"github.com/blackhawk42/vtmroll/pkg/vtmrollfmt"
	"github.com/charmbracelet/colorprofile"
)

// DiceStyles holds lipgloss styles for each die roll type.
type DiceStyles struct {
	NormalSuccessStyle lipgloss.Style

	NormalFailureStyle lipgloss.Style

	HungerSuccessStyle lipgloss.Style

	HungerFailureStyle lipgloss.Style

	HalfCriticalStyle lipgloss.Style

	HalfMessyCriticalStyle lipgloss.Style

	PossibleBestialFailureStyle lipgloss.Style
}

// NewDiceStyle returns a DiceStyles with all styles set to their zero value.
func NewDiceStyle() DiceStyles {
	return DiceStyles{
		NormalSuccessStyle:          lipgloss.NewStyle(),
		NormalFailureStyle:          lipgloss.NewStyle(),
		HungerSuccessStyle:          lipgloss.NewStyle(),
		HungerFailureStyle:          lipgloss.NewStyle(),
		HalfCriticalStyle:           lipgloss.NewStyle(),
		HalfMessyCriticalStyle:      lipgloss.NewStyle(),
		PossibleBestialFailureStyle: lipgloss.NewStyle(),
	}
}

// DieFormatFunction converts a single die roll and its RollType into a string
// representation.
type DieFormatFunction func(roll int, rollType vtmroll.RollType) string

// DiceFormatter formats VTMRollerResult dice by applying a DieFormatFunction
// and styling each result with the corresponding DiceStyles.
//
// The DiceFormatter includes a colorprofile.Profile (a companion library to lipgloss).
// This allows to easily downgrade or completely deactivate ANSI coloring sequencies.
type DiceFormatter struct {
	DiceStyles DiceStyles

	colorprofileWriter *colorprofile.Writer
	colorBuffer        *strings.Builder

	formatFunction DieFormatFunction
}

// NewDiceFormatter creates a DiceFormatter.
//
// If formatFunction is nil, BUILTIN_FORMATFUNCTION_NUMERIC is used.
func NewDiceFormatter(
	formatFunction DieFormatFunction,
	diceStyles DiceStyles,
	colorProfile colorprofile.Profile,
) *DiceFormatter {
	if formatFunction == nil {
		formatFunction = BUILTIN_FORMATFUNCTION_NUMERIC
	}

	colorBuffer := new(strings.Builder)

	colorprofileWriter := new(colorprofile.Writer)
	colorprofileWriter.Forward = colorBuffer
	colorprofileWriter.Profile = colorProfile

	return &DiceFormatter{
		DiceStyles: diceStyles,

		formatFunction:     formatFunction,
		colorBuffer:        colorBuffer,
		colorprofileWriter: colorprofileWriter,
	}
}

var _ vtmrollfmt.VTMRollResultDiceFormatter = (*DiceFormatter)(nil)

func (fmtter *DiceFormatter) FormatDice(vtmres vtmroll.VTMRollerResult) []string {
	results := make([]string, 0, vtmres.Len())

	var currentResult string

	for roll, rollType := range vtmres.Rolls() {
		// NOTE: Resetting colorBuffer while reusing colorprofileWriter across iterations
		// appears safe based on the implementation and observed behavior, but may need
		// further clarification if colorprofile.Writer is found to maintain internal state
		// that could be corrupted by resetting only the underlying buffer.
		fmtter.colorBuffer.Reset()

		currentResult = fmtter.formatFunction(roll, rollType)

		switch rollType {
		case vtmroll.NormalSuccess:
			currentResult = fmtter.DiceStyles.NormalSuccessStyle.Render(currentResult)
		case vtmroll.NormalFailure:
			currentResult = fmtter.DiceStyles.NormalFailureStyle.Render(currentResult)
		case vtmroll.HungerSuccess:
			currentResult = fmtter.DiceStyles.HungerSuccessStyle.Render(currentResult)
		case vtmroll.HungerFailure:
			currentResult = fmtter.DiceStyles.HungerFailureStyle.Render(currentResult)
		case vtmroll.HalfCritical:
			currentResult = fmtter.DiceStyles.HalfCriticalStyle.Render(currentResult)
		case vtmroll.HalfMessyCritical:
			currentResult = fmtter.DiceStyles.HalfMessyCriticalStyle.Render(currentResult)
		case vtmroll.PossibleBestialFailure:
			currentResult = fmtter.DiceStyles.PossibleBestialFailureStyle.Render(currentResult)
		}

		fmtter.colorprofileWriter.WriteString(currentResult)

		results = append(results, fmtter.colorBuffer.String())
	}

	return results
}

// SummaryStyles holds lipgloss styles for each summary message category.
type SummaryStyles struct {
	SuccessesMessageStyle        lipgloss.Style
	IsCriticalMessageStyle       lipgloss.Style
	IsTotalFailureMessageStyle   lipgloss.Style
	IsBestialFailureMessageStyle lipgloss.Style
	IsMessyCriticalMessageStyle  lipgloss.Style
}

// NewSummaryStyles returns a SummaryStyles with all styles set to their zero
// value.
func NewSummaryStyles() SummaryStyles {
	return SummaryStyles{
		SuccessesMessageStyle:        lipgloss.NewStyle(),
		IsCriticalMessageStyle:       lipgloss.NewStyle(),
		IsTotalFailureMessageStyle:   lipgloss.NewStyle(),
		IsBestialFailureMessageStyle: lipgloss.NewStyle(),
		IsMessyCriticalMessageStyle:  lipgloss.NewStyle(),
	}
}

// SummaryFormatFunction generates a VTMRollResultSummaryMessages from a
// VTMRollerResult.
type SummaryFormatFunction func(result vtmroll.VTMRollerResult) vtmrollfmt.VTMRollResultSummaryMessages

// ResultSummarizer formats a VTMRollerResult into styled summary messages.
//
// The ResultSummarizer includes a colorprofile.Profile (a companion library to lipgloss).
// This allows to easily downgrade or completely deactivate ANSI coloring sequencies.
type ResultSummarizer struct {
	SummaryStyles SummaryStyles

	summaryFormatFunction SummaryFormatFunction

	colorprofileWriter *colorprofile.Writer
	colorBuffer        *strings.Builder
}

// NewResultSummarizer creates a ResultSummarizer. If summaryFormatFunction is
// nil, BUILTIN_SUMMARYFORMATFUNCTION_SIMPLE is used.
func NewResultSummarizer(
	summarizeStyles SummaryStyles,
	summaryFormatFunction SummaryFormatFunction,
	colorProfile colorprofile.Profile,
) *ResultSummarizer {
	if summaryFormatFunction == nil {
		summaryFormatFunction = BUILTIN_SUMMARYFORMATFUNCTION_SIMPLE
	}

	colorBuffer := new(strings.Builder)

	colorprofileWriter := new(colorprofile.Writer)
	colorprofileWriter.Forward = colorBuffer
	colorprofileWriter.Profile = colorProfile

	return &ResultSummarizer{
		summaryFormatFunction: summaryFormatFunction,
		SummaryStyles:         summarizeStyles,
		colorprofileWriter:    colorprofileWriter,
		colorBuffer:           colorBuffer,
	}
}

var _ vtmrollfmt.VTMRollResultSummarizer = (*ResultSummarizer)(nil)

func (rsumm *ResultSummarizer) Summarize(vtmres vtmroll.VTMRollerResult) vtmrollfmt.VTMRollResultSummaryMessages {
	summary := rsumm.summaryFormatFunction(vtmres)

	rsumm.colorBuffer.Reset()
	rsumm.colorprofileWriter.WriteString(rsumm.SummaryStyles.SuccessesMessageStyle.Render(summary.SuccessesMessage))
	summary.SuccessesMessage = rsumm.colorBuffer.String()

	if summary.IsCriticalMessage != "" {
		rsumm.colorBuffer.Reset()
		rsumm.colorprofileWriter.WriteString(rsumm.SummaryStyles.IsCriticalMessageStyle.Render(summary.IsCriticalMessage))
		summary.IsCriticalMessage = rsumm.colorBuffer.String()
	}

	if summary.IsTotalFailureMessage != "" {
		rsumm.colorBuffer.Reset()
		rsumm.colorprofileWriter.WriteString(rsumm.SummaryStyles.IsTotalFailureMessageStyle.Render(summary.IsTotalFailureMessage))
		summary.IsTotalFailureMessage = rsumm.colorBuffer.String()
	}

	if summary.IsBestialFailureMessage != "" {
		rsumm.colorBuffer.Reset()
		rsumm.colorprofileWriter.WriteString(rsumm.SummaryStyles.IsBestialFailureMessageStyle.Render(summary.IsBestialFailureMessage))
		summary.IsBestialFailureMessage = rsumm.colorBuffer.String()
	}

	if summary.IsMessyCriticalMessage != "" {
		rsumm.colorBuffer.Reset()
		rsumm.colorprofileWriter.WriteString(rsumm.SummaryStyles.IsMessyCriticalMessageStyle.Render(summary.IsMessyCriticalMessage))
		summary.IsMessyCriticalMessage = rsumm.colorBuffer.String()
	}

	return summary
}
