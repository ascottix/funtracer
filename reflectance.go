// Copyright (c) 2018 Alessandro Scotti
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"math"
)

// SchlickReflectance returns a number between 0 and 1, representing how much
// light is reflected at the specified point (the other portion is refracted)
func SchlickReflectance(ii *IntersectionInfo) float64 {
	cosTheta := ii.Eyev.DotProduct(ii.Normalv) // θi is the angle of incidence

	if ii.N1 > ii.N2 {
		nRatio := ii.N1 / ii.N2
		sin2Thetat := nRatio * nRatio * (1 - cosTheta*cosTheta) // sin(θt)^2, where θt is the angle of refraction

		if sin2Thetat > 1 {
			return 1 // Total internal reflection
		}

		cosTheta = math.Sqrt(1 - sin2Thetat) // Replace cos(θi) with cos(θt)
	}

	r0 := (ii.N1 - ii.N2) / (ii.N1 + ii.N2)
	r0 = r0 * r0

	return r0 + (1-r0)*math.Pow((1-cosTheta), 5)
}

// See: https://www.scratchapixel.com/lessons/3d-basic-rendering/introduction-to-shading/reflection-refraction-fresnel
func FresnelReflectance(ii *IntersectionInfo) float64 {
	cosi := ii.Eyev.DotProduct(ii.Normalv) // θi is the angle of incidence
	etai := ii.N2
	etat := ii.N1

	if cosi > 0 {
		etai, etat = etat, etai
	}

	// Compute sini using Snell's law
	sint := etai / etat * math.Sqrt(1-cosi*cosi)

	if sint >= 1 {
		// Total internal reflection
		return 1
	} else {
		cost := math.Sqrt(1 - sint*sint)
		cosi = math.Abs(cosi)
		Rs := ((etat * cosi) - (etai * cost)) / ((etat * cosi) + (etai * cost))
		Rp := ((etai * cosi) - (etat * cost)) / ((etai * cosi) + (etat * cost))
		kr := (Rs*Rs + Rp*Rp) / 2

		return kr
	}
}
