// Copyright (c) 2018 Alessandro Scotti
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

type Scene struct {
	Name   string
	World  *World
	Camera *Camera
}

func NewScene() *Scene {
	return &Scene{
		"",
		NewWorld(),
		nil,
	}
}

func (s *Scene) SyncOptions(options *Options) {
	s.World.SetOptions(options)

	if s.Camera == nil {
		s.Camera = NewCamera(0, 0, Pi/2)
	}

	// Sync camera view size: user-specified has precedence, then scene file, then hardcoded defaults
	w := options.OutWidth
	if w == 0 {
		w = s.Camera.HSize
		if w == 0 {
			w = 400
		}
	}

	h := options.OutHeight
	if h == 0 {
		h = s.Camera.VSize
		if h == 0 {
			h = 225
		}
	}

	s.Camera.SetViewSize(w, h)

	options.OutWidth = w
	options.OutHeight = h
}
