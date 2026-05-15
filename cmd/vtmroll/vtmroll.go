package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/blackhawk42/vtmroll/pkg/vtmroll"
	"github.com/blackhawk42/vtmroll/pkg/vtmrollfmt"
	"github.com/blackhawk42/vtmroll/pkg/vtmrollfmt/fmtbuiltins"
	"github.com/charmbracelet/colorprofile"
)

func main() {
	roller := vtmroll.NewVTMRoller()

	// var formatter vtmrollfmt.VTMRollResultDiceFormatter = vtmrollfmt.VTMRollResultDiceFormatterFunc(func(vtmres vtmroll.VTMRollerResult) []string {
	// 	result := make([]string, 0, vtmres.Len())

	// 	for r, rollType := range vtmres.Rolls() {
	// 		result = append(result, fmt.Sprintf("%d (%s)", r, rollType))
	// 	}

	// 	return result
	// })

	var formatter vtmrollfmt.VTMRollResultDiceFormatter = fmtbuiltins.NewDiceFormatter(
		fmtbuiltins.BUILTIN_FORMATFUNCTION_ASCII,
		fmtbuiltins.BUILTIN_DICESTYLES_ANSI,
		colorprofile.Detect(os.Stdout, os.Environ()),
		// colorprofile.NoTTY,
	)

	var summ vtmrollfmt.VTMRollResultSummarizer = fmtbuiltins.NewResultSummarizer(
		fmtbuiltins.BUILTIN_SUMMARYSTLES_ANSI,
		fmtbuiltins.BUILTIN_SUMMARYFORMATFUNCTION_SIMPLE,
		colorprofile.NoTTY,
	)

	for _ = range 10 {
		result := roller.Roll(6, 3)

		fmt.Println(strings.Join(formatter.FormatDice(result), " "))

		summary := summ.Summarize(result)

		fmt.Println(summary.SuccessesMessage)
		if summary.IsTotalFailureMessage != "" {
			fmt.Println(summary.IsTotalFailureMessage)
		}
		if summary.IsBestialFailureMessage != "" {
			fmt.Println(summary.IsBestialFailureMessage)
		}
		if summary.IsMessyCriticalMessage != "" {
			fmt.Println(summary.IsMessyCriticalMessage)
		}
		if summary.IsCriticalMessage != "" {
			fmt.Println(summary.IsCriticalMessage)
		}

		fmt.Println()
	}
}
