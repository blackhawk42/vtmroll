package fmtbuiltins

import (
	"strconv"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/blackhawk42/vtmroll/pkg/vtmroll"
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
	colorprofileWriter *colorprofile.Writer
	colorBuffer        *strings.Builder

	formatFunction FormatFunction

	DiceStyles DiceStyles
}

type FormatFunction func(roll int, rollType vtmroll.RollType) string

func NewDiceFormatter(
	formatFunction FormatFunction,
	diceStyles DiceStyles,
	colorProfile colorprofile.Profile,
) *DiceFormatter {
	if formatFunction == nil {
		formatFunction = func(roll int, rollType vtmroll.RollType) string { return strconv.Itoa(roll) }
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
