package models

// EstimationMode enumerates the supported estimation dimensions. Card values
// are data, never hardcoded into UI logic, so new modes need no code changes
// beyond extending this catalog.
type EstimationMode string

const (
	ModeTokens EstimationMode = "TOKENS"
	ModeCost   EstimationMode = "COST"
	ModeDays   EstimationMode = "DAYS"
	ModeModel  EstimationMode = "MODEL"
)

// ModeDefinition describes a mode and the ordered set of cards it offers.
type ModeDefinition struct {
	Mode        EstimationMode `json:"mode"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Numeric     bool           `json:"numeric"`
	Cards       []string       `json:"cards"`
}

// EstimationModes is the canonical catalog served to clients via GET /api/modes.
var EstimationModes = []ModeDefinition{
	{
		Mode:        ModeTokens,
		Name:        "AI Tokens",
		Description: "Estimate how many tokens the work will consume.",
		Numeric:     true,
		Cards:       []string{"500K", "1M", "2M", "5M", "10M", "20M", "50M", "?"},
	},
	{
		Mode:        ModeCost,
		Name:        "AI Cost",
		Description: "Estimate the approximate AI spend for the work.",
		Numeric:     true,
		Cards:       []string{"$1", "$5", "$10", "$25", "$50", "$100", "$250", "?"},
	},
	{
		Mode:        ModeDays,
		Name:        "Engineering Days",
		Description: "Estimate AI-assisted engineering days to ship.",
		Numeric:     true,
		Cards:       []string{"1", "2", "3", "5", "8", "13", "21", "?"},
	},
	{
		Mode:        ModeModel,
		Name:        "Best AI Model",
		Description: "Pick the model best suited for the work.",
		Numeric:     false,
		Cards:       []string{"GPT", "Claude", "Gemini", "DeepSeek", "Cursor", "Codex", "Other"},
	},
}

// ModeByName returns the definition for a mode, and whether it exists.
func ModeByName(mode EstimationMode) (ModeDefinition, bool) {
	for _, d := range EstimationModes {
		if d.Mode == mode {
			return d, true
		}
	}
	return ModeDefinition{}, false
}

// IsValidCard reports whether a card belongs to a mode's catalog.
func (d ModeDefinition) IsValidCard(card string) bool {
	for _, c := range d.Cards {
		if c == card {
			return true
		}
	}
	return false
}
