// Copyright (c) 2019 Alessandro Scotti
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"flag"
	"os"
	"runtime"
)

const (
	DefaultSceneFileExt  = ".fun"
	DefaultSceneFileName = "have" + DefaultSceneFileExt
)

type Options struct {
	OutFilename     string `json:"o"`
	OutWidth        int    `json:"ow"`
	OutHeight       int    `json:"oh"`
	NumThreads      int    `json:"nt"`
	Supersampling   int    `json:"ss"`
	ReflectionDepth int    `json:"rd"`
}

func NewOptions() *Options {
	// Assign a default value to all options
	options := Options{
		OutFilename:     "fun.png",
		OutWidth:        0,
		OutHeight:       0,
		NumThreads:      runtime.GOMAXPROCS(0),
		Supersampling:   1,
		ReflectionDepth: 4,
	}

	return &options
}

func (options *Options) InitFlags() {
	flag.StringVar(&options.OutFilename, "o", options.OutFilename, "name of output image file")
	flag.IntVar(&options.OutWidth, "ow", options.OutWidth, "output image width")
	flag.IntVar(&options.OutHeight, "oh", options.OutHeight, "output image height")
	flag.IntVar(&options.NumThreads, "nt", options.NumThreads, "how many threads can be used for processing")
	flag.IntVar(&options.Supersampling, "ss", options.Supersampling, "supersampling level: each pixel is sampled n*n times")
	flag.IntVar(&options.ReflectionDepth, "rd", options.ReflectionDepth, "maximum depth of secondary rays")
}

func (options *Options) LoadFromJSON(filename string) {
	f, err := os.Open(filename)

	if err == nil {
		defer f.Close()

		dec := json.NewDecoder(f)

		dec.Decode(options)
	}
}
