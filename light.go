// Copyright (c) 2018 Alessandro Scotti
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"math"
)

type Light interface {
	LightenHit(ii *IntersectionInfo, shadowed bool) Color
	IsShadowed(rt *Raytracer, point Tuple) bool
}

type PointLight struct {
	Pos       Tuple
	Intensity Color
}

type DirectionalLight struct {
	Dir       Tuple
	Intensity Color
}

type SpotLight struct {
	Pos       Tuple
	Dir       Tuple
	AngleMin  float64
	AngleMax  float64
	Intensity Color
}

// IsShadowed returns true if there is an opaque object between the light position and the specified point
func IsShadowed(lightPos Tuple, rt *Raytracer, point Tuple) bool {
	v := lightPos.Sub(point)
	distance := v.Length()
	direction := v.Normalize()
	ray := NewRay(point, direction)
	hit := rt.HitForShadow(ray)

	return hit.Valid() && hit.T < distance
}

// OrenNayar implements the Van Ouwerkerks rewrite of Oren-Nayar model,
// see: http://shaderjvo.blogspot.com/2011/08/van-ouwerkerks-rewrite-of-oren-nayar.html
// In this model sigma represents the material roughness,
// if sigma = 0 then the material is not rough at all and the formula simplifies to the Lambert model
func OrenNayar(eyev, lightv, normalv Tuple, sigma float64) float64 {
	// Quick exit in case we're dealing with a simple Lambert shading
	if sigma == 0 {
		return 1
	}

	sigma2 := sigma * sigma
	A := 1 - (sigma2 / (2 * (sigma2 + 0.33)))
	B := 0.45 * sigma2 / (sigma2 + 0.09)

	L := math.Max(0, normalv.DotProduct(lightv)) // cosThetaI
	V := math.Max(0, normalv.DotProduct(eyev))   // cosThetaO

	if L >= (1-Epsilon) || V >= (1-Epsilon) {
		// cosPhi and sinTheta will be zero, exit now
		return A
	}

	lightPlane := lightv.Sub(normalv.Mul(L)).Normalize()
	viewPlane := eyev.Sub(normalv.Mul(V)).Normalize()
	P := math.Max(0, viewPlane.DotProduct(lightPlane)) // cosPhi

	sinTheta := math.Sqrt((1 - L*L) * (1 - V*V))
	den := math.Max(L, V)

	return A + B*P*sinTheta/den
}

// LightenHit computes the color of a point on a surface, for a specified light.
// It uses the Oren-Nayar model for diffuse and the Blinn-Phong model for specular
// (ambient contribution is computed elsewhere).
func LightenHit(lightv Tuple, lightIntensity Color, ii *IntersectionInfo) (result Color) {
	if cosTheta := lightv.DotProduct(ii.Normalv); cosTheta >= 0 { // Cosine of angle between light vector and surface normal
		// Light is on the same side of the surface, need to compute both diffuse and specular
		material := ii.O.Material()
		diffuseColor := ii.Mat.DiffuseColor.Blend(lightIntensity) // Combine light and surface colors

		// The energy of the light hitting the surface depends on the cosine of the angle
		// between the light incident direction and the surface normal (Lambert's cosine law)
		result = result.Add(diffuseColor.Mul(material.Diffuse * cosTheta * OrenNayar(ii.Eyev, lightv, ii.Normalv, material.Roughness)))

		// The Blinn-Phong model accounts for light that may be reflected directly towards the eye,
		// controlled by Specular (intensity of reflected light) and Shininess
		// (size of reflecting area, higher values yield a smaller area with harder reflection)
		halfv := lightv.Add(ii.Eyev).Normalize() // Half-vector, this would be reflectv := lightv.Neg().Reflect(normalv) in the standard Phong model

		if nDotH := ii.Normalv.DotProduct(halfv); nDotH > 0 { // Would be reflectv.DotProduct(eyev) in the standard Phong model
			f := math.Pow(nDotH, material.Shininess*4)                     // Multiply by 4 to keep "compatibility" with values tuned for the standard Phong model
			result = result.Add(lightIntensity.Mul(material.Specular * f)) // Add specular component
		}
		// ...else specular is black
	}
	// ...else light is on the other side of the surface: diffuse and specular are both black

	return result
}

// Lighten is used only for tests, it builds a dummy IntersectionInfo object then calls LightenHit to get the color
func Lighten(light Light, objectColor Color, object Hittable, point, eyev, normalv Tuple, shadowed bool) Color {
	ii := IntersectionInfo{Intersection: Intersection{O: object}, Point: point, Eyev: eyev, Normalv: normalv}
	ii.Mat.DiffuseColor = objectColor

	return light.LightenHit(&ii, shadowed)
}

func NewPointLight(pos Tuple, intensity Color) *PointLight {
	return &PointLight{pos, intensity}
}

func (light *PointLight) LightenHit(ii *IntersectionInfo, shadowed bool) (result Color) {
	if !shadowed {
		lightv := light.Pos.Sub(ii.Point).Normalize() // Direction to the light source
		result = LightenHit(lightv, light.Intensity, ii)
	}

	return
}

func (light *PointLight) IsShadowed(rt *Raytracer, point Tuple) bool {
	return IsShadowed(light.Pos, rt, point)
}

func NewDirectionalLight(dir Tuple, intensity Color) *DirectionalLight {
	return &DirectionalLight{dir.Normalize().Neg(), intensity}
}

func (light *DirectionalLight) IsShadowed(rt *Raytracer, point Tuple) bool {
	ray := NewRay(point, light.Dir)

	hit := rt.HitForShadow(ray)

	return hit.Valid()
}

func (light *DirectionalLight) LightenHit(ii *IntersectionInfo, shadowed bool) (result Color) {
	if !shadowed {
		result = LightenHit(light.Dir, light.Intensity, ii)
	}

	return
}

func NewSpotLight(pos, target Tuple, angleMin, angleMax float64, intensity Color) *SpotLight {
	return &SpotLight{pos, target.Sub(pos).Normalize(), angleMin, angleMax, intensity}
}

func (light *SpotLight) LightenHit(ii *IntersectionInfo, shadowed bool) (result Color) {
	if !shadowed {
		lightv := light.Pos.Sub(ii.Point).Normalize() // Direction to the light source

		cosSpotAngle := lightv.Neg().DotProduct(light.Dir)
		spotAngle := math.Acos(cosSpotAngle)

		a := math.Abs(spotAngle)

		if a < light.AngleMax {
			intensity := 1.0

			if a > light.AngleMin {
				// Modulate light so that it fades off gently
				t := (light.AngleMax - a) / (light.AngleMax - light.AngleMin) // Linear modulation: not good enough
				sqt := t * t                                                  // Quadratic modulation: much better! But still not perfect...
				intensity = sqt / (2*(sqt-t) + 1)
			}

			result = LightenHit(lightv, light.Intensity.Mul(intensity), ii)
		}
	}

	return result
}

func (light *SpotLight) IsShadowed(rt *Raytracer, point Tuple) bool {
	return IsShadowed(light.Pos, rt, point)
}
