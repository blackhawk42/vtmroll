package vtmrollfmt

import "github.com/blackhawk42/vtmroll/pkg/vtmroll"

type VTMRollResultDiceFormatter interface {
	FormatDice(vtmres vtmroll.VTMRollerResult) []string
}

type VTMRollResultDiceFormatterFunc func(vtmres vtmroll.VTMRollerResult) []string

func (f VTMRollResultDiceFormatterFunc) FormatDice(vtmres vtmroll.VTMRollerResult) []string {
	return f(vtmres)
}

type VTMRollResultDiceParser interface {
	Parse(rolls []string, roller *vtmroll.VTMRoller, hungerDice int) (vtmroll.VTMRollerResult, error)
}

type VTMRollResultDiceParserFunc func(rolls []string, roller *vtmroll.VTMRoller, hungerDice int) (vtmroll.VTMRollerResult, error)

func (f VTMRollResultDiceParserFunc) Parse(rolls []string, roller *vtmroll.VTMRoller, hungerDice int) (vtmroll.VTMRollerResult, error) {
	return f(rolls, roller)
}

type VTMRollResultSummaryMessages struct {
	SuccessesMessage        string
	IsCriticalMessage       string
	IsTotalFailureMessage   string
	IsBestialFailureMessage string
	IsMessyCriticalMessage  string
}

type VTMRollResultSummarizer interface {
	Summarize(vtmres vtmroll.VTMRollerResult) VTMRollResultSummaryMessages
}

type VTMRollResultSummarizerFunc func(vtmres vtmroll.VTMRollerResult) VTMRollResultSummaryMessages

func (f VTMRollResultSummarizerFunc) Summarize(vtmres vtmroll.VTMRollerResult) VTMRollResultSummaryMessages {
	return f(vtmres)
}
