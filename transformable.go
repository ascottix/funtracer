// Copyright (c) 2018 Alessandro Scotti
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

type Transformable interface {
	Transform() Matrix
	InverseTransform() Matrix
	SetTransform(transforms ...Matrix)
}

type Transformer struct {
	transform Matrix
	Tinverse  Matrix
	TinverseT Matrix
}

func (s *Transformer) SetTransform(transforms ...Matrix) {
	s.transform = Identity()

	for _, t := range transforms {
		s.transform = s.transform.Mul(t)
	}

	s.Tinverse = s.transform.Inverse() // Precompute inverse transform for performance (this helps _a lot_)
	s.TinverseT = s.Tinverse.Transpose()
}

func (s *Transformer) Transform() Matrix {
	return s.transform
}

func (s *Transformer) InverseTransform() Matrix {
	return s.Tinverse
}
