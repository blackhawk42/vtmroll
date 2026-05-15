package fmtbuiltins

import (
	"fmt"

	"charm.land/lipgloss/v2"
	"github.com/blackhawk42/vtmroll/pkg/vtmroll"
)

var DEFAULT_FORMATFUNCTION_ASCII FormatFunction = func(roll int, rollType vtmroll.RollType) string {
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

var DEFAULT_FORMATFUNCTION_CLASSIC_SIMPLE FormatFunction = func(roll int, rollType vtmroll.RollType) string {
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

var DEFAULT_FORMATFUNCTION_CLASSIC_DETAILED FormatFunction = func(roll int, rollType vtmroll.RollType) string {
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

var DEFAULT_DICESTYLE_ANSI DiceStyles

func init() {
	DEFAULT_DICESTYLE_ANSI = NewDiceStyle()
	DEFAULT_DICESTYLE_ANSI.NormalSuccessStyle = lipgloss.NewStyle().Foreground(lipgloss.Green)
	DEFAULT_DICESTYLE_ANSI.NormalFailureStyle = lipgloss.NewStyle().Foreground(lipgloss.BrightBlack)
	DEFAULT_DICESTYLE_ANSI.HungerSuccessStyle = lipgloss.NewStyle().Foreground(lipgloss.Red)
	DEFAULT_DICESTYLE_ANSI.HungerFailureStyle = lipgloss.NewStyle().Foreground(lipgloss.Red).Faint(true)
	DEFAULT_DICESTYLE_ANSI.HalfCriticalStyle = lipgloss.NewStyle().Foreground(lipgloss.BrightGreen).Bold(true)
	DEFAULT_DICESTYLE_ANSI.HalfMessyCriticalStyle = lipgloss.NewStyle().Foreground(lipgloss.BrightYellow).Bold(true)
	DEFAULT_DICESTYLE_ANSI.PossibleBestialFailureStyle = lipgloss.NewStyle().Foreground(lipgloss.BrightRed).Bold(true)
}
