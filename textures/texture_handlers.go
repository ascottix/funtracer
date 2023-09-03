// Copyright (c) 2023 Alessandro Scotti
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package textures

// Extension point to allow custom processing
type TextureOnMapUv func(u, v float64, ii *IntersectionInfo) (float64, float64) // Triggers after u and v are fetched, but before they are used
type TextureOnImage func(data []ColorRGBA, w, h int)
type TextureOnApply func(c ColorRGBA, ii *IntersectionInfo) // Triggers before color is applied to hit, replaces standard processing

func TextureMirrorV(u, v float64, ii *IntersectionInfo) (float64, float64) {
	return u, 1 - v
}
