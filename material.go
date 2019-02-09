// Copyright (c) 2018 Alessandro Scotti
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

type MaterialParams struct {
	DiffuseColor Color // Diffuse color
}

type Material struct {
	MaterialParams
	Texture      Texture
	Ambient      float64 // Ka
	Diffuse      float64 // Kd
	Roughness    float64 // Diffuse roughness (i.e. sigma in the Oren-Nayar model)
	Specular     float64
	Shininess    float64 // Ns, specular exponent
	Reflect      Color
	Refract      Color
	Ior          float64 // Ni, index of refraction
	ReflectColor Color
	ReflectLevel float64
	RefractColor Color
	RefractLevel float64
}

func NewMaterial() *Material {
	m := Material{
		MaterialParams: MaterialParams{
			DiffuseColor: White,
		},
		Ambient:      0.1,
		Diffuse:      0.9,
		Specular:     0.9,
		Shininess:    200.0,
		ReflectLevel: 0.0,
		ReflectColor: White,
		Ior:          1.0, // Index of refraction
	}

	m.SetReflect(0, White)
	m.SetRefract(0, White)

	return &m
}

func (m *Material) GetParamsAt(ii *IntersectionInfo) {
	ii.Mat = m.MaterialParams

	if m.Texture != nil {
		m.Texture.ApplyAtHit(ii)
	}
}

func (m *Material) SetDiffuseColor(c Color) *Material {
	m.DiffuseColor = c
	return m
}

func (m *Material) SetAmbient(v float64) *Material {
	// Ambient represents how much of the ambient light is reflected by the material
	m.Ambient = v
	return m
}

func (m *Material) SetDiffuse(v float64) *Material {
	// Diffuse represents how much light is scattered by the material,
	// it is used to compute the diffuse surface color
	m.Diffuse = v
	return m
}

func (m *Material) SetPattern(p Texture) *Material {
	// Pattern describes the appearance of the material
	m.Texture = p
	return m
}

func (m *Material) SetReflect(level float64, color Color) *Material {
	// Reflect level and color control how light is reflected by the material, for example
	// mirrors would have the highest reflectivity (1.0) and white color,
	// chrome gold would have high reflectivity (e.g. 0.95) and a goldish color (e.g. "ffc360")
	m.ReflectLevel = level
	m.ReflectColor = color
	m.Reflect = color.Mul(level)
	return m
}

func (m *Material) SetRefract(level float64, color Color) *Material {
	// Refract level and color control how light is refracted by the material,
	// the refract level goes from 0 (fully opaque) to 1 (fully transparent)
	m.RefractLevel = level
	m.RefractColor = color
	m.Refract = color.Mul(level)
	return m
}

func (m *Material) SetIor(v float64) *Material {
	m.Ior = v
	return m
}

func (m *Material) SetShininess(v float64) *Material {
	// Shininess controls the area where the specular highlight is distributed,
	// higher values yield a smaller area and a harder effect
	m.Shininess = v
	return m
}

func (m *Material) SetSpecular(v float64) *Material {
	// Specular controls the intensity of the light that is reflected directly towards the eye
	m.Specular = v
	return m
}

func (m *Material) ProxifyPatterns(g Patternable) *Material {
	// TODO: should we clone the material?!

	if m.Texture != nil {
		if t, ok := m.Texture.(Pattern); ok {
			p := NewProxyPattern(g, t)
			m.SetPattern(p)
		}
	}

	return m
}

// Compatibility with old API

func (m *Material) SetReflective(v float64) *Material {
	return m.SetReflect(v, White)
}

func (m *Material) SetRefractive(t, r float64) *Material {
	m.Ior = r
	return m.SetRefract(t, White)
}
