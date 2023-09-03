// Copyright (c) 2018 Alessandro Scotti
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package traits

import (
	"fmt"
)

type Namable interface {
	Name() string
	SetName(string)
}

type Namer struct {
	name string
}

var idcounter = 0 // Not thread-safe

func (n *Namer) Name() string {
	return n.name
}

func (n *Namer) SetName(name string) {
	n.name = name
}

func (n *Namer) SetNameForKind(kind string) {
	idcounter++
	n.name = fmt.Sprintf("%s_%d", kind, idcounter)
}
