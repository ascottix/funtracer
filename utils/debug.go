// Copyright (c) 2018 Alessandro Scotti
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package utils

import (
	"fmt"
)

func Debugf(format string, a ...interface{}) (int, error) {
	return fmt.Printf(format, a...)
}

func Debugln(a ...interface{}) (int, error) {
	return fmt.Println(a...)
}
