package models

import "testing"

func TestModeByNameKnownModes(t *testing.T) {
	for _, mode := range []EstimationMode{ModeTokens, ModeCost, ModeDays, ModeModel} {
		def, ok := ModeByName(mode)
		if !ok {
			t.Errorf("ModeByName(%q) reported not found", mode)
			continue
		}
		if def.Mode != mode {
			t.Errorf("ModeByName(%q).Mode = %q", mode, def.Mode)
		}
		if len(def.Cards) == 0 {
			t.Errorf("mode %q has no cards", mode)
		}
		if def.Name == "" {
			t.Errorf("mode %q has an empty display name", mode)
		}
	}
}

func TestModeByNameUnknown(t *testing.T) {
	if _, ok := ModeByName("BOGUS"); ok {
		t.Error("ModeByName(\"BOGUS\") should report not found")
	}
	if _, ok := ModeByName(""); ok {
		t.Error("ModeByName(\"\") should report not found")
	}
}

func TestIsValidCard(t *testing.T) {
	tokens, _ := ModeByName(ModeTokens)
	if !tokens.IsValidCard("500K") {
		t.Error("500K should be a valid TOKENS card")
	}
	if !tokens.IsValidCard("?") {
		t.Error("? (unsure) should be a valid TOKENS card")
	}
	if tokens.IsValidCard("999X") {
		t.Error("999X should not be a valid TOKENS card")
	}
	// Cross-mode: a COST card is not valid under TOKENS.
	if tokens.IsValidCard("$5") {
		t.Error("$5 should not be a valid TOKENS card")
	}

	model, _ := ModeByName(ModeModel)
	if !model.IsValidCard("Claude") {
		t.Error("Claude should be a valid MODEL card")
	}
	if model.IsValidCard("claude") {
		t.Error("card matching is case-sensitive; lowercase claude should be invalid")
	}
}

func TestEstimationModesCatalogIsComplete(t *testing.T) {
	// The four documented V1 modes must all be present, in a stable order.
	want := []EstimationMode{ModeTokens, ModeCost, ModeDays, ModeModel}
	if len(EstimationModes) != len(want) {
		t.Fatalf("catalog has %d modes, want %d", len(EstimationModes), len(want))
	}
	for i, m := range want {
		if EstimationModes[i].Mode != m {
			t.Errorf("EstimationModes[%d].Mode = %q, want %q", i, EstimationModes[i].Mode, m)
		}
	}
}

func TestNumericFlagMatchesModeSemantics(t *testing.T) {
	// Token/cost/day modes are numeric; model selection is not.
	numeric := map[EstimationMode]bool{
		ModeTokens: true, ModeCost: true, ModeDays: true, ModeModel: false,
	}
	for _, def := range EstimationModes {
		if def.Numeric != numeric[def.Mode] {
			t.Errorf("mode %q Numeric = %v, want %v", def.Mode, def.Numeric, numeric[def.Mode])
		}
	}
}
