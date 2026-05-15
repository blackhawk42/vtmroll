package fmtbuiltins

import (
	"fmt"
	"strconv"

	"charm.land/lipgloss/v2"
	"github.com/blackhawk42/vtmroll/pkg/vtmroll"
	"github.com/blackhawk42/vtmroll/pkg/vtmrollfmt"
)

var BUILTIN_FORMATFUNCTION_NUMERIC DieFormatFunction = func(roll int, rollType vtmroll.RollType) string { return strconv.Itoa(roll) }

var BUILTIN_FORMATFUNCTION_ASCII DieFormatFunction = func(roll int, rollType vtmroll.RollType) string {
	switch rollType {
	case vtmroll.NormalSuccess:
		return fmt.Sprintf("[%d]", roll)
	case vtmroll.NormalFailure:
		return fmt.Sprintf("%d", roll)
	case vtmroll.HungerSuccess:
		return fmt.Sprintf("{%d}", roll)
	case vtmroll.HungerFailure:
		return fmt.Sprintf("%d", roll)
	case vtmroll.HalfCritical:
		return fmt.Sprintf("*%d*", roll)
	case vtmroll.HalfMessyCritical:
		return fmt.Sprintf("*{%d}*", roll)
	case vtmroll.PossibleBestialFailure:
		return fmt.Sprintf("<%d>", roll)
	default:
		return fmt.Sprintf("%d", roll)
	}
}

var BUILTIN_FORMATFUNCTION_CLASSIC_SIMPLE DieFormatFunction = func(roll int, rollType vtmroll.RollType) string {
	switch rollType {
	case vtmroll.NormalSuccess:
		return "☥"
	case vtmroll.NormalFailure:
		return "○"
	case vtmroll.HungerSuccess:
		return "☥"
	case vtmroll.HungerFailure:
		return "○"
	case vtmroll.HalfCritical:
		return "٭☥٭"
	case vtmroll.HalfMessyCritical:
		return "٭☥٭"
	case vtmroll.PossibleBestialFailure:
		return "○"
	default:
		return fmt.Sprintf("%d", roll)
	}
}

var BUILTIN_FORMATFUNCTION_CLASSIC_DETAILED DieFormatFunction = func(roll int, rollType vtmroll.RollType) string {
	switch rollType {
	case vtmroll.NormalSuccess:
		return "☥"
	case vtmroll.NormalFailure:
		return "○"
	case vtmroll.HungerSuccess:
		return "{☥}"
	case vtmroll.HungerFailure:
		return "{○}"
	case vtmroll.HalfCritical:
		return "٭☥٭"
	case vtmroll.HalfMessyCritical:
		return "{٭☥٭}"
	case vtmroll.PossibleBestialFailure:
		return "<○>"
	default:
		return fmt.Sprintf("%d", roll)
	}
}

var BUILTIN_DICESTYLES_ANSI DiceStyles

var BUILTIN_SUMMARYSTLES_ANSI SummaryStyles

func init() {
	BUILTIN_DICESTYLES_ANSI = NewDiceStyle()
	BUILTIN_DICESTYLES_ANSI.NormalSuccessStyle = lipgloss.NewStyle().Foreground(lipgloss.Green)
	BUILTIN_DICESTYLES_ANSI.NormalFailureStyle = lipgloss.NewStyle().Foreground(lipgloss.BrightBlack)
	BUILTIN_DICESTYLES_ANSI.HungerSuccessStyle = lipgloss.NewStyle().Foreground(lipgloss.Red)
	BUILTIN_DICESTYLES_ANSI.HungerFailureStyle = lipgloss.NewStyle().Foreground(lipgloss.Red).Faint(true)
	BUILTIN_DICESTYLES_ANSI.HalfCriticalStyle = lipgloss.NewStyle().Foreground(lipgloss.BrightGreen).Bold(true)
	BUILTIN_DICESTYLES_ANSI.HalfMessyCriticalStyle = lipgloss.NewStyle().Foreground(lipgloss.BrightYellow).Bold(true)
	BUILTIN_DICESTYLES_ANSI.PossibleBestialFailureStyle = lipgloss.NewStyle().Foreground(lipgloss.BrightRed).Bold(true)

	BUILTIN_SUMMARYSTLES_ANSI.SuccessesMessageStyle = lipgloss.NewStyle()
	BUILTIN_SUMMARYSTLES_ANSI.IsCriticalMessageStyle = BUILTIN_DICESTYLES_ANSI.HalfCriticalStyle
	BUILTIN_SUMMARYSTLES_ANSI.IsTotalFailureMessageStyle = BUILTIN_DICESTYLES_ANSI.NormalFailureStyle
	BUILTIN_SUMMARYSTLES_ANSI.IsBestialFailureMessageStyle = BUILTIN_DICESTYLES_ANSI.PossibleBestialFailureStyle
	BUILTIN_SUMMARYSTLES_ANSI.IsMessyCriticalMessageStyle = BUILTIN_DICESTYLES_ANSI.HalfMessyCriticalStyle
}

var BUILTIN_SUMMARYFORMATFUNCTION_SIMPLE = func(result vtmroll.VTMRollerResult) vtmrollfmt.VTMRollResultSummaryMessages {
	summary := vtmrollfmt.VTMRollResultSummaryMessages{}

	summary.SuccessesMessage = fmt.Sprintf("Successes: %d", result.Successes())
	if result.IsCritical() {
		summary.IsCriticalMessage = "Critical!"
	}
	if result.IsTotalFailure() {
		summary.IsTotalFailureMessage = "Total failure!"
	}
	if result.IsBestialFailure() {
		summary.IsBestialFailureMessage = "Bestial failure!"
	}
	if result.IsMessyCritical() {
		summary.IsMessyCriticalMessage = "Messy critical!"
	}

	return summary
}
