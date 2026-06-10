package fmtbuiltins

import (
	"slices"
	"testing"

	"github.com/blackhawk42/vtmroll/pkg/vtmroll"
	"github.com/blackhawk42/vtmroll/pkg/vtmrollfmt"
	"github.com/charmbracelet/colorprofile"
)

func TestNewDiceStyle(t *testing.T) {
	ds := NewDiceStyle()
	if ds.NormalSuccessStyle.Render("x") != "x" {
		t.Error("NormalSuccessStyle not initialized properly")
	}
	if ds.NormalFailureStyle.Render("x") != "x" {
		t.Error("NormalFailureStyle not initialized properly")
	}
	if ds.HungerSuccessStyle.Render("x") != "x" {
		t.Error("HungerSuccessStyle not initialized properly")
	}
	if ds.HungerFailureStyle.Render("x") != "x" {
		t.Error("HungerFailureStyle not initialized properly")
	}
	if ds.HalfCriticalStyle.Render("x") != "x" {
		t.Error("HalfCriticalStyle not initialized properly")
	}
	if ds.HalfMessyCriticalStyle.Render("x") != "x" {
		t.Error("HalfMessyCriticalStyle not initialized properly")
	}
	if ds.PossibleBestialFailureStyle.Render("x") != "x" {
		t.Error("PossibleBestialFailureStyle not initialized properly")
	}
}

func TestNewSummaryStyles(t *testing.T) {
	ss := NewSummaryStyles()
	if ss.SuccessesMessageStyle.Render("x") != "x" {
		t.Error("SuccessesMessageStyle not initialized properly")
	}
	if ss.IsCriticalMessageStyle.Render("x") != "x" {
		t.Error("IsCriticalMessageStyle not initialized properly")
	}
	if ss.IsTotalFailureMessageStyle.Render("x") != "x" {
		t.Error("IsTotalFailureMessageStyle not initialized properly")
	}
	if ss.IsBestialFailureMessageStyle.Render("x") != "x" {
		t.Error("IsBestialFailureMessageStyle not initialized properly")
	}
	if ss.IsMessyCriticalMessageStyle.Render("x") != "x" {
		t.Error("IsMessyCriticalMessageStyle not initialized properly")
	}
}

