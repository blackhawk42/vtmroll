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
		fmtbuiltins.DEFAULT_FORMATFUNCTION_CLASSIC_DETAILED,
		fmtbuiltins.DEFAULT_DICESTYLE_ANSI,
		colorprofile.Detect(os.Stdout, os.Environ()),
		// colorprofile.NoTTY,
	)

	for _ = range 10 {
		result := roller.Roll(6, 3)

		fmt.Println(strings.Join(formatter.FormatDice(result), " "))
	}
}
