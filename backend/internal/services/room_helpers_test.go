package services

import (
	"strings"
	"testing"

	"freetokenspoker/internal/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestNormalizeCode(t *testing.T) {
	cases := map[string]string{
		"abcdef":    "ABCDEF",
		"  abc123 ": "ABC123",
		"ABCDEF":    "ABCDEF",
		"\tmix4d\n": "MIX4D",
		"":          "",
	}
	for in, want := range cases {
		if got := normalizeCode(in); got != want {
			t.Errorf("normalizeCode(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestIsMember(t *testing.T) {
	owner := primitive.NewObjectID()
	other := primitive.NewObjectID()
	stranger := primitive.NewObjectID()

	room := &models.Room{
		OwnerID: owner,
		Members: []models.Member{
			{UserID: owner, Name: "Owner"},
			{UserID: other, Name: "Other"},
		},
	}

	if !isMember(room, owner.Hex()) {
		t.Error("owner should be a member")
	}
	if !isMember(room, other.Hex()) {
		t.Error("added member should be a member")
	}
	if isMember(room, stranger.Hex()) {
		t.Error("a stranger should not be a member")
	}
	if isMember(room, "not-an-object-id") {
		t.Error("a garbage id should not match any member")
	}
}

func TestGenerateCodeShape(t *testing.T) {
	// codeAlphabet deliberately omits ambiguous characters so shared codes are
	// easy to read and type.
	const ambiguous = "01OIL"
	seen := map[string]bool{}

	for i := 0; i < 500; i++ {
		code := generateCode()
		if len(code) != codeLength {
			t.Fatalf("code %q has length %d, want %d", code, len(code), codeLength)
		}
		for _, r := range code {
			if !strings.ContainsRune(codeAlphabet, r) {
				t.Fatalf("code %q contains char %q outside the alphabet", code, r)
			}
			if strings.ContainsRune(ambiguous, r) {
				t.Fatalf("code %q contains ambiguous char %q", code, r)
			}
		}
		seen[code] = true
	}

	// 500 draws from a 31^6 space should almost never collide; a tiny variety
	// bar catches an accidentally-constant generator.
	if len(seen) < 400 {
		t.Errorf("generateCode produced only %d distinct codes in 500 draws; looks non-random", len(seen))
	}
}
