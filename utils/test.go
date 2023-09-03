package utils

import (
	"testing"

	. "ascottix/funtracer/shapes"
	. "ascottix/funtracer/textures"
)

const (
	SkipAllTestsWithImages = false
)

func TestWithImage(t *testing.T) {
	if SkipAllTestsWithImages {
		t.SkipNow()
	}
}

func ApplyTexture(s *Shape, filename string) *ImageTexture {
	txt := NewImageTexture()
	err := txt.LoadFromFile(filename)
	if err == nil {
		s.Material().SetPattern(txt).SetSpecular(0)
	} else {
		Debugln("Cannot load texture: ", filename)
	}

	return txt
}