func TestBuiltinFormatFunctionNumeric(t *testing.T) {
	tests := []struct {
		name     string
		roll     int
		rollType vtmroll.RollType
		want     string
	}{
		{"NormalSuccess", 7, vtmroll.NormalSuccess, "7"},
		{"NormalFailure", 3, vtmroll.NormalFailure, "3"},
		{"HungerSuccess", 8, vtmroll.HungerSuccess, "8"},
		{"HungerFailure", 2, vtmroll.HungerFailure, "2"},
		{"HalfCritical", 10, vtmroll.HalfCritical, "10"},
		{"HalfMessyCritical", 10, vtmroll.HalfMessyCritical, "10"},
		{"PossibleBestialFailure", 1, vtmroll.PossibleBestialFailure, "1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BUILTIN_FORMATFUNCTION_NUMERIC(tt.roll, tt.rollType)
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestBuiltinFormatFunctionASCII(t *testing.T) {
	tests := []struct {
		name     string
		roll     int
		rollType vtmroll.RollType
		want     string
	}{
		{"NormalSuccess", 7, vtmroll.NormalSuccess, "[7]"},
		{"NormalFailure", 3, vtmroll.NormalFailure, "3"},
		{"HungerSuccess", 8, vtmroll.HungerSuccess, "{8}"},
		{"HungerFailure", 2, vtmroll.HungerFailure, "2"},
		{"HalfCritical", 10, vtmroll.HalfCritical, "*10*"},
		{"HalfMessyCritical", 10, vtmroll.HalfMessyCritical, "*{10}*"},
		{"PossibleBestialFailure", 1, vtmroll.PossibleBestialFailure, "<1>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BUILTIN_FORMATFUNCTION_ASCII(tt.roll, tt.rollType)
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestBuiltinFormatFunctionClassicSimple(t *testing.T) {
	tests := []struct {
		name     string
		roll     int
		rollType vtmroll.RollType
		want     string
	}{
		{"NormalSuccess", 7, vtmroll.NormalSuccess, "☥"},
		{"NormalFailure", 3, vtmroll.NormalFailure, "○"},
		{"HungerSuccess", 8, vtmroll.HungerSuccess, "☥"},
		{"HungerFailure", 2, vtmroll.HungerFailure, "○"},
		{"HalfCritical", 10, vtmroll.HalfCritical, "٭☥٭"},
		{"HalfMessyCritical", 10, vtmroll.HalfMessyCritical, "٭☥٭"},
		{"PossibleBestialFailure", 1, vtmroll.PossibleBestialFailure, "●"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BUILTIN_FORMATFUNCTION_CLASSIC_SIMPLE(tt.roll, tt.rollType)
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestBuiltinFormatFunctionClassicDetailed(t *testing.T) {
	tests := []struct {
		name     string
		roll     int
		rollType vtmroll.RollType
		want     string
	}{
		{"NormalSuccess", 7, vtmroll.NormalSuccess, "☥"},
		{"NormalFailure", 3, vtmroll.NormalFailure, "○"},
		{"HungerSuccess", 8, vtmroll.HungerSuccess, "{☥}"},
		{"HungerFailure", 2, vtmroll.HungerFailure, "{○}"},
		{"HalfCritical", 10, vtmroll.HalfCritical, "٭☥٭"},
		{"HalfMessyCritical", 10, vtmroll.HalfMessyCritical, "{٭☥٭}"},
		{"PossibleBestialFailure", 1, vtmroll.PossibleBestialFailure, "<●>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BUILTIN_FORMATFUNCTION_CLASSIC_DETAILED(tt.roll, tt.rollType)
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestBuiltinSummaryFormatFunctionSimple(t *testing.T) {
	tests := []struct {
		name   string
		rolls  []int
		hunger int
		want   vtmrollfmt.VTMRollResultSummaryMessages
	}{
		{
			name:   "plain success",
			rolls:  []int{6, 7, 3, 4},
			hunger: 0,
			want: vtmrollfmt.VTMRollResultSummaryMessages{
				SuccessesMessage: "Successes: 2",
			},
		},
		{
			name:   "critical",
			rolls:  []int{10, 10, 3, 4},
			hunger: 0,
			want: vtmrollfmt.VTMRollResultSummaryMessages{
				SuccessesMessage:  "Successes: 4",
				IsCriticalMessage: "Critical!",
			},
		},
		{
			name:   "total failure",
			rolls:  []int{1, 2, 3, 4, 5},
			hunger: 0,
			want: vtmrollfmt.VTMRollResultSummaryMessages{
				SuccessesMessage:      "Successes: 0",
				IsTotalFailureMessage: "Total failure!",
			},
		},
		{
			name:   "bestial failure",
			rolls:  []int{1, 2, 3, 4, 5},
			hunger: 2,
			want: vtmrollfmt.VTMRollResultSummaryMessages{
				SuccessesMessage:        "Successes: 0",
				IsTotalFailureMessage:   "Total failure!",
				IsBestialFailureMessage: "Bestial failure!",
			},
		},
		{
			name:   "messy critical",
			rolls:  []int{10, 10, 6, 7},
			hunger: 2,
			want: vtmrollfmt.VTMRollResultSummaryMessages{
				SuccessesMessage:       "Successes: 6",
				IsCriticalMessage:      "Critical!",
				IsMessyCriticalMessage: "Messy critical!",
			},
		},
		{
			name:   "critical without messy (hunger present, crits on normal dice)",
			rolls:  []int{6, 7, 10, 10},
			hunger: 2,
			want: vtmrollfmt.VTMRollResultSummaryMessages{
				SuccessesMessage:  "Successes: 6",
				IsCriticalMessage: "Critical!",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			roller := vtmroll.NewVTMRoller()
			result := vtmroll.NewVTMRollerResult(tt.rolls, roller, tt.hunger)
			got := BUILTIN_SUMMARYFORMATFUNCTION_SIMPLE(result)
			if got != tt.want {
				t.Errorf("got %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestBuiltinParserNumericASCII(t *testing.T) {
	roller := vtmroll.NewVTMRoller()

	t.Run("round trip", func(t *testing.T) {
		result := vtmroll.NewVTMRollerResult([]int{10, 1, 8, 3, 6}, roller, 2)
		formatted := make([]string, 0, result.Len())
		for roll, rt := range result.Rolls() {
			formatted = append(formatted, BUILTIN_FORMATFUNCTION_ASCII(roll, rt))
		}
		parsed, err := BUILTIN_PARSER_NUMERIC_ASCII.Parse(formatted, roller, 2)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !slices.Equal(parsed.GetRolls(), result.GetRolls()) {
			t.Errorf("rolls: got %v, want %v", parsed.GetRolls(), result.GetRolls())
		}
	})

	t.Run("unrecognized format", func(t *testing.T) {
		_, err := BUILTIN_PARSER_NUMERIC_ASCII.Parse([]string{"xyz"}, roller, 0)
		if err == nil {
			t.Error("expected error for unrecognized format")
		}
	})

	t.Run("empty rolls", func(t *testing.T) {
		parsed, err := BUILTIN_PARSER_NUMERIC_ASCII.Parse([]string{}, roller, 0)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if parsed.Len() != 0 {
			t.Errorf("expected 0 rolls, got %d", parsed.Len())
		}
	})
}

func TestBuiltinParserClassic(t *testing.T) {
	roller := vtmroll.NewVTMRoller()

	t.Run("round trip with simple format", func(t *testing.T) {
		result := vtmroll.NewVTMRollerResult([]int{10, 1, 8, 3, 6}, roller, 2)
		formatted := make([]string, 0, result.Len())
		for roll, rt := range result.Rolls() {
			formatted = append(formatted, BUILTIN_FORMATFUNCTION_CLASSIC_SIMPLE(roll, rt))
		}
		parsed, err := BUILTIN_PARSER_CLASSIC.Parse(formatted, roller, 2)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if parsed.Successes() != result.Successes() {
			t.Errorf("Successes: got %d, want %d", parsed.Successes(), result.Successes())
		}
		if parsed.IsCritical() != result.IsCritical() {
			t.Errorf("IsCritical: got %v, want %v", parsed.IsCritical(), result.IsCritical())
		}
		if parsed.IsTotalFailure() != result.IsTotalFailure() {
			t.Errorf("IsTotalFailure: got %v, want %v", parsed.IsTotalFailure(), result.IsTotalFailure())
		}
		if parsed.IsBestialFailure() != result.IsBestialFailure() {
			t.Errorf("IsBestialFailure: got %v, want %v", parsed.IsBestialFailure(), result.IsBestialFailure())
		}
		if parsed.IsMessyCritical() != result.IsMessyCritical() {
			t.Errorf("IsMessyCritical: got %v, want %v", parsed.IsMessyCritical(), result.IsMessyCritical())
		}
		if parsed.HungerDice() != result.HungerDice() {
			t.Errorf("HungerDice: got %d, want %d", parsed.HungerDice(), result.HungerDice())
		}
	})

	t.Run("round trip with detailed format", func(t *testing.T) {
		result := vtmroll.NewVTMRollerResult([]int{10, 1, 8, 3, 6}, roller, 2)
		formatted := make([]string, 0, result.Len())
		for roll, rt := range result.Rolls() {
			formatted = append(formatted, BUILTIN_FORMATFUNCTION_CLASSIC_DETAILED(roll, rt))
		}
		parsed, err := BUILTIN_PARSER_CLASSIC.Parse(formatted, roller, 2)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if parsed.Successes() != result.Successes() {
			t.Errorf("Successes: got %d, want %d", parsed.Successes(), result.Successes())
		}
		if parsed.IsCritical() != result.IsCritical() {
			t.Errorf("IsCritical: got %v, want %v", parsed.IsCritical(), result.IsCritical())
		}
		if parsed.IsMessyCritical() != result.IsMessyCritical() {
			t.Errorf("IsMessyCritical: got %v, want %v", parsed.IsMessyCritical(), result.IsMessyCritical())
		}
	})

	t.Run("unrecognized symbol", func(t *testing.T) {
		_, err := BUILTIN_PARSER_CLASSIC.Parse([]string{"?"}, roller, 0)
		if err == nil {
			t.Error("expected error for unrecognized symbol")
		}
	})

	t.Run("empty rolls", func(t *testing.T) {
		parsed, err := BUILTIN_PARSER_CLASSIC.Parse([]string{}, roller, 0)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if parsed.Len() != 0 {
			t.Errorf("expected 0 rolls, got %d", parsed.Len())
		}
	})
}

func TestDiceFormatterDefaultFormatFunction(t *testing.T) {
	df := NewDiceFormatter(nil, NewDiceStyle(), colorprofile.ASCII)
	roller := vtmroll.NewVTMRoller()
	result := vtmroll.NewVTMRollerResult([]int{6, 10, 3}, roller, 1)
	got := df.FormatDice(result)
	want := []string{"6", "10", "3"}
	if !slices.Equal(got, want) {
		t.Errorf("FormatDice: got %v, want %v", got, want)
	}
}

func TestDiceFormatterFormatDice(t *testing.T) {
	customFn := func(roll int, rt vtmroll.RollType) string {
		if rt == vtmroll.HalfCritical || rt == vtmroll.HalfMessyCritical {
			return "CRIT"
		}
		return "DIE"
	}
	df := NewDiceFormatter(customFn, NewDiceStyle(), colorprofile.ASCII)
	roller := vtmroll.NewVTMRoller()
	result := vtmroll.NewVTMRollerResult([]int{10, 1, 8, 3, 6}, roller, 2)
	got := df.FormatDice(result)
	want := []string{"CRIT", "DIE", "DIE", "DIE", "DIE"}
	if !slices.Equal(got, want) {
		t.Errorf("FormatDice: got %v, want %v", got, want)
	}
}

func TestResultSummarizerDefaultSummaryFunction(t *testing.T) {
	rs := NewResultSummarizer(NewSummaryStyles(), nil, colorprofile.ASCII)
	roller := vtmroll.NewVTMRoller()
	result := vtmroll.NewVTMRollerResult([]int{6, 7, 8, 3, 4}, roller, 0)
	got := rs.Summarize(result)
	want := vtmrollfmt.VTMRollResultSummaryMessages{
		SuccessesMessage: "Successes: 3",
	}
	if got != want {
		t.Errorf("Summarize: got %+v, want %+v", got, want)
	}
}

func TestResultSummarizerSummarize(t *testing.T) {
	customFn := func(result vtmroll.VTMRollerResult) vtmrollfmt.VTMRollResultSummaryMessages {
		return vtmrollfmt.VTMRollResultSummaryMessages{
			SuccessesMessage:        "custom success",
			IsCriticalMessage:       "custom critical",
			IsTotalFailureMessage:   "custom fail",
			IsBestialFailureMessage: "custom beast",
			IsMessyCriticalMessage:  "custom messy",
		}
	}
	rs := NewResultSummarizer(NewSummaryStyles(), customFn, colorprofile.ASCII)
	roller := vtmroll.NewVTMRoller()
	result := vtmroll.NewVTMRollerResult([]int{6, 7, 8, 3, 4}, roller, 0)
	got := rs.Summarize(result)
	// With zero styles + ASCII profile, messages pass through unchanged
	want := vtmrollfmt.VTMRollResultSummaryMessages{
		SuccessesMessage:        "custom success",
		IsCriticalMessage:       "custom critical",
		IsTotalFailureMessage:   "custom fail",
		IsBestialFailureMessage: "custom beast",
		IsMessyCriticalMessage:  "custom messy",
	}
	if got != want {
		t.Errorf("Summarize: got %+v, want %+v", got, want)
	}
}

func TestBuiltinDiceStylesANSINonNil(t *testing.T) {
	if BUILTIN_DICESTYLES_ANSI.NormalSuccessStyle.Render("x") == "" {
		t.Error("NormalSuccessStyle should be initialized")
	}
	if BUILTIN_DICESTYLES_ANSI.NormalFailureStyle.Render("x") == "" {
		t.Error("NormalFailureStyle should be initialized")
	}
	if BUILTIN_DICESTYLES_ANSI.HungerSuccessStyle.Render("x") == "" {
		t.Error("HungerSuccessStyle should be initialized")
	}
	if BUILTIN_DICESTYLES_ANSI.HungerFailureStyle.Render("x") == "" {
		t.Error("HungerFailureStyle should be initialized")
	}
	if BUILTIN_DICESTYLES_ANSI.HalfCriticalStyle.Render("x") == "" {
		t.Error("HalfCriticalStyle should be initialized")
	}
	if BUILTIN_DICESTYLES_ANSI.HalfMessyCriticalStyle.Render("x") == "" {
		t.Error("HalfMessyCriticalStyle should be initialized")
	}
	if BUILTIN_DICESTYLES_ANSI.PossibleBestialFailureStyle.Render("x") == "" {
		t.Error("PossibleBestialFailureStyle should be initialized")
	}
}

func TestBuiltinSummaryStylesANSINonNil(t *testing.T) {
	if BUILTIN_SUMMARYSTYLES_ANSI.SuccessesMessageStyle.Render("x") == "" {
		t.Error("SuccessesMessageStyle should be initialized")
	}
	if BUILTIN_SUMMARYSTYLES_ANSI.IsCriticalMessageStyle.Render("x") == "" {
		t.Error("IsCriticalMessageStyle should be initialized")
	}
	if BUILTIN_SUMMARYSTYLES_ANSI.IsTotalFailureMessageStyle.Render("x") == "" {
		t.Error("IsTotalFailureMessageStyle should be initialized")
	}
	if BUILTIN_SUMMARYSTYLES_ANSI.IsBestialFailureMessageStyle.Render("x") == "" {
		t.Error("IsBestialFailureMessageStyle should be initialized")
	}
	if BUILTIN_SUMMARYSTYLES_ANSI.IsMessyCriticalMessageStyle.Render("x") == "" {
		t.Error("IsMessyCriticalMessageStyle should be initialized")
	}
}
