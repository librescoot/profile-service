package avatar

import (
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

func TestGenerate_CreatesPNG(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "avatar.png")

	err := Generate(42, path)
	if err != nil {
		t.Fatalf("Generate error: %v", err)
	}

	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("open error: %v", err)
	}
	defer f.Close()

	img, err := png.Decode(f)
	if err != nil {
		t.Fatalf("png.Decode error: %v", err)
	}

	bounds := img.Bounds()
	if bounds.Dx() != Size || bounds.Dy() != Size {
		t.Errorf("size = %dx%d, want %dx%d", bounds.Dx(), bounds.Dy(), Size, Size)
	}
}

func TestGenerate_DeterministicForSameSeed(t *testing.T) {
	dir := t.TempDir()
	path1 := filepath.Join(dir, "a1.png")
	path2 := filepath.Join(dir, "a2.png")

	Generate(42, path1)
	Generate(42, path2)

	data1, _ := os.ReadFile(path1)
	data2, _ := os.ReadFile(path2)

	if len(data1) != len(data2) {
		t.Error("same seed produced different file sizes")
	}
}

func TestGenerate_DifferentSeedsDifferentOutput(t *testing.T) {
	dir := t.TempDir()
	path1 := filepath.Join(dir, "a1.png")
	path2 := filepath.Join(dir, "a2.png")

	Generate(1, path1)
	Generate(2, path2)

	data1, _ := os.ReadFile(path1)
	data2, _ := os.ReadFile(path2)

	if string(data1) == string(data2) {
		t.Error("different seeds produced identical output")
	}
}
