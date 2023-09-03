package textures

// ColorRGBA is a RGB color with alpha, in linear space
type ColorRGBA struct {
	R, G, B, A float32
}

func (c ColorRGBA) Add(d ColorRGBA) ColorRGBA {
	return ColorRGBA{R: c.R + d.R, G: c.G + d.G, B: c.B + d.B, A: c.A + d.A}
}

func (c ColorRGBA) Mul(f float32) ColorRGBA {
	return ColorRGBA{R: c.R * f, G: c.G * f, B: c.B * f, A: c.A * f}
}

func (c ColorRGBA) RGB() Color {
	return RGB(float64(c.R), float64(c.G), float64(c.B))
}
