package names

import (
	"strings"
	"testing"
)

func TestGenerate_TwoWords(t *testing.T) {
	name := Generate()
	parts := strings.SplitN(name, " ", 2)
	if len(parts) != 2 {
		t.Errorf("expected 'Adjective Animal', got %q", name)
	}
}

func TestGenerate_NonEmpty(t *testing.T) {
	name := Generate()
	if name == "" {
		t.Error("expected non-empty name")
	}
}

func TestGenerate_Unique(t *testing.T) {
	seen := make(map[string]bool)
	dupes := 0
	for i := 0; i < 100; i++ {
		name := Generate()
		if seen[name] {
			dupes++
		}
		seen[name] = true
	}
	if dupes > 10 {
		t.Errorf("too many duplicates in 100 names: %d", dupes)
	}
}

func TestGenerateWithSeed_Deterministic(t *testing.T) {
	name1 := GenerateWithSeed(42)
	name2 := GenerateWithSeed(42)
	if name1 != name2 {
		t.Errorf("same seed produced different names: %q vs %q", name1, name2)
	}
}

func TestGenerateWithSeed_DifferentSeeds(t *testing.T) {
	name1 := GenerateWithSeed(1)
	name2 := GenerateWithSeed(2)
	if name1 == name2 {
		t.Errorf("different seeds produced same name: %q", name1)
	}
}
