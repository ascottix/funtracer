// Copyright (c) 2018 Alessandro Scotti
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import ()

type Matrix struct {
	Cols int
	Rows int
	A    []float64
}

func NewMatrix(cols, rows int, data ...float64) Matrix {
	s := cols * rows

	m := Matrix{cols, rows, make([]float64, s, s)}
	copy(m.A, data)

	return m
}

func NewIdentityMatrix4x4() Matrix {
	return NewMatrix(4, 4,
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1)
}

func (m Matrix) At(row, col int) float64 {
	return m.A[col+row*m.Cols]
}

func (m Matrix) SetAt(row, col int, val float64) {
	m.A[col+row*m.Cols] = val
}

func (m Matrix) Equals(o Matrix) bool {
	return SliceFloatEqual(m.A, o.A)
}

func (m Matrix) Mul(o Matrix) Matrix {
	// This is only implemented for 4x4 matrices!
	r := NewMatrix(4, 4)

	for row := 0; row < 4; row++ {
		for col := 0; col < 4; col++ {
			val := m.At(row, 0)*o.At(0, col) + m.At(row, 1)*o.At(1, col) + m.At(row, 2)*o.At(2, col) + m.At(row, 3)*o.At(3, col)
			r.SetAt(row, col, val)
		}
	}

	return r
}

func (m Matrix) MulT(t Tuple) Tuple {
	// This is only implemented for 4x4 matrices!
	return Tuple{
		m.A[0]*t.X + m.A[1]*t.Y + m.A[2]*t.Z + m.A[3]*t.W,
		m.A[4]*t.X + m.A[5]*t.Y + m.A[6]*t.Z + m.A[7]*t.W,
		m.A[8]*t.X + m.A[9]*t.Y + m.A[10]*t.Z + m.A[11]*t.W,
		t.W,
	}

	// Uncompressed version:
	// var r Tuple
	// r.X = m.At(0, 0)*t.X + m.At(0, 1)*t.Y + m.At(0, 2)*t.Z + m.At(0, 3)*t.W
	// r.Y = m.At(1, 0)*t.X + m.At(1, 1)*t.Y + m.At(1, 2)*t.Z + m.At(1, 3)*t.W
	// r.Z = m.At(2, 0)*t.X + m.At(2, 1)*t.Y + m.At(2, 2)*t.Z + m.At(2, 3)*t.W
	// r.W = m.At(3, 0)*t.X + m.At(3, 1)*t.Y + m.At(3, 2)*t.Z + m.At(3, 3)*t.W
	// return r
}

func (m Matrix) Transpose() Matrix {
	t := NewMatrix(m.Rows, m.Cols)

	for row := 0; row < m.Rows; row++ {
		for col := 0; col < m.Cols; col++ {
			t.SetAt(col, row, m.At(row, col))
		}
	}

	return t
}

func (m Matrix) Determinant() (d float64) {
	if m.Cols == 2 {
		d = m.A[0]*m.A[3] - m.A[1]*m.A[2]
	} else {
		for c := 0; c < m.Cols; c++ {
			d += m.At(0, c) * m.Cofactor(0, c)
		}
	}

	return
}

// Submatrix returns the (sub)matrix obtained by removing the specified column and row
func (m Matrix) Submatrix(row, col int) Matrix {
	s := NewMatrix(m.Cols-1, m.Rows-1)

	// Top left
	for c := 0; c < col; c++ {
		for r := 0; r < row; r++ {
			s.SetAt(r, c, m.At(r, c))
		}
	}

	// Top right
	for c := col + 1; c < m.Cols; c++ {
		for r := 0; r < row; r++ {
			s.SetAt(r, c-1, m.At(r, c))
		}
	}

	// Bottom left
	for c := 0; c < col; c++ {
		for r := row + 1; r < m.Rows; r++ {
			s.SetAt(r-1, c, m.At(r, c))
		}
	}

	// Bottom right
	for c := col + 1; c < m.Cols; c++ {
		for r := row + 1; r < m.Rows; r++ {
			s.SetAt(r-1, c-1, m.At(r, c))
		}
	}

	return s
}

func (m Matrix) Minor(row, col int) float64 {
	return m.Submatrix(row, col).Determinant()
}

func (m Matrix) Cofactor(row, col int) float64 {
	r := m.Minor(row, col)

	if (row+col)%2 != 0 {
		r = -r
	}

	return r
}

func (m Matrix) IsInvertible() bool {
	return m.Determinant() != 0
}

func (m Matrix) Inverse() Matrix {
	d := m.Determinant()

	r := NewMatrix(m.Rows, m.Cols)

	for row := 0; row < m.Rows; row++ {
		for col := 0; col < m.Cols; col++ {
			r.SetAt(col, row, m.Cofactor(row, col)/d) // Note: we set (col,row) instead of (row,col) to transpose the result
		}
	}

	return r
}
