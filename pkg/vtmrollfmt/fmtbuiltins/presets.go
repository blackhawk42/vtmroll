package fmtbuiltins

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/blackhawk42/vtmroll/pkg/vtmroll"
	"github.com/blackhawk42/vtmroll/pkg/vtmrollfmt"
)

// BUILTIN_FORMATFUNCTION_NUMERIC is a DieFormatFunction that formats every die
// as its plain numeric value.
var BUILTIN_FORMATFUNCTION_NUMERIC DieFormatFunction = func(roll int, rollType vtmroll.RollType) string { return strconv.Itoa(roll) }

// BUILTIN_FORMATFUNCTION_ASCII is a DieFormatFunction that uses ASCII
// delimiters to distinguish die types (e.g. [10] for normal success, {1} for
// hunger success, *10* for critical, etc.).
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

var numericFormatRegex = regexp.MustCompile(`^(\[\d+\]|\d+|\{\d+\}|\*\d+\*|\*\{\d+\}\*|<\d+>)$`)

// BUILTIN_PARSER_NUMERIC_ASCII parses rolls produced by
// BUILTIN_FORMATFUNCTION_ASCII back into a VTMRollerResult.
var BUILTIN_PARSER_NUMERIC_ASCII vtmrollfmt.VTMRollResultDiceParser = vtmrollfmt.VTMRollResultDiceParserFunc(func(rolls []string, roller *vtmroll.VTMRoller, hungerDice int) (vtmroll.VTMRollerResult, error) {
	rollsInt := make([]int, 0, len(rolls))

	for _, roll := range rolls {
		if !numericFormatRegex.MatchString(roll) {
			return vtmroll.VTMRollerResult{}, fmt.Errorf("unrecognized roll format: %s", roll)
		}

		rollInt, err := strconv.Atoi(strings.Trim(roll, "[]{}*<>"))
		if err != nil {
			return vtmroll.VTMRollerResult{}, fmt.Errorf("while trying to convert %s to a number: %w", roll, err)
		}

		rollsInt = append(rollsInt, rollInt)
	}

	return vtmroll.NewVTMRollerResult(rollsInt, roller, hungerDice), nil
})

// BUILTIN_FORMATFUNCTION_CLASSIC_SIMPLE is a DieFormatFunction that uses
// classic VtM symbols (ankh, circle) without distinguishing hunger visually
// from normal dice.
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
		return "●"
	default:
		return fmt.Sprintf("%d", roll)
	}
}

// BUILTIN_FORMATFUNCTION_CLASSIC_DETAILED is a DieFormatFunction that uses
// classic VtM symbols with hunger dice wrapped in braces and bestial failures
// in angle brackets for added distinction.
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
		return "<●>"
	default:
		return fmt.Sprintf("%d", roll)
	}
}

// BUILTIN_PARSER_CLASSIC parses rolls produced by either
// BUILTIN_FORMATFUNCTION_CLASSIC_SIMPLE or BUILTIN_FORMATFUNCTION_CLASSIC_DETAILED.
var BUILTIN_PARSER_CLASSIC vtmrollfmt.VTMRollResultDiceParser = vtmrollfmt.VTMRollResultDiceParserFunc(func(rolls []string, roller *vtmroll.VTMRoller, hungerDice int) (vtmroll.VTMRollerResult, error) {
	rollsInt := make([]int, 0, len(rolls))

	for _, roll := range rolls {
		switch {
		case strings.Contains(roll, "٭☥٭"):
			rollsInt = append(rollsInt, roller.RollUpperLimit)
		case strings.Contains(roll, "☥"):
			rollsInt = append(rollsInt, roller.SuccessThreshold)
		case strings.Contains(roll, "●"):
			rollsInt = append(rollsInt, roller.RollLowerLimit)
		default:
			if roll != "○" {
				return vtmroll.VTMRollerResult{}, fmt.Errorf("unrecognized value: %s", roll)
			}

			rollsInt = append(rollsInt, roller.SuccessThreshold-1)
		}
	}

	return vtmroll.NewVTMRollerResult(rollsInt, roller, hungerDice), nil
})

// BUILTIN_DICESTYLES_ANSI is a DiceStyles preset with ANSI color assignments
// for terminal display.
var BUILTIN_DICESTYLES_ANSI DiceStyles

// BUILTIN_SUMMARYSTYLES_ANSI is a SummaryStyles preset with ANSI color
// assignments for terminal display.
var BUILTIN_SUMMARYSTYLES_ANSI SummaryStyles

func init() {
	BUILTIN_DICESTYLES_ANSI = NewDiceStyle()
	BUILTIN_DICESTYLES_ANSI.NormalSuccessStyle = lipgloss.NewStyle().Foreground(lipgloss.Green)
	BUILTIN_DICESTYLES_ANSI.NormalFailureStyle = lipgloss.NewStyle().Foreground(lipgloss.BrightBlack)
	BUILTIN_DICESTYLES_ANSI.HungerSuccessStyle = lipgloss.NewStyle().Foreground(lipgloss.Red)
	BUILTIN_DICESTYLES_ANSI.HungerFailureStyle = lipgloss.NewStyle().Foreground(lipgloss.Red).Faint(true)
	BUILTIN_DICESTYLES_ANSI.HalfCriticalStyle = lipgloss.NewStyle().Foreground(lipgloss.BrightGreen).Bold(true)
	BUILTIN_DICESTYLES_ANSI.HalfMessyCriticalStyle = lipgloss.NewStyle().Foreground(lipgloss.BrightYellow).Bold(true)
	BUILTIN_DICESTYLES_ANSI.PossibleBestialFailureStyle = lipgloss.NewStyle().Foreground(lipgloss.BrightRed).Bold(true)

	BUILTIN_SUMMARYSTYLES_ANSI.SuccessesMessageStyle = lipgloss.NewStyle()
	BUILTIN_SUMMARYSTYLES_ANSI.IsCriticalMessageStyle = BUILTIN_DICESTYLES_ANSI.HalfCriticalStyle
	BUILTIN_SUMMARYSTYLES_ANSI.IsTotalFailureMessageStyle = BUILTIN_DICESTYLES_ANSI.NormalFailureStyle
	BUILTIN_SUMMARYSTYLES_ANSI.IsBestialFailureMessageStyle = BUILTIN_DICESTYLES_ANSI.PossibleBestialFailureStyle
	BUILTIN_SUMMARYSTYLES_ANSI.IsMessyCriticalMessageStyle = BUILTIN_DICESTYLES_ANSI.HalfMessyCriticalStyle
}

// BUILTIN_SUMMARYFORMATFUNCTION_SIMPLE generates a one-line summary with
// successes count and any special conditions (critical, total failure, bestial
// failure, messy critical).
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
