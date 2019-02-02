// Copyright (c) 2018 Alessandro Scotti
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"math"
)

type Light interface {
	LightenHit(ii *IntersectionInfo, rt *Raytracer) Color
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

type RectLight struct {
	Pos       Tuple
	Uv        Tuple
	Vv        Tuple
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

		// The energy of the light hitting the surface depends on the cosine of the angle
		// between the light direction and the surface normal (Lambert's cosine law) i.e. cosTheta
		result = ii.Mat.DiffuseColor.Mul(material.Diffuse * cosTheta * OrenNayar(ii.Eyev, lightv, ii.Normalv, material.Roughness))

		// The Blinn-Phong model accounts for light that may be reflected directly towards the eye,
		// controlled by Specular (intensity of reflected light) and Shininess
		// (size of reflecting area, higher values yield a smaller area with harder reflection)
		halfv := lightv.Add(ii.Eyev).Normalize() // Half-vector, this would be reflectv := lightv.Neg().Reflect(ii.Normalv) in the standard Phong model

		if nDotH := ii.Normalv.DotProduct(halfv); nDotH > Epsilon { // Would be reflectv.DotProduct(ii.Eyev) in the standard Phong model
			f := math.Pow(nDotH, 4*material.Shininess) // Multiply shininess by 4 to keep "compatibility" with values tuned for the standard Phong model

			result = result.Add(material.ReflectColor.Mul(material.Specular * f)) // Add specular component
		}

		// Phong
		// reflectv := lightv.Neg().Reflect(ii.Normalv)

		// if nDotH := reflectv.DotProduct(ii.Eyev); nDotH > 0 { // Would be reflectv.DotProduct(eyev) in the standard Phong model
		// 	f := math.Pow(nDotH, material.Shininess) // Multiply shininess by 4 to keep "compatibility" with values tuned for the standard Phong model

		// 	result = result.Add(material.ReflectColor.Mul(material.Specular * f)) // Add specular component
		// }

		// ...else specular is black

		result = result.Blend(lightIntensity)
	}
	// ...else light is on the other side of the surface: diffuse and specular are both black

	return result
}

// Lighten is used only for tests, it builds a dummy IntersectionInfo object then calls LightenHit to get the color
func Lighten(light *PointLight, object Hittable, point, eyev, normalv Tuple, shadowed bool) (Color, Color) {
	if shadowed {
		return Black, Black
	}

	ii := IntersectionInfo{Intersection: Intersection{O: object}, Point: point, Eyev: eyev, Normalv: normalv}
	object.Material().GetParamsAt(&ii)

	lightv := light.Pos.Sub(ii.Point).Normalize() // Direction to the light source

	return LightenHit(lightv, light.Intensity, &ii), ii.Mat.DiffuseColor
}

func NewPointLight(pos Tuple, intensity Color) *PointLight {
	return &PointLight{pos, intensity}
}

func (light *PointLight) LightenHit(ii *IntersectionInfo, rt *Raytracer) (result Color) {
	if !IsShadowed(light.Pos, rt, ii.OverPoint) {
		lightv := light.Pos.Sub(ii.Point).Normalize() // Direction to the light source
		result = LightenHit(lightv, light.Intensity, ii)
	}

	return
}

func angleToPoint(p1, p2, point Tuple) float64 {
	v1 := p1.Sub(point)
	v2 := p2.Sub(point)
	cosa := v1.DotProduct(v2) / (v1.Length() * v2.Length())
	return math.Acos(cosa)
}

func cosAngleToPoint(p1, p2, point Tuple) float64 {
	v1 := p1.Sub(point)
	v2 := p2.Sub(point)
	cosa := v1.DotProduct(v2) / (v1.Length() * v2.Length())
	return cosa
}

func cosAngleToPoint2(v1, p, point Tuple) float64 {
	v2 := point.Sub(p)
	cosa := v1.DotProduct(v2) / (v1.Length() * v2.Length())
	return cosa
}

func AnalyzeRectLight(light *RectLight, point Tuple) {
	p1 := light.Pos
	p2 := p1.Add(light.Uv)
	p3 := p1.Add(light.Vv)
	p4 := p2.Add(light.Vv)

	a1 := angleToPoint(p1, p4, point)
	a2 := angleToPoint(p2, p3, point)

	// vp := p1.Sub(point)
	// cosn := n.DotProduct(vp) / (vp.Length() * n.Length())

	Debugf("a1=%.2f, a2=%.2f\n", a1, a2)
}

func NewRectLight(intensity Color) *RectLight {
	return &RectLight{Pos: Point(0, 0, 0), Uv: Vector(1, 0, 0), Vv: Vector(0, 1, 0), Intensity: intensity}
}

func (light *RectLight) SetParams(pos, uvec, vvec Tuple) {
	light.Pos = pos
	light.Uv = uvec
	light.Vv = vvec
}

func (light *RectLight) SetSize(usize, vsize float64) {
	light.Uv = light.Uv.Normalize().Mul(usize)
	light.Vv = light.Vv.Normalize().Mul(vsize)
}

