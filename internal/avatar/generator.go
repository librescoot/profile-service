package avatar

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"math/rand/v2"
	"os"
	"path/filepath"
)

const Size = 128

func Generate(seed int, path string) error {
	r := rand.New(rand.NewPCG(uint64(seed), uint64(seed*31)))

	scheme := schemes[r.IntN(len(schemes))]
	shape := shapeType(r.IntN(int(shapeCount)))

	img := image.NewNRGBA(image.Rect(0, 0, Size, Size))
	fillBackground(img, scheme.Background)
	drawShape(img, shape, scheme.Primary, scheme.Secondary, r)

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("create avatar directory: %w", err)
	}

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create avatar file: %w", err)
	}
	defer f.Close()

	return png.Encode(f, img)
}

func fillBackground(img *image.NRGBA, c color.NRGBA) {
	for y := range Size {
		for x := range Size {
			img.SetNRGBA(x, y, c)
		}
	}
}

func drawShape(img *image.NRGBA, shape shapeType, primary, secondary color.NRGBA, r *rand.Rand) {
	cx, cy := Size/2, Size/2
	radius := Size/2 - 12

	switch shape {
	case shapeCircle:
		fillCircle(img, cx, cy, radius, primary)
		fillCircle(img, cx, cy, radius*2/3, secondary)
	case shapeDiamond:
		fillDiamond(img, cx, cy, radius, primary)
		fillDiamond(img, cx, cy, radius*2/3, secondary)
	case shapeRoundedSquare:
		fillRoundedRect(img, cx-radius, cy-radius, cx+radius, cy+radius, 16, primary)
		inner := radius * 2 / 3
		fillRoundedRect(img, cx-inner, cy-inner, cx+inner, cy+inner, 10, secondary)
	case shapeHexagon:
		fillHexagon(img, cx, cy, radius, primary)
		fillHexagon(img, cx, cy, radius*2/3, secondary)
	}

	dotCount := 2 + r.IntN(4)
	for i := range dotCount {
		angle := float64(i) * 2 * math.Pi / float64(dotCount)
		dx := cx + int(float64(radius)*0.45*math.Cos(angle))
		dy := cy + int(float64(radius)*0.45*math.Sin(angle))
		fillCircle(img, dx, dy, 6, primary)
	}
}

func fillCircle(img *image.NRGBA, cx, cy, r int, c color.NRGBA) {
	for y := cy - r; y <= cy+r; y++ {
		for x := cx - r; x <= cx+r; x++ {
			if x >= 0 && x < Size && y >= 0 && y < Size {
				dx, dy := x-cx, y-cy
				if dx*dx+dy*dy <= r*r {
					img.SetNRGBA(x, y, c)
				}
			}
		}
	}
}

func fillDiamond(img *image.NRGBA, cx, cy, r int, c color.NRGBA) {
	for y := cy - r; y <= cy+r; y++ {
		for x := cx - r; x <= cx+r; x++ {
			if x >= 0 && x < Size && y >= 0 && y < Size {
				dx := x - cx
				dy := y - cy
				if dx < 0 {
					dx = -dx
				}
				if dy < 0 {
					dy = -dy
				}
				if dx+dy <= r {
					img.SetNRGBA(x, y, c)
				}
			}
		}
	}
}

func fillRoundedRect(img *image.NRGBA, x0, y0, x1, y1, r int, c color.NRGBA) {
	for y := y0; y <= y1; y++ {
		for x := x0; x <= x1; x++ {
			if x >= 0 && x < Size && y >= 0 && y < Size {
				inside := false
				if x >= x0+r && x <= x1-r {
					inside = true
				} else if y >= y0+r && y <= y1-r {
					inside = true
				} else {
					ccx := x0 + r
					ccy := y0 + r
					if x > x1-r {
						ccx = x1 - r
					}
					if y > y1-r {
						ccy = y1 - r
					}
					ddx, ddy := x-ccx, y-ccy
					if ddx*ddx+ddy*ddy <= r*r {
						inside = true
					}
				}
				if inside {
					img.SetNRGBA(x, y, c)
				}
			}
		}
	}
}

func fillHexagon(img *image.NRGBA, cx, cy, r int, c color.NRGBA) {
	for y := cy - r; y <= cy+r; y++ {
		for x := cx - r; x <= cx+r; x++ {
			if x >= 0 && x < Size && y >= 0 && y < Size {
				dx := math.Abs(float64(x - cx))
				dy := math.Abs(float64(y - cy))
				if dy <= float64(r)*0.866 && dx <= float64(r)-dy*0.577 {
					img.SetNRGBA(x, y, c)
				}
			}
		}
	}
}
