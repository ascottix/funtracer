// Copyright (c) 2019 Alessandro Scotti
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"time"
)

func fail(err error) {
	fmt.Fprintf(os.Stderr, "*** Error: %s\n", err.Error())
}

func exists(filename string) bool {
	_, err := os.Stat(filename)

	return !os.IsNotExist(err)
}

func main() {
	options := NewOptions()

	// Load default options from JSON configuration file
	options.LoadFromJSON("config.json")

	// Allow user to override options using the command line
	options.InitFlags()

	flag.Parse()

	// Get name of scene file
	sceneFilename := flag.Arg(0)

	if sceneFilename == "" {
		// Use a default file if possible
		if exists(DefaultSceneFileName) {
			sceneFilename = DefaultSceneFileName
		} else {
			if len(os.Args) == 1 {
				fmt.Fprintln(os.Stderr, "Funtracer is a simple raytracer written in Go.")
			} else {
				fail(errors.New("scene file not specified"))
			}
			fmt.Fprintf(os.Stderr, "\nUsage: funtracer [options] scenefile[.fun]\n")
			flag.PrintDefaults()
			os.Exit(1)
		}
	}

	if !exists(sceneFilename) && exists(sceneFilename+DefaultSceneFileExt) {
		// Supply an extension if missing
		sceneFilename = sceneFilename + DefaultSceneFileExt
	}

	// Read scene and render!
	scene, err := ParseSbtSceneFromFile(sceneFilename)

	start := time.Now()

	if err == nil {
		scene.SyncOptions(options) // Sync options from outside with options from command line

		fmt.Printf("Options: %+v\n", *options)

		fmt.Printf("Rendering '%s' into '%s'...", sceneFilename, options.OutFilename)

		err = scene.World.RenderToPNG(scene.Camera, options.OutFilename)
	}

	elapsed := time.Now().Sub(start)

	if err != nil {
		fail(err)
		os.Exit(1)
	}

	fmt.Printf(" done in %s\n", elapsed.Round(time.Millisecond))
}
