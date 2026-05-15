package fmtbuiltins

import (
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/blackhawk42/vtmroll/pkg/vtmroll"
	"github.com/blackhawk42/vtmroll/pkg/vtmrollfmt"
	"github.com/charmbracelet/colorprofile"
)

type DiceStyles struct {
	NormalSuccessStyle lipgloss.Style

	NormalFailureStyle lipgloss.Style

	HungerSuccessStyle lipgloss.Style

	HungerFailureStyle lipgloss.Style

	HalfCriticalStyle lipgloss.Style

	HalfMessyCriticalStyle lipgloss.Style

	PossibleBestialFailureStyle lipgloss.Style
}

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

type DiceFormatter struct {
	DiceStyles DiceStyles

	colorprofileWriter *colorprofile.Writer
	colorBuffer        *strings.Builder

	formatFunction DieFormatFunction
}

type DieFormatFunction func(roll int, rollType vtmroll.RollType) string

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

type SummaryStyles struct {
	SuccessesMessageStyle        lipgloss.Style
	IsCriticalMessageStyle       lipgloss.Style
	IsTotalFailureMessageStyle   lipgloss.Style
	IsBestialFailureMessageStyle lipgloss.Style
	IsMessyCriticalMessageStyle  lipgloss.Style
}

func NewSummaryStyles() SummaryStyles {
	return SummaryStyles{
		SuccessesMessageStyle:        lipgloss.NewStyle(),
		IsCriticalMessageStyle:       lipgloss.NewStyle(),
		IsTotalFailureMessageStyle:   lipgloss.NewStyle(),
		IsBestialFailureMessageStyle: lipgloss.NewStyle(),
		IsMessyCriticalMessageStyle:  lipgloss.NewStyle(),
	}
}

type SummaryFormatFunction func(result vtmroll.VTMRollerResult) vtmrollfmt.VTMRollResultSummaryMessages

type ResultSummarizer struct {
	SummaryStyles SummaryStyles

	summaryFormatFunction SummaryFormatFunction

	colorprofileWriter *colorprofile.Writer
	colorBuffer        *strings.Builder
}

func NewResultSummarizer(
	summarizeStyles SummaryStyles,
	summarSummaryFormatFunction SummaryFormatFunction,
	colorProfile colorprofile.Profile,
) *ResultSummarizer {
	if summarSummaryFormatFunction == nil {
		summarSummaryFormatFunction = BUILTIN_SUMMARYFORMATFUNCTION_SIMPLE
	}

	colorBuffer := new(strings.Builder)

	colorprofileWriter := new(colorprofile.Writer)
	colorprofileWriter.Forward = colorBuffer
	colorprofileWriter.Profile = colorProfile

	return &ResultSummarizer{
		summaryFormatFunction: summarSummaryFormatFunction,
		SummaryStyles:         summarizeStyles,
		colorprofileWriter:    colorprofileWriter,
		colorBuffer:           colorBuffer,
	}
}

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
