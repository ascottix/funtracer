// Copyright (c) 2018 Alessandro Scotti
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"math"
)

type Interpolator func(float64) float64

func ErpLinear(t float64) float64 {
	return t
}

func ErpInverse(t float64) float64 {
	return 1 - t
}

// See: https://stackoverflow.com/questions/13462001/ease-in-and-ease-out-animation-formula
func ErpBezier(t float64) float64 {
	return t * t * (3 - 2*t)
}

// See: https://stackoverflow.com/questions/13462001/ease-in-and-ease-out-animation-formula
func ErpQuadratic(t float64) float64 {
	sqt := t * t
	return sqt / (2*(sqt-t) + 1)
}

// ErpToPerlinRange brings a value from [-1,1] to [0,1]
func ErpToPerlinRange(t float64) float64 {
	return t*2 - 1
}

// ErpClip forces a value to be between 0 and 1 as follows:
// t in (-Inf, 0) = 0
// t in [0, 1]    = t
// t in (1, +Inf) = 1
func ErpClip(t float64) float64 {
	if t < 0 {
		t = 0
	}
	if t > 1 {
		t = 1
	}
	return t
}

// Gamma correction uses a more or less standard value of 2.2
const Gamma = 2.2

func ErpLinearToGamma(t float64) float64 {
	return math.Pow(t, 1.0/Gamma)
}

func ErpGammaToLinear(t float64) float64 {
	return math.Pow(t, Gamma)
}

// sRGB convers a value from linear to sRGB
func ErpLinearTosRGB(t float64) float64 {
	if t <= 0.00313066844250063 {
		return t * 12.92
	} else {
		return 1.055*math.Pow(t, 1/2.4) - 0.055
	}
}

func ErpsRGBToLinear(t float64) float64 {
	if t <= 0.0404482362771082 {
		return t / 12.92
	} else {
		return math.Pow(((t + 0.055) / 1.055), 2.4)
	}
}