func (light *RectLight) SetDirection(pos, target Tuple) {
	// Compute the normal to the desired area light plane
	N := target.Sub(pos).Normalize()

	// Compute the tangent and bitangent vectors for the normal,
	// see: https://computergraphics.stackexchange.com/questions/5498/compute-sphere-tangent-for-normal-mapping
	A := Vector(0, 1, 0)
	T := A.CrossProduct(N).Normalize()
	B := T.CrossProduct(N)

	// Replace the light area vectors, but preserve the length
	light.Uv = T.Mul(light.Uv.Length())
	light.Vv = B.Mul(light.Vv.Length())

	// Set the corner so that the specified position falls in the center of the area
	light.Pos = pos.Sub(light.Uv.Mul(0.5)).Sub(light.Vv.Mul(0.5))
}

const LSamples float64 = 16

func (light *RectLight) LightenHit(ii *IntersectionInfo, rt *Raytracer) (result Color) {
	numSamplesU := LSamples
	numSamplesV := LSamples

	p1 := light.Pos
	p2 := light.Pos
	p3 := p1.Add(light.Uv)
	p4 := p2.Add(light.Vv)

	chits := 0
	if IsShadowed(p1, rt, ii.OverPoint) {
		chits++
	}
	if IsShadowed(p2, rt, ii.OverPoint) {
		chits++
	}
	if IsShadowed(p3, rt, ii.OverPoint) {
		chits++
	}
	if IsShadowed(p4, rt, ii.OverPoint) {
		chits++
	}

	if chits == 0 || chits == 4 {
		a1 := cosAngleToPoint(p1, p3, ii.Point)
		a2 := cosAngleToPoint(p2, p4, ii.Point)

		const apow = 7 // Seems to work with 6 too, this is a bit safer

		a1 = math.Abs(a1)
		a1 = 1 - math.Pow(a1, apow)
		numSamplesU *= a1
		numSamplesU = math.Round(numSamplesU + 0.5)

		a2 = math.Abs(a2)
		a2 = 1 - math.Pow(a2, apow)
		numSamplesV *= a2
		numSamplesV = math.Round(numSamplesV + 0.5)

		usize := 1.0 / math.Max(4, math.Min(numSamplesU, 8))
		vsize := 1.0 / math.Max(4, math.Min(numSamplesV, 8))

		hits := 0
		done := 0

		for u := Epsilon; u < 1; u += usize {
			for v := Epsilon; v < 1; v += vsize {
				pos := light.Pos.Add(light.Uv.Mul(u + rt.rand()*usize)).Add(light.Vv.Mul(v + rt.rand()*vsize))

				done++

				if !IsShadowed(pos, rt, ii.OverPoint) {
					hits++
				}
			}
		}

		if hits == done && chits == 0 {
			for u := Epsilon; u < 1; u += usize {
				for v := Epsilon; v < 1; v += vsize {
					pos := light.Pos.Add(light.Uv.Mul(u + rt.rand()*usize)).Add(light.Vv.Mul(v + rt.rand()*vsize))

					lightv := pos.Sub(ii.Point).Normalize() // Direction to the light source
					result = result.Add(LightenHit(lightv, light.Intensity, ii))
				}
			}
			return result.Mul(1 / float64(hits))
		}

		if hits == 0 {
			return Black
		}
	}

	result = Black
	numSamplesU = 16
	numSamplesV = 16
	usize := 1 / numSamplesU
	vsize := 1 / numSamplesV

	for u := Epsilon; u < 1; u += usize {
		for v := Epsilon; v < 1; v += vsize {
			pos := light.Pos.Add(light.Uv.Mul(u + rt.rand()*usize)).Add(light.Vv.Mul(v + rt.rand()*vsize))

			if !IsShadowed(pos, rt, ii.OverPoint) {
				lightv := pos.Sub(ii.Point).Normalize() // Direction to the light source
				result = result.Add(LightenHit(lightv, light.Intensity, ii))
			}
		}
	}

	return result.Mul(1 / (numSamplesU * numSamplesV))
}

func NewDirectionalLight(dir Tuple, intensity Color) *DirectionalLight {
	return &DirectionalLight{dir.Normalize().Neg(), intensity}
}

func (light *DirectionalLight) IsShadowed(rt *Raytracer, point Tuple) bool {
	ray := NewRay(point, light.Dir)

	hit := rt.HitForShadow(ray)

	return hit.Valid()
}

func (light *DirectionalLight) LightenHit(ii *IntersectionInfo, rt *Raytracer) (result Color) {
	if !light.IsShadowed(rt, ii.OverPoint) {
		result = LightenHit(light.Dir, light.Intensity, ii)
	}

	return
}

func NewSpotLight(pos, target Tuple, angleMin, angleMax float64, intensity Color) *SpotLight {
	return &SpotLight{pos, target.Sub(pos).Normalize(), angleMin, angleMax, intensity}
}

func (light *SpotLight) LightenHit(ii *IntersectionInfo, rt *Raytracer) (result Color) {
	if !IsShadowed(light.Pos, rt, ii.OverPoint) {
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
