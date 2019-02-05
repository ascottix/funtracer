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

// LightenHitWithJitteredStratified samples the area with a jittered stratified sampler
func (light *RectLight) LightenHitWithJitteredStratified(ii *IntersectionInfo, rt *Raytracer, samples float64) (result Color) {
	usamples, vsamples := samples, samples

	usize := 1 / usamples
	vsize := 1 / vsamples

	for u := Epsilon; u < 1; u += usize {
		for v := Epsilon; v < 1; v += vsize {
			pos := light.Pos.Add(light.Uv.Mul(u + rt.rand()*usize)).Add(light.Vv.Mul(v + rt.rand()*vsize))

			if !IsShadowed(pos, rt, ii.OverPoint) {
				lightv := pos.Sub(ii.Point).Normalize() // Direction to the light source
				result = result.Add(LightenHit(lightv, light.Intensity, ii))
			}
		}
	}

	return result.Mul(usize * vsize)
}

// LightenHitWithAdaptiveSampling recursively splits the area in two parts
// and tries to spend more samples in "difficult" zones where the variance is higher:
// it usually gives acceptable results very quickly and good results quite
// faster than the jittered sampler anyway
func (light *RectLight) LightenHitWithAdaptiveSampling(ii *IntersectionInfo, rt *Raytracer, minDepth, maxDepth int) (result Color) {
	// This evaluates a point light placed at position (u,v)
	// of the light surface and cast on the intersection point
	sample := func(u, v float64) Color {
		pos := light.Pos.Add(light.Uv.Mul(u)).Add(light.Vv.Mul(v))

		if IsShadowed(pos, rt, ii.OverPoint) {
			return Black
		}

		lightv := pos.Sub(ii.Point).Normalize() // Direction to the light source

		return LightenHit(lightv, light.Intensity, ii)
	}

	var estimateArea func(u, v, w, h float64, p0, p1, p2, p3 Color, depth int) Color

	estimateArea = func(u, v, w, h float64, p0, p1, p2, p3 Color, depth int) Color {
		// We use IsBlack() as a cheap IsShadowed() here
		fs := p0.IsBlack() + p1.IsBlack() + p2.IsBlack() + p3.IsBlack()

		// Interrupt the recursion if either at max depth or past the minimum depth with all samples in agreement
		if depth >= maxDepth || ((fs == 0 || fs == 4) && depth >= minDepth) {
			c := p0.Add(p1).Add(p2).Add(p3)

			// If there is already a majority, return now: this provided a small but measurable improvement in my tests
			if fs != 2 {
				return c.Mul(w * h / 4)
			}

			// Since the jury is split, get one last sample in the middle and make it count:
			// this is another idea that seemed to work quite well in tests
			c = c.Add(sample(u+w/2, v+h/2).Mul(4))

			return c.Mul(w * h / 8)
		} else {
			// Need more information, split the rectangle at the longest edge
			if w > h || (w == h && p0.IsBlack() == p2.IsBlack()) {
				w = w / 2
				pa := sample(u+w, v)
				pe := sample(u+w, v+h)
				c1 := estimateArea(u, v, w, h, p0, pa, p2, pe, depth+1)
				c2 := estimateArea(u+w, v, w, h, pa, p1, pe, p3, depth+1)
				return c1.Add(c2)
			} else {
				h = h / 2
				pb := sample(u, v+h)
				pd := sample(u+w, v+h)
				c3 := estimateArea(u, v, w, h, p0, p1, pb, pd, depth+1)
				c4 := estimateArea(u, v+h, w, h, pb, pd, p2, p3, depth+1)
				return c3.Add(c4)
			}
		}
	}

	result = estimateArea(0, 0, 1, 1, sample(0, 0), sample(0, 1), sample(1, 0), sample(1, 1), 0)

	return result
}

func (light *RectLight) LightenHit(ii *IntersectionInfo, rt *Raytracer) (result Color) {
	if false {
		return light.LightenHitWithJitteredStratified(ii, rt, float64(rt.world.Options.AreaLightSamples))
	} else {
		return light.LightenHitWithAdaptiveSampling(ii, rt, 5, 9)
	}
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
