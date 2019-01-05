// Copyright (c) 2018 Alessandro Scotti
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"math"
)

type Light interface {
	Lighten(objectColor Color, object Hittable, point, eyev, normalv Tuple, shadowed bool) Color
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

// PhongLight implements the Phong model of light, i.e. ambient (computed elsewhere) + diffuse + specular
func PhongLight(lightv Tuple, lightIntensity Color, objectColor Color, object Hittable, point, eyev, normalv Tuple) (result Color) {
	if cosTheta := lightv.DotProduct(normalv); cosTheta >= 0 { // Cosine of angle between light vector and surface normal
		// Light is on the same side of the surface, need to compute both diffuse and specular
		material := object.Material()
		effectiveColor := objectColor.Blend(lightIntensity) // Combine light and surface colors

		// The energy of the light hitting the surface depends on the cosine of the angle
		// between the light incident direction and the surface normal (Lambert's cosine law),
		// which is then multiplied by the light intensity and a material-dependent parameter
		result = result.Add(effectiveColor.Mul(material.Diffuse * cosTheta))

		// The Phong model accounts for light that may be reflected directly towards the eye,
		// controlled by Specular (intensity of reflected light) and Shininess
		// (size of reflecting area, higher values yield a smaller area with harder reflection)
		reflectv := lightv.Neg().Reflect(normalv)

		if reflectDotEye := reflectv.DotProduct(eyev); reflectDotEye > 0 {
			f := math.Pow(reflectDotEye, material.Shininess)
			result = result.Add(lightIntensity.Mul(material.Specular * f)) // Add specular component
		}
		// ...else specular is black
	}
	// ...else light is on the other side of the surface: diffuse and specular are both black

	return
}

func NewPointLight(pos Tuple, intensity Color) *PointLight {
	return &PointLight{pos, intensity}
}

func (light *PointLight) Lighten(objectColor Color, object Hittable, point, eyev, normalv Tuple, shadowed bool) (result Color) {
	if !shadowed {
		lightv := light.Pos.Sub(point).Normalize() // Direction to the light source
		result = PhongLight(lightv, light.Intensity, objectColor, object, point, eyev, normalv)
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

func (light *DirectionalLight) Lighten(objectColor Color, object Hittable, point, eyev, normalv Tuple, shadowed bool) (result Color) {
	if !shadowed {
		result = PhongLight(light.Dir, light.Intensity, objectColor, object, point, eyev, normalv)
	}

	return
}

func NewSpotLight(pos, target Tuple, angleMin, angleMax float64, intensity Color) *SpotLight {
	return &SpotLight{pos, target.Sub(pos).Normalize(), angleMin, angleMax, intensity}
}

func (light *SpotLight) Lighten(objectColor Color, object Hittable, point, eyev, normalv Tuple, shadowed bool) (result Color) {
	if !shadowed {
		lightv := light.Pos.Sub(point).Normalize() // Direction to the light source

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

			result = PhongLight(lightv, light.Intensity.Mul(intensity), objectColor, object, point, eyev, normalv)
		}
	}

	return result
}

func (light *SpotLight) IsShadowed(rt *Raytracer, point Tuple) bool {
	return IsShadowed(light.Pos, rt, point)
}
